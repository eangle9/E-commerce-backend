package dto

import (
	"Eccomerce-website/internal/constant"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type Product struct {
	ID            uuid.UUID              `json:"id"`
	CategoryID    uuid.UUID              `json:"category_id"`
	Brand         string                 `json:"brand"`
	Name          string                 `json:"name"`
	Status        constant.ProductStatus `json:"status"`
	Description   string                 `json:"description"`
	AverageRating float64                `json:"average_rating"`
	TotalReviews  int64                    `json:"total_reviews"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

type CreateProductParams struct {
	CategoryID  uuid.UUID              `json:"category_id"`
	Brand       string                 `json:"brand"`
	Name        string                 `json:"name"`
	Status      constant.ProductStatus `json:"status"`
	Description string                 `json:"description"`
}

func (c CreateProductParams) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.CategoryID,
			validation.Required.Error("category_id is required"),
			validation.By(func(value interface{}) error {
				ID, ok := value.(uuid.UUID)
				if !ok {
					return fmt.Errorf("category_id is not valid")
				}
				if ID == uuid.Nil {
					return fmt.Errorf("category_id can not be nil uuid")
				}
				return nil
			})),
		validation.Field(&c.Name, validation.Required.Error("name is required")),
		validation.Field(&c.Status, validation.Required.Error("status is required"), validation.In(
			constant.ProductStatusActive,
			constant.ProductStatusInActive,
			constant.ProductStatusOutOfStock,
		)),
	)
}
