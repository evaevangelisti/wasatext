package repositories

import (
	"database/sql"
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

	var messages []models.Message

	for rows.Next() {
		var messageID uuid.UUID

		if err := rows.Scan(&messageID); err != nil {
			return nil, errors.ErrInternal
		}

		message, err := repository.GetMessageByID(messageID)
		if err != nil {
			return nil, err
		}

		messages = append(messages, *message)
	}

	return messages, nil
}

func (repository *MessageRepository) GetMessageByID(messageID uuid.UUID) (*models.Message, error) {
	row := repository.Database.QueryRow("SELECT message_id, sender_id, content, attachment, sent_at, edited_at, conversation_id FROM messages WHERE message_id = ?", messageID.String())

	var message models.Message

	var senderID string
	var sentAt, editedAt string

	if err := row.Scan(&message.ID, &senderID, &message.Content, &message.Attachment, &sentAt, &editedAt, &message.ConversationID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.ErrInternal
	}

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
	} else if err != sql.ErrNoRows {
		return nil, errors.ErrInternal
	} else {
		message.IsForwarded = false
	}

	trackingRows, err := repository.Database.Query(`SELECT user_id, read_at FROM message_trackings WHERE message_id = ?`, message.ID.String())
	if err != nil {
		return nil, errors.ErrInternal
	}

	defer trackingRows.Close()

	message.Trackings.Read = make(map[string]string)

	for trackingRows.Next() {
		var userID, readAt string

		if err := trackingRows.Scan(&userID, &readAt); err != nil {
			return nil, errors.ErrInternal
		}

		if readAt != "" {
			message.Trackings.Read[userID] = readAt
		}
	}

	message.SentAt, err = globaltime.Parse(sentAt)
	if err != nil {
		return nil, errors.ErrInternal
	}

	if editedAt != "" {
		message.EditedAt, err = globaltime.Parse(editedAt)
		if err != nil {
			return nil, errors.ErrInternal
		}
	}

	return &message, nil
}

func (repository *MessageRepository) CreateMessage(conversationID, userID uuid.UUID, content, attachment string) (uuid.UUID, error) {
	messageID := uuid.New()
	sentAtTime := globaltime.Now()

	sentAtStr, err := globaltime.Format(sentAtTime)
	if err != nil {
		return uuid.Nil, errors.ErrInternal
	}

	_, err = repository.Database.Exec("INSERT INTO messages (message_id, conversation_id, sender_id, content, attachment, sent_at) VALUES (?, ?, ?, ?, ?, ?)", messageID.String(), conversationID.String(), userID.String(), content, attachment, sentAtStr)
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
	forwardedAtTime := globaltime.Now()

	forwardedAtStr, err := globaltime.Format(forwardedAtTime)
	if err != nil {
		return uuid.Nil, errors.ErrInternal
	}

	tx, err := repository.Database.Begin()
	if err != nil {
		return uuid.Nil, errors.ErrInternal
	}

	defer tx.Rollback()

	_, err = repository.Database.Exec("INSERT INTO messages (message_id, content, attachment, sent_at, conversation_id, sender_id) VALUES (?, ?, ?, ?, ?, ?)", forwardedMessageID.String(), originalMessage.Content, originalMessage.Attachment, forwardedAtStr, conversationID.String(), userID.String())
	if err != nil {
		return uuid.Nil, errors.ErrInternal
	}

	_, err = tx.Exec("INSERT INTO forwarded_messages (forwarded_message_id, forwarded_at, original_message_id, conversation_id, sender_id) VALUES (?, ?, ?, ?, ?)", forwardedMessageID.String(), forwardedAtStr, originalMessage.ID.String(), conversationID.String(), userID.String())
	if err != nil {
		return uuid.Nil, errors.ErrInternal
	}

	if err := tx.Commit(); err != nil {
		return uuid.Nil, errors.ErrInternal
	}

	return forwardedMessageID, nil
}

func (repository *MessageRepository) AddMessageTracking(messageID, userID uuid.UUID, readAtTime time.Time) error {
	readAtStr, err := globaltime.Format(readAtTime)
	if err != nil {
		return errors.ErrInternal
	}

	_, err = repository.Database.Exec(`INSERT INTO message_trackings (message_id, user_id, read_at) VALUES (?, ?, ?)`, messageID.String(), userID.String(), readAtStr)
	if err != nil {
		return errors.ErrInternal
	}

	return nil
}

func (repository *MessageRepository) UpdateMessage(messageID uuid.UUID, content string) error {
	editedAtTime := globaltime.Now()

	editedAtStr, err := globaltime.Format(editedAtTime)
	if err != nil {
		return errors.ErrInternal
	}

	_, err = repository.Database.Exec("UPDATE messages SET content = ?, edited_at = ? WHERE message_id = ?", content, editedAtStr, messageID.String())
	if err != nil {
		return errors.ErrInternal
	}

	return nil
}

func (repository *MessageRepository) DeleteMessage(messageID uuid.UUID) error {
	tx, err := repository.Database.Begin()
	if err != nil {
		return errors.ErrInternal
	}

	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM message_trackings WHERE message_id = ?", messageID.String())
	if err != nil {
		return errors.ErrInternal
	}

	_, err = tx.Exec("DELETE FROM forwarded_messages WHERE message_id = ?", messageID.String())
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

	return nil
}
