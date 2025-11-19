package publisher

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/segmentio/kafka-go"
	"golang-social-media/apps/ecommerce-service/internal/infrastructure/eventbus/publisher/contracts"
	"golang-social-media/pkg/events"
	"golang-social-media/pkg/logger"
)

var _ contracts.EcommercePublisher = (*KafkaPublisher)(nil)

type KafkaPublisher struct {
	productCreatedWriter      *kafka.Writer
	productStockUpdatedWriter *kafka.Writer
	orderCreatedWriter        *kafka.Writer
	orderItemAddedWriter      *kafka.Writer
	orderConfirmedWriter      *kafka.Writer
	orderCancelledWriter      *kafka.Writer
}

func NewKafkaPublisher(brokers []string) (*KafkaPublisher, error) {
	if len(brokers) == 0 {
		return nil, errors.New("kafka brokers must be provided")
	}

	// Product Created Writer
	productCreatedWriter := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    events.TopicProductCreated,
		Balancer: &kafka.LeastBytes{},
		BatchSize:    100,
		BatchBytes:   1048576,
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kafka.RequireOne,
		MaxAttempts:  10,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		WriteBackoffMin: 100 * time.Millisecond,
		WriteBackoffMax: 1 * time.Second,
		Async:      true,
		Compression: kafka.Snappy,
	}

	// Product Stock Updated Writer
	productStockUpdatedWriter := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    events.TopicProductStockUpdated,
		Balancer: &kafka.LeastBytes{},
		BatchSize:    100,
		BatchBytes:   1048576,
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kafka.RequireOne,
		MaxAttempts:  10,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		WriteBackoffMin: 100 * time.Millisecond,
		WriteBackoffMax: 1 * time.Second,
		Async:      true,
		Compression: kafka.Snappy,
	}

	// Order Created Writer
	orderCreatedWriter := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    events.TopicOrderCreated,
		Balancer: &kafka.LeastBytes{},
		BatchSize:    100,
		BatchBytes:   1048576,
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kafka.RequireOne,
		MaxAttempts:  10,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		WriteBackoffMin: 100 * time.Millisecond,
		WriteBackoffMax: 1 * time.Second,
		Async:      true,
		Compression: kafka.Snappy,
	}

	// Order Item Added Writer
	orderItemAddedWriter := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    events.TopicOrderItemAdded,
		Balancer: &kafka.LeastBytes{},
		BatchSize:    100,
		BatchBytes:   1048576,
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kafka.RequireOne,
		MaxAttempts:  10,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		WriteBackoffMin: 100 * time.Millisecond,
		WriteBackoffMax: 1 * time.Second,
		Async:      true,
		Compression: kafka.Snappy,
	}

	// Order Confirmed Writer
	orderConfirmedWriter := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    events.TopicOrderConfirmed,
		Balancer: &kafka.LeastBytes{},
		BatchSize:    100,
		BatchBytes:   1048576,
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kafka.RequireOne,
		MaxAttempts:  10,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		WriteBackoffMin: 100 * time.Millisecond,
		WriteBackoffMax: 1 * time.Second,
		Async:      true,
		Compression: kafka.Snappy,
	}

	// Order Cancelled Writer
	orderCancelledWriter := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    events.TopicOrderCancelled,
		Balancer: &kafka.LeastBytes{},
		BatchSize:    100,
		BatchBytes:   1048576,
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kafka.RequireOne,
		MaxAttempts:  10,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		WriteBackoffMin: 100 * time.Millisecond,
		WriteBackoffMax: 1 * time.Second,
		Async:      true,
		Compression: kafka.Snappy,
	}

	logger.Component("ecommerce.publisher").
		Info().
		Strs("brokers", brokers).
		Msg("kafka publisher initialized")

	return &KafkaPublisher{
		productCreatedWriter:      productCreatedWriter,
		productStockUpdatedWriter: productStockUpdatedWriter,
		orderCreatedWriter:        orderCreatedWriter,
		orderItemAddedWriter:      orderItemAddedWriter,
		orderConfirmedWriter:      orderConfirmedWriter,
		orderCancelledWriter:      orderCancelledWriter,
	}, nil
}

func (p *KafkaPublisher) PublishProductCreated(ctx context.Context, event contracts.ProductCreated) error {
	payload, err := json.Marshal(event)
	if err != nil {
		logger.Component("ecommerce.publisher").
			Error().
			Err(err).
			Msg("failed to marshal ProductCreated event")
		return err
	}

	if err := p.productCreatedWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.ProductID),
		Value: payload,
	}); err != nil {
		logger.Component("ecommerce.publisher").
			Error().
			Err(err).
			Str("product_id", event.ProductID).
			Msg("failed to publish ProductCreated event")
		return err
	}

	logger.Component("ecommerce.publisher").
		Info().
		Str("product_id", event.ProductID).
		Str("topic", events.TopicProductCreated).
		Msg("published ProductCreated event")

	return nil
}

func (p *KafkaPublisher) PublishProductStockUpdated(ctx context.Context, event contracts.ProductStockUpdated) error {
	payload, err := json.Marshal(event)
	if err != nil {
		logger.Component("ecommerce.publisher").
			Error().
			Err(err).
			Msg("failed to marshal ProductStockUpdated event")
		return err
	}

	if err := p.productStockUpdatedWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.ProductID),
		Value: payload,
	}); err != nil {
		logger.Component("ecommerce.publisher").
			Error().
			Err(err).
			Str("product_id", event.ProductID).
			Msg("failed to publish ProductStockUpdated event")
		return err
	}

	logger.Component("ecommerce.publisher").
		Info().
		Str("product_id", event.ProductID).
		Str("topic", events.TopicProductStockUpdated).
		Msg("published ProductStockUpdated event")

	return nil
}

func (p *KafkaPublisher) PublishOrderCreated(ctx context.Context, event contracts.OrderCreated) error {
	payload, err := json.Marshal(event)
	if err != nil {
		logger.Component("ecommerce.publisher").
			Error().
			Err(err).
			Msg("failed to marshal OrderCreated event")
		return err
	}

	if err := p.orderCreatedWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.OrderID),
		Value: payload,
	}); err != nil {
		logger.Component("ecommerce.publisher").
			Error().
			Err(err).
			Str("order_id", event.OrderID).
			Msg("failed to publish OrderCreated event")
		return err
	}

	logger.Component("ecommerce.publisher").
		Info().
		Str("order_id", event.OrderID).
		Str("topic", events.TopicOrderCreated).
		Msg("published OrderCreated event")

	return nil
}

func (p *KafkaPublisher) PublishOrderItemAdded(ctx context.Context, event contracts.OrderItemAdded) error {
	payload, err := json.Marshal(event)
	if err != nil {
		logger.Component("ecommerce.publisher").
			Error().
			Err(err).
			Msg("failed to marshal OrderItemAdded event")
		return err
	}

	if err := p.orderItemAddedWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.OrderID),
		Value: payload,
	}); err != nil {
		logger.Component("ecommerce.publisher").
			Error().
			Err(err).
			Str("order_id", event.OrderID).
			Msg("failed to publish OrderItemAdded event")
		return err
	}

	logger.Component("ecommerce.publisher").
		Info().
		Str("order_id", event.OrderID).
		Str("topic", events.TopicOrderItemAdded).
		Msg("published OrderItemAdded event")

	return nil
}

func (p *KafkaPublisher) PublishOrderConfirmed(ctx context.Context, event contracts.OrderConfirmed) error {
	payload, err := json.Marshal(event)
	if err != nil {
		logger.Component("ecommerce.publisher").
			Error().
			Err(err).
			Msg("failed to marshal OrderConfirmed event")
		return err
	}

	if err := p.orderConfirmedWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.OrderID),
		Value: payload,
	}); err != nil {
		logger.Component("ecommerce.publisher").
			Error().
			Err(err).
			Str("order_id", event.OrderID).
			Msg("failed to publish OrderConfirmed event")
		return err
	}

	logger.Component("ecommerce.publisher").
		Info().
		Str("order_id", event.OrderID).
		Str("topic", events.TopicOrderConfirmed).
		Msg("published OrderConfirmed event")

	return nil
}

func (p *KafkaPublisher) PublishOrderCancelled(ctx context.Context, event contracts.OrderCancelled) error {
	payload, err := json.Marshal(event)
	if err != nil {
		logger.Component("ecommerce.publisher").
			Error().
			Err(err).
			Msg("failed to marshal OrderCancelled event")
		return err
	}

	if err := p.orderCancelledWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.OrderID),
		Value: payload,
	}); err != nil {
		logger.Component("ecommerce.publisher").
			Error().
			Err(err).
			Str("order_id", event.OrderID).
			Msg("failed to publish OrderCancelled event")
		return err
	}

	logger.Component("ecommerce.publisher").
		Info().
		Str("order_id", event.OrderID).
		Str("topic", events.TopicOrderCancelled).
		Msg("published OrderCancelled event")

	return nil
}

func (p *KafkaPublisher) Close() error {
	var errs []error

	if err := p.productCreatedWriter.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := p.productStockUpdatedWriter.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := p.orderCreatedWriter.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := p.orderItemAddedWriter.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := p.orderConfirmedWriter.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := p.orderCancelledWriter.Close(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errors.New("failed to close some kafka writers")
	}

	return nil
}

