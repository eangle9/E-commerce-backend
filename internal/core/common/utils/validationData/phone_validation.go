package validationdata

import (
	"fmt"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	PhoneReg = regexp.MustCompile(`^(\+251|251|0)?[97]\d{8}$`)
)

func FormatPhoneNumber(phone string) (string, error) {
	if isValid := PhoneReg.MatchString(phone); !isValid {
		err := fmt.Errorf("invalid phone number format.Please enter a valid phone number")
		return "", err
	}
	reg := regexp.MustCompile(`[^\d]`)
	phone = reg.ReplaceAllString(phone, "")

	if phone[:1] == "0" {
		phone = phone[1:]
	}

	if phone[:3] != "251" {
		phone = "251" + phone
	}

	return phone, nil
}

func ValidatePhoneNumber(value interface{}) error {
	phoneNumber, exist := value.(string)
	if !exist {
		return validation.NewError("validation_phonenumber", "cannot be blank")
	}

	isValid := regexp.MustCompile(`^(\+251|251|0)?[97]\d{8}$`)
	if !isValid.MatchString(phoneNumber) {
		return validation.NewError("validation_phonenuber", "invalid phoneNumber format")
	}

	return nil
}
