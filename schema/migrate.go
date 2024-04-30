package schema

import (
	"Eccomerce-website/internal/core/port/repository"
	"fmt"
	"os"
)

func Migrate(db repository.Database) error {
	DB := db.GetDB()
	schema, err := os.ReadFile("../schema/schema.sql")
	if err != nil {
		err = fmt.Errorf("error reading schema file: %s", err)
		return err
	}

	if _, err := DB.Exec(string(schema)); err != nil {
		// err = errors.New("error executing schema SQL")
		return err
	}

	fmt.Println("schema created successfully")

	return nil
}
