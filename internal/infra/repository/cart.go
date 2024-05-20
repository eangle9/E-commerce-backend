package repository

import (
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"Eccomerce-website/internal/core/port/repository"
	"database/sql"

	"github.com/shopspring/decimal"
)

type cartRepository struct {
	db repository.Database
}

func NewCartRepository(db repository.Database) repository.CartRepository {
	return &cartRepository{
		db: db,
	}
}

func (c cartRepository) InsertCartItem(request request.CartRequest, userId uint) ([]response.CartResponse, error) {
	var cartResponses []response.CartResponse
	DB := c.db.GetDB()
	// userId := request.UserID
	productItemId := request.ProductItemID
	quantity := request.Quantity

	var unitPrice decimal.Decimal
	query := `SELECT price FROM product_item WHERE product_item_id = ? AND deleted_at IS NULL`
	if err := DB.QueryRow(query, productItemId).Scan(&unitPrice); err != nil {
		return nil, err
	}

	var cartID int
	if err := DB.QueryRow("SELECT cart_id FROM shopping_cart WHERE user_id = ? AND deleted_at IS NULL", userId).Scan(&cartID); err != nil {
		if err == sql.ErrNoRows {
			result, err := DB.Exec("INSERT INTO shopping_cart (user_id) VALUES (?)", userId)
			if err != nil {
				return nil, err
			}

			cart64, err := result.LastInsertId()
			if err != nil {
				return nil, err
			}
			cartID = int(cart64)
		} else {
			return nil, err
		}
	}

	var count int
	if err := DB.QueryRow("SELECT COUNT(*) FROM cart_item WHERE product_item_id = ? AND cart_id = ?", productItemId, cartID).Scan(&count); err != nil {
		return nil, err
	}

	if count > 0 {
		var oldQuantity int
		if err := DB.QueryRow("SELECT quantity FROM cart_item WHERE product_item_id = ? AND cart_id = ?", productItemId, cartID).Scan(&oldQuantity); err != nil {
			return nil, err
		}
		newQuantity := quantity + oldQuantity
		deciQuantity := decimal.NewFromInt(int64(newQuantity))
		total := deciQuantity.Mul(unitPrice)

		updateQuery := `UPDATE cart_item SET quantity = ?, total = ?, updated_at = NOW() WHERE product_item_id = ? AND cart_id = ?`
		if _, err := DB.Exec(updateQuery, newQuantity, total, productItemId, cartID); err != nil {
			return nil, err
		}
	} else {
		deciQuantity := decimal.NewFromInt(int64(quantity))
		total := deciQuantity.Mul(unitPrice)

		if _, err := DB.Exec("INSERT INTO cart_item (cart_id, product_item_id, quantity, total) VALUES (?, ?, ?, ?)", cartID, productItemId, quantity, total); err != nil {
			return nil, err
		}
	}

	if err := UpdateShoppingCart(DB, cartID); err != nil {
		return nil, err
	}

	joinQuery := `
	SELECT 
	    pimage.image_url,
	    p.product_name, 
	    p.description, 
	    pi.price AS unit_price,
		pi.qty_in_stock, 
		ci.item_id,
	    ci.quantity, 
	    ci.total AS sub_total, 
	    sc.sub_total AS cart_sub_total, 
        sc.total
    FROM 
	   users u
    JOIN 
	  shopping_cart sc ON u.user_id = sc.user_id
    JOIN 
	  cart_item ci ON sc.cart_id = ci.cart_id 
    JOIN 
	  product_item pi ON ci.product_item_id = pi.product_item_id
    JOIN 
	  product_image pimage ON pi.product_item_id = pimage.product_item_id
    JOIN 
	  product p ON pi.product_id = p.product_id
    WHERE 
	  u.user_id = ?
	`

	rows, err := DB.Query(joinQuery, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var cartResponse response.CartResponse
		if err := rows.Scan(
			&cartResponse.ImageUrl,
			&cartResponse.ProductName,
			&cartResponse.Description,
			&cartResponse.UnitPrice,
			&cartResponse.QtyInStock,
			&cartResponse.CartItemID,
			&cartResponse.Quantity,
			&cartResponse.SubTotal,
			&cartResponse.CartSubTotal,
			&cartResponse.Total,
		); err != nil {
			return nil, err
		}
		cartResponses = append(cartResponses, cartResponse)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cartResponses, nil

}

func UpdateShoppingCart(DB *sql.DB, cartID int) error {
	// var totalOld decimal.Decimal
	// var subTotalOld decimal.Decimal
	// var shipping decimal.Decimal
	// var serviceCharge decimal.Decimal
	// var vat decimal.Decimal
	// if err := DB.QueryRow("SELECT sub_total, total, shipping, service_charge, vat FROM shopping_cart WHERE cart_id = ? AND deleted_at IS NULL", cartID).Scan(&subTotalOld, &totalOld, &shipping, &serviceCharge, &vat); err != nil {
	// 	return err
	// }

	var subTotal decimal.Decimal
	if err := DB.QueryRow("SELECT SUM(total) FROM cart_item WHERE cart_id = ? AND deleted_at IS NULL", cartID).Scan(&subTotal); err != nil {
		return err
	}

	// subTotalNew := subTotalOld.Add(subTotal)
	// totalNew := totalOld.Add(shipping).Add(serviceCharge).Add(vat)
	if _, err := DB.Exec("UPDATE shopping_cart SET sub_total = ?, total = sub_total + shipping + service_charge + vat, updated_at = NOW() WHERE cart_id = ?", subTotal, cartID); err != nil {
		return err
	}

	return nil
}
