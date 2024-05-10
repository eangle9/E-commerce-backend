package repository

import (
	"Eccomerce-website/internal/core/common/utils"
	"Eccomerce-website/internal/core/dto"
	"Eccomerce-website/internal/core/port/repository"
	"errors"
	"fmt"
	"strings"
	"time"
)

type colorRepository struct {
	db repository.Database
}

func NewColorRepository(db repository.Database) repository.ColorRepository {
	return &colorRepository{
		db: db,
	}
}

func (c colorRepository) InsertColor(color dto.Color) (*int, error) {
	DB := c.db.GetDB()

	var count int
	if err := DB.QueryRow("SELECT COUNT(*) FROM color WHERE color_name = ? AND deleted_at IS NULL", color.Name).Scan(&count); err != nil {
		return nil, err
	}

	if count > 0 {
		err := fmt.Errorf("color with color_name '%s' already exists", color.Name)
		return nil, err
	}

	query := `INSERT INTO color(color_name) VALUES(?)`
	result, err := DB.Exec(query, color.Name)
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

func (c colorRepository) ListColors() ([]utils.Color, error) {
	var colors []utils.Color
	DB := c.db.GetDB()

	query := `SELECT color_id, color_name, created_at, updated_at, deleted_at FROM color WHERE deleted_at IS NULL`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var color utils.Color
		if err := rows.Scan(&color.ID, &color.ColorName, &color.CreatedAt, &color.UpdatedAt, &color.DeletedAt); err != nil {
			return nil, err
		}

		colors = append(colors, color)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return colors, nil
}

func (c colorRepository) GetColorById(id int) (utils.Color, error) {
	var color utils.Color
	DB := c.db.GetDB()

	query := `SELECT color_id, color_name, created_at, updated_at, deleted_at FROM color WHERE color_id = ? AND deleted_at IS NULL`
	if err := DB.QueryRow(query, id).Scan(&color.ID, &color.ColorName, &color.CreatedAt, &color.UpdatedAt, &color.DeletedAt); err != nil {
		err = fmt.Errorf("color with color_id '%d' not found", id)
		return utils.Color{}, err
	}

	return color, nil
}

func (c colorRepository) EditColorById(id int, color utils.UpdateColor) (utils.Color, error) {
	DB := c.db.GetDB()
	var updateFields []string
	var values []interface{}

	if color.ColorName != "" {
		updateFields = append(updateFields, "color_name = ?")
		values = append(values, color.ColorName)
	}

	if len(updateFields) == 0 {
		err := errors.New("failed to update color:No fields provided for update.Please provide at least one field to update")
		return utils.Color{}, err
	}

	if len(values) > 0 {
		updateFields = append(updateFields, "updated_at = ?")
		values = append(values, time.Now())
	}

	query := fmt.Sprintf("UPDATE color SET %s WHERE color_id = ? AND deleted_at IS NULL", strings.Join(updateFields, ", "))
	values = append(values, id)
	if _, err := DB.Exec(query, values...); err != nil {
		return utils.Color{}, err
	}

	updatedColor, err := c.GetColorById(id)
	if err != nil {
		return utils.Color{}, err
	}

	return updatedColor, nil
}
