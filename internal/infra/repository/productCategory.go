package repository

import (
	"Eccomerce-website/internal/core/common/utils"
	"Eccomerce-website/internal/core/dto"
	"Eccomerce-website/internal/core/port/repository"
	"fmt"
	"strings"
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

	query := `SELECT category_id, name, parent_id, created_at, updated_at, deleted_at FROM product_category WHERE category_id = ?`
	if err := DB.QueryRow(query, id).Scan(&category.ID, &category.Name, &category.ParentID, &category.CreatedAt, &category.UpdatedAt, &category.DeletedAt); err != nil {
		err = fmt.Errorf("product category with category_id %d not found", id)
		return utils.ProductCategory{}, err
	}

	return category, nil
}

func (p productCategoryRepository) EditProductCategoryById(id int, category utils.UpdateCategory) (utils.ProductCategory, error) {
	DB := p.db.GetDB()
	var updateFields []string
	var values []interface{}

	if category.Name != "" {
		updateFields = append(updateFields, "name = ?")
		values = append(values, category.Name)
	}

	if category.ParentID != 0 {
		updateFields = append(updateFields, "parent_id = ?")
		values = append(values, category.ParentID)
	}

	// oldProductCategory, err := p.GetProductCategoryById(id)
	// if err != nil {
	// 	return utils.ProductCategory{}, err
	// }

	// if oldProductCategory.Name == category.Name && oldProductCategory.ParentID == &category.ParentID {
	// 	err := errors.New("please provide updated information to proceed")
	// 	return utils.ProductCategory{}, err
	// }

	query := fmt.Sprintf("UPDATE product_category SET %s WHERE category_id = ?", strings.Join(updateFields, ", "))
	values = append(values, id)

	if _, err := DB.Exec(query, values...); err != nil {
		return utils.ProductCategory{}, err
	}

	productCategory, err := p.GetProductCategoryById(id)
	if err != nil {
		return utils.ProductCategory{}, err
	}

	return productCategory, nil
}
