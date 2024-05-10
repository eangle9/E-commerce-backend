package utils

import "time"

type ProductCategory struct {
	ID        int        `json:"category_id"`
	Name      string     `json:"name"`
	ParentID  *int       `json:"parent_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type UpdateCategory struct {
	Name     string `json:"name"`
	ParentID int    `json:"parent_id"`
}

type Product struct {
	ID          int        `json:"product_id"`
	CategoryID  int        `json:"category_id"`
	ProductName string     `json:"product_name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

// type Category struct {
// 	ID        int
// 	Name      string
// 	ParentID  *int
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// 	DeletedAt *time.Time
// }
