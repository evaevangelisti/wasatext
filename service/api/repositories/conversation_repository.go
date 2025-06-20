package repositories

import (
	"database/sql"

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
		SELECT c.conversation_id
		FROM conversations c
		LEFT JOIN participants p ON c.conversation_id = p.conversation_id
		LEFT JOIN members m ON c.conversation_id = m.conversation_id
		LEFT JOIN (
		    SELECT conversation_id, MAX(sent_at) AS last_message_at
		    FROM messages
		    GROUP BY conversation_id
		) lm ON c.conversation_id = lm.conversation_id
		WHERE p.user_id = ? OR m.user_id = ?
		GROUP BY c.conversation_id
		ORDER BY COALESCE(lm.last_message_at, c.created_at) DESC
	`

	rows, err := repository.Database.Query(query, userID.String(), userID.String())
	if err != nil {
		return nil, errors.ErrInternal
	}

	defer rows.Close()

	var conversations []models.Conversation

	for rows.Next() {
		var conversationID uuid.UUID

		if err := rows.Scan(&conversationID); err != nil {
			return nil, errors.ErrInternal
		}

		conversation, err := repository.GetConversationByID(conversationID)
		if err != nil {
			return nil, err
		}

		conversations = append(conversations, conversation)
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

		var name, photo string

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

		return &models.GroupConversation{
			ID:        conversationID,
			Type:      "group",
			Name:      name,
			Photo:     photo,
			Members:   members,
			Messages:  messages,
			CreatedAt: createdAtTime,
		}, nil

	default:
		return nil, errors.ErrInternal
	}
}

func (repository *ConversationRepository) GetPrivateConversationByParticipants(participantIDs []uuid.UUID) (*models.PrivateConversation, error) {
	query := `
        SELECT c.conversation_id
        FROM conversations c
        JOIN participants p1 ON c.conversation_id = p1.conversation_id AND p1.user_id = ?
        JOIN participants p2 ON c.conversation_id = p2.conversation_id AND p2.user_id = ?
        WHERE c.type = 'private'
        GROUP BY c.conversation_id
        HAVING COUNT(*) = 2
        LIMIT 1
    `

	row := repository.Database.QueryRow(query, participantIDs[0].String(), participantIDs[1].String())

	var conversationID uuid.UUID

	if err := row.Scan(&conversationID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.ErrInternal
	}

	conversation, err := repository.GetConversationByID(conversationID)
	if err != nil {
		return nil, err
	}

	privateConversation, ok := conversation.(*models.PrivateConversation)
	if !ok {
		return nil, errors.ErrInternal
	}

	return privateConversation, nil
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
		var createdAt string

		if err := rows.Scan(&participant.ID, &participant.Username, &participant.ProfilePicture, createdAt); err != nil {
			return nil, errors.ErrInternal
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
		var createdAt string

		if err := rows.Scan(&member.ID, &member.Username, &member.ProfilePicture, createdAt); err != nil {
			return nil, errors.ErrInternal
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
			return false, nil
		}

		return false, errors.ErrInternal
	}

	return true, nil
}

func (repository *ConversationRepository) CreatePrivateConversation(participantIDs []uuid.UUID) (uuid.UUID, error) {
	conversationID := uuid.New()
	createdAtTime := globaltime.Now()

	createdAtStr, err := globaltime.Format(createdAtTime)
	if err != nil {
		return uuid.Nil, errors.ErrInternal
	}

	tx, err := repository.Database.Begin()
	if err != nil {
		return uuid.Nil, errors.ErrInternal
	}

	defer func() {
		_ = tx.Rollback()
	}()

	_, err = tx.Exec("INSERT INTO conversations (conversation_id, type, created_at) VALUES (?, ?, ?)", conversationID.String(), "private", createdAtStr)
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
	createdAtTime := globaltime.Now()

	createdAtStr, err := globaltime.Format(createdAtTime)
	if err != nil {
		return uuid.Nil, errors.ErrInternal
	}

	tx, err := repository.Database.Begin()
	if err != nil {
		return uuid.Nil, errors.ErrInternal
	}

	defer func() {
		_ = tx.Rollback()
	}()

	_, err = tx.Exec("INSERT INTO conversations (conversation_id, type, created_at) VALUES (?, ?, ?)", conversationID.String(), "group", createdAtStr)
	if err != nil {
		return uuid.Nil, errors.ErrInternal
	}

	_, err = tx.Exec("INSERT INTO group_conversations (conversation_id, name, photo) VALUES (?, ?)", conversationID.String(), name)
	if err != nil {
		return uuid.Nil, errors.ErrInternal
	}

	for _, userID := range memberIDs {
		_, err = tx.Exec("INSERT INTO members (conversation_id, user_id, joined_at) VALUES (?, ?, ?)", conversationID.String(), userID.String(), createdAtStr)
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
	joinedAtTime := globaltime.Now()

	joinedAtStr, err := globaltime.Format(joinedAtTime)
	if err != nil {
		return uuid.Nil, errors.ErrInternal
	}

	_, err = repository.Database.Exec("INSERT INTO members (conversation_id, user_id, joined_at) VALUES (?, ?, ?)", conversationID.String(), userID.String(), joinedAtStr)
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
	_, err := repository.Database.Exec("UPDATE group_conversations SET photo = ? WHERE conversation_id = ?", photo, conversationID.String())
	if err != nil {
		return errors.ErrInternal
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
