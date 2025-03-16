package constant

type Domain string

const (
	Client   Domain = "client"
	Merchant Domain = "merchant"
	User     Domain = "user"
	System   Domain = "system"
)

type UserRole string

const (
	UserRoleAdmin    UserRole = "ADMIN"
	UserRoleVendor   UserRole = "VENDOR"
	UserRoleCustomer UserRole = "CUSTOMER"
)

type ProductStatus string

const (
	ProductStatusActive     ProductStatus = "ACTIVE"
	ProductStatusInActive   ProductStatus = "INACTIVE"
	ProductStatusOutOfStock ProductStatus = "OUT_OF_STOCK"
)

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "PENDING"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusShipped    OrderStatus = "SHIPPED"
	OrderStatusDelivered  OrderStatus = "DELIVERED"
	OrderStatusCancelled  OrderStatus = "CANCELED"
	OrderStatusReturned   OrderStatus = "RETURNED"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusCompleted PaymentStatus = "COMPLETED"
	PaymentStatusFailed    PaymentStatus = "FAILED"
	PaymentStatusRefunded  PaymentStatus = "REFUNDED"
	PaymentStatusCancelled PaymentStatus = "CANCELED"
)

type Currency string

const (
	CurrencyETB Currency = "ETB"
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
	CurrencyGBP Currency = "GBP"
)

type Rating string

const (
	RatingOne   Rating = "1"
	RatingTwo   Rating = "2"
	RatingThree Rating = "3"
	RatingFour  Rating = "4"
	RatingFive  Rating = "5"
)

type InventoryStatus string

const (
	InventoryStatusInStock  InventoryStatus = "IN_STOCK"
	InventoryStatusLowStock InventoryStatus = "LOW_STOCK"
	InventoryStatusOutStock InventoryStatus = "OUT_OF_STOCK"
)

type WishlistVisibility string

const (
	WishlistVisibilityPrivate WishlistVisibility = "PRIVATE"
	WishlistVisibilityPublic  WishlistVisibility = "PUBLIC"
)
