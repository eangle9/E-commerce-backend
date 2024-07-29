package repository

import (
	"Eccomerce-website/internal/core/common/utils"
	"Eccomerce-website/internal/core/dto"
	"Eccomerce-website/internal/core/entity"
	"Eccomerce-website/internal/core/port/repository"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
)

type productCategoryRepository struct {
	db       repository.Database
	dbLogger *zap.Logger
}

func NewProductCategoryRepository(db repository.Database, dbLogger *zap.Logger) repository.ProductCategoryRepository {
	return &productCategoryRepository{
		db:       db,
		dbLogger: dbLogger,
	}
}

func (p productCategoryRepository) InsertProductCategory(ctx context.Context, category dto.ProductCategory, requestID string) (*int, error) {
	DB := p.db.GetDB()

	name := category.Name
	parentId := category.ParentID

	query := "SELECT COUNT(*) FROM product_category WHERE name = ? And deleted_at IS NULL"
	var count int
	if err := DB.QueryRowContext(ctx, query, name).Scan(&count); err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "unable to read COUNT in the query").WithProperty(entity.StatusCode, 500)
		p.dbLogger.Error("failed to read 'COUNT' in the query",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "InsertProductCategory"),
			zap.String("requestID", requestID),
			zap.String("query", query),
			zap.String("categoryName", name),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, err
	}

	if count > 0 {
		err := fmt.Errorf("product category with category_name '%s' is already exists", name)
		errorResponse := entity.DuplicateEntry.Wrap(err, "conflict error").WithProperty(entity.StatusCode, 409)
		p.dbLogger.Error("duplicate entry",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "InsertProductCategory"),
			zap.String("requestID", requestID),
			zap.Any("requestData", category),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, errorResponse
	}

	if parentId == 0 {
		insertQuery := `INSERT INTO product_category(name) VALUES(?)`
		result, err := DB.ExecContext(ctx, insertQuery, name)
		if err != nil {
			errorResponse := entity.UnableToSave.Wrap(err, "failed to insert product category").WithProperty(entity.StatusCode, 500)
			p.dbLogger.Error("failed to create product category",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "databaseLayer"),
				zap.String("function", "InsertProductCategory"),
				zap.String("requestID", requestID),
				zap.String("query", insertQuery),
				zap.Any("requestData", category),
				zap.Error(errorResponse),
				zap.Stack("stacktrace"),
			)
			return nil, errorResponse
		}
		id64, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}
		id := int(id64)
		return &id, nil
	}

	insertCategoryQuery := `INSERT INTO product_category(name, parent_id) VALUES(?, ?)`
	result, err := DB.ExecContext(ctx, insertCategoryQuery, name, parentId)
	if err != nil {
		errorResponse := entity.UnableToSave.Wrap(err, "failed to insert product category").WithProperty(entity.StatusCode, 500)
		p.dbLogger.Error("failed to create product category",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "InsertProductCategory"),
			zap.String("requestID", requestID),
			zap.String("query", insertCategoryQuery),
			zap.Any("requestData", category),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, errorResponse
	}

	id64, err := result.LastInsertId()
	if err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "failed to get the inserted id").WithProperty(entity.StatusCode, 500)
		p.dbLogger.Error("unable to get lastInserted id",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "InsertProductCategory"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, errorResponse
	}
	id := int(id64)
	return &id, nil

}

func (p productCategoryRepository) ListProductCategory(ctx context.Context, offset, limit int, requestID string) ([]utils.ProductCategory, error) {
	var productCategories []utils.ProductCategory
	DB := p.db.GetDB()

	query := `SELECT category_id, name, parent_id, created_at, updated_at, deleted_at FROM product_category WHERE deleted_at IS NULL ORDER BY category_id LIMIT ? OFFSET ?`
	rows, err := DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		errorResponse := entity.UnableToFindResource.Wrap(err, "failed to get list of category").WithProperty(entity.StatusCode, 404)
		p.dbLogger.Error("category not found",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "ListProductCategory"),
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
		var productCategory utils.ProductCategory
		if err := rows.Scan(&productCategory.ID, &productCategory.Name, &productCategory.ParentID, &productCategory.CreatedAt, &productCategory.UpdatedAt, &productCategory.DeletedAt); err != nil {
			errorResponse := entity.UnableToRead.Wrap(err, "failed to scan category data").WithProperty(entity.StatusCode, 500)
			p.dbLogger.Error("unable to scan category data",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "databaseLayer"),
				zap.String("function", "ListProductCategory"),
				zap.String("requestID", requestID),
				zap.Error(errorResponse),
				zap.Stack("stacktrace"),
			)
			return nil, errorResponse
		}

		productCategories = append(productCategories, productCategory)
	}

	if err := rows.Err(); err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "db rows error").WithProperty(entity.StatusCode, 500)
		p.dbLogger.Error("db rows error",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "ListProductCategory"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, errorResponse
	}

	return productCategories, nil
}

func (p productCategoryRepository) GetProductCategoryById(ctx context.Context, id int, requestID string) (utils.ProductCategory, error) {
	var category utils.ProductCategory
	DB := p.db.GetDB()

	query := `SELECT category_id, name, parent_id, created_at, updated_at, deleted_at FROM product_category WHERE category_id = ? AND deleted_at IS NULL`
	if err := DB.QueryRowContext(ctx, query, id).Scan(&category.ID, &category.Name, &category.ParentID, &category.CreatedAt, &category.UpdatedAt, &category.DeletedAt); err != nil {
		errMessage := fmt.Sprintf("product category with category_id '%d' not found", id)
		errorResponse := entity.UnableToFindResource.Wrap(err, errMessage).WithProperty(entity.StatusCode, 404)
		p.dbLogger.Error("category not found",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "GetProductCategoryById"),
			zap.String("requestID", requestID),
			zap.String("query", query),
			zap.Int("categoryID", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return utils.ProductCategory{}, errorResponse
	}

	return category, nil
}

func (p productCategoryRepository) EditProductCategoryById(ctx context.Context, id int, category utils.UpdateCategory, requestID string) (utils.ProductCategory, error) {
	DB := p.db.GetDB()
	var updateFields []string
	var values []interface{}

	query := "SELECT COUNT(*) FROM product_category WHERE category_id = ? AND deleted_at IS NULL"
	var count int
	if err := DB.QueryRowContext(ctx, query, id).Scan(&count); err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "unable to read COUNT in the query").WithProperty(entity.StatusCode, 500)
		p.dbLogger.Error("failed to read 'COUNT' in the query",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "EditProductCategoryById"),
			zap.String("requestID", requestID),
			zap.String("query", query),
			zap.Int("categoryID", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return utils.ProductCategory{}, errorResponse
	}

	if count == 0 {
		err := fmt.Errorf("category with category_id '%d' not found", id)
		errorResponse := entity.UnableToFindResource.Wrap(err, "failed to get product category by category_id").WithProperty(entity.StatusCode, 404)
		p.dbLogger.Error("product category not found",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "EditProductCategoryById"),
			zap.String("requestID", requestID),
			zap.Int("categoryID", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return utils.ProductCategory{}, errorResponse
	}

	if category.Name != "" {
		updateFields = append(updateFields, "name = ?")
		values = append(values, category.Name)
	}

	if category.ParentID != 0 {
		updateFields = append(updateFields, "parent_id = ?")
		values = append(values, category.ParentID)
	}

	if len(updateFields) == 0 {
		err := errors.New("failed to update color:No fields provided for update.Please provide at least one field to update")
		errorResponse := entity.BadRequest.Wrap(err, "updated fields are required").WithProperty(entity.StatusCode, 400)
		p.dbLogger.Error("the updateCategory fields are empty",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "EditProductCategoryById"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return utils.ProductCategory{}, errorResponse
	}

	if len(values) > 0 {
		updateFields = append(updateFields, "updated_at = ?")
		values = append(values, time.Now())
	}

	updateQuery := fmt.Sprintf("UPDATE product_category SET %s WHERE category_id = ? AND deleted_at IS NULL", strings.Join(updateFields, ", "))
	values = append(values, id)

	if _, err := DB.ExecContext(ctx, updateQuery, values...); err != nil {
		errorResponse := entity.UnableToSave.Wrap(err, "failed to update category data").WithProperty(entity.StatusCode, 500)
		p.dbLogger.Error("failed to edit category data",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "EditProductCategoryById"),
			zap.String("requestID", requestID),
			zap.String("query", updateQuery),
			zap.Any("requestData", category),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return utils.ProductCategory{}, errorResponse
	}

	productCategory, err := p.GetProductCategoryById(ctx, id, requestID)
	if err != nil {
		return utils.ProductCategory{}, err
	}

	return productCategory, nil
}

func (p productCategoryRepository) DeleteProductCategoryById(ctx context.Context, id int, requestID string) error {
	DB := p.db.GetDB()

	query := "SELECT COUNT(*) FROM product_category WHERE category_id = ?"
	var count int
	if err := DB.QueryRowContext(ctx, query, id).Scan(&count); err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "unable to read COUNT in the query").WithProperty(entity.StatusCode, 500)
		p.dbLogger.Error("failed to read 'COUNT' in the query",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "DeleteProductCategoryById"),
			zap.String("requestID", requestID),
			zap.String("query", query),
			zap.Int("categoryID", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return errorResponse
	}

	if count == 0 {
		err := fmt.Errorf("product category with id '%d' not found", id)
		errorResponse := entity.UnableToFindResource.Wrap(err, "failed to get category by id").WithProperty(entity.StatusCode, 404)
		p.dbLogger.Error("category not found",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "DeleteProductCategoryById"),
			zap.String("requestID", requestID),
			zap.Int("categoryID", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return errorResponse
	}

	deleteQuery := `DELETE FROM product_category WHERE category_id = ?`
	if _, err := DB.ExecContext(ctx, deleteQuery, id); err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "failed to delete category by id").WithProperty(entity.StatusCode, 500)
		p.dbLogger.Error("unable to delete category",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "DeleteProductCategoryById"),
			zap.String("requestID", requestID),
			zap.String("query", deleteQuery),
			zap.Int("categoryID", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return errorResponse
	}

	return nil

}

// if err := DB.QueryRow("SELECT deleted_at FROM product_category WHERE category_id = ?", id).Scan(&deleted_at); err != nil {
// 	errType := entity.NotFoundError
// 	err := fmt.Errorf("product category with id '%d' not found", id)
// 	status := http.StatusNotFound
// 	return "", status, errType, err
// }

// if deleted_at != nil {
// 	errType := "CONFLICT_ERROR"
// 	err := errors.New("can't delete already deleted product category")
// 	status := http.StatusConflict
// 	return "", status, errType, err
// }

// query := `UPDATE product_category SET deleted_at = ? WHERE category_id = ?`
// if _, err := DB.Exec(query, time.Now(), id); err != nil {
// 	errType := entity.InternalError
// 	status := http.StatusInternalServerError
// 	return "", status, errType, err
// }
