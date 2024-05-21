package repository

import (
	"Eccomerce-website/internal/core/dto"
	"Eccomerce-website/internal/core/port/repository"
	"fmt"
)

type sizeRepository struct {
	db repository.Database
}

func NewSizeRepository(db repository.Database) repository.SizeRepository {
	return &sizeRepository{
		db: db,
	}
}

func (s sizeRepository) InsertSize(size dto.Size) (*int, error) {
	DB := s.db.GetDB()

	var count int
	if err := DB.QueryRow("SELECT COUNT(*) FROM size WHERE size_name = ?", size.SizeName).Scan(&count); err != nil {
		return nil, err
	}

	if count > 0 {
		err := fmt.Errorf("size with size_name '%s' already exists")
		return nil, err
	}

	query := `INSERT INTO size (size_name) VALUES (?)`
	result, err := DB.Exec(query, size.SizeName)
	if err != nil {
		return nil, err
	}

	id64, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	id := int(id64)

	return &id, nil
}
