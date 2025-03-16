package dto

import (
	"Eccomerce-website/internal/constant"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type Wishlist struct {
	ID            uuid.UUID                   `json:"id"`
	UserID        uuid.UUID                   `json:"user_id"`
	ProductItemID uuid.UUID                   `json:"product_item_id"`
	Visibility    constant.WishlistVisibility `json:"visibility"`
	CreatedAt     time.Time                   `json:"created_at"`
	UpdatedAt     time.Time                   `json:"updated_at"`
}

func (w Wishlist) Validate() error {
	return validation.ValidateStruct(&w,
		validation.Field(&w.UserID,
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
		validation.Field(&w.ProductItemID,
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
	)
}
