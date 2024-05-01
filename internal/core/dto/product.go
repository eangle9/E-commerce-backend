package dto

type ProductCategory struct {
	ID       int    `json:"category_id,omitempty"`
	ParentID int    `json:"parent_id,omitempty"`
	Name     string `json:"name,omitempty"`
}
