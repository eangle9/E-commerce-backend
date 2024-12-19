package dto

import "github.com/shopspring/decimal"

type Size struct {
	ID            int             `json:"size_id"`
	ProductItemID int             `json:"product_item_id"`
	SizeName      string          `json:"size_name"`
	Price         decimal.Decimal `json:"price"`
	Discount      decimal.Decimal `json:"discount"`
	QtyInStock    int             `json:"qty_in_stock"`
}
