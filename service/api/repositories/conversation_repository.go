package repositories

import (
	"database/sql"
	stdErrors "errors"
	"os"

	"github.com/evaevangelisti/wasatext/service/api/models"
	"github.com/evaevangelisti/wasatext/service/database"
	"github.com/evaevangelisti/wasatext/service/utils/errors"
	"github.com/evaevangelisti/wasatext/service/utils/globaltime"
	"github.com/google/uuid"
)

type ConversationRepository struct {
	Database database.Database
}

func (repository *ConversationRepository) GetConversationsByUserID(userID uuid.UUID) ([]models.Conversation, error) {
	query := `
		SELECT c.conversation_id, (
        	SELECT m.message_id
         	FROM messages m
          	WHERE m.conversation_id = c.conversation_id
           	ORDER BY m.sent_at DESC
            LIMIT 1
        ) AS last_message_id
		FROM conversations c
		LEFT JOIN participants p ON c.conversation_id = p.conversation_id
		LEFT JOIN members mbr ON c.conversation_id = mbr.conversation_id
		WHERE p.user_id = ? OR mbr.user_id = ?
		ORDER BY COALESCE((
        	SELECT m.sent_at
         	FROM messages m
          	WHERE m.conversation_id = c.conversation_id
           	ORDER BY m.sent_at DESC
            LIMIT 1
        ), c.created_at) DESC
	`

	rows, err := repository.Database.Query(query, userID.String(), userID.String())
	if err != nil {
		return nil, errors.ErrInternal
	}

	defer rows.Close()

	conversations := []models.Conversation{}

	for rows.Next() {
		var (
			conversationID string
			lastMessageID  sql.NullString
		)

		if err := rows.Scan(&conversationID, &lastMessageID); err != nil {
			return nil, errors.ErrInternal
		}

		cid, err := uuid.Parse(conversationID)
		if err != nil {
			return nil, errors.ErrInternal
		}

		conversation, err := repository.GetConversationByID(cid)
		if err != nil {
			return nil, err
		}

		var lastMessage *models.Message

		if lastMessageID.Valid {
			mid, err := uuid.Parse(lastMessageID.String)
			if err != nil {
				return nil, errors.ErrInternal
			}

			messageRepository := MessageRepository{Database: repository.Database}

			lastMessage, err = messageRepository.GetMessageByID(mid)
			if err != nil {
				return nil, err
			}
		}

		switch conversation := conversation.(type) {
		case *models.PrivateConversation:
			conversation.Messages = nil
			conversation.LastMessage = lastMessage

			conversations = append(conversations, conversation)

		case *models.GroupConversation:
			conversation.Messages = nil
			conversation.LastMessage = lastMessage

			conversations = append(conversations, conversation)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, errors.ErrInternal
	}

	return conversations, nil
}

func (repository *ConversationRepository) GetConversationByID(conversationID uuid.UUID) (models.Conversation, error) {
	row := repository.Database.QueryRow("SELECT type, created_at FROM conversations WHERE conversation_id = ?", conversationID.String())

	var typ, createdAtStr string

	if err := row.Scan(&typ, &createdAtStr); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.ErrInternal
	}

	createdAtTime, err := globaltime.Parse(createdAtStr)
	if err != nil {
		return nil, errors.ErrInternal
	}

	messageRepository := MessageRepository{Database: repository.Database}

	messages, err := messageRepository.GetMessagesByConversationID(conversationID)
	if err != nil {
		return nil, err
	}

	switch typ {
	case "private":
		participants, err := repository.GetParticipants(conversationID)
		if err != nil {
			return nil, err
		}

		return &models.PrivateConversation{
			ID:           conversationID,
			Type:         "private",
			Participants: participants,
			Messages:     messages,
			CreatedAt:    createdAtTime,
		}, nil

	case "group":
		row := repository.Database.QueryRow("SELECT name, photo FROM group_conversations WHERE conversation_id = ?", conversationID.String())

		var (
			name  string
			photo sql.NullString
		)

		if err := row.Scan(&name, &photo); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, errors.ErrInternal
		}

		members, err := repository.GetMembers(conversationID)
		if err != nil {
			return nil, err
		}

		groupConversation := &models.GroupConversation{
			ID:        conversationID,
			Type:      "group",
			Name:      name,
			Members:   members,
			Messages:  messages,
			CreatedAt: createdAtTime,
		}

		if photo.Valid {
			groupConversation.Photo = photo.String
		}

		return groupConversation, nil

	default:
		return nil, errors.ErrInternal
	}
}

func (repository *ConversationRepository) GetPrivateConversationByParticipants(participantIDs []uuid.UUID) (*models.PrivateConversation, error) {
	query := `
		SELECT c.conversation_id
        FROM conversations c
        JOIN participants p ON c.conversation_id = p.conversation_id
        WHERE c.type = 'private' AND (p.user_id = ? OR p.user_id = ?)
        GROUP BY c.conversation_id
        HAVING COUNT(DISTINCT p.user_id) = 2
        LIMIT 1
    `

	row := repository.Database.QueryRow(query, participantIDs[0].String(), participantIDs[1].String())

	var conversationID string

	if err := row.Scan(&conversationID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.ErrInternal
	}

	cid, err := uuid.Parse(conversationID)
	if err != nil {
		return nil, errors.ErrInternal
	}

	conversation, err := repository.GetConversationByID(cid)
	if err != nil {
		return nil, err
	}

	privateConversation, ok := conversation.(*models.PrivateConversation)
	if !ok {
		return nil, errors.ErrInternal
	}

	return privateConversation, nil
}

func (repository *ConversationRepository) GetConversationByMessageID(messageID uuid.UUID) (models.Conversation, error) {
	row := repository.Database.QueryRow("SELECT conversation_id FROM messages WHERE message_id = ?", messageID.String())

	var conversationID string

	if err := row.Scan(&conversationID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.ErrInternal
	}

	cid, err := uuid.Parse(conversationID)
	if err != nil {
		return nil, errors.ErrInternal
	}

	return repository.GetConversationByID(cid)
}

func (repository *ConversationRepository) GetConversationByCommentID(commentID uuid.UUID) (models.Conversation, error) {
	row := repository.Database.QueryRow("SELECT message_id FROM comments WHERE comment_id = ?", commentID.String())

	var messageID string

	if err := row.Scan(&messageID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.ErrInternal
	}

	mid, err := uuid.Parse(messageID)
	if err != nil {
		return nil, errors.ErrInternal
	}

	return repository.GetConversationByMessageID(mid)
}

func (repository *ConversationRepository) GetParticipants(conversationID uuid.UUID) ([]models.User, error) {
	rows, err := repository.Database.Query("SELECT u.user_id, u.username, u.profile_picture, u.created_at FROM participants p JOIN users u ON p.user_id = u.user_id WHERE p.conversation_id = ?", conversationID.String())
	if err != nil {
		return nil, errors.ErrInternal
	}

	defer rows.Close()

	var participants []models.User

	for rows.Next() {
		var participant models.User

		var (
			participantID, createdAt string
			profilePicture           sql.NullString
		)

		if err := rows.Scan(&participantID, &participant.Username, &profilePicture, &createdAt); err != nil {
			return nil, errors.ErrInternal
		}

		pid, err := uuid.Parse(participantID)
		if err != nil {
			return nil, errors.ErrInternal
		}

		participant.ID = pid

		if profilePicture.Valid {
			participant.ProfilePicture = profilePicture.String
		}

		participant.CreatedAt, err = globaltime.Parse(createdAt)
		if err != nil {
			return nil, errors.ErrInternal
		}

		participants = append(participants, participant)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.ErrInternal
	}

	return participants, nil
}

func (repository *ConversationRepository) GetMembers(conversationID uuid.UUID) ([]models.User, error) {
	rows, err := repository.Database.Query("SELECT u.user_id, u.username, u.profile_picture, u.created_at FROM members m JOIN users u ON m.user_id = u.user_id WHERE m.conversation_id = ?", conversationID.String())
	if err != nil {
		return nil, errors.ErrInternal
	}

	defer rows.Close()

	var members []models.User

	for rows.Next() {
		var member models.User

		var (
			memberID, createdAt string
			profilePicture      sql.NullString
		)

		if err := rows.Scan(&memberID, &member.Username, &profilePicture, &createdAt); err != nil {
			return nil, errors.ErrInternal
		}

		mid, err := uuid.Parse(memberID)
		if err != nil {
			return nil, errors.ErrInternal
		}

		member.ID = mid

		if profilePicture.Valid {
			member.ProfilePicture = profilePicture.String
		}

		member.CreatedAt, err = globaltime.Parse(createdAt)
		if err != nil {
			return nil, errors.ErrInternal
		}

		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.ErrInternal
	}

	return members, nil
}

func (repository *ConversationRepository) IsUserInConversation(conversationID, userID uuid.UUID) (bool, error) {
	query := `
        SELECT 1
        FROM conversations c
        LEFT JOIN participants p ON c.conversation_id = p.conversation_id AND p.user_id = ?
        LEFT JOIN members m ON c.conversation_id = m.conversation_id AND m.user_id = ?
        WHERE c.conversation_id = ? AND (p.user_id IS NOT NULL OR m.user_id IS NOT NULL)
        LIMIT 1
    `

	row := repository.Database.QueryRow(query, userID.String(), userID.String(), conversationID.String())

	var exists int

	if err := row.Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			return false, errors.ErrNotFound
		}

		return false, errors.ErrInternal
	}

	return true, nil
}

func (repository *ConversationRepository) CreatePrivateConversation(participantIDs []uuid.UUID) (uuid.UUID, error) {
	conversationID := uuid.New()
	createdAt := globaltime.Now()

	tx, err := repository.Database.Begin()
	if err != nil {
		return uuid.Nil, errors.ErrInternal
	}

	defer func() {
		_ = tx.Rollback()
	}()

	_, err = tx.Exec("INSERT INTO conversations (conversation_id, type, created_at) VALUES (?, ?, ?)", conversationID.String(), "private", globaltime.Format(createdAt))
	if err != nil {
		return uuid.Nil, errors.ErrInternal
	}

	_, err = tx.Exec("INSERT INTO private_conversations (conversation_id) VALUES (?)", conversationID.String())
	if err != nil {
		return uuid.Nil, errors.ErrInternal
	}

	for _, userID := range participantIDs {
		_, err = tx.Exec("INSERT INTO participants (conversation_id, user_id) VALUES (?, ?)", conversationID.String(), userID.String())
		if err != nil {
			return uuid.Nil, errors.ErrInternal
		}
	}

	if err := tx.Commit(); err != nil {
		return uuid.Nil, errors.ErrInternal
	}

	return conversationID, nil
}

func (repository *ConversationRepository) CreateGroupConversation(name string, memberIDs []uuid.UUID) (uuid.UUID, error) {
	conversationID := uuid.New()
	createdAt := globaltime.Now()

	tx, err := repository.Database.Begin()
	if err != nil {
		return uuid.Nil, errors.ErrInternal
	}

	defer func() {
		_ = tx.Rollback()
	}()

	_, err = tx.Exec("INSERT INTO conversations (conversation_id, type, created_at) VALUES (?, ?, ?)", conversationID.String(), "group", globaltime.Format(createdAt))
	if err != nil {
		return uuid.Nil, errors.ErrInternal
	}

	_, err = tx.Exec("INSERT INTO group_conversations (conversation_id, name) VALUES (?, ?)", conversationID.String(), name)
	if err != nil {
		return uuid.Nil, errors.ErrInternal
	}

	for _, userID := range memberIDs {
		_, err = tx.Exec("INSERT INTO members (conversation_id, user_id) VALUES (?, ?)", conversationID.String(), userID.String())
		if err != nil {
			return uuid.Nil, errors.ErrInternal
		}
	}

	if err := tx.Commit(); err != nil {
		return uuid.Nil, errors.ErrInternal
	}

	return conversationID, nil
}

func (repository *ConversationRepository) AddMember(conversationID, userID uuid.UUID) (uuid.UUID, error) {
	joinedAt := globaltime.Now()

	_, err := repository.Database.Exec("INSERT INTO members (conversation_id, user_id, joined_at) VALUES (?, ?, ?)", conversationID.String(), userID.String(), globaltime.Format(joinedAt))
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func (repository *ConversationRepository) UpdateGroupName(conversationID uuid.UUID, name string) error {
	_, err := repository.Database.Exec("UPDATE group_conversations SET name = ? WHERE conversation_id = ?", name, conversationID.String())
	if err != nil {
		return errors.ErrInternal
	}

	return nil
}

func (repository *ConversationRepository) UpdateGroupPhoto(conversationID uuid.UUID, photo string) error {
	var oldPhoto sql.NullString

	err := repository.Database.QueryRow("SELECT photo FROM group_conversations WHERE conversation_id = ?", conversationID.String()).Scan(&oldPhoto)
	if err != nil && !stdErrors.Is(err, sql.ErrNoRows) {
		return errors.ErrInternal
	}

	_, err = repository.Database.Exec("UPDATE group_conversations SET photo = ? WHERE conversation_id = ?", sql.NullString{String: photo, Valid: photo != ""}, conversationID.String())
	if err != nil {
		return errors.ErrInternal
	}

	if oldPhoto.Valid && oldPhoto.String != "" {
		if err := os.Remove("." + oldPhoto.String); err != nil && !os.IsNotExist(err) {
			return errors.ErrInternal
		}
	}

	return nil
}

func (repository *ConversationRepository) DeleteGroupConversation(conversationID uuid.UUID) error {
	var groupPhoto sql.NullString

	err := repository.Database.QueryRow("SELECT photo FROM group_conversations WHERE conversation_id = ?", conversationID.String()).Scan(&groupPhoto)
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

	_, err = tx.Exec("DELETE FROM messages WHERE conversation_id = ?", conversationID.String())
	if err != nil {
		return errors.ErrInternal
	}

	_, err = tx.Exec("DELETE FROM group_conversations WHERE conversation_id = ?", conversationID.String())
	if err != nil {
		return errors.ErrInternal
	}

	_, err = tx.Exec("DELETE FROM members WHERE conversation_id = ?", conversationID.String())
	if err != nil {
		return errors.ErrInternal
	}

	_, err = tx.Exec("DELETE FROM conversations WHERE conversation_id = ?", conversationID.String())
	if err != nil {
		return errors.ErrInternal
	}

	if err := tx.Commit(); err != nil {
		return errors.ErrInternal
	}

	if groupPhoto.Valid && groupPhoto.String != "" {
		if err := os.Remove("." + groupPhoto.String); err != nil && !os.IsNotExist(err) {
			return errors.ErrInternal
		}
	}

	return nil
}

func (repository *ConversationRepository) RemoveMember(conversationID, userID uuid.UUID) error {
	_, err := repository.Database.Exec("DELETE FROM members WHERE conversation_id = ? AND user_id = ?", conversationID.String(), userID.String())
	if err != nil {
		return errors.ErrInternal
	}

	return nil
}
