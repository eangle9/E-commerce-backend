package state

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type TokenKey struct {
	SymmetricKey string `json:"symmetric_key"`
	Issuer       string `json:"issuer"`
	Footer       string `json:"footer"`
	KeyLength    int    `json:"key_length"`
	// this is the duration in minute
	ExpireAt time.Duration `json:"expire_at"`
}

func (m TokenKey) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.SymmetricKey, validation.Required.Error("Symmetric-key is required")),
		validation.Field(&m.Issuer, validation.Required.Error("issuer is required")),
		validation.Field(&m.Footer, validation.Required.Error("footer short code is required")),
		validation.Field(&m.KeyLength, validation.Required.Error("key length is required")),
		validation.Field(&m.ExpireAt, validation.Required.Error("expire at is required")),
	)
}
