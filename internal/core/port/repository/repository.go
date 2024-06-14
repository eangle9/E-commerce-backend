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
	Authentication(request.LoginRequest) (utils.User, error)
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
	InsertProductItem(item request.ProductItemRequest) (*int, string, error)
	ListProductItems() ([]utils.ProductItem, error)
	GetProductItemById(id int) (utils.ProductItem, error)
	EditProductItemById(id int, productItem utils.UpdateProductItem) (utils.ProductItem, error)
	DeleteProductItemById(id int) (string, int, string, error)
}

type ProductImageRepository interface {
	InsertProductImage(request request.ProductImageRequest) (*int, string, error)
}

type CartRepository interface {
	InsertCartItem(request request.CartRequest, userId uint) ([]response.CartResponse, error)
}

type SizeRepository interface {
	InsertSize(size dto.Size) (*int, error)
	ListSizes() ([]utils.Size, error)
	GetSizeById(id int) (utils.Size, error)
	EditSizeById(id int, size utils.UpdateSize) (utils.Size, error)
	DeleteSizeById(id int) (string, int, string, error)
}

type GetProducts interface {
	ListAllProducts() ([]utils.ListProduct, error)
	GetSingleProductById(id int) (utils.SingleProduct, error)
}

type ReviewRepository interface {
	InsertReview(review dto.Review) (*int, error)
	ListReviews() ([]utils.Review, error)
}
