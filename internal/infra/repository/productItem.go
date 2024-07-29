package repository

import (
	"Eccomerce-website/internal/core/common/utils"
	cloudinaryupload "Eccomerce-website/internal/core/common/utils/cloudinary_upload"
	"Eccomerce-website/internal/core/entity"
	"context"

	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/port/repository"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type productItemRepository struct {
	db       repository.Database
	dbLogger *zap.Logger
}

func NewProductItemRepository(db repository.Database, dbLogger *zap.Logger) repository.ProductItemRepository {
	return &productItemRepository{
		db:       db,
		dbLogger: dbLogger,
	}
}

func (p productItemRepository) InsertProductItem(ctx context.Context, item request.ProductItemRequest, requestID string) (*int, string, error) {
	DB := p.db.GetDB()
	productId := item.ProductID
	colorId := item.ColorID
	price := item.Price
	discount := item.Discount
	qty := item.QtyInStock
	file := item.File

	checkProductQuery := "SELECT COUNT(*) FROM product WHERE product_id = ?"
	var productCount int
	if err := DB.QueryRowContext(ctx, checkProductQuery, productId).Scan(&productCount); err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "unable to read COUNT in the query").WithProperty(entity.StatusCode, 500)
		p.dbLogger.Error("failed to read 'COUNT' in the query",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "InsertProductItem"),
			zap.String("requestID", requestID),
			zap.String("query", checkProductQuery),
			zap.Int("productID", productId),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)

		return nil, "", errorResponse
	}
	if productCount == 0 {
		err := fmt.Errorf("product_id '%d' does not exist in the product table", productId)
		errorResponse := entity.UnableToFindResource.Wrap(err, "failed to get product with product_id").WithProperty(entity.StatusCode, 404)
		p.dbLogger.Error("product not found",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "InsertProductItem"),
			zap.String("requestID", requestID),
			zap.Int("productID", productId),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)

		return nil, "", errorResponse
	}

	checkDuplicateItemQuery := "SELECT COUNT(*) FROM product_item WHERE product_id = ? AND color_id = ? AND deleted_at IS NULL"
	var count int
	if colorId != nil {
		if err := DB.QueryRowContext(ctx, checkDuplicateItemQuery, productId, colorId).Scan(&count); err != nil {
			errorResponse := entity.UnableToRead.Wrap(err, "unable to read COUNT in the query").WithProperty(entity.StatusCode, 500)
			p.dbLogger.Error("failed to read 'COUNT' in the query",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "databaseLayer"),
				zap.String("function", "InsertProductItem"),
				zap.String("requestID", requestID),
				zap.String("query", checkDuplicateItemQuery),
				zap.Int("productID", productId),
				zap.Int("colorID", *colorId),
				zap.Error(errorResponse),
				zap.Stack("stacktrace"),
			)

			return nil, "", errorResponse
		}
	}

	if count > 0 {
		err := fmt.Errorf("product_item with product_id '%d' and color_id '%d' already exists", productId, *colorId)
		errorResponse := entity.DuplicateEntry.Wrap(err, "conflict error").WithProperty(entity.StatusCode, 409)
		p.dbLogger.Error("duplicate entry",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "InsertProductItem"),
			zap.String("requestID", requestID),
			zap.Any("requestData", item),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, "", errorResponse
	}

	image_url, err := cloudinaryupload.UploadToCloudinary(file, p.dbLogger, requestID)
	if err != nil {
		return nil, "", err
	}

	query := `INSERT INTO product_item(product_id, color_id, image_url, price, discount, qty_in_stock) VALUES(?, ?, ?, ?, ?, ?)`
	result, err := DB.ExecContext(ctx, query, productId, colorId, image_url, price, discount, qty)
	if err != nil {
		errorResponse := entity.UnableToSave.Wrap(err, "failed to insert product item").WithProperty(entity.StatusCode, 500)
		p.dbLogger.Error("failed to create product item",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "InsertProductItem"),
			zap.String("requestID", requestID),
			zap.String("query", query),
			zap.Any("requestData", item),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, "", errorResponse
	}

	id64, err := result.LastInsertId()
	if err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "failed to get the inserted id").WithProperty(entity.StatusCode, 500)
		p.dbLogger.Error("unable to get lastInserted id",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "InsertProductItem"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, "", errorResponse
	}

	id := int(id64)

	return &id, image_url, nil
}

func (p productItemRepository) ListProductItems(ctx context.Context, offset, limit int, requestID string) ([]utils.ProductItem, error) {
	var items []utils.ProductItem
	DB := p.db.GetDB()

	query := `SELECT product_item_id, product_id, color_id, image_url, price, discount, qty_in_stock, created_at, updated_at, deleted_at FROM product_item WHERE deleted_at IS NULL ORDER BY product_item_id LIMIT ? OFFSET ?`
	rows, err := DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		errorResponse := entity.UnableToFindResource.Wrap(err, "failed to get list of product items").WithProperty(entity.StatusCode, 404)
		p.dbLogger.Error("product items not found",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "ListProductItems"),
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
		var item utils.ProductItem
		if err := rows.Scan(&item.ID, &item.ProductID, &item.ColorID, &item.ImageUrl, &item.Price, &item.Discount, &item.QtyInStock, &item.CreatedAt, &item.UpdatedAt, &item.DeletedAt); err != nil {
			errorResponse := entity.UnableToRead.Wrap(err, "failed to scan productItem data").WithProperty(entity.StatusCode, 500)
			p.dbLogger.Error("unable to scan productItem data",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "databaseLayer"),
				zap.String("function", "ListProductItems"),
				zap.String("requestID", requestID),
				zap.Error(errorResponse),
				zap.Stack("stacktrace"),
			)
			return nil, errorResponse
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "db rows error").WithProperty(entity.StatusCode, 500)
		p.dbLogger.Error("db rows error",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "ListProductItems"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, errorResponse
	}

	return items, nil

}

func (p productItemRepository) GetProductItemById(ctx context.Context, id int, requestID string) (utils.ProductItem, error) {
	DB := p.db.GetDB()
	var productItem utils.ProductItem

	query := `SELECT product_item_id, product_id, color_id, image_url, price, discount, qty_in_stock, created_at, updated_at, deleted_at FROM product_item WHERE product_item_id = ? AND deleted_at IS NULL`
	if err := DB.QueryRowContext(ctx, query, id).Scan(&productItem.ID, &productItem.ProductID, &productItem.ColorID, &productItem.ImageUrl, &productItem.Price, &productItem.Discount, &productItem.QtyInStock, &productItem.CreatedAt, &productItem.UpdatedAt, &productItem.DeletedAt); err != nil {
		errorMessage := fmt.Sprintf("product item with product_item_id '%d' not found", id)
		errorResponse := entity.UnableToFindResource.Wrap(err, errorMessage).WithProperty(entity.StatusCode, 404)
		p.dbLogger.Error("productItem not found",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "GetProductItemById"),
			zap.String("requestID", requestID),
			zap.String("query", query),
			zap.Int("productItemID", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return utils.ProductItem{}, errorResponse
	}

	return productItem, nil
}

func (p productItemRepository) EditProductItemById(ctx context.Context, id int, productItem utils.UpdateProductItem, requestID string) (utils.ProductItem, error) {
	DB := p.db.GetDB()
	zeroDecimal := decimal.NewFromInt(0)
	var updateFields []string
	var values []interface{}

	query := "SELECT COUNT(*) FROM product_item WHERE product_item_id = ? AND deleted_at IS NULL"
	var count int
	if err := DB.QueryRowContext(ctx, query, id).Scan(&count); err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "unable to read COUNT in the query").WithProperty(entity.StatusCode, 500)
		p.dbLogger.Error("failed to read 'COUNT' in the query",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "EditProductItemById"),
			zap.String("requestID", requestID),
			zap.String("query", query),
			zap.Int("productItemID", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return utils.ProductItem{}, errorResponse
	}

	if count == 0 {
		err := fmt.Errorf("product item with product_item_id '%d' not found", id)
		errorResponse := entity.UnableToFindResource.Wrap(err, "failed to get product item by id").WithProperty(entity.StatusCode, 404)
		p.dbLogger.Error("product item not found",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "EditProductItemById"),
			zap.String("requestID", requestID),
			zap.Int("productItemID", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return utils.ProductItem{}, errorResponse
	}

	if productItem.ProductID != nil {
		updateFields = append(updateFields, "product_id = ?")
		values = append(values, productItem.ProductID)
	}

	if productItem.ColorID != nil {
		updateFields = append(updateFields, "color_id = ?")
		values = append(values, productItem.ColorID)
	}

	if !productItem.Price.Equal(zeroDecimal) {
		updateFields = append(updateFields, "price = ?")
		values = append(values, productItem.Price)
	}

	if !productItem.Discount.Equal(zeroDecimal) {
		updateFields = append(updateFields, "discount = ?")
		values = append(values, productItem.Discount)
	}

	if productItem.QtyInStock != nil {
		updateFields = append(updateFields, "qty_in_stock = ?")
		values = append(values, productItem.QtyInStock)
	}

	if productItem.File != nil {
		ImageUrl, err := cloudinaryupload.UploadToCloudinary(productItem.File, p.dbLogger, requestID)
		if err != nil {
			return utils.ProductItem{}, err
		}
		updateFields = append(updateFields, "image_url = ?")
		values = append(values, ImageUrl)
	}

	if len(updateFields) == 0 {
		err := errors.New("failed to update color:No fields provided for update.Please provide at least one field to update")
		errorResponse := entity.BadRequest.Wrap(err, "updated fields are required").WithProperty(entity.StatusCode, 400)
		p.dbLogger.Error("the updateProductItem fields are empty",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "EditProductItemById"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return utils.ProductItem{}, errorResponse
	}

	if len(values) > 0 {
		updateFields = append(updateFields, "updated_at = ?")
		values = append(values, time.Now())
	}

	updateQuery := fmt.Sprintf("UPDATE product_item SET %s WHERE product_item_id = ? AND deleted_at IS NULL", strings.Join(updateFields, ", "))
	values = append(values, id)

	if _, err := DB.ExecContext(ctx, updateQuery, values...); err != nil {
		errorResponse := entity.UnableToSave.Wrap(err, "failed to update product item data").WithProperty(entity.StatusCode, 500)
		p.dbLogger.Error("failed to edit productItem data",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "EditProductItemById"),
			zap.String("requestID", requestID),
			zap.String("query", query),
			zap.Any("requestData", productItem),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return utils.ProductItem{}, errorResponse
	}

	updatedItem, err := p.GetProductItemById(ctx, id, requestID)
	if err != nil {
		return utils.ProductItem{}, err
	}

	return updatedItem, nil
}

func (p productItemRepository) DeleteProductItemById(ctx context.Context, id int, requestID string) error {
	DB := p.db.GetDB()

	query := "SELECT COUNT(*) FROM product_item WHERE product_item_id = ?"
	var count int
	if err := DB.QueryRowContext(ctx, query, id).Scan(&count); err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "unable to read COUNT in the query").WithProperty(entity.StatusCode, 500)
		p.dbLogger.Error("failed to read 'COUNT' in the query",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "DeleteProductItemById"),
			zap.String("requestID", requestID),
			zap.String("query", query),
			zap.Int("productID", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return errorResponse
	}

	if count == 0 {
		err := fmt.Errorf("product item with id '%d' not found", id)
		errorResponse := entity.UnableToFindResource.Wrap(err, "failed to get product item by id").WithProperty(entity.StatusCode, 404)
		p.dbLogger.Error("product item not found",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "DeleteProductItemById"),
			zap.String("requestID", requestID),
			zap.Int("productID", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return errorResponse
	}

	deleteQuery := `DELETE FROM product_item WHERE product_item_id = ?`
	if _, err := DB.ExecContext(ctx, deleteQuery, id); err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "failed to delete product item by id").WithProperty(entity.StatusCode, 500)
		p.dbLogger.Error("unable to delete product item",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "DeleteProductItemById"),
			zap.String("requestID", requestID),
			zap.String("query", deleteQuery),
			zap.Int("productID", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return errorResponse
	}

	return nil

	// var productID int
	// if err := DB.QueryRow("SELECT product_id FROM product_item WHERE product_item_id = ?", id).Scan(&productID); err != nil {
	// 	status := http.StatusInternalServerError
	// 	errType := entity.InternalError
	// 	return "", status, errType, err
	// }

	// var cartID int
	// if err := DB.QueryRow("SELECT cart_id FROM cart_item WHERE product_item_id = ?", id).Scan(&cartID); err != nil {
	// 	status := http.StatusInternalServerError
	// 	errType := entity.InternalError
	// 	return "", status, errType, err
	// }

	// var productCount int
	// if err := DB.QueryRow("SELECT COUNT(*) FROM product_item WHERE product_id = ?", productID).Scan(&productCount); err != nil {
	// 	status := http.StatusInternalServerError
	// 	errType := entity.InternalError
	// 	return "", status, errType, err
	// }

	// if productCount == 0 {
	// 	var pdtCount int
	// 	if err := DB.QueryRow("SELECT COUNT(*) FROM product WHERE product_id = ?", productID).Scan(&pdtCount); err != nil {
	// 		status := http.StatusInternalServerError
	// 		errType := entity.InternalError
	// 		return "", status, errType, err
	// 	}
	// 	if pdtCount > 0 {
	// 		pdtQuery := `DELETE FROM product WHERE product_id = ?`
	// 		if _, err := DB.Exec(pdtQuery, productID); err != nil {
	// 			status := http.StatusInternalServerError
	// 			errType := entity.InternalError
	// 			return "", status, errType, err
	// 		}
	// 	}
	// }

	// var cartCount int
	// if err := DB.QueryRow("SELECT COUNT(*) FROM cart_item WHERE cart_id = ?", cartID).Scan(&cartCount); err != nil {
	// 	status := http.StatusInternalServerError
	// 	errType := entity.InternalError
	// 	return "", status, errType, err
	// }

	// if cartCount == 0 {
	// 	var ctCount int
	// 	if err := DB.QueryRow("SELECT COUNT(*) FROM shopping_cart WHERE cart_id = ?", cartID).Scan(&ctCount); err != nil {
	// 		status := http.StatusInternalServerError
	// 		errType := entity.InternalError
	// 		return "", status, errType, err
	// 	}
	// 	if ctCount > 0 {
	// 		cartQuery := `DELETE FROM shopping_cart WHERE cart_id = ?`
	// 		if _, err := DB.Exec(cartQuery, cartID); err != nil {
	// 			status := http.StatusInternalServerError
	// 			errType := entity.InternalError
	// 			return "", status, errType, err
	// 		}
	// 	}
	// }

	// if cartCount > 0 {
	// 	if err := UpdateShoppingCart(DB, cartID); err != nil {
	// 		status := http.StatusInternalServerError
	// 		errType := entity.InternalError
	// 		return "", status, errType, err
	// 	}
	// }

}
