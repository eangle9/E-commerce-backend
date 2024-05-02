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
	ListProductCategory() ([]utils.CategoryList, error)
}
