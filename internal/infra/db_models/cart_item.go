package dbmodels

import "time"

type CartItem struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	ProductID uint
	AddedAt   time.Time
}
