package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Ban struct {
	UserId       uuid.UUID `json:"userId" validate:"required, uuid4"`
	BannedUserId uuid.UUID `json:"bannedUserId" validate:"required, uuid4"`
	BannedAt     time.Time `json:"bannedAt" validate:"required"`
}
