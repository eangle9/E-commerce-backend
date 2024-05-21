package request

import (
	"mime/multipart"

	"github.com/shopspring/decimal"
)

type SignUpRequest struct {
	Username    string `json:"username" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required"`
	FirstName   string `json:"firstName" validate:"required"`
	LastName    string `json:"lastName" validate:"required"`
	PhoneNumber string `json:"phoneNumber" validate:"required"`
	// Address     string `json:"address" validate:"required"`
	// Role string `json:"role" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ProductCategoryRequest struct {
	ParentID int    `json:"parent_id"`
	Name     string `json:"name" validate:"required"`
}

type ColorRequest struct {
	ColorName string `json:"color_name" validate:"required"`
}

type ProductRequest struct {
	CategoryID  int    `json:"category_id" validate:"required"`
	Brand       string `json:"brand"`
	ProductName string `json:"product_name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type ProductItemRequest struct {
	ProductID  int             `json:"product_id" validate:"required"`
	ColorID    *int            `json:"color_id"`
	Price      decimal.Decimal `json:"price" validate:"required"`
	QtyInStock int             `json:"qty_in_stock"`
}

type ProductImageRequest struct {
	ProductItemId int `json:"product_item_id" validate:"required"`
	File          *multipart.FileHeader
}

// type CartItem struct {
// 	ProductItemID int `json:"product_item_id" validate:"required"`
// 	Quantity      int `json:"quantity" validate:"required"`
// }

type CartRequest struct {
	// UserID        int `json:"user_id" validate:"required"`
	ProductItemID int `json:"product_item_id" validate:"required"`
	Quantity      int `json:"quantity" validate:"required"`
}

type SizeRequest struct {
	SizeName string `json:"size_name" validate:"required"`
}
