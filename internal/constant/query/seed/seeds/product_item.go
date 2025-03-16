package seeds

import (
	"Eccomerce-website/internal/constant/model/dto"
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

func CreateProductItem(db *pgxpool.Pool, param *dto.ProductItem) error {
	log.Println(param)
	_, err := db.Exec(context.Background(),
		`INSERT INTO product_items(
			id,
			product_id,
			color_id,
			size_id,
			sku,
			status,
			image_url,
			price,
			discount
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		param.ID,
		param.ProductID,
		param.ColorID,
		param.SizeID,
		param.SKU,
		param.Status,
		param.ImageURL,
		param.Price,
		param.Discount,
	)
	return err
}
