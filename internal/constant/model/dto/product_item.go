package dto

import (
	"Eccomerce-website/internal/constant"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type ProductItem struct {
	ID        uuid.UUID              `json:"id"`
	ProductID uuid.UUID              `json:"product_id"`
	ColorID   uuid.UUID              `json:"color_id"`
	SizeID    uuid.UUID              `json:"size_id"`
	SKU       string                 `json:"sku"`
	Status    constant.ProductStatus `json:"status"`
	ImageURL  string                 `json:"image_url"`
	Price     float64                `json:"price"`
	Discount  float64                `json:"discount"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

func (p ProductItem) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.ProductID,
			validation.Required.Error("product_id is required"),
			validation.By(func(value interface{}) error {
				ID, ok := value.(uuid.UUID)
				if !ok {
					return fmt.Errorf("product_id is not valid")
				}
				if ID == uuid.Nil {
					return fmt.Errorf("product_id can not be nil uuid")
				}
				return nil
			})),
		validation.Field(&p.ColorID,
			validation.Required.Error("color_id is required"),
			validation.By(func(value interface{}) error {
				ID, ok := value.(uuid.UUID)
				if !ok {
					return fmt.Errorf("color_id is not valid")
				}
				if ID == uuid.Nil {
					return fmt.Errorf("color_id can not be nil uuid")
				}
				return nil
			})),
		validation.Field(&p.SizeID,
			validation.Required.Error("size_id is required"),
			validation.By(func(value interface{}) error {
				ID, ok := value.(uuid.UUID)
				if !ok {
					return fmt.Errorf("size_id is not valid")
				}
				if ID == uuid.Nil {
					return fmt.Errorf("size_id can not be nil uuid")
				}
				return nil
			})),
		validation.Field(&p.SKU, validation.Required.Error("sku is required")),
		validation.Field(&p.Status, validation.Required.Error("status is required"), validation.In(
			constant.ProductStatusActive,
			constant.ProductStatusInActive,
			constant.ProductStatusOutOfStock,
		)),
		validation.Field(&p.ImageURL, validation.Required.Error("image_url is required")),
		validation.Field(&p.Price, validation.Required.Error("price is required")),
	)
}

type Size struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s Size) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name, validation.Required.Error("name is required")),
	)
}

type Color struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c Color) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required.Error("name is required")),
	)
}
