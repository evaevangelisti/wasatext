package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID          uuid.UUID `json:"commentId" validate:"required"`
	Commenter   User      `json:"commenter" validate:"required"`
	Emoji       string    `json:"emoji" validate:"required,min=1,max=16"`
	CommentedAt time.Time `json:"commentedAt" validate:"required"`
}
