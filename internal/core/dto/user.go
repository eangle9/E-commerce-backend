package dto

type User struct {
	Username       string `json:"username,omitempty"`
	Email          string `json:"email,omitempty"`
	Password       string `json:"password,omitempty"`
	FirstName      string `json:"firstName,omitempty"`
	LastName       string `json:"lastName,omitempty"`
	PhoneNumber    string `json:"phoneNumber,omitempty"`
	Address        string `json:"address,omitempty"`
	Role           string `json:"role,omitempty"`
	EmailVerified  bool   `json:"emailVerified"`
	ProfilePicture string `json:"profilePicture,omitempty"`
	// CreatedAt      time.Time      `json:"createdAt,omitempty"`
	// UpdatedAt      time.Time      `json:"updatedAt,omitempty"`
	// DeletedAt      gorm.DeletedAt `json:"deletedAt,omitempty"`
}
