package repository

import (
	"Eccomerce-website/internals/core/common/utils"
	"Eccomerce-website/internals/core/dto"
	"Eccomerce-website/internals/core/entity"
	"Eccomerce-website/internals/core/port/repository"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
)

type colorRepository struct {
	db       repository.Database
	dbLogger *zap.Logger
}

func NewColorRepository(db repository.Database, dbLogger *zap.Logger) repository.ColorRepository {
	return &colorRepository{
		db:       db,
		dbLogger: dbLogger,
	}
}

func (c colorRepository) InsertColor(ctx context.Context, color dto.Color, requestID string) (*int, error) {
	DB := c.db.GetDB()

	query := "SELECT COUNT(*) FROM color WHERE color_name = ? AND deleted_at IS NULL"
	var count int
	if err := DB.QueryRowContext(ctx, query, color.Name).Scan(&count); err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "unable to read COUNT in the query").WithProperty(entity.StatusCode, 500)
		c.dbLogger.Error("failed to read 'COUNT' in the query",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "InsertColor"),
			zap.String("requestID", requestID),
			zap.String("query", query),
			zap.String("colorName", color.Name),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, errorResponse
	}

	if count > 0 {
		err := fmt.Errorf("color with color_name '%s' already exists", color.Name)
		errorResponse := entity.DuplicateEntry.Wrap(err, "conflict error").WithProperty(entity.StatusCode, 409)
		c.dbLogger.Error("duplicate entry",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "InsertColor"),
			zap.String("requestID", requestID),
			zap.Any("requestData", color),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, errorResponse
	}

	insertQuery := `INSERT INTO color(color_name) VALUES(?)`
	result, err := DB.ExecContext(ctx, insertQuery, color.Name)
	if err != nil {
		errorResponse := entity.UnableToSave.Wrap(err, "failed to insert color").WithProperty(entity.StatusCode, 500)
		c.dbLogger.Error("failed to create product color",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "InsertColor"),
			zap.String("requestID", requestID),
			zap.String("query", insertQuery),
			zap.Any("requestData", color),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, errorResponse
	}

	id64, err := result.LastInsertId()
	if err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "failed to get the inserted id").WithProperty(entity.StatusCode, 500)
		c.dbLogger.Error("unable to get lastInserted id",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "InsertColor"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, errorResponse
	}

	id := int(id64)

	return &id, nil
}

func (c colorRepository) ListColors(ctx context.Context, offset, limit int, requestID string) ([]utils.Color, error) {
	var colors []utils.Color
	DB := c.db.GetDB()

	query := `SELECT color_id, color_name, created_at, updated_at, deleted_at FROM color WHERE deleted_at IS NULL ORDER BY color_id LIMIT ? OFFSET ?`
	rows, err := DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		errorResponse := entity.UnableToFindResource.Wrap(err, "failed to get list of colors").WithProperty(entity.StatusCode, 404)
		c.dbLogger.Error("colors not found",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "ListColors"),
			zap.String("requestID", requestID),
			zap.String("query", query),
			zap.Int("offset", offset),
			zap.Int("limit", limit),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, errorResponse
	}

	defer rows.Close()

	for rows.Next() {
		var color utils.Color
		if err := rows.Scan(&color.ID, &color.ColorName, &color.CreatedAt, &color.UpdatedAt, &color.DeletedAt); err != nil {
			errorResponse := entity.UnableToRead.Wrap(err, "failed to scan color data").WithProperty(entity.StatusCode, 500)
			c.dbLogger.Error("unable to scan color data",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "databaseLayer"),
				zap.String("function", "ListColors"),
				zap.String("requestID", requestID),
				zap.Error(errorResponse),
				zap.Stack("stacktrace"),
			)
			return nil, errorResponse
		}

		colors = append(colors, color)
	}

	if err := rows.Err(); err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "db rows error").WithProperty(entity.StatusCode, 500)
		c.dbLogger.Error("db rows error",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "ListColors"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, errorResponse
	}

	return colors, nil
}

func (c colorRepository) GetColorById(ctx context.Context, id int, requestID string) (utils.Color, error) {
	var color utils.Color
	DB := c.db.GetDB()

	query := `SELECT color_id, color_name, created_at, updated_at, deleted_at FROM color WHERE color_id = ? AND deleted_at IS NULL`
	if err := DB.QueryRowContext(ctx, query, id).Scan(&color.ID, &color.ColorName, &color.CreatedAt, &color.UpdatedAt, &color.DeletedAt); err != nil {
		errMessage := fmt.Sprintf("color with color_id '%d' not found", id)
		errorResponse := entity.UnableToFindResource.Wrap(err, errMessage).WithProperty(entity.StatusCode, 404)
		c.dbLogger.Error("size not found",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "GetColorById"),
			zap.String("requestID", requestID),
			zap.String("query", query),
			zap.Int("colorID", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return utils.Color{}, errorResponse
	}

	return color, nil
}

func (c colorRepository) EditColorById(ctx context.Context, id int, color utils.UpdateColor, requestID string) (utils.Color, error) {
	DB := c.db.GetDB()
	var updateFields []string
	var values []interface{}

	query := "SELECT COUNT(*) FROM color WHERE color_id = ? AND deleted_at IS NULL"
	var count int
	if err := DB.QueryRowContext(ctx, query, id).Scan(&count); err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "unable to read COUNT in the query").WithProperty(entity.StatusCode, 500)
		c.dbLogger.Error("failed to read 'COUNT' in the query",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "EditColorById"),
			zap.String("requestID", requestID),
			zap.String("query", query),
			zap.Int("colorID", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return utils.Color{}, errorResponse
	}

	if count == 0 {
		err := fmt.Errorf("color with color_id '%d' not found", id)
		errorResponse := entity.UnableToFindResource.Wrap(err, "failed to get product color by color_id").WithProperty(entity.StatusCode, 404)
		c.dbLogger.Error("product color not found",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "EditColorById"),
			zap.String("requestID", requestID),
			zap.Int("colorID", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return utils.Color{}, errorResponse
	}

	if color.ColorName != "" {
		updateFields = append(updateFields, "color_name = ?")
		values = append(values, color.ColorName)
	}

	if len(updateFields) == 0 {
		err := errors.New("failed to update color:No fields provided for update.Please provide at least one field to update")
		errorResponse := entity.BadRequest.Wrap(err, "updated fields are required").WithProperty(entity.StatusCode, 400)
		c.dbLogger.Error("the updateColor fields are empty",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "EditColorById"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return utils.Color{}, errorResponse
	}

	if len(values) > 0 {
		updateFields = append(updateFields, "updated_at = ?")
		values = append(values, time.Now())
	}

	updateQuery := fmt.Sprintf("UPDATE color SET %s WHERE color_id = ? AND deleted_at IS NULL", strings.Join(updateFields, ", "))
	values = append(values, id)
	if _, err := DB.ExecContext(ctx, updateQuery, values...); err != nil {
		errorResponse := entity.UnableToSave.Wrap(err, "failed to update color data").WithProperty(entity.StatusCode, 500)
		c.dbLogger.Error("failed to edit color data",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "EditColorById"),
			zap.String("requestID", requestID),
			zap.String("query", updateQuery),
			zap.Any("requestData", color),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return utils.Color{}, errorResponse
	}

	updatedColor, err := c.GetColorById(ctx, id, requestID)
	if err != nil {
		return utils.Color{}, err
	}

	return updatedColor, nil
}

func (c colorRepository) DeleteColorById(ctx context.Context, id int, requestID string) error {
	DB := c.db.GetDB()

	query := "SELECT COUNT(*) FROM color WHERE color_id = ?"
	var count int
	if err := DB.QueryRowContext(ctx, query, id).Scan(&count); err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "unable to read COUNT in the query").WithProperty(entity.StatusCode, 500)
		c.dbLogger.Error("failed to read 'COUNT' in the query",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "DeleteColorById"),
			zap.String("requestID", requestID),
			zap.String("query", query),
			zap.Int("colorID", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return errorResponse
	}
	if count == 0 {
		err := fmt.Errorf("color with id '%d' not found", id)
		errorResponse := entity.UnableToFindResource.Wrap(err, "failed to get color by id").WithProperty(entity.StatusCode, 404)
		c.dbLogger.Error("color not found",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "DeleteColorById"),
			zap.String("requestID", requestID),
			zap.Int("colorID", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return errorResponse
	}

	deleteQuery := `DELETE FROM color WHERE color_id = ?`
	if _, err := DB.ExecContext(ctx, deleteQuery, id); err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "failed to delete color by id").WithProperty(entity.StatusCode, 500)
		c.dbLogger.Error("unable to delete color",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "DeleteColorById"),
			zap.String("requestID", requestID),
			zap.String("query", deleteQuery),
			zap.Int("colorID", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return errorResponse
	}

	return nil
}
