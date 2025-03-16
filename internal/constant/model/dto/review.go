package dto

import (
	"Eccomerce-website/internal/constant"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type Review struct {
	ID        uuid.UUID       `json:"id"`
	ProductID uuid.UUID       `json:"product_id"`
	UserID    uuid.UUID       `json:"user_id"`
	OrderID   uuid.UUID       `json:"order_id"`
	Rating    constant.Rating `json:"rating"`
	Comment   string          `json:"comment"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

func (r Review) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ProductID,
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
		validation.Field(&r.UserID,
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
		validation.Field(&r.OrderID,
			validation.Required.Error("order_id is required"),
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
		validation.Field(&r.Rating,
			validation.Required.Error("rating is required"),
			validation.In(
				constant.RatingOne,
				constant.RatingTwo,
				constant.RatingThree,
				constant.RatingFour,
				constant.RatingFive,
			)),
	)
}
