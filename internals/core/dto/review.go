package dto

type Review struct {
	ID        int    `json:"review_id"`
	UserID    int    `json:"user_id"`
	ProductID int    `json:"product_id"`
	Rating    uint   `json:"rating"`
	Comment   string `json:"comment"`
}
