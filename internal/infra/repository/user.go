package repository

import (
	"Eccomerce-website/internal/core/common/utils"
	validationdata "Eccomerce-website/internal/core/common/utils/validationData"
	"Eccomerce-website/internal/core/dto"
	"Eccomerce-website/internal/core/entity"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/port/repository"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
)

type userRepository struct {
	db       repository.Database
	dbLogger *zap.Logger
}

func NewUserRepository(db repository.Database, dbLogger *zap.Logger) repository.UserRepository {
	return &userRepository{
		db:       db,
		dbLogger: dbLogger,
	}
}

func (u userRepository) InsertUser(ctx context.Context, user dto.User, requestID string) (int, error) {
	user.Email = strings.ToLower(user.Email)
	DB := u.db.GetDB()

	var count int
	if err := DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE username = ? OR email = ?", user.Username, user.Email).Scan(&count); err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "unable to read COUNT in the query").WithProperty(entity.StatusCode, 500)
		u.dbLogger.Error("failed to read 'COUNT' in the query",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "InsertUser"),
			zap.String("requestID", requestID),
			zap.String("query", "SELECT COUNT(*) FROM users WHERE username = ? OR email = ?"),
			zap.String("username", user.Username),
			zap.String("email", user.Email),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return 0, errorResponse
	}
	if count > 0 {
		err := errors.New("user already exists")
		errorResponse := entity.DuplicateEntry.Wrap(err, "conflict error").WithProperty(entity.StatusCode, 409)
		u.dbLogger.Error("duplicate entry",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "InsertUser"),
			zap.String("requestID", requestID),
			zap.Any("requestData", user),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return 0, errorResponse
	}

	query := `INSERT INTO users(username, email, password, first_name, last_name, phone_number, email_verified, profile_picture)
	          VALUES(?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := DB.ExecContext(ctx, query, user.Username, user.Email, user.Password, user.FirstName, user.LastName, user.PhoneNumber, user.EmailVerified, user.ProfilePicture)
	if err != nil {
		errorResponse := entity.UnableToSave.Wrap(err, "failed to insert user").WithProperty(entity.StatusCode, 500)
		u.dbLogger.Error("failed to create user",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "InsertUser"),
			zap.String("requestID", requestID),
			zap.String("query", "INSERT INTO users(username, email, password, first_name, last_name, phone_number, email_verified, profile_picture)VALUES(?, ?, ?, ?, ?, ?, ?, ?)"),
			zap.Any("requestData", user),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return 0, errorResponse
	}

	id, err := result.LastInsertId()
	if err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "failed to get the inserted id").WithProperty(entity.StatusCode, 500)
		u.dbLogger.Error("unable to get lastInserted id",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "InsertUser"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return 0, errorResponse
	}

	return int(id), nil
}

func (u userRepository) Authentication(ctx context.Context, request request.LoginRequest, requestID string) (utils.User, error) {
	var user utils.User
	email := strings.ToLower(request.Email)
	password := request.Password
	DB := u.db.GetDB()

	query := `
	 SELECT user_id, username, email, password, first_name,
	 last_name, phone_number, profile_picture, email_verified, 
	 role, created_at, updated_at FROM users WHERE email = ? AND deleted_at IS NULL
	 `

	if err := DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.FirstName,
		&user.LastName, &user.PhoneNumber, &user.ProfilePicture,
		&user.EmailVerified, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
		errorMessage := fmt.Sprintf("user with email %s is not found", email)
		errorResponse := entity.UnableToFindResource.Wrap(err, errorMessage).WithProperty(entity.StatusCode, 404)
		u.dbLogger.Error("user not found",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "Authentication"),
			zap.String("requestID", requestID),
			zap.String("query", "SELECT user_id, username, email, password, first_name, last_name, phone_number, profile_picture, email_verified,  role, created_at, updated_at FROM users WHERE email = ? AND deleted_at IS NULL"),
			zap.Any("requestData", request),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return utils.User{}, errorResponse
	}

	match, err := validationdata.MatchPassword(user.Password, password, u.dbLogger, requestID)
	if err != nil || !match {
		return utils.User{}, err
	}

	return user, nil

}

func (u userRepository) ListUsers(ctx context.Context, offset, perPage int, requestID string) ([]utils.User, error) {
	var users []utils.User
	DB := u.db.GetDB()

	query := `SELECT user_id, username, email, password, first_name,
	last_name, phone_number, profile_picture, email_verified, 
	role, created_at, updated_at FROM users WHERE deleted_at IS NULL ORDER BY user_id LIMIT ? OFFSET ?`

	rows, err := DB.QueryContext(ctx, query, perPage, offset)
	if err != nil {
		errorResponse := entity.UnableToFindResource.Wrap(err, "failed to get list of users").WithProperty(entity.StatusCode, 404)
		u.dbLogger.Error("unable to get users list",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "ListUsers"),
			zap.String("requestID", requestID),
			zap.String("query", "SELECT user_id, username, email, password, first_name,last_name, phone_number, profile_picture, email_verified, role, created_at, updated_at FROM users WHERE deleted_at IS NULL ORDER BY user_id LIMIT ? OFFSET ?"),
			zap.Int("offset", offset),
			zap.Int("perPage", perPage),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, errorResponse
	}
	defer rows.Close()

	for rows.Next() {
		var user utils.User
		if err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.Password, &user.FirstName,
			&user.LastName, &user.PhoneNumber, &user.ProfilePicture,
			&user.EmailVerified, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
			errorResponse := entity.UnableToRead.Wrap(err, "failed to scan users data").WithProperty(entity.StatusCode, 500)
			u.dbLogger.Error("unable to scan user data",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "databaseLayer"),
				zap.String("function", "ListUsers"),
				zap.String("requestID", requestID),
				zap.Error(errorResponse),
				zap.Stack("stacktrace"),
			)
			return nil, errorResponse
		}

		users = append(users, user)

	}

	if err := rows.Err(); err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "db rows error").WithProperty(entity.StatusCode, 500)
		u.dbLogger.Error("db rows error",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "ListUsers"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, errorResponse
	}

	return users, nil
}

func (u userRepository) GetUserById(ctx context.Context, id int, requestID string) (utils.User, error) {
	var user utils.User
	DB := u.db.GetDB()
	query := `SELECT user_id, username, email, password, first_name,
	last_name, phone_number, profile_picture, email_verified, 
	role, created_at, updated_at FROM users WHERE user_id = ? AND deleted_at IS NULL`

	if err := DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.FirstName,
		&user.LastName, &user.PhoneNumber, &user.ProfilePicture,
		&user.EmailVerified, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	); err != nil {
		errorMessage := fmt.Sprintf("user with id %d not found", id)
		errorResponse := entity.UnableToFindResource.Wrap(err, errorMessage).WithProperty(entity.StatusCode, 404)
		u.dbLogger.Error("user not found",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "GetUserById"),
			zap.String("requestID", requestID),
			zap.String("query", "SELECT user_id, username, email, password, first_name,last_name, phone_number, profile_picture, email_verified, role, created_at, updated_at FROM users WHERE user_id = ? AND deleted_at IS NULL"),
			zap.Int("id", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return utils.User{}, errorResponse
	}

	return user, nil
}

func (u userRepository) EditUserById(ctx context.Context, id int, user request.UpdateUser, requestID string) (utils.User, error) {
	var updateFields []string
	var values []interface{}
	DB := u.db.GetDB()

	if user.Username != "" {
		updateFields = append(updateFields, "username = ?")
		values = append(values, user.Username)
	}
	if user.Email != "" {
		updateFields = append(updateFields, "email = ?")
		values = append(values, user.Email)
	}
	if user.Password != "" {
		updateFields = append(updateFields, "password = ?")
		values = append(values, user.Password)
	}
	if user.FirstName != "" {
		updateFields = append(updateFields, "first_name = ?")
		values = append(values, user.FirstName)
	}
	if user.LastName != "" {
		updateFields = append(updateFields, "last_name = ?")
		values = append(values, user.LastName)
	}
	if user.PhoneNumber != "" {
		updateFields = append(updateFields, "phone_number = ?")
		values = append(values, user.PhoneNumber)
	}
	if user.ProfilePicture != "" {
		updateFields = append(updateFields, "profile_picture = ?")
		values = append(values, user.ProfilePicture)
	}

	if len(updateFields) == 0 {
		err := errors.New("failed to update user:No fields provided for update.Please provide at least one field to update")
		errorResponse := entity.BadRequest.Wrap(err, "updated fields are required").WithProperty(entity.StatusCode, 400)
		u.dbLogger.Error("the updata fields are empty",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "EditUserById"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return utils.User{}, errorResponse
	}

	if len(values) > 0 {
		updateFields = append(updateFields, "updated_at = ?")
		values = append(values, time.Now())
	}

	query := fmt.Sprintf("UPDATE users SET %s WHERE user_id = ? AND deleted_at IS NULL", strings.Join(updateFields, ", "))
	values = append(values, id)

	if _, err := DB.ExecContext(ctx, query, values...); err != nil {
		errorResponse := entity.UnableToSave.Wrap(err, "failed to update user data").WithProperty(entity.StatusCode, 500)
		u.dbLogger.Error("failed to edit user data",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "EditUserById"),
			zap.String("requestID", requestID),
			zap.String("query", "UPDATE users SET %s WHERE user_id = ? AND deleted_at IS NULL"),
			zap.Any("requestData", user),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return utils.User{}, errorResponse
	}

	updateUser, err := u.GetUserById(ctx, id, requestID)
	if err != nil {
		return utils.User{}, err
	}

	return updateUser, nil
}

func (u userRepository) DeleteUserById(ctx context.Context, id int, requestID string) error {
	DB := u.db.GetDB()

	var count int
	if err := DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE user_id = ?", id).Scan(&count); err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "unable to read COUNT in the query").WithProperty(entity.StatusCode, 500)
		u.dbLogger.Error("failed to read 'COUNT' in the query",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "DeleteUserById"),
			zap.String("requestID", requestID),
			zap.String("query", "SELECT COUNT(*) FROM users WHERE user_id = ?"),
			zap.Int("id", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return errorResponse
	}
	if count == 0 {
		err := fmt.Errorf("user with user_id '%d' not found", id)
		errorResponse := entity.UnableToFindResource.Wrap(err, "failed to get user by id").WithProperty(entity.StatusCode, 404)
		u.dbLogger.Error("user not found",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "DeleteUserById"),
			zap.String("requestID", requestID),
			zap.Int("id", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return errorResponse
	}

	query := `DELETE FROM users WHERE user_id = ?`
	if _, err := DB.ExecContext(ctx, query, id); err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "failed to delete user by id").WithProperty(entity.StatusCode, 500)
		u.dbLogger.Error("unable to delete user",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "DeleteUserById"),
			zap.String("requestID", requestID),
			zap.String("query", "DELETE FROM users WHERE user_id = ?"),
			zap.Int("id", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return errorResponse
	}

	return nil
}
