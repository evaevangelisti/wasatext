package services

import (
	"github.com/evaevangelisti/wasatext/service/api/models"
	"github.com/evaevangelisti/wasatext/service/api/repositories"
	"github.com/evaevangelisti/wasatext/service/utils"
	"github.com/evaevangelisti/wasatext/service/utils/errors"
	"github.com/google/uuid"
)

type CommentService struct {
	Repository *repositories.CommentRepository
}

func (service *CommentService) CreateComment(messageID, userID uuid.UUID, emoji string) (*models.Comment, error) {
	messageRepository := &repositories.MessageRepository{Database: service.Repository.Database}

	conversationID, err := utils.GetConversationIDFromMessage(messageRepository, messageID)
	if err != nil {
		return nil, err
	}

	conversationRepository := &repositories.ConversationRepository{Database: service.Repository.Database}

	hasAccess, err := utils.IsUserInConversation(conversationRepository, conversationID, userID)
	if err != nil {
		return nil, err
	}

	if !hasAccess {
		return nil, errors.ErrForbidden
	}

	comments, err := service.Repository.GetCommentsByMessageID(messageID)
	if err != nil {
		return nil, err
	}

	for _, comment := range comments {
		if comment.Commenter.ID == userID {
			return nil, errors.ErrConflict
		}
	}

	commentID, err := service.Repository.CreateComment(messageID, userID, emoji)
	if err != nil {
		return nil, err
	}

	comment, err := service.Repository.GetCommentByID(commentID)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (service *CommentService) DeleteComment(commentID, authenticatedUserID uuid.UUID) error {
	messageRepository := &repositories.MessageRepository{Database: service.Repository.Database}

	conversationID, err := utils.GetConversationIDFromComment(service.Repository, messageRepository, commentID)
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

	comment, err := service.Repository.GetCommentByID(commentID)
	if err != nil {
		return err
	}

	if comment == nil {
		return errors.ErrNotFound
	}

	if comment.Commenter.ID != authenticatedUserID {
		return errors.ErrForbidden
	}

	return service.Repository.DeleteComment(commentID)
}
