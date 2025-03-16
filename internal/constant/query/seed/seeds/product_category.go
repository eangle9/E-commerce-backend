package seeds

import (
	"Eccomerce-website/internal/constant/model/dto"
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

func CreateProductCategory(db *pgxpool.Pool, param *dto.ProductCategory) error {
	log.Println(param)
	_, err := db.Exec(context.Background(),
		`INSERT INTO product_category(
			id,
			name,
			description,
			image
			) VALUES ($1, $2, $3, $4)`,
		param.ID,
		param.Name,
		param.Description,
		param.Image,
	)
	return err
}
