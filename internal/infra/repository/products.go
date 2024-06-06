package repository

import (
	"Eccomerce-website/internal/core/common/utils"
	"Eccomerce-website/internal/core/port/repository"
)

type productsRepository struct {
	db repository.Database
}

func NewProductsRepository(db repository.Database) repository.GetProducts {
	return &productsRepository{
		db: db,
	}
}

func (p productsRepository) ListAllProducts() ([]utils.SingleProduct, error) {
	DB := p.db.GetDB()
	var products []utils.SingleProduct

	productRows, err := DB.Query("SELECT product_id, product_name FROM product")
	if err != nil {
		return nil, err
	}

	defer productRows.Close()

	for productRows.Next() {
		var singleProduct utils.SingleProduct

		if err := productRows.Scan(&singleProduct.ProductID, &singleProduct.Product); err != nil {
			return nil, err
		}

		query := `
		SELECT 
		    p.product_item_id, 
			c.color_name, 
			p.image_url, 
			p.price, 
			p.qty_in_stock 
		FROM 
		    product_item p 
		LEFT JOIN 
		    color c ON p.color_id = c.color_id
		WHERE 
		    p.product_id = ?
		 `
		productItemRows, err := DB.Query(query, singleProduct.ProductID)
		if err != nil {
			return nil, err
		}

		defer productItemRows.Close()

		var productItems []utils.ProductVariant
		for productItemRows.Next() {
			var productItem utils.ProductVariant
			if err := productItemRows.Scan(&productItem.ItemID, &productItem.Color, &productItem.ImageUrl, &productItem.Price, &productItem.InStock); err != nil {
				return nil, err
			}
			productItems = append(productItems, productItem)
		}

		singleProduct.ProductItems = productItems
		products = append(products, singleProduct)
	}

	return products, nil
}
