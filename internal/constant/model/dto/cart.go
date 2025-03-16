package dto

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type ShoppingCart struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	SubTotal      float64   `json:"sub_total"`
	DeliveryFee   float64   `json:"delivery_fee"`
	ServiceCharge float64   `json:"service_charge"`
	Vat           float64   `json:"vat"`
	Total         float64   `json:"total"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (s ShoppingCart) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.UserID,
			validation.Required.Error("user_id is required"),
			validation.By(func(value interface{}) error {
				ID, ok := value.(uuid.UUID)
				if !ok {
					return fmt.Errorf("user_id is not valid")
				}
				if ID == uuid.Nil {
					return fmt.Errorf("user_id can not be nil uuid")
				}
				return nil
			})),
		validation.Field(&s.SubTotal, validation.Required.Error("sub_total is required")),
		validation.Field(&s.Total, validation.Required.Error("total is required")),
	)
}

type CartItems struct {
	ID            uuid.UUID `json:"id"`
	CartID        uuid.UUID `json:"cart_id"`
	ProductItemID uuid.UUID `json:"product_item_id"`
	Price         float64   `json:"price"`
	Quantity      int64       `json:"quantity"`
	TotalPrice    float64   `json:"total_price"`
}

func (c CartItems) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.CartID,
			validation.Required.Error("cart_id is required"),
			validation.By(func(value interface{}) error {
				ID, ok := value.(uuid.UUID)
				if !ok {
					return fmt.Errorf("cart_id is not valid")
				}
				if ID == uuid.Nil {
					return fmt.Errorf("cart_id can not be nil uuid")
				}
				return nil
			})),
		validation.Field(&c.ProductItemID,
			validation.Required.Error("product_item_id is required"),
			validation.By(func(value interface{}) error {
				ID, ok := value.(uuid.UUID)
				if !ok {
					return fmt.Errorf("product_item_id is not valid")
				}
				if ID == uuid.Nil {
					return fmt.Errorf("product_item_id can not be nil uuid")
				}
				return nil
			})),
		validation.Field(&c.Price, validation.Required.Error("price is required")),
		validation.Field(&c.Quantity, validation.Required.Error("quantity is required")),
		validation.Field(&c.TotalPrice, validation.Required.Error("total_price is required")),
	)
}
