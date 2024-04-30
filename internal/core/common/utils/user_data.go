package utils

import "time"

type User struct {
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
}

type UpdateUser struct {
	Username       string
	Email          string
	Password       string
	FirstName      string
	LastName       string
	PhoneNumber    string
	Address        string
	ProfilePicture string
	Role           string
}
