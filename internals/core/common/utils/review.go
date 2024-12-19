package utils

import "time"

type Review struct {
	ID        int        `json:"review_id"`
	UserID    int        `json:"user_id"`
	ProductID int        `json:"product_id"`
	Rating    uint       `json:"rating"`
	Comment   string     `json:"comment"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
