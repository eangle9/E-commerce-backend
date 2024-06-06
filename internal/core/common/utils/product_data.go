package utils

import (
	"time"

	"github.com/shopspring/decimal"
)

type ProductCategory struct {
	ID        int        `json:"category_id"`
	Name      string     `json:"name"`
	ParentID  *int       `json:"parent_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type UpdateCategory struct {
	Name     string `json:"name"`
	ParentID int    `json:"parent_id"`
}

type Product struct {
	ID          int        `json:"product_id"`
	CategoryID  int        `json:"category_id"`
	ProductName string     `json:"product_name"`
	Brand       string     `json:"brand"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type UpdateProduct struct {
	CategoryID  int    `json:"category_id"`
	Brand       string `json:"brand"`
	ProductName string `json:"product_name"`
	Description string `json:"description"`
}

type ProductItem struct {
	ID         int             `json:"product_item_id"`
	ProductID  int             `json:"product_id"`
	ColorID    *int            `json:"color_id"`
	ImageUrl   string          `json:"image_url"`
	Price      decimal.Decimal `json:"price"`
	QtyInStock int             `json:"qty_in_stock"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
	DeletedAt  *time.Time      `json:"deleted_at"`
}

type UpdateProductItem struct {
	ProductID  int             `json:"product_id"`
	ColorID    int             `json:"color_id"`
	Price      decimal.Decimal `json:"price"`
	QtyInStock *int            `json:"qty_in_stock"`
}

type SingleProduct struct {
	ProductID    int              `json:"product_id"`
	Product      string           `json:"product"`
	ProductItems []ProductVariant `json:"product_items"`
}

type ProductVariant struct {
	ItemID   int             `json:"item_id"`
	Color    *string         `json:"color"`
	ImageUrl string          `json:"image_url"`
	Price    decimal.Decimal `json:"price"`
	InStock  *int            `json:"in stock"`
	// Sizes    []string `json:"sizes"`
}

// type Category struct {
// 	ID        int
// 	Name      string
// 	ParentID  *int
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// 	DeletedAt *time.Time
// }
