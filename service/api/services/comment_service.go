package services

import (
	"github.com/evaevangelisti/wasatext/service/api/models"
	"github.com/evaevangelisti/wasatext/service/api/repositories"
	"github.com/evaevangelisti/wasatext/service/utils/errors"
	"github.com/google/uuid"
)

type CommentService struct {
	Repository *repositories.CommentRepository
}

func (service *CommentService) CreateComment(messageID, userID uuid.UUID, emoji string) (*models.Comment, error) {
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

func (service *CommentService) DeleteComment(commentID, userID uuid.UUID) error {
	conversationRepository := &repositories.ConversationRepository{Database: service.Repository.Database}

	conversation, err := conversationRepository.GetConversationByCommentID(commentID)
	if err != nil {
		return err
	}

	hasAccess, err := conversationRepository.IsUserInConversation(conversation.GetID(), userID)
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

	if comment.Commenter.ID != userID {
		return errors.ErrForbidden
	}

	return service.Repository.DeleteComment(commentID)
}
