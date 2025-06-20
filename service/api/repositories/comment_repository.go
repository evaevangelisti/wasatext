package repositories

import (
	"database/sql"

	"github.com/evaevangelisti/wasatext/service/api/models"
	"github.com/evaevangelisti/wasatext/service/database"
	"github.com/evaevangelisti/wasatext/service/utils/errors"
	"github.com/evaevangelisti/wasatext/service/utils/globaltime"
	"github.com/google/uuid"
)

type CommentRepository struct {
	Database database.Database
}

func (repository *CommentRepository) GetCommentsByMessageID(messageID uuid.UUID) ([]models.Comment, error) {
	rows, err := repository.Database.Query("SELECT c.comment_id, c.emoji, c.commented_at, c.message_id, u.user_id, u.username, u.profile_picture, u.created_at FROM comments c JOIN users u ON c.commenter_id = u.user_id WHERE c.conversation_id = ? ORDER BY c.commented_at ASC", messageID.String())
	if err != nil {
		return nil, errors.ErrInternal
	}

	defer rows.Close()

	var comments []models.Comment

	for rows.Next() {
		var comment models.Comment

		var commenter models.User
		var commentedAt string

		if err := rows.Scan(&comment.ID, &comment.Emoji, &commentedAt, &comment.MessageID, &commenter.ID, &commenter.Username, &commenter.ProfilePicture, &commenter.CreatedAt); err != nil {
			return nil, errors.ErrInternal
		}

		comment.Commenter = commenter

		comment.CommentedAt, err = globaltime.Parse(commentedAt)
		if err != nil {
			return nil, errors.ErrInternal
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

func (repository *CommentRepository) GetCommentByID(commentID uuid.UUID) (*models.Comment, error) {
	row := repository.Database.QueryRow("SELECT comment_id, commenter_id, emoji, commented_at, message_id FROM comments WHERE comment_id = ?", commentID.String())

	var comment models.Comment

	var commenterID, commentedAt string

	if err := row.Scan(&comment.ID, &commenterID, &comment.Emoji, &commentedAt, &comment.MessageID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.ErrInternal
	}

	uid, err := uuid.Parse(commenterID)
	if err != nil {
		return nil, errors.ErrInternal
	}

	userRepository := UserRepository{Database: repository.Database}

	commenter, err := userRepository.GetUserByID(uid)
	if err != nil {
		return nil, err
	}

	comment.Commenter = *commenter

	comment.CommentedAt, err = globaltime.Parse(commentedAt)
	if err != nil {
		return nil, errors.ErrInternal
	}

	return &comment, nil
}

func (repository *CommentRepository) CreateComment(messageID, userID uuid.UUID, emoji string) (uuid.UUID, error) {
	commentID := uuid.New()
	commentedAtTime := globaltime.Now()

	commentedAtStr, err := globaltime.Format(commentedAtTime)
	if err != nil {
		return uuid.Nil, errors.ErrInternal
	}

	_, err = repository.Database.Exec("INSERT INTO comments (comment_id, emoji, commented_at, message_id, commenter_id) VALUES (?, ?, ?, ?, ?)", commentID.String(), emoji, commentedAtStr, messageID.String(), userID.String())
	if err != nil {
		return uuid.Nil, errors.ErrInternal
	}

	return commentID, nil
}

func (repository *CommentRepository) DeleteComment(commentID uuid.UUID) error {
	_, err := repository.Database.Exec("DELETE FROM comments WHERE comment_id = ?", commentID.String())
	if err != nil {
		return errors.ErrInternal
	}

	return nil
}
