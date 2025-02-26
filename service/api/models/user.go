package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	UserID    uuid.UUID `json:"userId" validate:"required, uuid4"`
	Username  string    `json:"username" validate:"required, min=3, max=16, alphanum"`
	CreatedAt time.Time `json:"createdAt" validate:"required"`
}
