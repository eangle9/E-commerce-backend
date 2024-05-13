package repository

import (
	"Eccomerce-website/internal/core/dto"
	"Eccomerce-website/internal/core/port/repository"
	"fmt"
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
