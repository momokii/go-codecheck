package user

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/momokii/go-codecheck/backend/core/models"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) FindByToken(tx *sql.Tx, token string) (*models.User, error) {
	var user models.User

	query := "SELECT id, username, password, is_completed_setup, session_token, session_expired, created_at, updated_at FROM users WHERE session_token = ?"

	if err := tx.QueryRow(query, token).Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.IsCompletedSetup,
		&user.SessionToken,
		&user.SessionExpired,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return &user, nil // User not found
		}
		return nil, err // Other error
	}

	return &user, nil
}

func (r *UserRepository) FindById(tx *sql.Tx, id int) (*models.User, error) {
	var user models.User

	query := "SELECT id, username, password, is_completed_setup, session_token, session_expired, created_at, updated_at FROM users WHERE id = ?"

	if err := tx.QueryRow(query, id).Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.IsCompletedSetup,
		&user.SessionToken,
		&user.SessionExpired,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return &user, nil // User not found
		}
		return nil, err // Other error
	}

	return &user, nil
}

func (r *UserRepository) FindByUsername(tx *sql.Tx, username string) (*models.User, error) {

	var user models.User

	query := "SELECT id, username, password, is_completed_setup, session_token, session_expired, created_at, updated_at FROM users WHERE username = ?"

	if err := tx.QueryRow(query, username).Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.IsCompletedSetup,
		&user.SessionToken,
		&user.SessionExpired,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return &user, nil // User not found
		}
		return nil, err // Other error
	}

	return &user, nil
}

func (r *UserRepository) Update(tx *sql.Tx, userUpdate *models.UserUpdate) error {
	setParts := []string{}
	args := []interface{}{}

	// Profile fields
	if userUpdate.Username != "" {
		setParts = append(setParts, "username = ?")
		args = append(args, userUpdate.Username)
	}

	if userUpdate.Password != "" {
		setParts = append(setParts, "password = ?")
		args = append(args, userUpdate.Password)
	}

	if userUpdate.IsCompletedSetup != nil {
		setParts = append(setParts, "is_completed_setup = ?")
		args = append(args, *userUpdate.IsCompletedSetup)
	}

	// Session fields (include in unified update)
	if userUpdate.SessionToken != nil {
		setParts = append(setParts, "session_token = ?")
		args = append(args, *userUpdate.SessionToken)
	}

	if userUpdate.SessionExpired != nil {
		setParts = append(setParts, "session_expired = ?")
		args = append(args, *userUpdate.SessionExpired)
	}

	if len(setParts) == 0 {
		return errors.New("no fields to update")
	}

	setParts = append(setParts, "updated_at = datetime('now')")
	args = append(args, userUpdate.Id)

	query := fmt.Sprintf("UPDATE users SET %s WHERE id = ?", strings.Join(setParts, ", "))
	_, err := tx.Exec(query, args...)
	return err
}
