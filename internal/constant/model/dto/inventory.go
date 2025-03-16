package dto

import (
	"Eccomerce-website/internal/constant"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type Inventory struct {
	ID                uuid.UUID                `json:"id"`
	ProductItemID     uuid.UUID                `json:"product_item_id"`
	SKU               string                   `json:"sku"`
	AvaliableQuantity int                      `json:"avaliable_quantity"`
	ReservedQuantity  int                      `json:"reserved_quantity"`
	MinimumQuantity   int                      `json:"minimum_quantity"`
	Status            constant.InventoryStatus `json:"status"`
	RestockDate       time.Time                `json:"restock_date"`
	CreatedAt         time.Time                `json:"created_at"`
	UpdatedAt         time.Time                `json:"updated_at"`
}

func (i Inventory) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.ProductItemID, validation.Required.Error("product_item_id is required"),
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
		validation.Field(&i.SKU, validation.Required.Error("sku is required")),
		validation.Field(&i.AvaliableQuantity, validation.Required.Error("avaliable_quantity is required")),
		validation.Field(&i.Status,
			validation.Required.Error("status is required"), validation.In(
				constant.InventoryStatusInStock,
				constant.InventoryStatusOutStock,
				constant.InventoryStatusLowStock,
			)),
	)
}
