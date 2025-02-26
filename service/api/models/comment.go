package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Comment struct {
	CommentId   uuid.UUID `json:"commentId" validate:"required, uuid4"`
	Content     string    `json:"content" validate:"required, min=1, max=256"`
	CommentedAt time.Time `json:"commentedAt" validate:"required"`
	UserId      uuid.UUID `json:"userId" validate:"required, uuid4"`
	PhotoId     uuid.UUID `json:"photoId" validate:"required, uuid4"`
}
