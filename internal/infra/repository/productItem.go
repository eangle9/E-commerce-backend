package repository

import (
	"Eccomerce-website/internal/core/common/utils"
	cloudinaryupload "Eccomerce-website/internal/core/common/utils/cloudinary_upload"

	// "Eccomerce-website/internal/core/dto"
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/port/repository"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type productItemRepository struct {
	db repository.Database
}

func NewProductItemRepository(db repository.Database) repository.ProductItemRepository {
	return &productItemRepository{
		db: db,
	}
}

func (p productItemRepository) InsertProductItem(item request.ProductItemRequest) (*int, string, error) {
	DB := p.db.GetDB()
	productId := item.ProductID
	colorId := item.ColorID
	sizeId := item.SizeID
	price := item.Price
	discount := item.Discount
	qty := item.QtyInStock
	file := item.File

	var count int
	if colorId != nil && sizeId != nil {
		if err := DB.QueryRow("SELECT COUNT(*) FROM product_item WHERE product_id = ? AND color_id = ? AND size_id = ? AND deleted_at IS NULL", productId, colorId, sizeId).Scan(&count); err != nil {
			return nil, "", err
		}
	}

	// if colorId == nil {
	// 	if err := DB.QueryRow("SELECT COUNT(*) FROM product_item WHERE product_id = ? AND color_id IS NULL AND deleted_at IS NULL", productId).Scan(&count); err != nil {
	// 		return nil, "", err
	// 	}
	// } else {
	// 	if err := DB.QueryRow("SELECT COUNT(*) FROM product_item WHERE product_id = ? AND color_id = ? AND deleted_at IS NULL", productId, colorId).Scan(&count); err != nil {
	// 		return nil, "", err
	// 	}
	// }

	if count > 0 {
		err := fmt.Errorf("product_item with product_id '%d', color_id '%v' and size_id '%v' already exists", productId, colorId, sizeId)
		return nil, "", err
	}

	image_url, err := cloudinaryupload.UploadToCloudinary(file)
	if err != nil {
		return nil, "", err
	}

	query := `INSERT INTO product_item(product_id, color_id, size_id, image_url, price, discount, qty_in_stock) VALUES(?, ?, ?, ?, ?, ?, ?)`
	result, err := DB.Exec(query, productId, colorId, sizeId, image_url, price, discount, qty)
	if err != nil {
		return nil, "", err
	}

	id64, err := result.LastInsertId()
	if err != nil {
		return nil, "", err
	}

	id := int(id64)

	return &id, image_url, nil
}

func (p productItemRepository) ListProductItems() ([]utils.ProductItem, error) {
	var items []utils.ProductItem
	DB := p.db.GetDB()

	query := `SELECT product_item_id, product_id, color_id, size_id, image_url, price, discount, qty_in_stock, created_at, updated_at, deleted_at FROM product_item WHERE deleted_at IS NULL`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var item utils.ProductItem
		if err := rows.Scan(&item.ID, &item.ProductID, &item.ColorID, &item.SizeID, &item.ImageUrl, &item.Price, &item.Discount, &item.QtyInStock, &item.CreatedAt, &item.UpdatedAt, &item.DeletedAt); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil

}

func (p productItemRepository) GetProductItemById(id int) (utils.ProductItem, error) {
	DB := p.db.GetDB()
	var productItem utils.ProductItem

	query := `SELECT product_item_id, product_id, color_id, size_id, image_url, price, discount, qty_in_stock, created_at, updated_at, deleted_at FROM product_item WHERE product_item_id = ? AND deleted_at IS NULL`
	if err := DB.QueryRow(query, id).Scan(&productItem.ID, &productItem.ProductID, &productItem.ColorID, &productItem.SizeID, &productItem.ImageUrl, &productItem.Price, &productItem.Discount, &productItem.QtyInStock, &productItem.CreatedAt, &productItem.UpdatedAt, &productItem.DeletedAt); err != nil {
		return utils.ProductItem{}, err
	}

	return productItem, nil
}

func (p productItemRepository) EditProductItemById(id int, productItem utils.UpdateProductItem) (utils.ProductItem, error) {
	DB := p.db.GetDB()
	zeroDecimal := decimal.NewFromInt(0)
	var updateFields []string
	var values []interface{}

	if productItem.ProductID != nil {
		updateFields = append(updateFields, "product_id = ?")
		values = append(values, productItem.ProductID)
	}

	if productItem.ColorID != nil {
		updateFields = append(updateFields, "color_id = ?")
		values = append(values, productItem.ColorID)
	}

	if productItem.SizeID != nil {
		updateFields = append(updateFields, "size_id = ?")
		values = append(values, productItem.SizeID)
	}

	if !productItem.Price.Equal(zeroDecimal) {
		updateFields = append(updateFields, "price = ?")
		values = append(values, productItem.Price)
	}

	if productItem.Discount != zeroDecimal {
		updateFields = append(updateFields, "discount = ?")
		values = append(values, productItem.Discount)
	}

	if productItem.QtyInStock != nil {
		updateFields = append(updateFields, "qty_in_stock = ?")
		values = append(values, productItem.QtyInStock)
	}

	if productItem.File != nil {
		ImageUrl, err := cloudinaryupload.UploadToCloudinary(productItem.File)
		if err != nil {
			return utils.ProductItem{}, err
		}
		updateFields = append(updateFields, "image_url")
		values = append(values, ImageUrl)
	}

	if len(updateFields) == 0 {
		err := errors.New("failed to update color:No fields provided for update.Please provide at least one field to update")
		return utils.ProductItem{}, err
	}

	if len(values) > 0 {
		updateFields = append(updateFields, "updated_at = ?")
		values = append(values, time.Now())
	}

	query := fmt.Sprintf("UPDATE product_item SET %s WHERE product_item_id = ? AND deleted_at IS NULL", strings.Join(updateFields, ", "))
	values = append(values, id)

	if _, err := DB.Exec(query, values...); err != nil {
		return utils.ProductItem{}, err
	}

	updatedItem, err := p.GetProductItemById(id)
	if err != nil {
		return utils.ProductItem{}, err
	}

	return updatedItem, nil
}

func (p productItemRepository) DeleteProductItemById(id int) (string, int, string, error) {
	DB := p.db.GetDB()
	// var deleted_at *time.Time

	var count int
	if err := DB.QueryRow("SELECT COUNT(*) FROM product_item WHERE product_item_id = ?", id).Scan(&count); err != nil {
		status := http.StatusInternalServerError
		errType := errorcode.InternalError
		return "", status, errType, err
	}
	if count == 0 {
		status := http.StatusNotFound
		errType := errorcode.NotFoundError
		err := fmt.Errorf("product item with id '%d' not found", id)
		return "", status, errType, err
	}

	query := `DELETE FROM product_item WHERE product_item_id = ?`
	if _, err := DB.Exec(query, id); err != nil {
		status := http.StatusInternalServerError
		errType := errorcode.InternalError
		return "", status, errType, err
	}

	status := http.StatusOK
	errType := errorcode.Success
	resp := fmt.Sprintf("product item with id '%d' deleted successfully!", id)

	return resp, status, errType, nil

	// var productID int
	// if err := DB.QueryRow("SELECT product_id FROM product_item WHERE product_item_id = ?", id).Scan(&productID); err != nil {
	// 	status := http.StatusInternalServerError
	// 	errType := errorcode.InternalError
	// 	return "", status, errType, err
	// }

	// var cartID int
	// if err := DB.QueryRow("SELECT cart_id FROM cart_item WHERE product_item_id = ?", id).Scan(&cartID); err != nil {
	// 	status := http.StatusInternalServerError
	// 	errType := errorcode.InternalError
	// 	return "", status, errType, err
	// }

	// var productCount int
	// if err := DB.QueryRow("SELECT COUNT(*) FROM product_item WHERE product_id = ?", productID).Scan(&productCount); err != nil {
	// 	status := http.StatusInternalServerError
	// 	errType := errorcode.InternalError
	// 	return "", status, errType, err
	// }

	// if productCount == 0 {
	// 	var pdtCount int
	// 	if err := DB.QueryRow("SELECT COUNT(*) FROM product WHERE product_id = ?", productID).Scan(&pdtCount); err != nil {
	// 		status := http.StatusInternalServerError
	// 		errType := errorcode.InternalError
	// 		return "", status, errType, err
	// 	}
	// 	if pdtCount > 0 {
	// 		pdtQuery := `DELETE FROM product WHERE product_id = ?`
	// 		if _, err := DB.Exec(pdtQuery, productID); err != nil {
	// 			status := http.StatusInternalServerError
	// 			errType := errorcode.InternalError
	// 			return "", status, errType, err
	// 		}
	// 	}
	// }

	// var cartCount int
	// if err := DB.QueryRow("SELECT COUNT(*) FROM cart_item WHERE cart_id = ?", cartID).Scan(&cartCount); err != nil {
	// 	status := http.StatusInternalServerError
	// 	errType := errorcode.InternalError
	// 	return "", status, errType, err
	// }

	// if cartCount == 0 {
	// 	var ctCount int
	// 	if err := DB.QueryRow("SELECT COUNT(*) FROM shopping_cart WHERE cart_id = ?", cartID).Scan(&ctCount); err != nil {
	// 		status := http.StatusInternalServerError
	// 		errType := errorcode.InternalError
	// 		return "", status, errType, err
	// 	}
	// 	if ctCount > 0 {
	// 		cartQuery := `DELETE FROM shopping_cart WHERE cart_id = ?`
	// 		if _, err := DB.Exec(cartQuery, cartID); err != nil {
	// 			status := http.StatusInternalServerError
	// 			errType := errorcode.InternalError
	// 			return "", status, errType, err
	// 		}
	// 	}
	// }

	// if cartCount > 0 {
	// 	if err := UpdateShoppingCart(DB, cartID); err != nil {
	// 		status := http.StatusInternalServerError
	// 		errType := errorcode.InternalError
	// 		return "", status, errType, err
	// 	}
	// }

}
