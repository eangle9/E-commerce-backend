package seeds

import (
	"Eccomerce-website/internal/constant/model/dto"
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

func CreateUser(db *pgxpool.Pool, param *dto.User) error {
	log.Println(param)
	_, err := db.Exec(context.Background(),
		`INSERT INTO users(
			id,
			username,
			email,
			phone_number,
			password, 
			first_name,
			last_name,
			profile_picture,
			email_verified,
			role
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		param.ID,
		param.Username,
		param.Email,
		param.PhoneNumber,
		param.Password,
		param.FirstName,
		param.LastName,
		param.ProfilePicture,
		param.EmailVerified,
		param.Role,
	)
	return err
}
