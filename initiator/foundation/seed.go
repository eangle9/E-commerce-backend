package foundation

import (
	"context"
	"fmt"

	"Eccomerce-website/internal/constant/query/seed/seeds"

	"github.com/jackc/pgx/v4/pgxpool"
)

func SeedDB(conn *pgxpool.Pool) error {

	// Seed tables if they are not seeded. And create a seed record.
	for _, seed := range seeds.All() {
		// check if table is seeded
		if !IsSeeded(conn, seed.Name) {
			// seed table
			err := seed.Run(conn)
			if err != nil {
				return fmt.Errorf(err.Error(), seed.Name)
			}
			// create seed record
			err = CreateSeed(conn, seed.Name)
			if err != nil {
				return fmt.Errorf(err.Error(), seed.Name)
			}
		}
	}
	return nil
}

func CreateSeed(conn *pgxpool.Pool, name string) error {
	_, err := conn.Exec(context.Background(), "INSERT INTO seeds (name) VALUES ($1)", name)
	return err
}

func IsSeeded(conn *pgxpool.Pool, name string) bool {
	var id int64
	err := conn.QueryRow(context.Background(), "SELECT * FROM seeds WHERE name = $1", name).Scan(&name, &id)
	if err != nil {
		return false
	}
	if name == "" {
		return false
	}
	return true
}
