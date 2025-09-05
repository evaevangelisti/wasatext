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
	rows, err := repository.Database.Query("SELECT c.comment_id FROM comments c WHERE c.message_id = ? ORDER BY c.commented_at ASC", messageID.String())
	if err != nil {
		return nil, errors.ErrInternal
	}

	defer rows.Close()

	comments := []models.Comment{}

	for rows.Next() {
		var commentID string

		if err := rows.Scan(&commentID); err != nil {
			return nil, errors.ErrInternal
		}

		cid, err := uuid.Parse(commentID)
		if err != nil {
			return nil, errors.ErrInternal
		}

		comment, err := repository.GetCommentByID(cid)
		if err != nil {
			return nil, err
		}

		comments = append(comments, *comment)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.ErrInternal
	}

	return comments, nil
}

func (repository *CommentRepository) GetCommentByID(commentID uuid.UUID) (*models.Comment, error) {
	row := repository.Database.QueryRow("SELECT user_id, emoji, commented_at FROM comments WHERE comment_id = ?", commentID.String())

	var comment models.Comment
	var commenterID, commentedAt string

	if err := row.Scan(&commenterID, &comment.Emoji, &commentedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.ErrInternal
	}

	comment.ID = commentID

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
	commentedAt := globaltime.Now()

	_, err := repository.Database.Exec("INSERT INTO comments (comment_id, emoji, commented_at, message_id, user_id) VALUES (?, ?, ?, ?, ?)", commentID.String(), emoji, globaltime.Format(commentedAt), messageID.String(), userID.String())
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
