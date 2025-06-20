package utils

import (
	"github.com/evaevangelisti/wasatext/service/api/repositories"
	"github.com/evaevangelisti/wasatext/service/utils/errors"
	"github.com/google/uuid"
)

func IsUserInConversation(conversationRepository *repositories.ConversationRepository, conversationID, userID uuid.UUID) (bool, error) {
	return conversationRepository.IsUserInConversation(conversationID, userID)
}

func GetConversationIDFromMessage(messageRepository *repositories.MessageRepository, messageID uuid.UUID) (uuid.UUID, error) {
	message, err := messageRepository.GetMessageByID(messageID)
	if err != nil || message == nil {
		return uuid.Nil, errors.ErrNotFound
	}

	return message.ConversationID, nil
}

func GetConversationIDFromComment(commentRepository *repositories.CommentRepository, messageRepository *repositories.MessageRepository, commentID uuid.UUID) (uuid.UUID, error) {
	comment, err := commentRepository.GetCommentByID(commentID)
	if err != nil || comment == nil {
		return uuid.Nil, errors.ErrNotFound
	}

	return GetConversationIDFromMessage(messageRepository, comment.MessageID)
}
