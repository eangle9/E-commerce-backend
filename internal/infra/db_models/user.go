package dbmodels

import "time"

type User struct {
	// gorm.Model
	ID             int
	Username       string
	Email          string
	Password       string
	FirstName      string
	LastName       string
	PhoneNumber    string
	Address        string
	ProfilePicture string
	EmailVerified  bool
	Role           string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time
	// PaymentInfo    *PaymentInfo
	// Wishlist       *[]Product
	// Cart           *[]CartItem
}
