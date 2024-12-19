package repository

import (
	"Eccomerce-website/internals/core/common/utils"
	"Eccomerce-website/internals/core/dto"
	"Eccomerce-website/internals/core/model/request"
	"Eccomerce-website/internals/core/model/response"
	"context"
	// dbmodels "Eccomerce-website/internal/infra/db_models"
)

type UserRepository interface {
	InsertUser(ctx context.Context, user dto.User, requestID string) (int, error)
	Authentication(ctx context.Context, request request.LoginRequest, requestID string) (utils.User, error)
	ListUsers(ctx context.Context, offset, perPage int, requestID string) ([]utils.User, error)
	GetUserById(ctx context.Context, id int, requestID string) (utils.User, error)
	EditUserById(ctx context.Context, id int, user request.UpdateUser, requestID string) (utils.User, error)
	DeleteUserById(ctx context.Context, id int, requestID string) error
}

type ProductCategoryRepository interface {
	InsertProductCategory(ctx context.Context, category dto.ProductCategory, requestID string) (*int, error)
	ListProductCategory(ctx context.Context, offset, limit int, requestID string) ([]utils.ProductCategory, error)
	GetProductCategoryById(ctx context.Context, id int, requestID string) (utils.ProductCategory, error)
	EditProductCategoryById(ctx context.Context, id int, category utils.UpdateCategory, requestID string) (utils.ProductCategory, error)
	DeleteProductCategoryById(ctx context.Context, id int, requestID string) error
}

type ColorRepository interface {
	InsertColor(ctx context.Context, color dto.Color, requestID string) (*int, error)
	ListColors(ctx context.Context, offset, limit int, requestID string) ([]utils.Color, error)
	GetColorById(ctx context.Context, id int, requestID string) (utils.Color, error)
	EditColorById(ctx context.Context, id int, color utils.UpdateColor, requestID string) (utils.Color, error)
	DeleteColorById(ctx context.Context, id int, requestID string) error
}

type ProductRepository interface {
	InsertProduct(ctx context.Context, product dto.Product, requestID string) (*int, error)
	ListProducts(ctx context.Context, offset, limit int, requestID string) ([]utils.Product, error)
	GetProductById(ctx context.Context, id int, requestID string) (utils.Product, error)
	EditProductById(ctx context.Context, id int, product utils.UpdateProduct, requestID string) (utils.Product, error)
	DeleteProductById(ctx context.Context, id int, requestID string) error
}

type ProductItemRepository interface {
	InsertProductItem(ctx context.Context, item request.ProductItemRequest, requestID string) (*int, string, error)
	ListProductItems(ctx context.Context, offset, limit int, requestID string) ([]utils.ProductItem, error)
	GetProductItemById(ctx context.Context, id int, requestID string) (utils.ProductItem, error)
	EditProductItemById(ctx context.Context, id int, productItem utils.UpdateProductItem, requestID string) (utils.ProductItem, error)
	DeleteProductItemById(ctx context.Context, id int, requestID string) error
}

type ProductImageRepository interface {
	InsertProductImage(ctx context.Context, request request.ProductImageRequest, requestID string) (*int, string, error)
}

type CartRepository interface {
	InsertCartItem(request request.CartRequest, userId uint) ([]response.CartResponse, error)
}

type SizeRepository interface {
	InsertSize(ctx context.Context, size dto.Size, requestID string) (*int, error)
	ListSizes(ctx context.Context, offset, limit int, requestID string) ([]utils.Size, error)
	GetSizeById(ctx context.Context, id int, requestID string) (utils.Size, error)
	EditSizeById(ctx context.Context, id int, size request.UpdateSize, requestID string) (utils.Size, error)
	DeleteSizeById(ctx context.Context, id int, requestID string) error
}

type GetProducts interface {
	ListAllProducts(ctx context.Context, offset, limit int, filters map[string]string, sort string, requestID string) ([]utils.ListProduct, error)
	GetSingleProductById(ctx context.Context, id int, requestID string) (utils.SingleProduct, error)
}

type ReviewRepository interface {
	InsertReview(ctx context.Context, review dto.Review, requestID string) (*int, error)
	ListReviews(ctx context.Context, offset, limit int, requestID string) ([]utils.Review, error)
}

type ChapaRepository interface{}
