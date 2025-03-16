package seed

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type Seed struct {
	Name string
	Run  func(*pgxpool.Pool) error
}
