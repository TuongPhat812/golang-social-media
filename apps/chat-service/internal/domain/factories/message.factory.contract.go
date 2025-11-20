package factories

import "golang-social-media/apps/chat-service/internal/domain/message"

// MessageFactory defines the contract for creating Message entities
type MessageFactory interface {
	CreateMessage(senderID, receiverID, content string) (*message.Message, error)
}


