package outbox

import (
	"context"
	"encoding/json"
	"time"

	"golang-social-media/apps/auth-service/internal/infrastructure/eventbus/publisher"
	"golang-social-media/apps/auth-service/internal/infrastructure/eventbus/publisher/contracts"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/postgres"
	"golang-social-media/pkg/logger"

	"github.com/rs/zerolog"
)

// Processor processes outbox events and publishes them
type Processor struct {
	outboxRepo *postgres.OutboxRepository
	publisher  *publisher.KafkaPublisher
	log        *zerolog.Logger
	batchSize  int
	interval   time.Duration
}

// NewProcessor creates a new OutboxProcessor
func NewProcessor(outboxRepo *postgres.OutboxRepository, kafkaPublisher *publisher.KafkaPublisher) *Processor {
	return &Processor{
		outboxRepo: outboxRepo,
		publisher:  kafkaPublisher,
		log:        logger.Component("auth.outbox.processor"),
		batchSize:  10,
		interval:   5 * time.Second,
	}
}

// Start starts processing outbox events
func (p *Processor) Start(ctx context.Context) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	p.log.Info().
		Int("batch_size", p.batchSize).
		Dur("interval", p.interval).
		Msg("outbox processor started")

	for {
		select {
		case <-ctx.Done():
			p.log.Info().Msg("outbox processor stopped")
			return
		case <-ticker.C:
			if err := p.processBatch(ctx); err != nil {
				p.log.Error().
					Err(err).
					Msg("failed to process outbox batch")
			}
		}
	}
}

// processBatch processes a batch of pending events
func (p *Processor) processBatch(ctx context.Context) error {
	events, err := p.outboxRepo.GetPendingEvents(ctx, p.batchSize)
	if err != nil {
		return err
	}

	if len(events) == 0 {
		return nil
	}

	p.log.Debug().
		Int("event_count", len(events)).
		Msg("processing outbox batch")

	for _, event := range events {
		if err := p.processEvent(ctx, event); err != nil {
			p.log.Error().
				Err(err).
				Str("event_id", event.ID).
				Str("event_type", event.EventType).
				Msg("failed to process event")
			// Mark as failed
			if markErr := p.outboxRepo.MarkAsFailed(ctx, event.ID, err.Error()); markErr != nil {
				p.log.Error().
					Err(markErr).
					Str("event_id", event.ID).
					Msg("failed to mark event as failed")
			}
		}
	}

	return nil
}

// processEvent processes a single outbox event
func (p *Processor) processEvent(ctx context.Context, event postgres.OutboxModel) error {
	// Parse payload
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(event.Payload), &payload); err != nil {
		return err
	}

	// Publish to Kafka using the appropriate publisher method
	// Convert payload map to contract struct
	var publishErr error
	switch event.EventType {
	case "UserCreated":
		var eventContract contracts.UserCreatedPayload
		if err := mapToStruct(payload, &eventContract); err != nil {
			return err
		}
		publishErr = p.publisher.PublishUserCreated(ctx, eventContract)
	default:
		p.log.Warn().
			Str("event_type", event.EventType).
			Msg("unknown event type, skipping")
		return nil
	}

	if publishErr != nil {
		// Increment retry count
		if retryErr := p.outboxRepo.IncrementRetry(ctx, event.ID); retryErr != nil {
			p.log.Error().
				Err(retryErr).
				Str("event_id", event.ID).
				Msg("failed to increment retry count")
		}
		return publishErr
	}

	// Mark as published
	if err := p.outboxRepo.MarkAsPublished(ctx, event.ID); err != nil {
		return err
	}

	p.log.Debug().
		Str("event_id", event.ID).
		Str("event_type", event.EventType).
		Msg("event published successfully")

	return nil
}

// mapToStruct converts a map to a struct
func mapToStruct(m map[string]interface{}, v interface{}) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

