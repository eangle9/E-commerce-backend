package repository

import (
	"Eccomerce-website/internal/core/common/utils"
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

func (p productCategoryRepository) InsertProductCategory(category dto.ProductCategory) (*int, error) {
	DB := p.db.GetDB()

	name := category.Name
	parentId := category.ParentID

	var count int
	if err := DB.QueryRow("SELECT COUNT(*) FROM product_category WHERE name = ?", name).Scan(&count); err != nil {
		return nil, err
	}

	if count > 0 {
		err := fmt.Errorf("product category with category_name '%s' is already exists", name)
		return nil, err
	}

	if parentId == 0 {
		query := `INSERT INTO product_category(name) VALUES(?)`
		result, err := DB.Exec(query, name)
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

	query := `INSERT INTO product_category(name, parent_id) VALUES(?, ?)`
	result, err := DB.Exec(query, name, parentId)
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

func (p productCategoryRepository) ListProductCategory() ([]utils.ProductCategory, error) {
	var productCategories []utils.ProductCategory
	DB := p.db.GetDB()

	query := `SELECT category_id, name, parent_id FROM product_category`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var productCategory utils.ProductCategory
		if err := rows.Scan(&productCategory.ID, &productCategory.Name, &productCategory.ParentID); err != nil {
			return nil, err
		}

		productCategories = append(productCategories, productCategory)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return productCategories, nil
}

func (p productCategoryRepository) GetProductCategoryById(id int) (utils.ProductCategory, error) {
	var category utils.ProductCategory
	DB := p.db.GetDB()

	query := `SELECT category_id, name, parent_id FROM product_category WHERE category_id = ?`
	if err := DB.QueryRow(query, id).Scan(&category.ID, &category.Name, &category.ParentID); err != nil {
		err = fmt.Errorf("product category with category_id %d not found", id)
		return utils.ProductCategory{}, err
	}

	return category, nil
}
