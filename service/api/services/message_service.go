package services

import (
	"github.com/evaevangelisti/wasatext/service/api/models"
	"github.com/evaevangelisti/wasatext/service/api/repositories"
	"github.com/evaevangelisti/wasatext/service/utils/errors"
	"github.com/google/uuid"
)

type MessageService struct {
	Repository *repositories.MessageRepository
}

func (service *MessageService) CreateMessage(conversationID, userID uuid.UUID, content, attachment string, replyToMessageID uuid.UUID) (*models.Message, error) {
	conversationRepository := &repositories.ConversationRepository{Database: service.Repository.Database}

	hasAccess, err := conversationRepository.IsUserInConversation(conversationID, userID)
	if err != nil {
		return nil, err
	}

	if !hasAccess {
		return nil, errors.ErrForbidden
	}

	if content == "" && attachment == "" {
		return nil, errors.ErrBadRequest
	}

	if replyToMessageID != uuid.Nil {
		replyMessage, err := service.Repository.GetMessageByID(replyToMessageID)

		if err != nil {
			return nil, err
		}

		if replyMessage == nil || replyMessage.Sender.ID == uuid.Nil {
			return nil, errors.ErrBadRequest
		}
	}

	messageID, err := service.Repository.CreateMessage(conversationID, userID, content, attachment, replyToMessageID)
	if err != nil {
		return nil, err
	}

	message, err := service.Repository.GetMessageByID(messageID)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (service *MessageService) CreateForwardedMessage(conversationID, userID, originalMessageID uuid.UUID) (*models.Message, error) {
	conversationRepository := &repositories.ConversationRepository{Database: service.Repository.Database}

	hasAccess, err := conversationRepository.IsUserInConversation(conversationID, userID)
	if err != nil {
		return nil, err
	}

	if !hasAccess {
		return nil, errors.ErrForbidden
	}

	originalConversation, err := conversationRepository.GetConversationByMessageID(originalMessageID)
	if err != nil {
		return nil, err
	}

	if originalConversation == nil {
		return nil, errors.ErrNotFound
	}

	hasAccess, err = conversationRepository.IsUserInConversation(originalConversation.GetID(), userID)
	if err != nil {
		return nil, err
	}

	if !hasAccess {
		return nil, errors.ErrForbidden
	}

	messageID, err := service.Repository.CreateForwardedMessage(conversationID, userID, originalMessageID)
	if err != nil {
		return nil, err
	}

	message, err := service.Repository.GetMessageByID(messageID)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (service *MessageService) UpdateMessage(messageID, userID uuid.UUID, content string) (*models.Message, error) {
	conversationRepository := &repositories.ConversationRepository{Database: service.Repository.Database}

	conversation, err := conversationRepository.GetConversationByMessageID(messageID)
	if err != nil {
		return nil, err
	}

	if conversation == nil {
		return nil, errors.ErrNotFound
	}

	hasAccess, err := conversationRepository.IsUserInConversation(conversation.GetID(), userID)
	if err != nil {
		return nil, err
	}

	if !hasAccess {
		return nil, errors.ErrForbidden
	}

	message, err := service.Repository.GetMessageByID(messageID)
	if err != nil {
		return nil, err
	}

	if message == nil {
		return nil, errors.ErrNotFound
	}

	if message.IsForwarded {
		return nil, errors.ErrBadRequest
	}

	if message.Sender.ID != userID {
		return nil, errors.ErrForbidden
	}

	if content == "" {
		return nil, errors.ErrBadRequest
	}

	err = service.Repository.UpdateMessage(messageID, content)
	if err != nil {
		return nil, err
	}

	updatedMessage, err := service.Repository.GetMessageByID(messageID)
	if err != nil {
		return nil, err
	}

	return updatedMessage, nil
}

func (service *MessageService) DeleteMessage(messageID, userID uuid.UUID) error {
	conversationRepository := &repositories.ConversationRepository{Database: service.Repository.Database}

	conversation, err := conversationRepository.GetConversationByMessageID(messageID)
	if err != nil {
		return err
	}

	if conversation == nil {
		return errors.ErrNotFound
	}

	hasAccess, err := conversationRepository.IsUserInConversation(conversation.GetID(), userID)
	if err != nil {
		return err
	}

	if !hasAccess {
		return errors.ErrForbidden
	}

	message, err := service.Repository.GetMessageByID(messageID)
	if err != nil {
		return err
	}

	if message == nil {
		return errors.ErrNotFound
	}

	if message.Sender.ID != userID {
		return errors.ErrForbidden
	}

	return service.Repository.DeleteMessage(messageID)
}
