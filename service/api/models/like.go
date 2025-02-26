package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Like struct {
	UserId  uuid.UUID `json:"userId" validate:"required, uuid4"`
	PhotoId uuid.UUID `json:"photoId" validate:"required, uuid4"`
	LikedAt time.Time `json:"likedAt" validate:"required"`
}
