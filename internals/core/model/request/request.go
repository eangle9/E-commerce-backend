package request

import (
	validationdata "Eccomerce-website/internals/core/common/utils/validationData"
	"encoding/json"
	"mime/multipart"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/shopspring/decimal"
)

type SignUpRequest struct {
	Username    string `json:"username" `
	Email       string `json:"email"`
	Password    string `json:"password" `
	FirstName   string `json:"firstName" `
	LastName    string `json:"lastName" `
	PhoneNumber Phone  `json:"phoneNumber" `
}

type Phone struct {
	Number string
}

func (p *Phone) UnmarshalJSON(data []byte) error {
	// if string(data) == "null" || string(data) == `""` {
	// 	err := errors.New("phoneNumber is required")
	// 	return err
	// }

	var phoneNumber string
	if err := json.Unmarshal(data, &phoneNumber); err != nil {
		return err
	}

	formatedPhone, err := validationdata.FormatPhoneNumber(phoneNumber)
	if err != nil {
		return err
	}

	p.Number = formatedPhone

	return nil
}

func (p Phone) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Number, validation.Required, validation.By(validationdata.ValidatePhoneNumber)),
	)
}

func (s SignUpRequest) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Username, validation.Required, validation.Length(5, 0), is.Alphanumeric),
		validation.Field(&s.Email, validation.Required, is.Email),
		validation.Field(&s.Password, validation.Length(8, 0), validation.By(validationdata.ValidatePassword)),
		validation.Field(&s.FirstName, validation.Required, validation.Length(2, 0)),
		validation.Field(&s.LastName, validation.Required, validation.Length(2, 0)),
		validation.Field(&s.PhoneNumber, validation.Required, validation.By(func(value interface{}) error {
			phone, _ := value.(Phone)
			return phone.Validate()
		})),
	)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l LoginRequest) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Email, validation.Required, is.Email),
		validation.Field(&l.Password, validation.Required, validation.By(validationdata.ValidatePassword)),
	)
}

type PaginationQuery struct {
	Page    int `form:"page"`
	PerPage int `form:"per_page"`
}

func (pq PaginationQuery) Validate() error {
	return validation.ValidateStruct(&pq,
		validation.Field(&pq.Page, validation.Required.When(pq.Page != 0), validation.Min(0)),
		validation.Field(&pq.PerPage, validation.Required.When(pq.PerPage != 0), validation.Min(0)),
	)
}

type SearchQuery struct {
	Name string `form:"name"`
}

type CategoryQuery struct {
	Category string `form:"category"`
}

type SortQuery struct {
	Sort string `form:"sort"`
}

type UpdateUser struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	PhoneNumber    string `json:"phone_number"`
	ProfilePicture string `json:"profile_picture"`
}

func (s UpdateUser) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Username, validation.By(func(value interface{}) error {
			if s.Username != "" {
				return validation.Validate(value, validation.Length(5, 0), is.Alphanumeric)
			}
			return nil
		})),
		validation.Field(&s.Email, validation.By(func(value interface{}) error {
			if s.Email != "" {
				return validation.Validate(value, is.Email)
			}
			return nil
		})),
		validation.Field(&s.Password, validation.By(func(value interface{}) error {
			if s.Password != "" {
				return validation.Validate(value, validation.Length(8, 0), validation.By(validationdata.ValidatePassword))
			}
			return nil
		})),
		validation.Field(&s.FirstName, validation.By(func(value interface{}) error {
			if s.FirstName != "" {
				return validation.Validate(value, validation.Length(2, 0))
			}
			return nil
		})),
		validation.Field(&s.LastName, validation.By(func(value interface{}) error {
			if s.LastName != "" {
				return validation.Validate(value, validation.Length(2, 0))
			}
			return nil
		})),
		validation.Field(&s.PhoneNumber, validation.By(func(value interface{}) error {
			if s.PhoneNumber != "" {
				return validation.Validate(value, validation.By(validationdata.ValidatePhoneNumber))
			}
			return nil
		})),
	)
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (r RefreshRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.RefreshToken, validation.Required),
	)
}

type ProductCategoryRequest struct {
	ParentID int    `json:"parent_id"`
	Name     string `json:"name" `
}

func (p ProductCategoryRequest) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Name, validation.Required),
	)
}

type ColorRequest struct {
	ColorName string `json:"color_name"`
}

func (c ColorRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.ColorName, validation.Required),
	)
}

type ProductRequest struct {
	CategoryID  int    `json:"category_id"`
	Brand       string `json:"brand"`
	ProductName string `json:"product_name"`
	Description string `json:"description"`
}

func (p ProductRequest) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.CategoryID, validation.Required),
		validation.Field(&p.ProductName, validation.Required),
		validation.Field(&p.Description, validation.Required),
	)
}

type ProductItemRequest struct {
	ProductID  int                   `json:"product_id"`
	ColorID    *int                  `json:"color_id"`
	Price      decimal.Decimal       `json:"price"`
	Discount   decimal.Decimal       `json:"discount"`
	QtyInStock int                   `json:"qty_in_stock"`
	File       *multipart.FileHeader `form:"file" binding:"required"`
}

func (p ProductItemRequest) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.ProductID, validation.Required),
		validation.Field(&p.Price, validation.Required),
	)
}

type ProductImageRequest struct {
	ProductItemId int                   `json:"product_item_id"`
	File          *multipart.FileHeader `json:"file"`
}

func (p ProductImageRequest) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.ProductItemId, validation.Required),
	)
}

type CartRequest struct {
	ProductItemID int `json:"product_item_id"`
	Quantity      int `json:"quantity"`
}

func (c CartRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.ProductItemID, validation.Required),
		validation.Field(&c.Quantity, validation.Required),
	)
}

// size request
type SizeRequest struct {
	ProductItemID int             `json:"product_item_id"`
	SizeName      string          `json:"size_name"`
	Price         decimal.Decimal `json:"price"`
	Discount      decimal.Decimal `json:"discount"`
	QtyInStock    int             `json:"qty_in_stock"`
}

func (s SizeRequest) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.ProductItemID, validation.Required),
		validation.Field(&s.SizeName, validation.Required),
		validation.Field(&s.Price, validation.Required),
		validation.Field(&s.QtyInStock, validation.Required),
	)
}

type UpdateSize struct {
	SizeName   string          `json:"size_name"`
	Price      decimal.Decimal `json:"price"`
	Discount   decimal.Decimal `json:"discount"`
	QtyInStock int             `json:"qty_in_stock"`
}

// review request
type ReviewRequest struct {
	UserID    int    `json:"user_id"`
	ProductID int    `json:"product_id"`
	Rating    uint   `json:"rating"`
	Comment   string `json:"comment"`
}

func (r ReviewRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.UserID, validation.Required),
		validation.Field(&r.ProductID, validation.Required),
		validation.Field(&r.Rating, validation.Required),
		validation.Field(&r.Comment, validation.Required),
	)
}
