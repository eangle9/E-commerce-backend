package dto

type ProductCategory struct {
	ID       int    `json:"category_id,omitempty"`
	Name     string `json:"name,omitempty"`
	ParentID int    `json:"parent_id,omitempty"`
}

type Product struct {
	ID          int    `json:"product_id"`
	CategoryID  int    `json:"category_id"`
	ProductName string `json:"product_name"`
	Description string `json:"description"`
}

type ProductItem struct {
	ID         int  `json:"product_item_id"`
	ProductID  int  `json:"product_id"`
	ColorID    *int `json:"color_id"`
	Price      int  `json:"price"`
	QtyInStock int  `json:"qty_in_stock"`
}
