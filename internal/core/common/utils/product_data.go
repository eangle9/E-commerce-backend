package utils

import (
	"mime/multipart"
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
	SizeID     *int            `json:"size_id"`
	ImageUrl   string          `json:"image_url"`
	Price      decimal.Decimal `json:"price"`
	Discount   decimal.Decimal `json:"discount"`
	QtyInStock int             `json:"qty_in_stock"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
	DeletedAt  *time.Time      `json:"deleted_at"`
}

type UpdateProductItem struct {
	ProductID  *int            `json:"product_id"`
	ColorID    *int            `json:"color_id"`
	SizeID     *int            `json:"size_id"`
	Price      decimal.Decimal `json:"price"`
	Discount   decimal.Decimal `json:"discount"`
	QtyInStock *int            `json:"qty_in_stock"`
	// ImageUrl   *string         `json:"image_url"`
	File *multipart.FileHeader
}

type ListProduct struct {
	ProductID    int              `json:"product_id"`
	Name         string           `json:"name"`
	ProductItems []ProductVariant `json:"product_items"`
	Reviews      []ProductReview  `json:"reviews"`
	// Sizes        []ProductSize    `json:"sizes"`
}

type ProductVariant struct {
	ItemID int     `json:"item_id"`
	Color  *string `json:"color"`
	// Size     *string         `json:"size"`
	ImageUrl string          `json:"image_url"`
	Price    decimal.Decimal `json:"price"`
	Discount decimal.Decimal `json:"discount"`
	InStock  *int            `json:"in_stock"`
	Sizes    []ProductSize   `json:"sizes"`
}

type SingleProduct struct {
	ProductID   int             `json:"product_id"`
	Name        string          `json:"name"`
	Category    string          `json:"category"`
	Brand       string          `json:"brand"`
	Description string          `json:"description"`
	Items       []ItemVariant   `json:"items"`
	Reviews     []ProductReview `json:"reviews"`
	// Sizes       []ProductSize   `json:"sizes"`
	// Sizes       []string        `json:"sizes"`
}

type ItemVariant struct {
	ItemID int     `json:"item_id"`
	Color  *string `json:"color"`
	// Size     *string         `json:"size"`
	ImageUrl string          `json:"image_url"`
	Price    decimal.Decimal `json:"price"`
	Discount decimal.Decimal `json:"discount"`
	InStock  *int            `json:"in_stock"`
	Sizes    []ProductSize   `json:"sizes"`
}

type ProductReview struct {
	ReviewID  int       `json:"review_id"`
	UserID    *int      `json:"user_id"`
	ProductID *int      `json:"product_id"`
	Rating    uint      `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	User      Reviewer  `json:"user"`
}

type ProductSize struct {
	ID         int             `json:"id"`
	Size       string          `json:"size"`
	Price      decimal.Decimal `json:"price"`
	Discount   decimal.Decimal `json:"discount"`
	QtyInStock int             `json:"qty_in_stock"`
}

type Reviewer struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

// type Category struct {
// 	ID        int
// 	Name      string
// 	ParentID  *int
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// 	DeletedAt *time.Time
// }
