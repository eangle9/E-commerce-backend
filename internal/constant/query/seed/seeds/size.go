package seeds

import (
	"Eccomerce-website/internal/constant/model/dto"
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

func CreateSize(db *pgxpool.Pool, param *dto.Size) error {
	log.Println(param)
	_, err := db.Exec(context.Background(),
		`INSERT INTO size(
			id,
			name
			) VALUES ($1, $2)`,
		param.ID,
		param.Name,
	)
	return err
}
