package repositories

import (
	"database/sql"
	stdErrors "errors"
	"os"
	"strings"
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
	rows, err := repository.Database.Query(
		`SELECT message_id, sender_id, content, attachment, sent_at, edited_at
		 FROM messages
		 WHERE conversation_id = ?
		 ORDER BY sent_at ASC`, conversationID.String())

	if err != nil {
		return nil, errors.ErrInternal
	}

	defer rows.Close()

	type rawMessage struct {
		ID         uuid.UUID
		SenderID   uuid.UUID
		Content    string
		Attachment string
		SentAt     string
		EditedAt   string
	}

	rawMessages := []rawMessage{}
	messageIDs := []uuid.UUID{}
	senderIDs := map[uuid.UUID]struct{}{}

	for rows.Next() {
		var (
			messageID, senderID, sentAt   string
			content, attachment, editedAt sql.NullString
		)

		if err := rows.Scan(&messageID, &senderID, &content, &attachment, &sentAt, &editedAt); err != nil {
			return nil, errors.ErrInternal
		}

		mid, err := uuid.Parse(messageID)
		if err != nil {
			return nil, errors.ErrInternal
		}

		sid, err := uuid.Parse(senderID)
		if err != nil {
			return nil, errors.ErrInternal
		}

		rawMessages = append(rawMessages, rawMessage{
			ID:         mid,
			SenderID:   sid,
			Content:    content.String,
			Attachment: attachment.String,
			SentAt:     sentAt,
			EditedAt:   editedAt.String,
		})

		messageIDs = append(messageIDs, mid)
		senderIDs[sid] = struct{}{}
	}

	if err := rows.Err(); err != nil {
		return nil, errors.ErrInternal
	}

	if len(rawMessages) == 0 {
		return []models.Message{}, nil
	}

	placeholders := make([]string, len(messageIDs))
	args := make([]interface{}, len(messageIDs))
	for i, mid := range messageIDs {
		placeholders[i] = "?"
		args[i] = mid.String()
	}

	commentRows, err := repository.Database.Query(
		`SELECT comment_id, emoji, commented_at, message_id, user_id
		 FROM comments
		 WHERE message_id IN (`+strings.Join(placeholders, ",")+`)
		 ORDER BY commented_at ASC`, args...)

	if err != nil {
		return nil, errors.ErrInternal
	}

	defer commentRows.Close()

	type rawComment struct {
		ID          uuid.UUID
		Emoji       string
		CommentedAt string
		MessageID   uuid.UUID
		UserID      uuid.UUID
	}

	rawComments := []rawComment{}
	commenterIDs := map[uuid.UUID]struct{}{}

	for commentRows.Next() {
		var (
			commentID, emoji, commentedAt, messageID, userID string
		)

		if err := commentRows.Scan(&commentID, &emoji, &commentedAt, &messageID, &userID); err != nil {
			return nil, errors.ErrInternal
		}

		cid, err := uuid.Parse(commentID)
		if err != nil {
			return nil, errors.ErrInternal
		}

		mid, err := uuid.Parse(messageID)
		if err != nil {
			return nil, errors.ErrInternal
		}

		uid, err := uuid.Parse(userID)
		if err != nil {
			return nil, errors.ErrInternal
		}

		rawComments = append(rawComments, rawComment{
			ID:          cid,
			Emoji:       emoji,
			CommentedAt: commentedAt,
			MessageID:   mid,
			UserID:      uid,
		})

		commenterIDs[uid] = struct{}{}
	}

	if err := commentRows.Err(); err != nil {
		return nil, errors.ErrInternal
	}

	trackingRows, err := repository.Database.Query(
		`SELECT message_id, user_id, read_at
		 FROM message_trackings
		 WHERE message_id IN (`+strings.Join(placeholders, ",")+`)`, args...)

	if err != nil {
		return nil, errors.ErrInternal
	}

	defer trackingRows.Close()

	type rawTracking struct {
		MessageID uuid.UUID
		UserID    uuid.UUID
		ReadAt    string
	}

	rawTrackings := []rawTracking{}
	trackingUserIDs := map[uuid.UUID]struct{}{}

	for trackingRows.Next() {
		var messageID, userID, readAt string
		if err := trackingRows.Scan(&messageID, &userID, &readAt); err != nil {
			return nil, errors.ErrInternal
		}

		mid, err := uuid.Parse(messageID)
		if err != nil {
			return nil, errors.ErrInternal
		}

		uid, err := uuid.Parse(userID)
		if err != nil {
			return nil, errors.ErrInternal
		}

		rawTrackings = append(rawTrackings, rawTracking{
			MessageID: mid,
			UserID:    uid,
			ReadAt:    readAt,
		})

		trackingUserIDs[uid] = struct{}{}
	}

	if err := trackingRows.Err(); err != nil {
		return nil, errors.ErrInternal
	}

	allUserIDs := map[uuid.UUID]struct{}{}

	for uid := range senderIDs {
		allUserIDs[uid] = struct{}{}
	}

	for uid := range commenterIDs {
		allUserIDs[uid] = struct{}{}
	}

	for uid := range trackingUserIDs {
		allUserIDs[uid] = struct{}{}
	}

	userIDList := []uuid.UUID{}

	for uid := range allUserIDs {
		userIDList = append(userIDList, uid)
	}

	userPlaceholders := make([]string, len(userIDList))
	userArgs := make([]interface{}, len(userIDList))
	for i, uid := range userIDList {
		userPlaceholders[i] = "?"
		userArgs[i] = uid.String()
	}

	userRows, err := repository.Database.Query(
		`SELECT user_id, username, profile_picture, created_at
		 FROM users
		 WHERE user_id IN (`+strings.Join(userPlaceholders, ",")+`)`, userArgs...)

	if err != nil {
		return nil, errors.ErrInternal
	}

	defer userRows.Close()

	userMap := map[uuid.UUID]models.User{}

	for userRows.Next() {
		var userID, username, createdAt string
		var profilePicture sql.NullString
		if err := userRows.Scan(&userID, &username, &profilePicture, &createdAt); err != nil {
			return nil, errors.ErrInternal
		}

		uid, err := uuid.Parse(userID)
		if err != nil {
			return nil, errors.ErrInternal
		}

		createdAtTime, err := globaltime.Parse(createdAt)
		if err != nil {
			return nil, errors.ErrInternal
		}

		user := models.User{
			ID:             uid,
			Username:       username,
			ProfilePicture: "",
			CreatedAt:      createdAtTime,
		}

		if profilePicture.Valid {
			user.ProfilePicture = profilePicture.String
		}

		userMap[uid] = user
	}

	if err := userRows.Err(); err != nil {
		return nil, errors.ErrInternal
	}

	commentsByMessage := map[uuid.UUID][]models.Comment{}

	for _, rc := range rawComments {
		commentedAtTime, err := globaltime.Parse(rc.CommentedAt)
		if err != nil {
			return nil, errors.ErrInternal
		}

		comment := models.Comment{
			ID:          rc.ID,
			Commenter:   userMap[rc.UserID],
			Emoji:       rc.Emoji,
			CommentedAt: commentedAtTime,
		}

		commentsByMessage[rc.MessageID] = append(commentsByMessage[rc.MessageID], comment)
	}

	trackingByMessage := map[uuid.UUID]map[uuid.UUID]time.Time{}

	for _, rt := range rawTrackings {
		readAtTime, err := globaltime.Parse(rt.ReadAt)
		if err != nil {
			return nil, errors.ErrInternal
		}

		if trackingByMessage[rt.MessageID] == nil {
			trackingByMessage[rt.MessageID] = make(map[uuid.UUID]time.Time)
		}

		trackingByMessage[rt.MessageID][rt.UserID] = readAtTime
	}

	forwardRows, err := repository.Database.Query(
		`SELECT forwarded_message_id, original_message_id
		 FROM forwarded_messages
		 WHERE forwarded_message_id IN (`+strings.Join(placeholders, ",")+`)`, args...)

	if err != nil {
		return nil, errors.ErrInternal
	}

	defer forwardRows.Close()

	forwardedMap := map[uuid.UUID]uuid.UUID{}

	for forwardRows.Next() {
		var fmid, omid string

		if err := forwardRows.Scan(&fmid, &omid); err != nil {
			return nil, errors.ErrInternal
		}

		fmidUUID, err := uuid.Parse(fmid)
		if err != nil {
			return nil, errors.ErrInternal
		}

		omidUUID, err := uuid.Parse(omid)
		if err != nil {
			return nil, errors.ErrInternal
		}

		forwardedMap[fmidUUID] = omidUUID
	}

	if err := forwardRows.Err(); err != nil {
		return nil, errors.ErrInternal
	}

	messages := make([]models.Message, 0, len(rawMessages))

	for _, rm := range rawMessages {
		sentAtTime, err := globaltime.Parse(rm.SentAt)
		if err != nil {
			return nil, errors.ErrInternal
		}

		var editedAtTime time.Time
		if rm.EditedAt != "" {
			editedAtTime, err = globaltime.Parse(rm.EditedAt)
			if err != nil {
				return nil, errors.ErrInternal
			}
		}

		msg := models.Message{
			ID:          rm.ID,
			Sender:      userMap[rm.SenderID],
			Content:     rm.Content,
			Attachment:  rm.Attachment,
			Comments:    commentsByMessage[rm.ID],
			IsForwarded: false,
			Trackings: struct {
				Read map[uuid.UUID]time.Time `json:"read,omitempty" validate:"omitempty"`
			}{
				Read: trackingByMessage[rm.ID],
			},
			SentAt:   sentAtTime,
			EditedAt: editedAtTime,
		}

		if omid, ok := forwardedMap[rm.ID]; ok {
			msg.IsForwarded = true
			msg.OriginalMessageID = omid
		}

		if msg.Trackings.Read == nil {
			msg.Trackings.Read = make(map[uuid.UUID]time.Time)
		}

		messages = append(messages, msg)
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
