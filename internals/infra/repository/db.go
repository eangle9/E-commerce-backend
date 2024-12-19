package repository

import (
	"Eccomerce-website/internals/core/entity"
	"Eccomerce-website/internals/core/port/repository"
	"Eccomerce-website/internals/infra/config"
	"database/sql"
	"time"

	"go.uber.org/zap"
)

type database struct {
	DB *sql.DB
	// dbLogger *zap.Logger
}

func NewDatabase(conf config.DatabaseConfig, dbLogger *zap.Logger) (repository.Database, error) {
	db, err := newDatabase(conf, dbLogger)
	if err != nil {
		return nil, err
	}

	return &database{
		DB: db,
	}, nil
}

func newDatabase(conf config.DatabaseConfig, dbLogger *zap.Logger) (*sql.DB, error) {
	db, err := sql.Open(conf.Driver, conf.Url)
	if err != nil {
		errorResponse := entity.ConnectionError.Wrap(err, "failed db connection").WithProperty(entity.StatusCode, 500)
		dbLogger.Error("unable to connect database",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "newDatabse"),
			zap.String("error", errorResponse.Error()),
			zap.Stack("stacktrace"),
		)
		return nil, errorResponse
	}

	return db, nil
}

func (da database) GetDB() *sql.DB {
	db := da.DB
	return db
}
