package repository

import (
	"Eccomerce-website/internal/core/common/utils"
	"Eccomerce-website/internal/core/dto"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
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

type ProductRepository interface {
	InsertProduct(product dto.Product) (*int, error)
	ListProducts() ([]utils.Product, error)
	GetProductById(id int) (utils.Product, error)
	EditProductById(id int, product utils.UpdateProduct) (utils.Product, error)
	DeleteProductById(id int) (string, int, string, error)
}

type ProductItemRepository interface {
	InsertProductItem(item dto.ProductItem) (*int, error)
	ListProductItems() ([]utils.ProductItem, error)
	GetProductItemById(id int) (utils.ProductItem, error)
	EditProductItemById(id int, productItem utils.UpdateProductItem) (utils.ProductItem, error)
	DeleteProductItemById(id int) (string, int, string, error)
}

type ProductImageRepository interface {
	InsertProductImage(request request.ProductImageRequest) (*int, string, error)
}

type CartRepository interface {
	InsertCartItem(request request.CartRequest) ([]response.CartResponse, error)
}
