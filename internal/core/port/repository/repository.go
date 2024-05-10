package repository

import (
	"Eccomerce-website/internal/core/common/utils"
	"Eccomerce-website/internal/core/dto"
	// dbmodels "Eccomerce-website/internal/infra/db_models"
)

type UserRepository interface {
	InsertUser(user dto.User) (int, error)
	Authentication(email string, password string) (utils.User, error)
	ListUsers() ([]utils.User, error)
	GetUserById(id int) (utils.User, error)
	EditUserById(id int, user utils.UpdateUser) (utils.User, error)
	DeleteUserById(id int) (string, int, string, error)
}

type ProductCategoryRepository interface {
	InsertProductCategory(category dto.ProductCategory) (*int, error)
	ListProductCategory() ([]utils.ProductCategory, error)
	GetProductCategoryById(id int) (utils.ProductCategory, error)
	EditProductCategoryById(id int, category utils.UpdateCategory) (utils.ProductCategory, error)
	DeleteProductCategoryById(id int) (string, int, string, error)
}

type ColorRepository interface {
	InsertColor(color dto.Color) (*int, error)
	ListColors() ([]utils.Color, error)
	GetColorById(id int) (utils.Color, error)
	EditColorById(id int, color utils.UpdateColor) (utils.Color, error)
	DeleteColorById(id int) (string, int, string, error)
}
