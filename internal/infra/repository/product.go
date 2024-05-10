package repository

import (
	"Eccomerce-website/internal/core/dto"
	"Eccomerce-website/internal/core/port/repository"
	"fmt"
)

type productRepository struct {
	db repository.Database
}

func NewProductRepository(db repository.Database) repository.ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (p productRepository) InsertProduct(product dto.Product) (*int, error) {
	DB := p.db.GetDB()
	categoryId := product.CategoryID
	name := product.ProductName
	description := product.Description

	var count int
	if err := DB.QueryRow("SELECT COUNT(*) FROM product WHERE product_name = ? AND deleted_at IS NULL", name).Scan(&count); err != nil {
		return nil, err
	}

	if count > 0 {
		err := fmt.Errorf("product with product_id '%s' is already exists", name)
		return nil, err
	}

	query := `INSERT INTO product(category_id, product_name, description) VALUES(?, ?, ?)`
	result, err := DB.Exec(query, categoryId, name, description)
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
