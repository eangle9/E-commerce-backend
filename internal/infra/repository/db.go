package repository

import (
	"Eccomerce-website/internal/core/port/repository"
	"Eccomerce-website/internal/infra/config"
	"database/sql"
)

type database struct {
	DB *sql.DB
}

func NewDatabase(conf config.DatabaseConfig) (repository.Database, error) {
	db, err := newDatabase(conf)
	if err != nil {
		return nil, err
	}

	return &database{
		DB: db,
	}, nil
}

func newDatabase(conf config.DatabaseConfig) (*sql.DB, error) {
	// dsn := conf.Url
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db, err := sql.Open(conf.Driver, conf.Url)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (da database) GetDB() *sql.DB {
	db := da.DB
	return db
}
