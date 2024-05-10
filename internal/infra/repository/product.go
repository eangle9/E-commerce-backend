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

func (p productRepository) ListProducts() ([]utils.Product, error) {
	var products []utils.Product
	DB := p.db.GetDB()

	query := `SELECT product_id, category_id, product_name, description, created_at, updated_at, deleted_at FROM product WHERE deleted_at IS NULL`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var product utils.Product

		if err := rows.Scan(&product.ID, &product.CategoryID, &product.ProductName, &product.Description, &product.CreatedAt, &product.UpdatedAt, &product.DeletedAt); err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (p productRepository) GetProductById(id int) (utils.Product, error) {
	var product utils.Product
	DB := p.db.GetDB()

	query := `SELECT product_id, category_id, product_name, description, created_at, updated_at, deleted_at FROM product WHERE product_id = ? AND deleted_at IS NULL`
	if err := DB.QueryRow(query, id).Scan(&product.ID, &product.CategoryID, &product.ProductName, &product.Description, &product.CreatedAt, &product.UpdatedAt, &product.DeletedAt); err != nil {
		err := fmt.Errorf("product with product_id '%d' not found", id)
		return utils.Product{}, err
	}

	return product, nil
}

func (p productRepository) EditProductById(id int, product utils.UpdateProduct) (utils.Product, error) {
	DB := p.db.GetDB()
	var updateFields []string
	var values []interface{}

	if product.CategoryID != 0 {
		updateFields = append(updateFields, "category_id = ?")
		values = append(values, product.CategoryID)
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
		return utils.Product{}, err
	}

	if len(values) > 0 {
		updateFields = append(updateFields, "updated_at = ?")
		values = append(values, time.Now())
	}

	query := fmt.Sprintf("UPDATE product SET %s WHERE product_id = ? AND deleted_at IS NULL", strings.Join(updateFields, ", "))
	values = append(values, id)
	if _, err := DB.Exec(query, values...); err != nil {
		return utils.Product{}, err
	}

	updatedProduct, err := p.GetProductById(id)
	if err != nil {
		return utils.Product{}, err
	}

	return updatedProduct, nil

}
