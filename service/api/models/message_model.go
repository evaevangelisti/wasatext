package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID                uuid.UUID `json:"messageId" validate:"required,uuid4"`
	ConversationID    uuid.UUID `json:"conversationId" validate:"required,uuid4"`
	Sender            User      `json:"sender" validate:"required"`
	Content           string    `json:"content,omitempty" validate:"omitempty,min=1,max=1000"`
	Attachment        string    `json:"attachment,omitempty" validate:"omitempty,url,min=11,max=255"`
	Comments          []Comment `json:"comments,omitempty" validate:"omitempty,max=100"`
	IsForwarded       bool      `json:"isForwarded" validate:"required"`
	OriginalMessageID uuid.UUID `json:"originalMessageId,omitempty" validate:"omitempty,uuid4"`
	Trackings         struct {
		Read map[string]string `json:"read,omitempty" validate:"omitempty"`
	} `json:"trackings,omitempty" validate:"omitempty"`
	SentAt   time.Time `json:"sentAt" validate:"required"`
	EditedAt time.Time `json:"editedAt,omitempty" validate:"omitempty"`
}
