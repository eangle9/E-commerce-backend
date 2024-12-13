package repository

import (
	"Eccomerce-website/internal/core/common/utils"
	"Eccomerce-website/internal/core/entity"
	"Eccomerce-website/internal/core/port/repository"
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type productsRepository struct {
	db       repository.Database
	dbLogger *zap.Logger
}

func NewProductsRepository(db repository.Database, dbLogger *zap.Logger) repository.GetProducts {
	return &productsRepository{
		db:       db,
		dbLogger: dbLogger,
	}
}

func (p productsRepository) ListAllProducts(ctx context.Context, offset, limit int, filters map[string]string, sort string, requestID string) ([]utils.ListProduct, error) {
	DB := p.db.GetDB()
	var products []utils.ListProduct
	var args []interface{}

	query := `
	SELECT 
	    p.product_id, 
		p.product_name, 
		c.name,
		p.brand,
		p.description
	FROM 
	   product p
	LEFT JOIN
	   product_category c ON p.category_id = c.category_id   
     WHERE 
	   1=1
	   `

	// Search Filtering
	if name, ok := filters["name"]; ok {
		if name != "" {
			query += " AND p.product_name LIKE ?"
			args = append(args, "%"+name+"%")
		}
	}

	// Category Filtering
	if category, ok := filters["category"]; ok {
		if category != "" {
			var categoryId int
			categoryQuery := "SELECT category_id FROM product_category WHERE name = ?"
			if err := DB.QueryRowContext(ctx, categoryQuery, category).Scan(&categoryId); err != nil {
				errorMessage := fmt.Sprintf("categoryID with category_name '%s' not found", category)
				errorResponse := entity.UnableToFindResource.Wrap(err, errorMessage).WithProperty(entity.StatusCode, 404)
				p.dbLogger.Error("categoryID not found",
					zap.String("timestamp", time.Now().Format(time.RFC3339)),
					zap.String("layer", "databaseLayer"),
					zap.String("function", "ListAllProducts"),
					zap.String("requestID", requestID),
					zap.String("query", categoryQuery),
					zap.String("categoryName", category),
					zap.Error(errorResponse),
					zap.Stack("stacktrace"),
				)

				return nil, errorResponse
			}

			query += " AND p.category_id = ?"
			args = append(args, categoryId)
		}
	}

	// Sorting
	if sort != "" {
		query += " ORDER BY " + sort
	} else {
		query += " ORDER BY p.created_at DESC"
	}

	// Pagination
	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	// query := "SELECT product_id, product_name FROM product ORDER BY product_id LIMIT ? OFFSET ?"
	productRows, err := DB.QueryContext(ctx, query, args...)
	if err != nil {
		errorResponse := entity.UnableToFindResource.Wrap(err, "failed to get list of products").WithProperty(entity.StatusCode, 404)
		p.dbLogger.Error("products not found",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "ListAllProducts"),
			zap.String("requestID", requestID),
			zap.String("query", query),
			zap.Int("limit", limit),
			zap.Int("offset", offset),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)

		return nil, errorResponse
	}

	defer productRows.Close()

	for productRows.Next() {
		var singleProduct utils.ListProduct

		if err := productRows.Scan(&singleProduct.ProductID, &singleProduct.Name, &singleProduct.Category, &singleProduct.Brand, &singleProduct.Description); err != nil {
			errorResponse := entity.UnableToRead.Wrap(err, "failed to scan ListProduct data").WithProperty(entity.StatusCode, 500)
			p.dbLogger.Error("unable to scan list product data",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "databaseLayer"),
				zap.String("function", "ListAllProducts"),
				zap.String("requestID", requestID),
				zap.Error(errorResponse),
				zap.Stack("stacktrace"),
			)

			return nil, errorResponse
		}

		query := `
		SELECT 
		    p.product_item_id, 
			c.color_name, 
			p.image_url, 
			p.price, 
			p.discount,
			p.qty_in_stock
		FROM 
		    product_item p 
		LEFT JOIN 
		    color c ON p.color_id = c.color_id
		WHERE 
		    p.product_id = ?
		 `
		productItemRows, err := DB.QueryContext(ctx, query, singleProduct.ProductID)
		if err != nil {
			errorResponse := entity.UnableToFindResource.Wrap(err, "failed to get list of product items by product_id").WithProperty(entity.StatusCode, 404)
			p.dbLogger.Error("product items not found",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "databaseLayer"),
				zap.String("function", "ListAllProducts"),
				zap.String("requestID", requestID),
				zap.String("query", query),
				zap.Int("productID", singleProduct.ProductID),
				zap.Error(errorResponse),
				zap.Stack("stacktrace"),
			)

			return nil, errorResponse
		}

		defer productItemRows.Close()

		var productItems []utils.ProductVariant
		for productItemRows.Next() {
			var productItem utils.ProductVariant
			if err := productItemRows.Scan(&productItem.ItemID, &productItem.Color, &productItem.ImageUrl, &productItem.Price, &productItem.Discount, &productItem.InStock); err != nil {
				errorResponse := entity.UnableToRead.Wrap(err, "failed to scan ProductVariant data").WithProperty(entity.StatusCode, 500)
				p.dbLogger.Error("unable to scan ProductVariant data",
					zap.String("timestamp", time.Now().Format(time.RFC3339)),
					zap.String("layer", "databaseLayer"),
					zap.String("function", "ListAllProducts"),
					zap.String("requestID", requestID),
					zap.Error(errorResponse),
					zap.Stack("stacktrace"),
				)

				return nil, errorResponse
			}

			sizeQuery := `SELECT size_id, size_name, price, discount, qty_in_stock FROM size WHERE product_item_id = ?`
			sizeRows, err := DB.QueryContext(ctx, sizeQuery, productItem.ItemID)
			if err != nil {
				errorResponse := entity.UnableToFindResource.Wrap(err, "failed to get list of product sizes  by product_item_id").WithProperty(entity.StatusCode, 404)
				p.dbLogger.Error("product sizes not found",
					zap.String("timestamp", time.Now().Format(time.RFC3339)),
					zap.String("layer", "databaseLayer"),
					zap.String("function", "ListAllProducts"),
					zap.String("requestID", requestID),
					zap.String("query", query),
					zap.Int("productItemID", productItem.ItemID),
					zap.Error(errorResponse),
					zap.Stack("stacktrace"),
				)

				return nil, errorResponse
			}
			defer sizeRows.Close()

			var sizes []utils.ProductSize
			for sizeRows.Next() {
				var size utils.ProductSize

				if err := sizeRows.Scan(&size.ID, &size.Size, &size.Price, &size.Discount, &size.QtyInStock); err != nil {
					errorResponse := entity.UnableToRead.Wrap(err, "failed to scan ProductVariant data").WithProperty(entity.StatusCode, 500)
					p.dbLogger.Error("unable to scan ProductVariant data",
						zap.String("timestamp", time.Now().Format(time.RFC3339)),
						zap.String("layer", "databaseLayer"),
						zap.String("function", "ListAllProducts"),
						zap.String("requestID", requestID),
						zap.Error(errorResponse),
						zap.Stack("stacktrace"),
					)

					return nil, errorResponse
				}

				sizes = append(sizes, size)
			}

			if err := sizeRows.Err(); err != nil {
				errorResponse := entity.UnableToRead.Wrap(err, "db sizeRows error").WithProperty(entity.StatusCode, 500)
				p.dbLogger.Error("db sizeRows error",
					zap.String("timestamp", time.Now().Format(time.RFC3339)),
					zap.String("layer", "databaseLayer"),
					zap.String("function", "ListAllProducts"),
					zap.String("requestID", requestID),
					zap.Error(errorResponse),
					zap.Stack("stacktrace"),
				)
				return nil, errorResponse
			}

			productItem.Sizes = sizes
			productItems = append(productItems, productItem)
		}

		if err := productItemRows.Err(); err != nil {
			errorResponse := entity.UnableToRead.Wrap(err, "db productItemRows error").WithProperty(entity.StatusCode, 500)
			p.dbLogger.Error("db productItemRows error",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "databaseLayer"),
				zap.String("function", "ListAllProducts"),
				zap.String("requestID", requestID),
				zap.Error(errorResponse),
				zap.Stack("stacktrace"),
			)
			return nil, errorResponse
		}

		singleProduct.ProductItems = productItems

		reviewQuery := `
    SELECT 
        r.review_id,
        r.user_id,
        r.product_id,
        r.rating,
        r.comment,
        r.created_at,
        u.user_id,
        u.username,
        u.first_name,
        u.last_name,
        u.email,
        u.phone_number
    FROM 
        review r
    LEFT JOIN
        users u ON r.user_id = u.user_id
    WHERE
        r.product_id = ?
`

		reviewRows, err := DB.QueryContext(ctx, reviewQuery, singleProduct.ProductID)
		if err != nil {
			errorResponse := entity.UnableToFindResource.Wrap(err, "failed to get list of product reviews  by product_id").WithProperty(entity.StatusCode, 404)
			p.dbLogger.Error("product reviews not found",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "databaseLayer"),
				zap.String("function", "ListAllProducts"),
				zap.String("requestID", requestID),
				zap.String("query", reviewQuery),
				zap.Int("productID", singleProduct.ProductID),
				zap.Error(errorResponse),
				zap.Stack("stacktrace"),
			)

			return nil, errorResponse
		}

		defer reviewRows.Close()

		var reviews []utils.ProductReview
		for reviewRows.Next() {
			var review utils.ProductReview
			if err := reviewRows.Scan(&review.ReviewID, &review.UserID, &review.ProductID, &review.Rating, &review.Comment, &review.CreatedAt, &review.User.ID, &review.User.Username, &review.User.FirstName, &review.User.LastName, &review.User.Email, &review.User.PhoneNumber); err != nil {
				errorResponse := entity.UnableToRead.Wrap(err, "failed to scan ProductReview data").WithProperty(entity.StatusCode, 500)
				p.dbLogger.Error("unable to scan ProductReview data",
					zap.String("timestamp", time.Now().Format(time.RFC3339)),
					zap.String("layer", "databaseLayer"),
					zap.String("function", "ListAllProducts"),
					zap.String("requestID", requestID),
					zap.Error(errorResponse),
					zap.Stack("stacktrace"),
				)

				return nil, errorResponse
			}
			reviews = append(reviews, review)
		}

		if err := reviewRows.Err(); err != nil {
			errorResponse := entity.UnableToRead.Wrap(err, "db reviewRows error").WithProperty(entity.StatusCode, 500)
			p.dbLogger.Error("db reviewRows error",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "databaseLayer"),
				zap.String("function", "ListAllProducts"),
				zap.String("requestID", requestID),
				zap.Error(errorResponse),
				zap.Stack("stacktrace"),
			)
			return nil, errorResponse
		}

		singleProduct.Reviews = reviews

		products = append(products, singleProduct)
	}

	return products, nil
}

func (p productsRepository) GetSingleProductById(ctx context.Context, id int, requestID string) (utils.SingleProduct, error) {
	DB := p.db.GetDB()
	var singleProduct utils.SingleProduct

	prouductQuery := `
	SELECT
	   p.product_id,
	   c.name,
	   p.brand,
	   p.product_name,
	   p.description
	FROM 
	  product p
	LEFT JOIN
	  product_category c ON p.category_id = c.category_id
	WHERE
	  p.product_id = ?    
	`

	if err := DB.QueryRowContext(ctx, prouductQuery, id).Scan(&singleProduct.ProductID, &singleProduct.Category, &singleProduct.Brand, &singleProduct.Name, &singleProduct.Description); err != nil {
		errMessage := fmt.Sprintf("products with product_id '%d' not found", id)
		errorResponse := entity.UnableToFindResource.Wrap(err, errMessage).WithProperty(entity.StatusCode, 404)
		p.dbLogger.Error("product not found",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "GetSingleProductById"),
			zap.String("requestID", requestID),
			zap.String("query", prouductQuery),
			zap.Int("productID", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)

		return utils.SingleProduct{}, errorResponse
	}

	itemQuery := `
	SELECT 
	   p.product_item_id,
	   c.color_name,
	   p.image_url,
	   p.price,
	   p.discount,
	   qty_in_stock
	FROM 
	  product_item p
	LEFT JOIN
	  color c ON p.color_id = c.color_id
	WHERE
	  p.product_id = ?    
	`

	itemRows, err := DB.QueryContext(ctx, itemQuery, id)
	if err != nil {
		errMessage := fmt.Sprintf("product items with product_id '%d' not found", id)
		errorResponse := entity.UnableToFindResource.Wrap(err, errMessage).WithProperty(entity.StatusCode, 404)
		p.dbLogger.Error("product items not found",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "GetSingleProductById"),
			zap.String("requestID", requestID),
			zap.String("query", itemQuery),
			zap.Int("productID", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)

		return utils.SingleProduct{}, errorResponse

	}

	defer itemRows.Close()

	var items []utils.ItemVariant
	for itemRows.Next() {
		var item utils.ItemVariant

		if err := itemRows.Scan(&item.ItemID, &item.Color, &item.ImageUrl, &item.Price, &item.Discount, &item.InStock); err != nil {
			errorResponse := entity.UnableToRead.Wrap(err, "failed to scan ItemVariant data").WithProperty(entity.StatusCode, 500)
			p.dbLogger.Error("unable to scan ItemVariant data",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "databaseLayer"),
				zap.String("function", "GetSingleProductById"),
				zap.String("requestID", requestID),
				zap.Error(errorResponse),
				zap.Stack("stacktrace"),
			)
			return utils.SingleProduct{}, errorResponse
		}

		sizeQuery := `SELECT size_id, size_name, price, discount, qty_in_stock FROM size WHERE product_item_id = ?`
		sizeRows, err := DB.QueryContext(ctx, sizeQuery, item.ItemID)
		if err != nil {
			errMessage := fmt.Sprintf("product item with product_item_id '%d' not found", item.ItemID)
			errorResponse := entity.UnableToFindResource.Wrap(err, errMessage).WithProperty(entity.StatusCode, 404)
			p.dbLogger.Error("product items not found",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "databaseLayer"),
				zap.String("function", "GetSingleProductById"),
				zap.String("requestID", requestID),
				zap.String("query", sizeQuery),
				zap.Int("productItemID", item.ItemID),
				zap.Error(errorResponse),
				zap.Stack("stacktrace"),
			)

			return utils.SingleProduct{}, errorResponse
		}
		defer sizeRows.Close()

		var sizes []utils.ProductSize
		for sizeRows.Next() {
			var size utils.ProductSize

			if err := sizeRows.Scan(&size.ID, &size.Size, &size.Price, &size.Discount, &size.QtyInStock); err != nil {
				errorResponse := entity.UnableToRead.Wrap(err, "failed to scan ProductSize data").WithProperty(entity.StatusCode, 500)
				p.dbLogger.Error("unable to scan ProductSize data",
					zap.String("timestamp", time.Now().Format(time.RFC3339)),
					zap.String("layer", "databaseLayer"),
					zap.String("function", "GetSingleProductById"),
					zap.String("requestID", requestID),
					zap.Error(errorResponse),
					zap.Stack("stacktrace"),
				)
				return utils.SingleProduct{}, errorResponse
			}

			sizes = append(sizes, size)
		}

		if err := sizeRows.Err(); err != nil {
			errorResponse := entity.UnableToRead.Wrap(err, "db sizeRows error").WithProperty(entity.StatusCode, 500)
			p.dbLogger.Error("db sizeRows error",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "databaseLayer"),
				zap.String("function", "GetSingleProductById"),
				zap.String("requestID", requestID),
				zap.Error(errorResponse),
				zap.Stack("stacktrace"),
			)
			return utils.SingleProduct{}, errorResponse
		}

		item.Sizes = sizes
		items = append(items, item)
	}

	if err := itemRows.Err(); err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "db itemRows error").WithProperty(entity.StatusCode, 500)
		p.dbLogger.Error("db itemRows error",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "GetSingleProductById"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return utils.SingleProduct{}, errorResponse
	}

	singleProduct.Items = items

	reviewQuery := `
	SELECT
	   r.review_id,
	   r.user_id,
	   r.product_id,
	   r.rating,
	   r.comment,
	   r.created_at,
	   u.user_id,
	   u.username,
	   u.first_name,
	   u.last_name,
	   u.email,
	   u.phone_number
	FROM 
	  review r
	LEFT JOIN 
	  users u ON r.user_id = u.user_id
	WHERE
	  r.product_id = ?    
	`
	reviewRows, err := DB.QueryContext(ctx, reviewQuery, id)
	if err != nil {
		errMessage := fmt.Sprintf("product review with product_id '%d' not found", id)
		errorResponse := entity.UnableToFindResource.Wrap(err, errMessage).WithProperty(entity.StatusCode, 404)
		p.dbLogger.Error("product review not found",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "GetSingleProductById"),
			zap.String("requestID", requestID),
			zap.String("query", reviewQuery),
			zap.Int("productID", id),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)

		return utils.SingleProduct{}, errorResponse
	}

	defer reviewRows.Close()

	var reviews []utils.ProductReview
	for reviewRows.Next() {
		var review utils.ProductReview
		if err := reviewRows.Scan(&review.ReviewID, &review.UserID, &review.ProductID, &review.Rating, &review.Comment, &review.CreatedAt, &review.User.ID, &review.User.Username, &review.User.FirstName, &review.User.LastName, &review.User.Email, &review.User.PhoneNumber); err != nil {
			errorResponse := entity.UnableToRead.Wrap(err, "failed to scan ProductReview data").WithProperty(entity.StatusCode, 500)
			p.dbLogger.Error("unable to scan ProductReview data",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "databaseLayer"),
				zap.String("function", "GetSingleProductById"),
				zap.String("requestID", requestID),
				zap.Error(errorResponse),
				zap.Stack("stacktrace"),
			)
			return utils.SingleProduct{}, errorResponse
		}
		reviews = append(reviews, review)
	}

	if err := reviewRows.Err(); err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "db reviewRows error").WithProperty(entity.StatusCode, 500)
		p.dbLogger.Error("db reviewRows error",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "GetSingleProductById"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return utils.SingleProduct{}, errorResponse
	}

	singleProduct.Reviews = reviews

	return singleProduct, nil
}
