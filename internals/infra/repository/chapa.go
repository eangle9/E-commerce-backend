package repository

import (
	"Eccomerce-website/internals/core/port/repository"

	"go.uber.org/zap"
)

type chapaRepository struct {
	db       repository.Database
	dbLogger *zap.Logger
}

func NewChapaRepository(db repository.Database, dbLogger *zap.Logger) repository.ChapaRepository {
	return &chapaRepository{
		db:       db,
		dbLogger: dbLogger,
	}
}
