package services

import (
	"github.com/evaevangelisti/wasatext/service/api/models"
	"github.com/evaevangelisti/wasatext/service/api/repositories"
	"github.com/evaevangelisti/wasatext/service/utils"
	"github.com/evaevangelisti/wasatext/service/utils/errors"
	"github.com/evaevangelisti/wasatext/service/utils/globaltime"
	"github.com/google/uuid"
)

type ConversationService struct {
	Repository *repositories.ConversationRepository
}

func (service *ConversationService) GetConversationsByUserID(userID uuid.UUID) ([]models.Conversation, error) {
	return service.Repository.GetConversationsByUserID(userID)
}

func (service *ConversationService) GetConversationByID(conversationID, authenticatedUserID uuid.UUID) (models.Conversation, error) {
	messageRepository := &repositories.MessageRepository{Database: service.Repository.Database}

	messages, err := messageRepository.GetMessagesByConversationID(conversationID)
	if err != nil {
		return nil, err
	}

	readAt := globaltime.Now()

	for _, message := range messages {
		if _, alreadyRead := message.Trackings.Read[authenticatedUserID.String()]; !alreadyRead {
			err = messageRepository.AddMessageTracking(message.ID, authenticatedUserID, readAt)
			if err != nil {
				return nil, err
			}
		}
	}

	return service.Repository.GetConversationByID(conversationID)
}

func (service *ConversationService) CreatePrivateConversation(participantIDs []uuid.UUID) (*models.PrivateConversation, error) {
	if len(participantIDs) != 2 || participantIDs[0] == participantIDs[1] {
		return nil, errors.ErrBadRequest
	}

	existingConversation, err := service.Repository.GetPrivateConversationByParticipants(participantIDs)
	if err != nil {
		return nil, err
	}

	if existingConversation != nil {
		return existingConversation, errors.ErrConflict
	}

	conversationID, err := service.Repository.CreatePrivateConversation(participantIDs)
	if err != nil {
		return nil, err
	}

	conversation, err := service.Repository.GetConversationByID(conversationID)
	if err != nil {
		return nil, err
	}

	privateConversation, ok := conversation.(*models.PrivateConversation)
	if !ok {
		return nil, errors.ErrInternal
	}

	return privateConversation, nil
}

func (service *ConversationService) CreateGroupConversation(name string, memberIDs []uuid.UUID) (*models.GroupConversation, error) {
	if len(memberIDs) == 0 {
		return nil, errors.ErrBadRequest
	}

	conversationID, err := service.Repository.CreateGroupConversation(name, memberIDs)
	if err != nil {
		return nil, err
	}

	conversation, err := service.Repository.GetConversationByID(conversationID)
	if err != nil {
		return nil, err
	}

	groupConversation, ok := conversation.(*models.GroupConversation)
	if !ok {
		return nil, errors.ErrInternal
	}

	return groupConversation, nil
}

func (service *ConversationService) AddMember(conversationID, authenticatedUserID, userID uuid.UUID) (*models.GroupConversation, error) {
	conversation, err := service.Repository.GetConversationByID(conversationID)
	if err != nil {
		return nil, err
	}

	groupConversation, ok := conversation.(*models.GroupConversation)
	if !ok {
		return nil, errors.ErrBadRequest
	}

	if len(groupConversation.Members) >= 100 {

	}

	hasAccess, err := utils.IsUserInConversation(service.Repository, conversationID, authenticatedUserID)
	if err != nil {
		return nil, err
	}

	if !hasAccess {
		return nil, errors.ErrForbidden
	}

	for _, member := range groupConversation.Members {
		if member.ID == userID {
			return nil, errors.ErrConflict
		}
	}

	_, err = service.Repository.AddMember(conversationID, userID)
	if err != nil {
		return nil, err
	}

	updatedConversation, err := service.Repository.GetConversationByID(conversationID)
	if err != nil {
		return nil, err
	}

	updatedGroupConversation, ok := updatedConversation.(*models.GroupConversation)
	if !ok {
		return nil, errors.ErrInternal
	}

	return updatedGroupConversation, nil
}

func (service *ConversationService) UpdateGroupName(conversationID, authenticatedUserID uuid.UUID, name string) (*models.GroupConversation, error) {
	conversation, err := service.Repository.GetConversationByID(conversationID)
	if err != nil {
		return nil, err
	}

	_, ok := conversation.(*models.GroupConversation)
	if !ok {
		return nil, errors.ErrBadRequest
	}

	hasAccess, err := utils.IsUserInConversation(service.Repository, conversationID, authenticatedUserID)
	if err != nil {
		return nil, err
	}

	if !hasAccess {
		return nil, errors.ErrForbidden
	}

	err = service.Repository.UpdateGroupName(conversationID, name)
	if err != nil {
		return nil, err
	}

	updatedConversation, err := service.Repository.GetConversationByID(conversationID)
	if err != nil {
		return nil, err
	}

	updatedGroupConversation, ok := updatedConversation.(*models.GroupConversation)
	if !ok {
		return nil, errors.ErrInternal
	}

	return updatedGroupConversation, nil
}

func (service *ConversationService) UpdateGroupPhoto(conversationID, authenticatedUserID uuid.UUID, photo string) (*models.GroupConversation, error) {
	conversation, err := service.Repository.GetConversationByID(conversationID)
	if err != nil {
		return nil, err
	}

	_, ok := conversation.(*models.GroupConversation)
	if !ok {
		return nil, errors.ErrBadRequest
	}

	hasAccess, err := utils.IsUserInConversation(service.Repository, conversationID, authenticatedUserID)
	if err != nil {
		return nil, err
	}

	if !hasAccess {
		return nil, errors.ErrForbidden
	}

	err = service.Repository.UpdateGroupPhoto(conversationID, photo)
	if err != nil {
		return nil, err
	}

	updatedConversation, err := service.Repository.GetConversationByID(conversationID)
	if err != nil {
		return nil, err
	}

	updatedGroupConversation, ok := updatedConversation.(*models.GroupConversation)
	if !ok {
		return nil, errors.ErrInternal
	}

	return updatedGroupConversation, nil
}

func (service *ConversationService) RemoveMember(conversationID, userID uuid.UUID) error {
	hasAccess, err := utils.IsUserInConversation(service.Repository, conversationID, userID)
	if err != nil {
		return err
	}

	if !hasAccess {
		return errors.ErrNotFound
	}

	err = service.Repository.RemoveMember(conversationID, userID)
	if err != nil {
		return err
	}

	return nil
}
