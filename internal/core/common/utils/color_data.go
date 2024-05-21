package utils

import "time"

type Color struct {
	ID        int        `json:"color_id,omitempty"`
	ColorName string     `json:"color_name,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type UpdateColor struct {
	ColorName string `json:"color_name"`
}

type Size struct {
	ID        int        `json:"size_id,omitempty"`
	SizeName  string     `json:"size_name,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type UpdateSize struct {
	SizeName string `json:"size_name"`
}
