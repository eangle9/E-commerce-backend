package dto

import (
	"Eccomerce-website/internal/constant"
	"fmt"
	"regexp"
	"time"

	"github.com/dongri/phonenumber"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID         `json:"id"`
	Username       string            `json:"username"`
	Email          string            `json:"email"`
	PhoneNumber    string            `json:"phone_number"`
	Password       string            `json:"password"`
	FirstName      string            `json:"first_name"`
	LastName       string            `json:"last_name"`
	ProfilePicture string            `json:"profile_picture"`
	EmailVerified  bool              `json:"email_verified"`
	Role           constant.UserRole `json:"role"`
	LastLogin      time.Time         `json:"last_login"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

type CreateUserParams struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phone_number"`
	Password       string `json:"password"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	ProfilePicture string `json:"profile_picture"`
	EmailVerified  bool   `json:"email_verified"`
}

func (c CreateUserParams) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Username,
			validation.Required.Error("username is required"),
			validation.Length(3, 30).Error("username must be between 3 and 30 characters")),
		validation.Field(&c.Email,
			validation.Required.Error("email is required"),
			is.Email.Error("email is not valid")),
		validation.Field(&c.PhoneNumber,
			validation.Required.Error("phone number is required"),
			validation.By(ValidatePhone)),
		validation.Field(&c.Password,
			validation.Required.Error("password is required"),
			validation.Length(5, 30).Error("password must be between 5 and 30 characters"),
			validation.Match(regexp.MustCompile(`[a-zA-Z0-9]`)).Error("password must contain at least one letter and one number")),
		validation.Field(&c.FirstName, validation.Required.Error("first name is required")),
		validation.Field(&c.LastName, validation.Required.Error("last name is required")),
	)
}

func ValidatePhone(phone any) error {
	str := phonenumber.Parse(fmt.Sprintf("%v", phone), "ET")
	if str == "" {
		return fmt.Errorf("invalid phone number")
	}
	return nil
}

type UpdateUserParams struct {
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	PhoneNumber    string    `json:"phone_number"`
	Password       string    `json:"password"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	ProfilePicture string    `json:"profile_picture"`
	LastLogin      time.Time `json:"last_login"`
}

type UserAddress struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Country     string    `json:"country"`
	City        string    `json:"city"`
	SubCity     string    `json:"sub_city"`
	Woreda      string    `json:"woreda"`
	Kebele      string    `json:"kebele"`
	Street      string    `json:"street"`
	PhoneNumber string    `json:"phone_number"`
	IsPrimary   bool      `json:"is_primary"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (c UserAddress) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.UserID, validation.Required.Error("user id is required"), validation.By(func(value interface{}) error {
			ID, ok := value.(uuid.UUID)
			if !ok {
				return fmt.Errorf("invalid user id")
			}
			if ID == uuid.Nil {
				return fmt.Errorf("user_id can not be nil uuid")
			}
			return nil
		})),
		validation.Field(&c.Country, validation.Required.Error("country is required")),
		validation.Field(&c.City, validation.Required.Error("city is required")),
		validation.Field(&c.SubCity, validation.Required.Error("sub city is required")),
		validation.Field(&c.Street, validation.Required.Error("street is required")),
		validation.Field(&c.PhoneNumber, validation.Required.Error("phone number is required")),
	)
}
