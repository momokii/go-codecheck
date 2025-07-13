package handlers

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/momokii/go-codecheck/backend/core/databases"
	"github.com/momokii/go-codecheck/backend/core/models"
	"github.com/momokii/go-codecheck/backend/core/repository/user"
	"github.com/momokii/go-codecheck/backend/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

const (
	SALT_SIZE             = 16
	TOKEN_EXPIRATION_TIME = time.Hour * 8
)

type UserHandler struct {
	sqliteDB *databases.SQLiteDB
	UserRepo *user.UserRepository
}

func NewUserHandler(sqliteDB *databases.SQLiteDB, userRepo *user.UserRepository) *UserHandler {
	return &UserHandler{
		sqliteDB: sqliteDB,
		UserRepo: userRepo,
	}
}

func (h *UserHandler) GetAndValidateUserDataByToken(token string) (models.User, error) {

	var user models.User

	if token == "" {
		return user, fmt.Errorf("Invalid session token")
	}

	if err, _ := h.sqliteDB.Transaction(context.Background(), func(tx *sql.Tx) (err error, statusCode int) {

		// hash token from user
		hash := sha256.Sum256([]byte(token))
		hashedToken := hex.EncodeToString(hash[:])

		userCheck, err := h.UserRepo.FindByToken(tx, hashedToken)
		if err != nil {
			return fmt.Errorf("Error finding user by token: %w", err), http.StatusInternalServerError
		}

		if userCheck.Id == 0 {
			return fmt.Errorf("Invalid Token/ Token Expired"), http.StatusUnauthorized
		}

		// check if token is expired or not
		expiredTime, err := time.Parse(time.RFC3339, userCheck.SessionExpired)
		if err != nil {
			return fmt.Errorf("Error parsing session expired time: %w", err), http.StatusInternalServerError
		}

		if time.Now().After(expiredTime) {
			// if expired, delete the session data with update it the user data
			emptyString := ""
			if err := h.UserRepo.Update(
				tx,
				&models.UserUpdate{
					Id:             userCheck.Id,
					SessionToken:   &emptyString,
					SessionExpired: &emptyString,
				},
			); err != nil {
				return fmt.Errorf("Error updating user session: %w", err), http.StatusInternalServerError
			}

			return fmt.Errorf("Token Expired"), http.StatusUnauthorized
		}

		// if not expired, so give the user information for basic data
		user = models.User{
			Id:        userCheck.Id,
			Username:  userCheck.Username,
			CreatedAt: userCheck.CreatedAt,
			UpdatedAt: userCheck.UpdatedAt,
		}

		return nil, http.StatusOK

	}); err != nil {
		return user, fmt.Errorf("Error starting transaction: %w", err)
	}

	return user, nil
}

func (h *UserHandler) UserLogin(userLogin models.UserLogin) (userToken string, err error) {

	if err := utils.ValidateStruct(userLogin); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Username":
				return userToken, fmt.Errorf("Login Error: %s", err)
			case "Password":
				return userToken, fmt.Errorf("Login Error: %s", err)
			}
		}
	}

	// start checking and tx
	if err, _ := h.sqliteDB.Transaction(context.Background(), func(tx *sql.Tx) (err error, statusCode int) {

		// check user by username
		userCheck, err := h.UserRepo.FindByUsername(tx, userLogin.Username)
		if err != nil {
			return fmt.Errorf("Error checking user: %w", err), http.StatusInternalServerError
		}

		if userCheck.Id == 0 {
			return fmt.Errorf("Username/Password Wrong"), http.StatusBadRequest
		}

		// check password validation
		if err := bcrypt.CompareHashAndPassword([]byte(userCheck.Password), []byte(userLogin.Password)); err != nil {
			return fmt.Errorf("Username/Password Wrong"), http.StatusBadRequest
		}

		// generate session token
		rawSessionToken, err := utils.GenerateTokenSession()
		if err != nil {
			return fmt.Errorf("Error generating session token: %w", err), http.StatusInternalServerError
		}

		// hash token
		hashSessionToken := sha256.Sum256([]byte(rawSessionToken))
		hashedSessionToken := hex.EncodeToString(hashSessionToken[:])

		// expired token time
		expiredTokenTime := time.Now().Add(TOKEN_EXPIRATION_TIME).Format(time.RFC3339)

		// update user to add the new session data
		if err := h.UserRepo.Update(
			tx,
			&models.UserUpdate{
				Id:             userCheck.Id,
				SessionToken:   &hashedSessionToken,
				SessionExpired: &expiredTokenTime,
			},
		); err != nil {
			return fmt.Errorf("Error updating user session: %w", err), http.StatusInternalServerError
		}

		// give the raw token for user FE saved it
		userToken = rawSessionToken

		return nil, http.StatusOK

	}); err != nil {
		return "", err
	}

	return userToken, nil
}

func (h *UserHandler) UserUpdatePassword(userUpdate models.UserUpdate, isCompletedSetup bool) error {

	if err := utils.ValidateStruct(userUpdate); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Password":
				return fmt.Errorf("Error: %s", err)
			}
		}
	}

	if err, _ := h.sqliteDB.Transaction(context.Background(), func(tx *sql.Tx) (err error, statusCode int) {

		// user check
		userCheck, err := h.UserRepo.FindById(tx, userUpdate.Id)
		if err != nil {
			return fmt.Errorf("Error: %w", err), http.StatusInternalServerError
		}

		if userCheck.Id == 0 {
			return fmt.Errorf("User with ID %d not found", userUpdate.Id), http.StatusNotFound
		}

		// hash new pass
		hashedNewPass, err := bcrypt.GenerateFromPassword([]byte(userUpdate.Password), SALT_SIZE)
		if err != nil {
			return fmt.Errorf("Error hashing password: %w", err), http.StatusInternalServerError
		}

		// update user data
		if err := h.UserRepo.Update(
			tx,
			&models.UserUpdate{
				Id:       userUpdate.Id,
				Password: string(hashedNewPass),
			},
		); err != nil {
			return fmt.Errorf("Error updating user password: %w", err), http.StatusInternalServerError
		}

		return nil, http.StatusOK

	}); err != nil {
		return err
	}

	return nil
}

func (h *UserHandler) UserUpdateUsername(userUpdate models.UserUpdate) error {

	if err := utils.ValidateStruct(userUpdate); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Username":
				return fmt.Errorf("Error: %s", err)
			}
		}
	}

	// start checking and tx
	if err, _ := h.sqliteDB.Transaction(context.Background(), func(tx *sql.Tx) (err error, statusCode int) {

		// just in case if needed in future improvement
		// checkk user username
		userCheck, err := h.UserRepo.FindById(tx, userUpdate.Id)
		if err != nil {
			return fmt.Errorf("Error checking user: %w", err), http.StatusInternalServerError
		}

		if userCheck.Id != 0 && userCheck.Id != userUpdate.Id {
			return fmt.Errorf("Username %s is already taken", userUpdate.Username), http.StatusConflict
		}

		// change username
		if err := h.UserRepo.Update(
			tx,
			&models.UserUpdate{
				Id:       userUpdate.Id,
				Username: userUpdate.Username,
			},
		); err != nil {
			return fmt.Errorf("Error updating username: %w", err), http.StatusInternalServerError
		}

		return nil, http.StatusOK

	}); err != nil {
		return err
	}

	return nil
}

func (h *UserHandler) UserDeleteSession(userId int) error {

	if userId <= 0 {
		return fmt.Errorf("Invalid user ID")
	}

	// for now and simplicity, the logout process will just clear the session token
	if err, _ := h.sqliteDB.Transaction(context.Background(), func(tx *sql.Tx) (err error, statusCode int) {

		// check user by user id
		userCheck, err := h.UserRepo.FindById(tx, userId)
		if err != nil {
			return fmt.Errorf("Error checking user: %w", err), http.StatusInternalServerError
		}

		if userCheck.Id == 0 {
			return fmt.Errorf("User with ID %d not found", userId), http.StatusNotFound
		}

		// valid user, delete the session
		emptyString := ""
		if err := h.UserRepo.Update(
			tx,
			&models.UserUpdate{
				Id:             userCheck.Id,
				SessionToken:   &emptyString,
				SessionExpired: &emptyString,
			},
		); err != nil {
			return fmt.Errorf("Error updating user session: %w", err), http.StatusInternalServerError
		}

		return nil, http.StatusOK

	}); err != nil {
		return err
	}

	return nil
}
