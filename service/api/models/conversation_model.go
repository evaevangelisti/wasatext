package models

import (
	"time"

	"github.com/google/uuid"
)

type Conversation interface {
	GetID() uuid.UUID
	GetType() string
}

type PrivateConversation struct {
	ID           uuid.UUID `json:"conversationId" validate:"required,uuid4"`
	Type         string    `json:"type" validate:"required,oneof=private group"`
	Participants []User    `json:"participants" validate:"required,min=2,max=2"`
	Messages     []Message `json:"messages,omitempty" validate:"omitempty,max=1000"`
	CreatedAt    time.Time `json:"createdAt" validate:"required"`
}

func (conversation *PrivateConversation) GetID() uuid.UUID { return conversation.ID }
func (conversation *PrivateConversation) GetType() string  { return conversation.Type }

type GroupConversation struct {
	ID        uuid.UUID `json:"conversationId" validate:"required,uuid4"`
	Type      string    `json:"type" validate:"required,oneof=private group"`
	Name      string    `json:"name" validate:"required,min=1,max=50"`
	Photo     string    `json:"photo,omitempty" validate:"omitempty,url,min=11,max=255"`
	Members   []User    `json:"members" validate:"required,min=1,max=100"`
	Messages  []Message `json:"messages,omitempty" validate:"omitempty,max=1000"`
	CreatedAt time.Time `json:"createdAt" validate:"required"`
}

func (conversation *GroupConversation) GetID() uuid.UUID { return conversation.ID }
func (conversation *GroupConversation) GetType() string  { return conversation.Type }
