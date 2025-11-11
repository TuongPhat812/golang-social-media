package events

import (
	"time"

	"github.com/myself/golang-social-media/common/domain/message"
)

type ChatCreated struct {
	Message   message.Message
	CreatedAt time.Time
}
