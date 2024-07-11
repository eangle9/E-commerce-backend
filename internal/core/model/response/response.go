package response

import "github.com/shopspring/decimal"

type Response struct {
	Data       interface{} `json:"data"`
	StatusCode int         `json:"statusCode"`
	Message    interface{} `json:"message"`
}

type CartResponse struct {
	ImageUrl     string          `json:"image_url"`
	ProductName  string          `json:"product_name"`
	Color        string          `json:"color"`
	Size         string          `json:"size"`
	Description  string          `json:"description"`
	UnitPrice    decimal.Decimal `json:"unit_price"`
	QtyInStock   string          `json:"qty_in_stock"`
	CartItemID   int             `json:"cart_item_id"`
	Quantity     int             `json:"quantity"`
	SubTotal     decimal.Decimal `json:"sub_total"`
	CartSubTotal decimal.Decimal `json:"cart_sub_total"`
	Total        decimal.Decimal `json:"total"`
}
