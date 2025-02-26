package models

import (
	"time"

	"github.com/google/uuid"
)

type Photo struct {
	PhotoID    uuid.UUID `json:"photoId" validate:"required, uuid4"`
	Photo      string    `json:"photo" validate:"required"`
	UploadedAt time.Time `json:"uploadedAt" validate:"required"`
	UserId     uuid.UUID `json:"userId" validate:"required, uuid4"`
}
