package request

type SignUpRequest struct {
	Username    string `json:"username" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required"`
	FirstName   string `json:"firstName" validate:"required"`
	LastName    string `json:"lastName" validate:"required"`
	PhoneNumber string `json:"phoneNumber" validate:"required"`
	Address     string `json:"address" validate:"required"`
	Role        string `json:"role" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ProductCategoryRequest struct {
	ParentID int    `json:"parent_id"`
	Name     string `json:"name" validate:"required"`
}

type ColorRequest struct {
	ColorName string `json:"color_name" validate:"required"`
}
