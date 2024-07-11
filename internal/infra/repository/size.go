package repository

import (
	"Eccomerce-website/internal/core/common/utils"
	"Eccomerce-website/internal/core/dto"
	"Eccomerce-website/internal/core/entity"
	"Eccomerce-website/internal/core/port/repository"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/shopspring/decimal"
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
	if err := DB.QueryRow("SELECT COUNT(*) FROM size WHERE size_name = ? AND product_item_id = ?", size.SizeName, size.ProductItemID).Scan(&count); err != nil {
		return nil, err
	}

	if count > 0 {
		err := fmt.Errorf("size with size_name '%s' and product_item_id '%d' already exists", size.SizeName, size.ProductItemID)
		return nil, err
	}

	query := `INSERT INTO size (product_item_id, size_name, price, discount, qty_in_stock) VALUES (?, ?, ?, ?, ?)`
	result, err := DB.Exec(query, size.ProductItemID, size.SizeName, size.Price, size.Discount, size.QtyInStock)
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

func (s sizeRepository) ListSizes() ([]utils.Size, error) {
	var sizes []utils.Size
	DB := s.db.GetDB()

	query := `SELECT size_id, size_name, price, discount, qty_in_stock, created_at, updated_at, deleted_at FROM size WHERE deleted_at IS NULL`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var size utils.Size

		if err := rows.Scan(&size.ID, &size.SizeName, &size.Price, &size.Discount, &size.QtyInStock, &size.CreatedAt, &size.UpdatedAt, &size.DeletedAt); err != nil {
			return nil, err
		}

		sizes = append(sizes, size)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sizes, nil

}

func (s sizeRepository) GetSizeById(id int) (utils.Size, error) {
	var size utils.Size
	DB := s.db.GetDB()

	query := `SELECT size_id, size_name, price, discount, qty_in_stock, created_at, updated_at, deleted_at FROM size WHERE size_id = ? AND deleted_at IS NULL`
	if err := DB.QueryRow(query, id).Scan(&size.ID, &size.SizeName, &size.Price, &size.Discount, &size.QtyInStock, &size.CreatedAt, &size.UpdatedAt, &size.DeletedAt); err != nil {
		err = fmt.Errorf("size with size_id '%d' not found", id)
		return utils.Size{}, err
	}

	return size, nil
}

func (s sizeRepository) EditSizeById(id int, size utils.UpdateSize) (utils.Size, error) {
	DB := s.db.GetDB()
	zeroDecimal := decimal.NewFromInt(0)
	var updateFields []string
	var values []interface{}

	var count int
	if err := DB.QueryRow("SELECT COUNT(*) FROM size WHERE size_id = ? AND deleted_at IS NULL", id).Scan(&count); err != nil {
		return utils.Size{}, err
	}

	if count == 0 {
		err := fmt.Errorf("size with size_id '%d' not found", id)
		return utils.Size{}, err
	}

	if size.SizeName != "" {
		updateFields = append(updateFields, "size_name = ?")
		values = append(values, strings.ToUpper(size.SizeName))
	}

	if !size.Price.Equal(zeroDecimal) {
		updateFields = append(updateFields, "price = ?")
		values = append(values, size.Price)
	}

	if !size.Discount.Equal(zeroDecimal) {
		updateFields = append(updateFields, "discount = ?")
		values = append(values, size.Discount)
	}

	if size.QtyInStock != 0 {
		updateFields = append(updateFields, "qty_in_stock = ?")
		values = append(values, size.QtyInStock)
	}

	if len(updateFields) == 0 {
		err := errors.New("failed to update size:No fields provided for update.Please provide at least one field to update")
		return utils.Size{}, err
	}

	if len(values) > 0 {
		updateFields = append(updateFields, "updated_at = ?")
		values = append(values, time.Now())
	}

	query := fmt.Sprintf("UPDATE size SET %s WHERE size_id = ? AND deleted_at IS NULL", strings.Join(updateFields, ", "))
	values = append(values, id)

	if _, err := DB.Exec(query, values...); err != nil {
		return utils.Size{}, err
	}

	updatedSize, err := s.GetSizeById(id)
	if err != nil {
		return utils.Size{}, err
	}

	return updatedSize, nil
}

func (s sizeRepository) DeleteSizeById(id int) (string, int, string, error) {
	DB := s.db.GetDB()

	var count int
	if err := DB.QueryRow("SELECT COUNT(*) FROM size WHERE size_id = ? AND deleted_at IS NULL", id).Scan(&count); err != nil {
		status := http.StatusInternalServerError
		errType := entity.InternalError
		return "", status, errType, err
	}

	if count == 0 {
		status := http.StatusNotFound
		errType := entity.NotFoundError
		err := fmt.Errorf("size with size_id '%d' not found", id)
		return "", status, errType, err
	}

	query := `DELETE FROM size WHERE size_id = ? AND deleted_at IS NULL`
	if _, err := DB.Exec(query, id); err != nil {
		status := http.StatusInternalServerError
		errType := entity.InternalError
		return "", status, errType, err
	}

	status := http.StatusOK
	errType := entity.Success
	resp := fmt.Sprintf("size with size_id '%d' deleted successfully!", id)

	return resp, status, errType, nil
}
