package seeds

import (
	"Eccomerce-website/internal/constant/model/dto"
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

func CreateProduct(db *pgxpool.Pool, param *dto.Product) error {
	log.Println(param)
	_, err := db.Exec(context.Background(),
		`INSERT INTO products(
			id,
			category_id,
			brand,
			name,
			status,
			description,
			average_rating,
			total_reviews
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		param.ID,
		param.CategoryID,
		param.Brand,
		param.Name,
		param.Status,
		param.Description,
		param.AverageRating,
		param.TotalReviews,
	)
	return err
}
