package repositories

import (
	"database/sql"
	stdErrors "errors"
	"os"
	"time"

	"github.com/evaevangelisti/wasatext/service/api/models"
	"github.com/evaevangelisti/wasatext/service/database"
	"github.com/evaevangelisti/wasatext/service/utils/errors"
	"github.com/evaevangelisti/wasatext/service/utils/globaltime"
	"github.com/google/uuid"
)

type MessageRepository struct {
	Database database.Database
}

func (repository *MessageRepository) GetMessagesByConversationID(conversationID uuid.UUID) ([]models.Message, error) {
	rows, err := repository.Database.Query("SELECT message_id FROM messages WHERE conversation_id = ? ORDER BY sent_at ASC", conversationID.String())
	if err != nil {
		return nil, errors.ErrInternal
	}

	defer rows.Close()

	messages := []models.Message{}

	for rows.Next() {
		var messageID string

		if err := rows.Scan(&messageID); err != nil {
			return nil, errors.ErrInternal
		}

		mid, err := uuid.Parse(messageID)
		if err != nil {
			return nil, errors.ErrInternal
		}

		message, err := repository.GetMessageByID(mid)
		if err != nil {
			return nil, err
		}

		messages = append(messages, *message)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.ErrInternal
	}

	return messages, nil
}

func (repository *MessageRepository) GetMessageByID(messageID uuid.UUID) (*models.Message, error) {
	row := repository.Database.QueryRow("SELECT sender_id, content, attachment, sent_at, edited_at FROM messages WHERE message_id = ?", messageID.String())

	var message models.Message

	var (
		senderID, sentAt              string
		content, attachment, editedAt sql.NullString
	)

	if err := row.Scan(&senderID, &content, &attachment, &sentAt, &editedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.ErrInternal
	}

	message.ID = messageID

	uid, err := uuid.Parse(senderID)
	if err != nil {
		return nil, errors.ErrInternal
	}

	userRepository := UserRepository{Database: repository.Database}

	sender, err := userRepository.GetUserByID(uid)
	if err != nil {
		return nil, err
	}

	message.Sender = *sender

	if content.Valid {
		message.Content = content.String
	}

	if attachment.Valid {
		message.Attachment = attachment.String
	}

	commentRepository := CommentRepository{Database: repository.Database}

	comments, err := commentRepository.GetCommentsByMessageID(message.ID)
	if err != nil {
		return nil, err
	}

	message.Comments = comments

	forwardRow := repository.Database.QueryRow(`SELECT original_message_id FROM forwarded_messages WHERE forwarded_message_id = ?`, message.ID.String())

	var originalMessageID string

	if err := forwardRow.Scan(&originalMessageID); err == nil {
		message.IsForwarded = true

		message.OriginalMessageID, err = uuid.Parse(originalMessageID)
		if err != nil {
			return nil, errors.ErrInternal
		}
	} else if !stdErrors.Is(err, sql.ErrNoRows) {
		return nil, errors.ErrInternal
	} else {
		message.IsForwarded = false
	}

	trackingRows, err := repository.Database.Query(`SELECT user_id, read_at FROM message_trackings WHERE message_id = ?`, message.ID.String())
	if err != nil {
		return nil, errors.ErrInternal
	}

	defer trackingRows.Close()

	message.Trackings.Read = make(map[uuid.UUID]time.Time)

	for trackingRows.Next() {
		var userID, readAtStr string

		if err := trackingRows.Scan(&userID, &readAtStr); err != nil {
			return nil, errors.ErrInternal
		}

		uid, err = uuid.Parse(userID)
		if err != nil {
			return nil, errors.ErrInternal
		}

		readAtTime, err := globaltime.Parse(readAtStr)
		if err != nil {
			return nil, errors.ErrInternal
		}

		if readAtStr != "" {
			message.Trackings.Read[uid] = readAtTime
		}
	}

	if err := trackingRows.Err(); err != nil {
		return nil, errors.ErrInternal
	}

	message.SentAt, err = globaltime.Parse(sentAt)
	if err != nil {
		return nil, errors.ErrInternal
	}

	if editedAt.Valid && editedAt.String != "" {
		message.EditedAt, err = globaltime.Parse(editedAt.String)
		if err != nil {
			return nil, errors.ErrInternal
		}
	}

	return &message, nil
}

func (repository *MessageRepository) CreateMessage(conversationID, userID uuid.UUID, content, attachment string) (uuid.UUID, error) {
	messageID := uuid.New()
	sentAt := globaltime.Now()

	_, err := repository.Database.Exec("INSERT INTO messages (message_id, conversation_id, sender_id, content, attachment, sent_at) VALUES (?, ?, ?, ?, ?, ?)", messageID.String(), conversationID.String(), userID.String(), sql.NullString{String: content, Valid: content != ""}, sql.NullString{String: attachment, Valid: attachment != ""}, globaltime.Format(sentAt))
	if err != nil {
		return uuid.Nil, errors.ErrInternal
	}

	return messageID, nil
}

func (repository *MessageRepository) CreateForwardedMessage(conversationID, userID, originalMessageID uuid.UUID) (uuid.UUID, error) {
	messageRepository := MessageRepository{Database: repository.Database}

	originalMessage, err := messageRepository.GetMessageByID(originalMessageID)
	if err != nil {
		return uuid.Nil, err
	}

	forwardedMessageID := uuid.New()
	forwardedAt := globaltime.Now()

	tx, err := repository.Database.Begin()
	if err != nil {
		return uuid.Nil, errors.ErrInternal
	}

	defer func() {
		_ = tx.Rollback()
	}()

	_, err = repository.Database.Exec("INSERT INTO messages (message_id, content, attachment, sent_at, conversation_id, sender_id) VALUES (?, ?, ?, ?, ?, ?)", forwardedMessageID.String(), sql.NullString{String: originalMessage.Content, Valid: originalMessage.Content != ""}, sql.NullString{String: originalMessage.Attachment, Valid: originalMessage.Attachment != ""}, globaltime.Format(forwardedAt), conversationID.String(), userID.String())
	if err != nil {
		return uuid.Nil, errors.ErrInternal
	}

	_, err = tx.Exec("INSERT INTO forwarded_messages (forwarded_message_id, forwarded_at, original_message_id, conversation_id, sender_id) VALUES (?, ?, ?, ?, ?)", forwardedMessageID.String(), globaltime.Format(forwardedAt), originalMessage.ID.String(), conversationID.String(), userID.String())
	if err != nil {
		return uuid.Nil, errors.ErrInternal
	}

	if err := tx.Commit(); err != nil {
		return uuid.Nil, errors.ErrInternal
	}

	return forwardedMessageID, nil
}

func (repository *MessageRepository) AddMessageTracking(messageID, userID uuid.UUID, readAt time.Time) error {
	_, err := repository.Database.Exec(`INSERT INTO message_trackings (message_id, user_id, read_at) VALUES (?, ?, ?)`, messageID.String(), userID.String(), globaltime.Format(readAt))
	if err != nil {
		return errors.ErrInternal
	}

	return nil
}

func (repository *MessageRepository) UpdateMessage(messageID uuid.UUID, content string) error {
	editedAt := globaltime.Now()

	_, err := repository.Database.Exec("UPDATE messages SET content = ?, edited_at = ? WHERE message_id = ?", content, globaltime.Format(editedAt), messageID.String())
	if err != nil {
		return errors.ErrInternal
	}

	return nil
}

func (repository *MessageRepository) DeleteMessage(messageID uuid.UUID) error {
	var attachment sql.NullString

	err := repository.Database.QueryRow("SELECT attachment FROM messages WHERE message_id = ?", messageID.String()).Scan(&attachment)
	if err != nil && !stdErrors.Is(err, sql.ErrNoRows) {
		return errors.ErrInternal
	}

	tx, err := repository.Database.Begin()
	if err != nil {
		return errors.ErrInternal
	}

	defer func() {
		_ = tx.Rollback()
	}()

	_, err = tx.Exec("DELETE FROM message_trackings WHERE message_id = ?", messageID.String())
	if err != nil {
		return errors.ErrInternal
	}

	_, err = tx.Exec("DELETE FROM forwarded_messages WHERE original_message_id = ?", messageID.String())
	if err != nil {
		return errors.ErrInternal
	}

	_, err = tx.Exec("DELETE FROM comments WHERE message_id = ?", messageID.String())
	if err != nil {
		return errors.ErrInternal
	}

	_, err = tx.Exec("DELETE FROM messages WHERE message_id = ?", messageID.String())
	if err != nil {
		return errors.ErrInternal
	}

	if err = tx.Commit(); err != nil {
		return errors.ErrInternal
	}

	if attachment.Valid && attachment.String != "" {
		if err := os.Remove("." + attachment.String); err != nil && !os.IsNotExist(err) {
			return errors.ErrInternal
		}
	}

	return nil
}
