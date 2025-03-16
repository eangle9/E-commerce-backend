package dto

import (
	"Eccomerce-website/internal/constant"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type Payments struct {
	ID            uuid.UUID              `json:"id"`
	OrderID       uuid.UUID              `json:"order_id"`
	UserID        uuid.UUID              `json:"user_id"`
	PaymentMethod string                 `json:"payment_method"`
	Status        constant.PaymentStatus `json:"status"`
	TransactionID string                 `json:"transaction_id"`
	TxRef         string                 `json:"tx_ref"`
	Reference     string                 `json:"reference"`
	Type          string                 `json:"type"`
	TotalAmount   float64                `json:"total_amount"`
	CurrencyCode  constant.Currency      `json:"currency_code"`
	CreatedAT     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

func (p Payments) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.OrderID, validation.Required.Error("order_id is required"),
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
		validation.Field(&p.UserID, validation.Required.Error("user_id is required"),
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
		validation.Field(&p.PaymentMethod, validation.Required.Error("payment_method is required")),
		validation.Field(&p.Status,
			validation.Required.Error("status is required"),
			validation.In(
				constant.PaymentStatusPending,
				constant.PaymentStatusCompleted,
				constant.PaymentStatusFailed,
				constant.PaymentStatusRefunded,
				constant.PaymentStatusCancelled,
			)),
		validation.Field(&p.TxRef, validation.Required.Error("tx_ref is required")),
		validation.Field(&p.TotalAmount, validation.Required.Error("total_amount is required")),
		validation.Field(&p.CurrencyCode, validation.Required.Error("currency_code is required")),
	)
}
