package services

import (
	"github.com/evaevangelisti/wasatext/service/api/models"
	"github.com/evaevangelisti/wasatext/service/api/repositories"
	"github.com/evaevangelisti/wasatext/service/utils"
	"github.com/evaevangelisti/wasatext/service/utils/errors"
	"github.com/google/uuid"
)

type MessageService struct {
	Repository *repositories.MessageRepository
}

func (service *MessageService) CreateMessage(conversationID, userID uuid.UUID, content, attachment string) (*models.Message, error) {
	conversationRepository := &repositories.ConversationRepository{Database: service.Repository.Database}

	hasAccess, err := utils.IsUserInConversation(conversationRepository, conversationID, userID)
	if err != nil {
		return nil, err
	}

	if !hasAccess {
		return nil, errors.ErrForbidden
	}

	if content == "" && attachment == "" {
		return nil, errors.ErrBadRequest
	}

	messageID, err := service.Repository.CreateMessage(conversationID, userID, content, attachment)
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

	hasAccess, err := utils.IsUserInConversation(conversationRepository, conversationID, userID)
	if err != nil {
		return nil, err
	}

	if !hasAccess {
		return nil, errors.ErrForbidden
	}

	messageRepository := &repositories.MessageRepository{Database: service.Repository.Database}

	originalConversationID, err := utils.GetConversationIDFromMessage(messageRepository, originalMessageID)
	if err != nil {
		return nil, err
	}

	hasAccess, err = utils.IsUserInConversation(conversationRepository, originalConversationID, userID)
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

func (service *MessageService) UpdateMessage(messageID, authenticatedUserID uuid.UUID, content string) (*models.Message, error) {
	messageRepository := &repositories.MessageRepository{Database: service.Repository.Database}

	conversationID, err := utils.GetConversationIDFromMessage(messageRepository, messageID)
	if err != nil {
		return nil, err
	}

	conversationRepository := &repositories.ConversationRepository{Database: service.Repository.Database}

	hasAccess, err := utils.IsUserInConversation(conversationRepository, conversationID, authenticatedUserID)
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

	if message.Sender.ID != authenticatedUserID {
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

func (service *MessageService) DeleteMessage(messageID, authenticatedUserID uuid.UUID) error {
	messageRepository := &repositories.MessageRepository{Database: service.Repository.Database}

	conversationID, err := utils.GetConversationIDFromMessage(messageRepository, messageID)
	if err != nil {
		return err
	}

	conversationRepository := &repositories.ConversationRepository{Database: service.Repository.Database}

	hasAccess, err := utils.IsUserInConversation(conversationRepository, conversationID, authenticatedUserID)
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

	if message.Sender.ID != authenticatedUserID {
		return errors.ErrForbidden
	}

	return service.Repository.DeleteMessage(messageID)
}
