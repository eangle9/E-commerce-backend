package utils

type CategoryList struct {
	ID       int    `json:"category_id,omitempty"`
	Name     string `json:"name,omitempty"`
	ParentID *int   `json:"parent_id"`
}
