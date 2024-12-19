package repository

import (
	"Eccomerce-website/internals/core/common/utils"
	"Eccomerce-website/internals/core/dto"
	"Eccomerce-website/internals/core/entity"
	"Eccomerce-website/internals/core/model/request"
	"Eccomerce-website/internals/core/port/repository"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type sizeRepository struct {
	db       repository.Database
	dbLogger *zap.Logger
}

func NewSizeRepository(db repository.Database, dbLogger *zap.Logger) repository.SizeRepository {
	return &sizeRepository{
		db:       db,
		dbLogger: dbLogger,
	}
}

func (s sizeRepository) InsertSize(ctx context.Context, size dto.Size, requestID string) (*int, error) {
	DB := s.db.GetDB()

	query := "SELECT COUNT(*) FROM size WHERE size_name = ? AND product_item_id = ?"
	var count int
	if err := DB.QueryRowContext(ctx, query, size.SizeName, size.ProductItemID).Scan(&count); err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "unable to read COUNT in the query").WithProperty(entity.StatusCode, 500)
		s.dbLogger.Error("failed to read 'COUNT' in the query",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "InsertSize"),
			zap.String("requestID", requestID),
			zap.String("query", query),
			zap.String("sizeName", size.SizeName),
			zap.Int("productItemId", size.ProductItemID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, errorResponse
	}

	if count > 0 {
		err := fmt.Errorf("size with size_name '%s' and product_item_id '%d' already exists", size.SizeName, size.ProductItemID)
		errorResponse := entity.DuplicateEntry.Wrap(err, "conflict error").WithProperty(entity.StatusCode, 409)
		s.dbLogger.Error("duplicate entry",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "InsertSize"),
			zap.String("requestID", requestID),
			zap.Any("requestData", size),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, errorResponse
	}

	insertQuery := `INSERT INTO size (product_item_id, size_name, price, discount, qty_in_stock) VALUES (?, ?, ?, ?, ?)`
	result, err := DB.ExecContext(ctx, insertQuery, size.ProductItemID, size.SizeName, size.Price, size.Discount, size.QtyInStock)
	if err != nil {
		errorResponse := entity.UnableToSave.Wrap(err, "failed to insert product size").WithProperty(entity.StatusCode, 500)
		s.dbLogger.Error("failed to create product size",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "InsertSize"),
			zap.String("requestID", requestID),
			zap.String("query", insertQuery),
			zap.Any("requestData", size),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, errorResponse
	}

	id64, err := result.LastInsertId()
	if err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "failed to get the inserted id").WithProperty(entity.StatusCode, 500)
		s.dbLogger.Error("unable to get lastInserted id",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "InsertSize"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, errorResponse
	}

	id := int(id64)

	return &id, nil
}

func (s sizeRepository) ListSizes(ctx context.Context, offset, limit int, requestID string) ([]utils.Size, error) {
	var sizes []utils.Size
	DB := s.db.GetDB()

	query := `SELECT size_id, product_item_id, size_name, price, discount, qty_in_stock, created_at, updated_at, deleted_at FROM size WHERE deleted_at IS NULL ORDER BY size_id LIMIT ? OFFSET ?`

	rows, err := DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		errorResponse := entity.UnableToFindResource.Wrap(err, "failed to get list of sizes").WithProperty(entity.StatusCode, 404)
		s.dbLogger.Error("sizes not found",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "ListSizes"),
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
		var size utils.Size

		if err := rows.Scan(&size.ID, &size.ProductItemId, &size.SizeName, &size.Price, &size.Discount, &size.QtyInStock, &size.CreatedAt, &size.UpdatedAt, &size.DeletedAt); err != nil {
			errorResponse := entity.UnableToRead.Wrap(err, "failed to scan size data").WithProperty(entity.StatusCode, 500)
			s.dbLogger.Error("unable to scan size data",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "databaseLayer"),
				zap.String("function", "ListSizes"),
				zap.String("requestID", requestID),
				zap.Error(errorResponse),
				zap.Stack("stacktrace"),
			)
			return nil, errorResponse
		}

		sizes = append(sizes, size)
	}

	if err := rows.Err(); err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "db rows error").WithProperty(entity.StatusCode, 500)
		s.dbLogger.Error("db rows error",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "ListSizes"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, errorResponse
	}

	return sizes, nil

}

func (s sizeRepository) GetSizeById(ctx context.Context, id int, requestID string) (utils.Size, error) {
	var size utils.Size
	DB := s.db.GetDB()

	query := `SELECT size_id, product_item_id, size_name, price, discount, qty_in_stock, created_at, updated_at, deleted_at FROM size WHERE size_id = ? AND deleted_at IS NULL`
	if err := DB.QueryRowContext(ctx, query, id).Scan(&size.ID, &size.ProductItemId, &size.SizeName, &size.Price, &size.Discount, &size.QtyInStock, &size.CreatedAt, &size.UpdatedAt, &size.DeletedAt); err != nil {
		errorMessage := fmt.Sprintf("size with size_id '%d' not found", id)
		errorResponse := entity.UnableToFindResource.Wrap(err, errorMessage).WithProperty(entity.StatusCode, 404)
		s.dbLogger.Error("size not found",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "GetSizeById"),
			zap.String("requestID", requestID),
			zap.String("query", query),
			zap.Int("sizeID", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return utils.Size{}, errorResponse
	}

	return size, nil
}

func (s sizeRepository) EditSizeById(ctx context.Context, id int, size request.UpdateSize, requestID string) (utils.Size, error) {
	DB := s.db.GetDB()
	zeroDecimal := decimal.NewFromInt(0)
	var updateFields []string
	var values []interface{}

	query := "SELECT COUNT(*) FROM size WHERE size_id = ? AND deleted_at IS NULL"
	var count int
	if err := DB.QueryRowContext(ctx, query, id).Scan(&count); err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "unable to read COUNT in the query").WithProperty(entity.StatusCode, 500)
		s.dbLogger.Error("failed to read 'COUNT' in the query",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "EditSizeById"),
			zap.String("requestID", requestID),
			zap.String("query", query),
			zap.Int("sizeID", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return utils.Size{}, errorResponse
	}

	if count == 0 {
		err := fmt.Errorf("size with size_id '%d' not found", id)
		errorResponse := entity.UnableToFindResource.Wrap(err, "failed to get product size by size_id").WithProperty(entity.StatusCode, 404)
		s.dbLogger.Error("product size not found",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "EditSizeById"),
			zap.String("requestID", requestID),
			zap.Int("sizeID", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return utils.Size{}, errorResponse
	}

	if size.SizeName != "" {
		updateFields = append(updateFields, "size_name = ?")
		values = append(values, strings.ToUpper(size.SizeName))
	}

	if !size.Price.Equal(zeroDecimal) {
		updateFields = append(updateFields, "price = ?")
		values = append(values, size.Price)
	}

	if !size.Discount.Equal(zeroDecimal) {
		updateFields = append(updateFields, "discount = ?")
		values = append(values, size.Discount)
	}

	if size.QtyInStock != 0 {
		updateFields = append(updateFields, "qty_in_stock = ?")
		values = append(values, size.QtyInStock)
	}

	if len(updateFields) == 0 {
		err := errors.New("failed to update size:No fields provided for update.Please provide at least one field to update")
		errorResponse := entity.BadRequest.Wrap(err, "updated fields are required").WithProperty(entity.StatusCode, 400)
		s.dbLogger.Error("the updateSize fields are empty",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "EditSizeById"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return utils.Size{}, errorResponse
	}

	if len(values) > 0 {
		updateFields = append(updateFields, "updated_at = ?")
		values = append(values, time.Now())
	}

	updateQuery := fmt.Sprintf("UPDATE size SET %s WHERE size_id = ? AND deleted_at IS NULL", strings.Join(updateFields, ", "))
	values = append(values, id)

	if _, err := DB.ExecContext(ctx, updateQuery, values...); err != nil {
		errorResponse := entity.UnableToSave.Wrap(err, "failed to update size data").WithProperty(entity.StatusCode, 500)
		s.dbLogger.Error("failed to edit size data",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "EditSizeById"),
			zap.String("requestID", requestID),
			zap.String("query", updateQuery),
			zap.Any("requestData", size),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return utils.Size{}, errorResponse
	}

	updatedSize, err := s.GetSizeById(ctx, id, requestID)
	if err != nil {
		return utils.Size{}, err
	}

	return updatedSize, nil
}

func (s sizeRepository) DeleteSizeById(ctx context.Context, id int, requestID string) error {
	DB := s.db.GetDB()

	query := "SELECT COUNT(*) FROM size WHERE size_id = ? AND deleted_at IS NULL"
	var count int
	if err := DB.QueryRowContext(ctx, query, id).Scan(&count); err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "unable to read COUNT in the query").WithProperty(entity.StatusCode, 500)
		s.dbLogger.Error("failed to read 'COUNT' in the query",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "DeleteSizeById"),
			zap.String("requestID", requestID),
			zap.String("query", query),
			zap.Int("sizeID", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return errorResponse
	}

	if count == 0 {
		err := fmt.Errorf("size with size_id '%d' not found", id)
		errorResponse := entity.UnableToFindResource.Wrap(err, "failed to get size by id").WithProperty(entity.StatusCode, 404)
		s.dbLogger.Error("size not found",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "DeleteSizeById"),
			zap.String("requestID", requestID),
			zap.Int("sizeID", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return errorResponse
	}

	deleteSizeQuery := `DELETE FROM size WHERE size_id = ? AND deleted_at IS NULL`
	if _, err := DB.ExecContext(ctx, deleteSizeQuery, id); err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "failed to delete size by id").WithProperty(entity.StatusCode, 500)
		s.dbLogger.Error("unable to delete size",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "DeleteSizeById"),
			zap.String("requestID", requestID),
			zap.String("query", deleteSizeQuery),
			zap.Int("sizeID", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return errorResponse
	}

	return nil
}
