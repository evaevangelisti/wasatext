package models

import (
	"time"

	"github.com/google/uuid"
)

type Follow struct {
	FollowerId  uuid.UUID `json:"followerId" validate:"required, uuid4"`
	FollowingId uuid.UUID `json:"followingId" validate:"required, uuid4"`
	FollowedAt  time.Time `json:"followedAt" validate:"required"`
}
