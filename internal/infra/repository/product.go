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

type productRepository struct {
	db       repository.Database
	dbLogger *zap.Logger
}

func NewProductRepository(db repository.Database, dbLogger *zap.Logger) repository.ProductRepository {
	return &productRepository{
		db:       db,
		dbLogger: dbLogger,
	}
}

func (p productRepository) InsertProduct(ctx context.Context, product dto.Product, requestID string) (*int, error) {
	DB := p.db.GetDB()
	categoryId := product.CategoryID
	brand := product.Brand
	name := product.ProductName
	description := product.Description

	productQuery := "SELECT COUNT(*) FROM product WHERE product_name = ? AND deleted_at IS NULL"
	var count int
	if err := DB.QueryRowContext(ctx, productQuery, name).Scan(&count); err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "failed to read COUNT in the query").WithProperty(entity.StatusCode, 500)
		p.dbLogger.Error("unable to read COUNT in the query",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "InsertProduct"),
			zap.String("requestID", requestID),
			zap.String("query", productQuery),
			zap.String("productName", name),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)

		return nil, errorResponse
	}

	if count > 0 {
		err := fmt.Errorf("product with product_name '%s' is already exists", name)
		errorResponse := entity.DuplicateEntry.Wrap(err, "conflict error").WithProperty(entity.StatusCode, 409)
		p.dbLogger.Error("duplicate entry",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "InsertProduct"),
			zap.String("requestID", requestID),
			zap.Any("requestData", product),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, errorResponse
	}

	query := `INSERT INTO product(category_id, brand, product_name, description) VALUES(?, ?, ?, ?)`
	result, err := DB.ExecContext(ctx, query, categoryId, brand, name, description)
	if err != nil {
		errorResponse := entity.UnableToSave.Wrap(err, "failed to insert product").WithProperty(entity.StatusCode, 500)
		p.dbLogger.Error("failed to create product",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "InsertProduct"),
			zap.String("requestID", requestID),
			zap.String("query", query),
			zap.Any("requestData", product),
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
			zap.String("function", "InsertProduct"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, errorResponse
	}

	id := int(id64)

	return &id, nil

}

func (p productRepository) ListProducts(ctx context.Context, offset, limit int, requestID string) ([]utils.Product, error) {
	var products []utils.Product
	DB := p.db.GetDB()

	query := `SELECT product_id, category_id, brand, product_name, description, created_at, updated_at, deleted_at FROM product WHERE deleted_at IS NULL ORDER BY product_id LIMIT ? OFFSET ?`

	rows, err := DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		errorResponse := entity.UnableToFindResource.Wrap(err, "failed to get list of products").WithProperty(entity.StatusCode, 404)
		p.dbLogger.Error("products not found",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "ListProducts"),
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
		var product utils.Product

		if err := rows.Scan(&product.ID, &product.CategoryID, &product.Brand, &product.ProductName, &product.Description, &product.CreatedAt, &product.UpdatedAt, &product.DeletedAt); err != nil {
			errorResponse := entity.UnableToRead.Wrap(err, "failed to scan product data").WithProperty(entity.StatusCode, 500)
			p.dbLogger.Error("unable to scan the product data",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "databaseLayer"),
				zap.String("function", "ListProducts"),
				zap.String("requestID", requestID),
				zap.Error(errorResponse),
				zap.Stack("stacktrace"),
			)

			return nil, errorResponse
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "db rows error").WithProperty(entity.StatusCode, 500)
		p.dbLogger.Error("db rows error",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "ListProducts"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, errorResponse
	}

	return products, nil
}

func (p productRepository) GetProductById(ctx context.Context, id int, requestID string) (utils.Product, error) {
	var product utils.Product
	DB := p.db.GetDB()

	query := `SELECT product_id, category_id, brand, product_name, description, created_at, updated_at, deleted_at FROM product WHERE product_id = ? AND deleted_at IS NULL`
	if err := DB.QueryRowContext(ctx, query, id).Scan(&product.ID, &product.CategoryID, &product.Brand, &product.ProductName, &product.Description, &product.CreatedAt, &product.UpdatedAt, &product.DeletedAt); err != nil {
		errMessage := fmt.Sprintf("product with product_id '%d' not found", id)
		errorResponse := entity.UnableToFindResource.Wrap(err, errMessage).WithProperty(entity.StatusCode, 404)
		p.dbLogger.Error("product not found",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "GetProductById"),
			zap.String("requestID", requestID),
			zap.String("query", query),
			zap.Int("productID", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return utils.Product{}, errorResponse
	}

	return product, nil
}

func (p productRepository) EditProductById(ctx context.Context, id int, product utils.UpdateProduct, requestID string) (utils.Product, error) {
	DB := p.db.GetDB()
	var updateFields []string
	var values []interface{}

	var count int
	productUpdateQuery := "SELECT COUNT(*) FROM product WHERE product_id = ? AND deleted_at IS NULL"
	if err := DB.QueryRowContext(ctx, productUpdateQuery, id).Scan(&count); err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "unable to read COUNT in the query").WithProperty(entity.StatusCode, 500)
		p.dbLogger.Error("failed to read 'COUNT' in the query",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "EditProductById"),
			zap.String("requestID", requestID),
			zap.String("query", productUpdateQuery),
			zap.Int("sizeID", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return utils.Product{}, errorResponse
	}

	if count == 0 {
		err := fmt.Errorf("product with product_id '%d' not found", id)
		errorResponse := entity.UnableToFindResource.Wrap(err, "failed to get product by product_id").WithProperty(entity.StatusCode, 404)
		p.dbLogger.Error("product not found",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "EditProductById"),
			zap.String("requestID", requestID),
			zap.Int("productID", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return utils.Product{}, errorResponse
	}

	if product.CategoryID != 0 {
		updateFields = append(updateFields, "category_id = ?")
		values = append(values, product.CategoryID)
	}

	if product.Brand != "" {
		updateFields = append(updateFields, "brand = ?")
		values = append(values, product.Brand)
	}

	if product.Description != "" {
		updateFields = append(updateFields, "description = ?")
		values = append(values, product.Description)
	}

	if product.ProductName != "" {
		updateFields = append(updateFields, "product_name = ?")
		values = append(values, product.ProductName)
	}

	if len(updateFields) == 0 {
		err := errors.New("failed to update product:No fields provided for update.Please provide at least one field to update")
		errorResponse := entity.BadRequest.Wrap(err, "updated fields are required").WithProperty(entity.StatusCode, 400)
		p.dbLogger.Error("the updateProduct fields are empty",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "EditProductById"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return utils.Product{}, errorResponse
	}

	if len(values) > 0 {
		updateFields = append(updateFields, "updated_at = ?")
		values = append(values, time.Now())
	}

	query := fmt.Sprintf("UPDATE product SET %s WHERE product_id = ? AND deleted_at IS NULL", strings.Join(updateFields, ", "))
	values = append(values, id)
	if _, err := DB.ExecContext(ctx, query, values...); err != nil {
		errorResponse := entity.UnableToSave.Wrap(err, "failed to update product").WithProperty(entity.StatusCode, 500)
		p.dbLogger.Error("failed to edit product data",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "EditProductById"),
			zap.String("requestID", requestID),
			zap.String("query", query),
			zap.Any("requestData", product),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return utils.Product{}, errorResponse
	}

	updatedProduct, err := p.GetProductById(ctx, id, requestID)
	if err != nil {
		return utils.Product{}, err
	}

	return updatedProduct, nil

}

func (p productRepository) DeleteProductById(ctx context.Context, id int, requestID string) error {
	DB := p.db.GetDB()

	query := "SELECT COUNT(*) FROM product WHERE product_id = ?"
	var count int
	if err := DB.QueryRowContext(ctx, query, id).Scan(&count); err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "unable to read COUNT in the query").WithProperty(entity.StatusCode, 500)
		p.dbLogger.Error("failed to read 'COUNT' in the query",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "DeleteProductById"),
			zap.String("requestID", requestID),
			zap.String("query", query),
			zap.Int("productID", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return errorResponse
	}
	if count == 0 {
		err := fmt.Errorf("product with product_id '%d' not found", id)
		errorResponse := entity.UnableToFindResource.Wrap(err, "failed to get product by product_id").WithProperty(entity.StatusCode, 404)
		p.dbLogger.Error("product not found",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "DeleteProductById"),
			zap.String("requestID", requestID),
			zap.Int("productID", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return errorResponse
	}

	deleteProductQuery := `DELETE FROM product WHERE product_id = ?`
	if _, err := DB.ExecContext(ctx, query, id); err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "failed to delete product by product_id").WithProperty(entity.StatusCode, 500)
		p.dbLogger.Error("unable to delete product",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "DeleteProductById"),
			zap.String("requestID", requestID),
			zap.String("query", deleteProductQuery),
			zap.Int("sizeID", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)

		return errorResponse
	}

	return nil

}

// if err := DB.QueryRow("SELECT deleted_at FROM product WHERE product_id = ?", id).Scan(&deleted_at); err != nil {
// 	errType := errorcode.NotFoundError
// 	status := http.StatusNotFound
// 	err := fmt.Errorf("product with product_id '%d' not found", id)

// 	return "", status, errType, err
// }

// if deleted_at != nil {
// 	errType := "CONFLICT_ERROR"
// 	status := http.StatusConflict
// 	err := errors.New("can't delete already deleted product")

// 	return "", status, errType, err
// }

// query := `UPDATE product SET deleted_at = ? WHERE product_id = ?`
// if _, err := DB.Exec(query, time.Now(), id); err != nil {
// 	errType := errorcode.InternalError
// 	status := http.StatusInternalServerError

// 	return "", status, errType, err
// }

// errType := entity.Success
// status := http.StatusOK
// resp := fmt.Sprintf("product with product_id '%d' deleted successfully!", id)

// return resp, status, errType, nil
