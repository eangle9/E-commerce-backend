package dto

import (
	"Eccomerce-website/internal/constant"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type OrderDetails struct {
	ID            uuid.UUID            `json:"id"`
	UserID        uuid.UUID            `json:"user_id"`
	SubTotal      float64              `json:"sub_total"`
	DeliveryFee   float64              `json:"delivery_fee"`
	ServiceCharge float64              `json:"service_charge"`
	Vat           float64              `json:"vat"`
	Total         float64              `json:"total"`
	Status        constant.OrderStatus `json:"status"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
}

func (o OrderDetails) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.UserID, validation.Required.Error("user_id is required"),
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
		validation.Field(&o.SubTotal, validation.Required.Error("sub_total is required")),
		validation.Field(&o.Total, validation.Required.Error("total is required")),
		validation.Field(&o.Status, validation.Required.Error("status is required"), validation.In(
			constant.OrderStatusPending,
			constant.OrderStatusProcessing,
			constant.OrderStatusShipped,
			constant.OrderStatusDelivered,
			constant.OrderStatusCancelled,
			constant.OrderStatusReturned,
		)),
	)
}

type OrderItems struct {
	ID            uuid.UUID `json:"id"`
	OrderID       uuid.UUID `json:"order_id"`
	ProductItemID uuid.UUID `json:"product_item_id"`
	Price         float64   `json:"price"`
	Quantity      int64       `json:"quantity"`
	TotalPrice    float64   `json:"total_price"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (o OrderItems) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.OrderID, validation.Required.Error("order_id is required"),
			validation.By(func(value interface{}) error {
				ID, ok := value.(uuid.UUID)
				if !ok {
					return fmt.Errorf("order_id is not valid")
				}
				if ID == uuid.Nil {
					return fmt.Errorf("order_id can not be nil uuid")
				}
				return nil
			})),
		validation.Field(&o.ProductItemID, validation.Required.Error("product_item_id is required"),
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
		validation.Field(&o.Price, validation.Required.Error("price is required")),
		validation.Field(&o.Quantity, validation.Required.Error("quantity is required")),
		validation.Field(&o.TotalPrice, validation.Required.Error("total_price is required")),
	)
}
