package dbmodels

type PaymentInfo struct {
	ID             uint `gorm:"primaryKey"`
	UserID         uint
	CardNumber     *string
	ExpiryDate     *string
	BillingAddress *string
}
