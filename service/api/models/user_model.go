package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `json:"userId" validate:"required,uuid4"`
	Username       string    `json:"username" validate:"required,min=3,max=16"`
	ProfilePicture string    `json:"profilePicture,omitempty" validate:"omitempty,url,min=11,max=255"`
	CreatedAt      time.Time `json:"createdAt" validate:"required"`
}
