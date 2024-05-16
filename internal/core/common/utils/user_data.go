package utils

import "time"

type User struct {
	ID          int    `json:"user_id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	// Address        string    `json:"address"`
	ProfilePicture string    `json:"profile_picture"`
	EmailVerified  bool      `json:"email_verified"`
	Role           string    `json:"role"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      time.Time `json:"deleted_at"`
}

type UpdateUser struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	// Address        string `json:"address"`
	ProfilePicture string `json:"profile_picture"`
	// Role           string `json:"role"`
}
