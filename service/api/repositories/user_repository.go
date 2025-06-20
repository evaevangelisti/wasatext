package repositories

import (
	"database/sql"

	"github.com/evaevangelisti/wasatext/service/api/models"
	"github.com/evaevangelisti/wasatext/service/database"
	"github.com/evaevangelisti/wasatext/service/utils/errors"
	"github.com/evaevangelisti/wasatext/service/utils/globaltime"
	"github.com/google/uuid"
)

type UserRepository struct {
	Database database.Database
}

func (repository *UserRepository) GetUsers(q string, authenticatedUserID uuid.UUID) ([]models.User, error) {
	query := "SELECT user_id, username, profile_picture, created_at FROM users WHERE user_id != ?"
	args := []interface{}{authenticatedUserID}

	if q != "" {
		query += " WHERE username LIKE ?"
		args = append(args, q+"%")
	}

	rows, err := repository.Database.Query(query, args...)
	if err != nil {
		return nil, errors.ErrInternal
	}

	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User

		if err := rows.Scan(&user.ID, &user.Username, &user.ProfilePicture, &user.CreatedAt); err != nil {
			return nil, errors.ErrInternal
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.ErrInternal
	}

	return users, nil
}

func (repository *UserRepository) GetUserByID(userID uuid.UUID) (*models.User, error) {
	row := repository.Database.QueryRow("SELECT user_id, username, profile_picture, created_at FROM users WHERE user_id = ?", userID.String())

	var user models.User

	if err := row.Scan(&user.ID, &user.Username, &user.ProfilePicture, &user.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.ErrInternal
	}

	return &user, nil
}

func (repository *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	row := repository.Database.QueryRow("SELECT user_id, username, profile_picture, created_at FROM users WHERE username = ?", username)

	var user models.User

	if err := row.Scan(&user.ID, &user.Username, &user.ProfilePicture, &user.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.ErrInternal
	}

	return &user, nil
}

func (repository *UserRepository) CreateUser(username string) (uuid.UUID, error) {
	userID := uuid.New()
	createdAtTime := globaltime.Now()

	createdAtStr, err := globaltime.Format(createdAtTime)
	if err != nil {
		return uuid.Nil, errors.ErrInternal
	}

	_, err = repository.Database.Exec("INSERT INTO users (user_id, username, created_at) VALUES (?, ?, ?)", userID.String(), username, createdAtStr)
	if err != nil {
		return uuid.Nil, errors.ErrInternal
	}

	return userID, nil
}

func (repository *UserRepository) UpdateUsername(userID uuid.UUID, username string) error {
	_, err := repository.Database.Exec("UPDATE users SET username = ? WHERE user_id = ?", username, userID.String())
	if err != nil {
		return errors.ErrInternal
	}

	return nil
}

func (repository *UserRepository) UpdateProfilePicture(userID uuid.UUID, profilePicture string) error {
	_, err := repository.Database.Exec("UPDATE users SET profile_picture = ? WHERE user_id = ?", profilePicture, userID.String())
	if err != nil {
		return errors.ErrInternal
	}

	return nil
}
