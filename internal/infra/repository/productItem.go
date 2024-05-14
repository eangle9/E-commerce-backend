package repository

import (
	"Eccomerce-website/internal/core/common/utils"
	"Eccomerce-website/internal/core/dto"
	"Eccomerce-website/internal/core/port/repository"
	"errors"
	"fmt"
	"strings"
	"time"
)

type productItemRepository struct {
	db repository.Database
}

func NewProductItemRepository(db repository.Database) repository.ProductItemRepository {
	return &productItemRepository{
		db: db,
	}
}

func (p productItemRepository) InsertProductItem(item dto.ProductItem) (*int, error) {
	DB := p.db.GetDB()
	productId := item.ProductID
	colorId := item.ColorID
	price := item.Price
	qty := item.QtyInStock

	var count int
	if colorId == nil {
		if err := DB.QueryRow("SELECT COUNT(*) FROM product_item WHERE product_id = ? AND color_id IS NULL AND deleted_at IS NULL", productId).Scan(&count); err != nil {
			return nil, err
		}
	} else {
		if err := DB.QueryRow("SELECT COUNT(*) FROM product_item WHERE product_id = ? AND color_id = ? AND deleted_at IS NULL", productId, colorId).Scan(&count); err != nil {
			return nil, err
		}
	}

	if count > 0 {
		err := fmt.Errorf("product_item with product_id '%d' and color_id '%v' already exists", productId, colorId)
		return nil, err
	}

	query := `INSERT INTO product_item(product_id, color_id, price, qty_in_stock) VALUES(?, ?, ?, ?)`
	result, err := DB.Exec(query, productId, colorId, price, qty)
	if err != nil {
		return nil, err
	}

	id64, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	id := int(id64)

	return &id, nil
}

func (p productItemRepository) ListProductItems() ([]utils.ProductItem, error) {
	var items []utils.ProductItem
	DB := p.db.GetDB()

	query := `SELECT product_item_id, product_id, color_id, price, qty_in_stock, created_at, updated_at, deleted_at FROM product_item WHERE deleted_at IS NULL`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var item utils.ProductItem
		if err := rows.Scan(&item.ID, &item.ProductID, &item.ColorID, &item.Price, &item.QtyInStock, &item.CreatedAt, &item.UpdatedAt, &item.DeletedAt); err != nil {
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

	query := `SELECT product_item_id, product_id, color_id, price, qty_in_stock, created_at, updated_at, deleted_at FROM product_item WHERE product_item_id = ? AND deleted_at IS NULL`
	if err := DB.QueryRow(query, id).Scan(&productItem.ID, &productItem.ProductID, &productItem.ColorID, &productItem.Price, &productItem.QtyInStock, &productItem.CreatedAt, &productItem.UpdatedAt, &productItem.DeletedAt); err != nil {
		return utils.ProductItem{}, err
	}

	return productItem, nil
}

func (p productItemRepository) EditProductItemById(id int, productItem utils.UpdateProductItem) (utils.ProductItem, error) {
	DB := p.db.GetDB()
	var updateFields []string
	var values []interface{}

	if productItem.ProductID != 0 {
		updateFields = append(updateFields, "product_id = ?")
		values = append(values, productItem.ProductID)
	}

	if productItem.ColorID != 0 {
		updateFields = append(updateFields, "color_id = ?")
		values = append(values, productItem.ColorID)
	}

	if productItem.Price != 0 {
		updateFields = append(updateFields, "price = ?")
		values = append(values, productItem.Price)
	}

	if productItem.QtyInStock != nil {
		updateFields = append(updateFields, "qty_in_stock = ?")
		values = append(values, productItem.QtyInStock)
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

// func (p productItemRepository) DeleteProductItemById(id int)  error {
// 	DB := p.db.GetDB()
// 	var deleted_at *time.Time

// 	if err := DB.QueryRow("SELECT deleted_at FROM product_item WHERE product_item_id = ?").Scan(&deleted_at); err != nil {
// 		err = fmt.Errorf("product item with product_item_id '%d' not found", id)
// 		return utils.ProductItem{}, err
// 	}

// 	if deleted_at != nil {
// 		err := errors.New("you can't delete already deleted product item")
// 		return utils.ProductItem{}, err
// 	}

// 	query := `SELECT `
// }
