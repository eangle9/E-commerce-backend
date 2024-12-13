package dto

import (
	"fmt"

	"github.com/dongri/phonenumber"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/shopspring/decimal"
)

type PaymentRequest struct {
	Amount         float64                `json:"amount"`
	Currency       string                 `json:"currency"`
	Email          string                 `json:"email"`
	FirstName      string                 `json:"first_name"`
	LastName       string                 `json:"last_name"`
	Phone          string                 `json:"phone"`
	CallbackURL    string                 `json:"callback_url"`
	TransactionRef string                 `json:"tx_ref"`
	Customization  map[string]interface{} `json:"customization"`
}

func (p PaymentRequest) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Amount,
			validation.Required.Error("amount is required"),
			// validation.By(ValidateDecimalMin(decimal.NewFromInt(1),
			// 	fmt.Sprintf("value must be greater than or equal to %v", 1))),
		),
		validation.Field(&p.Currency, validation.Required),
		validation.Field(&p.Phone, validation.Required),
		validation.Field(&p.TransactionRef, validation.Required),
		validation.Field(&p.Phone, validation.Required, validation.By(ValidatePhoneNumber)),
	)
}

type PaymentResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Data    struct {
		CheckoutURL string `json:"checkout_url"`
	}
}

type VerifyResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Data    struct {
		TransactionFee float64 `json:"charge"`
	}
}

func ValidateDecimalMin(minValue decimal.Decimal, message string) validation.RuleFunc {
	return func(value interface{}) error {
		val, ok := value.(decimal.Decimal)
		if !ok {
			return fmt.Errorf("value must be a decimal type: %T", value)
		}

		if val.LessThan(minValue) {
			return validation.NewError("400", message)
		}
		return nil
	}
}

func ValidatePhoneNumber(phone any) error {
	str := phonenumber.Parse(fmt.Sprintf("%v", phone), "ET")
	if str == "" {
		return fmt.Errorf("invalid phone number")
	}
	return nil
}
