package factories

import (
	"time"

	"golang-social-media/apps/chat-service/internal/domain/message"
	"github.com/google/uuid"
)

// MessageFactory creates Message entities with proper initialization
type MessageFactory struct{}

// NewMessageFactory creates a new MessageFactory
func NewMessageFactory() *MessageFactory {
	return &MessageFactory{}
}

// CreateMessage creates a new Message with proper initialization
// This factory encapsulates the complex creation logic
func (f *MessageFactory) CreateMessage(senderID, receiverID, content string) (*message.Message, error) {
	if senderID == "" {
		return nil, &MessageFactoryError{Message: "sender ID cannot be empty"}
	}
	if receiverID == "" {
		return nil, &MessageFactoryError{Message: "receiver ID cannot be empty"}
	}
	if content == "" {
		return nil, &MessageFactoryError{Message: "content cannot be empty"}
	}

	now := time.Now().UTC()
	msg := &message.Message{
		ID:         uuid.NewString(),
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		CreatedAt:  now,
	}

	// Validate the created message
	if err := msg.Validate(); err != nil {
		return nil, &MessageFactoryError{
			Message: "failed to validate message",
			Cause:   err,
		}
	}

	// Domain logic: create message (this adds domain events internally)
	msg.Create()

	return msg, nil
}

// MessageFactoryError represents an error in message factory
type MessageFactoryError struct {
	Message string
	Cause   error
}

func (e *MessageFactoryError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

func (e *MessageFactoryError) Unwrap() error {
	return e.Cause
}

