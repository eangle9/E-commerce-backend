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

func (p productsRepository) ListAllProducts() ([]utils.ListProduct, error) {
	DB := p.db.GetDB()
	var products []utils.ListProduct

	productRows, err := DB.Query("SELECT product_id, product_name FROM product")
	if err != nil {
		return nil, err
	}

	defer productRows.Close()

	for productRows.Next() {
		var singleProduct utils.ListProduct

		if err := productRows.Scan(&singleProduct.ProductID, &singleProduct.Name); err != nil {
			return nil, err
		}

		query := `
		SELECT 
		    p.product_item_id, 
			c.color_name, 
			s.size_name,
			p.image_url, 
			p.price, 
			p.discount,
			p.qty_in_stock
		FROM 
		    product_item p 
		LEFT JOIN 
		    color c ON p.color_id = c.color_id
		LEFT JOIN
		    size s ON p.size_id = s.size_id	
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
			if err := productItemRows.Scan(&productItem.ItemID, &productItem.Color, &productItem.Size, &productItem.ImageUrl, &productItem.Price, &productItem.Discount, &productItem.InStock); err != nil {
				return nil, err
			}
			productItems = append(productItems, productItem)
		}

		if err := productItemRows.Err(); err != nil {
			return nil, err
		}

		singleProduct.ProductItems = productItems

		reviewQuery := `
    SELECT 
        r.review_id,
        r.user_id,
        r.product_id,
        r.rating,
        r.comment,
        r.created_at,
        u.user_id,
        u.username,
        u.first_name,
        u.last_name,
        u.email,
        u.phone_number
    FROM 
        review r
    LEFT JOIN
        users u ON r.user_id = u.user_id
    WHERE
        r.product_id = ?
`

		reviewRows, err := DB.Query(reviewQuery, singleProduct.ProductID)
		if err != nil {
			return nil, err
		}

		defer reviewRows.Close()

		var reviews []utils.ProductReview
		for reviewRows.Next() {
			var review utils.ProductReview
			if err := reviewRows.Scan(&review.ReviewID, &review.UserID, &review.ProductID, &review.Rating, &review.Comment, &review.CreatedAt, &review.User.ID, &review.User.Username, &review.User.FirstName, &review.User.LastName, &review.User.Email, &review.User.PhoneNumber); err != nil {
				return nil, err
			}
			reviews = append(reviews, review)
		}

		if err := reviewRows.Err(); err != nil {
			return nil, err
		}

		singleProduct.Reviews = reviews

		products = append(products, singleProduct)
	}

	return products, nil
}

func (p productsRepository) GetSingleProductById(id int) (utils.SingleProduct, error) {
	DB := p.db.GetDB()
	var singleProduct utils.SingleProduct

	prouductQuery := `
	SELECT
	   p.product_id,
	   c.name,
	   p.brand,
	   p.product_name,
	   p.description
	FROM 
	  product p
	LEFT JOIN
	  product_category c ON p.category_id = c.category_id
	WHERE
	  p.product_id = ?    
	`

	if err := DB.QueryRow(prouductQuery, id).Scan(&singleProduct.ProductID, &singleProduct.Category, &singleProduct.Brand, &singleProduct.Name, &singleProduct.Description); err != nil {
		return utils.SingleProduct{}, err
	}

	itemQuery := `
	SELECT 
	   p.product_item_id,
	   c.color_name,
	   s.size_name,
	   p.image_url,
	   p.price,
	   p.discount,
	   qty_in_stock
	FROM 
	  product_item p
	LEFT JOIN
	  color c ON p.color_id = c.color_id
	LEFT JOIN
	  size s ON p.size_id = s.size_id
	WHERE
	  p.product_id = ?    
	`

	itemRows, err := DB.Query(itemQuery, id)
	if err != nil {
		return utils.SingleProduct{}, err
	}

	defer itemRows.Close()

	var items []utils.ItemVariant
	for itemRows.Next() {
		var item utils.ItemVariant

		if err := itemRows.Scan(&item.ItemID, &item.Color, &item.Size, &item.ImageUrl, &item.Price, &item.Discount, &item.InStock); err != nil {
			return utils.SingleProduct{}, err
		}
		items = append(items, item)
	}

	if err := itemRows.Err(); err != nil {
		return utils.SingleProduct{}, err
	}

	singleProduct.Items = items

	reviewQuery := `
	SELECT
	   r.review_id,
	   r.user_id,
	   r.product_id,
	   r.rating,
	   r.comment,
	   r.created_at,
	   u.user_id,
	   u.username,
	   u.first_name,
	   u.last_name,
	   u.email,
	   u.phone_number
	FROM 
	  review r
	LEFT JOIN 
	  users u ON r.user_id = u.user_id
	WHERE
	  r.product_id = ?    
	`
	reviewRows, err := DB.Query(reviewQuery, id)
	if err != nil {
		return utils.SingleProduct{}, err
	}

	defer reviewRows.Close()

	var reviews []utils.ProductReview
	for reviewRows.Next() {
		var review utils.ProductReview
		if err := reviewRows.Scan(&review.ReviewID, &review.UserID, &review.ProductID, &review.Rating, &review.Comment, &review.CreatedAt, &review.User.ID, &review.User.Username, &review.User.FirstName, &review.User.LastName, &review.User.Email, &review.User.PhoneNumber); err != nil {
			return utils.SingleProduct{}, err
		}
		reviews = append(reviews, review)
	}

	if err := reviewRows.Err(); err != nil {
		return utils.SingleProduct{}, err
	}

	singleProduct.Reviews = reviews

	return singleProduct, nil
}
