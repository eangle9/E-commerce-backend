package repository

import (
	"Eccomerce-website/internal/core/dto"
	"Eccomerce-website/internal/core/port/repository"
	"fmt"
)

type productCategoryRepository struct {
	db repository.Database
}

func NewProductCategoryRepository(db repository.Database) repository.ProductCategoryRepository {
	return &productCategoryRepository{
		db: db,
	}
}

func (p productCategoryRepository) InsertProductCategory(category dto.ProductCategory) (int, error) {
	DB := p.db.GetDB()
	name := category.Name
	parentId := category.ParentID

	var count int
	if err := DB.QueryRow("SELECT COUNT(*) FROM product_category WHERE name = ?", name).Scan(&count); err != nil {
		return 0, err
	}

	if count > 0 {
		err := fmt.Errorf("product category with name %s is already exists", name)
		return 0, err
	}

	if parentId == 0 {
		query := `INSERT INTO product_category(name) VALUES(?)`
		result, err := DB.Exec(query, name)
		if err != nil {
			return 0, err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return 0, err
		}

		return int(id), nil
	}

	query := `INSERT INTO product_category(parent_id, name) VALUES(?, ?)`
	result, err := DB.Exec(query, parentId, name)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil

}
