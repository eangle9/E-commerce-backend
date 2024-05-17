package response

import "github.com/shopspring/decimal"

type Response struct {
	Data         interface{} `json:"data"`
	Status       int         `json:"status"`
	ErrorType    string      `json:"type"`
	ErrorMessage interface{} `json:"message"`
}

type CartResponse struct {
	ImageUrl     string          `json:"image_url"`
	ProductName  string          `json:"product_name"`
	Description  string          `json:"description"`
	UnitPrice    decimal.Decimal `json:"unit_price"`
	CartItemID   int             `json:"cart_item_id"`
	Quantity     int             `json:"quantity"`
	SubTotal     decimal.Decimal `json:"sub_total"`
	CartSubTotal decimal.Decimal `json:"cart_sub_total"`
	Total        decimal.Decimal `json:"total"`
}
