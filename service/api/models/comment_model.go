package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID          uuid.UUID `json:"commentId" validate:"required,uuid4"`
	MessageID   uuid.UUID `json:"messageId" validate:"required,uuid4"`
	Commenter   User      `json:"commenter" validate:"required"`
	Emoji       string    `json:"emoji" validate:"required,emoji"`
	CommentedAt time.Time `json:"commentedAt" validate:"required"`
}
