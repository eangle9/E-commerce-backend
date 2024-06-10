package dto

import "github.com/shopspring/decimal"

type ProductCategory struct {
	ID       int    `json:"category_id,omitempty"`
	Name     string `json:"name,omitempty"`
	ParentID int    `json:"parent_id,omitempty"`
}

type Product struct {
	ID          int    `json:"product_id"`
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
	Price      decimal.Decimal `json:"price"`
	Discount   decimal.Decimal `json:"discount"`
	QtyInStock *int            `json:"qty_in_stock"`
	ImageUrl   string          `json:"image_url"`
}

type ProductImage struct {
	ID            int    `json:"image_id"`
	ProductItemID int    `json:"product_item_id"`
	ImageUrl      string `json:"image_url"`
}
