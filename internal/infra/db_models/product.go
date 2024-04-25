package dbmodels

import "time"

type Product struct {
	ID          uint `gorm:"primaryKey"`
	UserID      uint
	ProductID   uint
	Name        string
	Description string
	Price       float64
	Quantity    uint
	Catagory    string
	Brand       string
	ImageUrl    string
	AddedAt     *time.Time
}
