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

type UserRepository struct {
	Database database.Database
}

func (repository *UserRepository) GetUsers(q string, authenticatedUserID uuid.UUID) ([]models.User, error) {
	query := "SELECT user_id, username, profile_picture, created_at FROM users WHERE user_id != ?"

	args := []interface{}{authenticatedUserID}
	if q != "" {
		query += " AND username LIKE ?"
		args = append(args, q+"%")
	}

	rows, err := repository.Database.Query(query, args...)
	if err != nil {
		return nil, errors.ErrInternal
	}

	defer rows.Close()

	users := []models.User{}

	for rows.Next() {
		var user models.User

		var (
			userID, createdAt string
			profilePicture    sql.NullString
		)

		if err := rows.Scan(&userID, &user.Username, &profilePicture, &createdAt); err != nil {
			return nil, errors.ErrInternal
		}

		uid, err := uuid.Parse(userID)
		if err != nil {
			return nil, errors.ErrInternal
		}

		user.ID = uid

		if profilePicture.Valid {
			user.ProfilePicture = profilePicture.String
		}

		user.CreatedAt, err = globaltime.Parse(createdAt)
		if err != nil {
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
	row := repository.Database.QueryRow("SELECT username, profile_picture, created_at FROM users WHERE user_id = ?", userID.String())

	var user models.User

	var (
		createdAt      string
		profilePicture sql.NullString
	)

	if err := row.Scan(&user.Username, &profilePicture, &createdAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.ErrInternal
	}

	user.ID = userID

	if profilePicture.Valid {
		user.ProfilePicture = profilePicture.String
	}

	var err error

	user.CreatedAt, err = globaltime.Parse(createdAt)
	if err != nil {
		return nil, errors.ErrInternal
	}

	return &user, nil
}

func (repository *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	row := repository.Database.QueryRow("SELECT user_id, username, profile_picture, created_at FROM users WHERE username = ?", username)

	var user models.User

	var (
		userID, createdAt string
		profilePicture    sql.NullString
	)

	if err := row.Scan(&userID, &user.Username, &profilePicture, &createdAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.ErrInternal
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.ErrInternal
	}

	user.ID = uid

	if profilePicture.Valid {
		user.ProfilePicture = profilePicture.String
	}

	user.CreatedAt, err = globaltime.Parse(createdAt)
	if err != nil {
		return nil, errors.ErrInternal
	}

	return &user, nil
}

func (repository *UserRepository) CreateUser(username string) (uuid.UUID, error) {
	userID := uuid.New()
	createdAt := globaltime.Now()

	_, err := repository.Database.Exec("INSERT INTO users (user_id, username, created_at) VALUES (?, ?, ?)", userID.String(), username, globaltime.Format(createdAt))
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
	var oldProfilePicture sql.NullString

	err := repository.Database.QueryRow("SELECT profile_picture FROM users WHERE user_id = ?", userID.String()).Scan(&oldProfilePicture)
	if err != nil && !stdErrors.Is(err, sql.ErrNoRows) {
		return errors.ErrInternal
	}

	_, err = repository.Database.Exec("UPDATE users SET profile_picture = ? WHERE user_id = ?", sql.NullString{String: profilePicture, Valid: profilePicture != ""}, userID.String())
	if err != nil {
		return errors.ErrInternal
	}

	if oldProfilePicture.Valid && oldProfilePicture.String != "" {
		if err := os.Remove("." + oldProfilePicture.String); err != nil && !os.IsNotExist(err) {
			return errors.ErrInternal
		}
	}

	return nil
}
