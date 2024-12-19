package repository

import (
	cloudinaryupload "Eccomerce-website/internals/core/common/utils/cloudinary_upload"
	"Eccomerce-website/internals/core/entity"
	"Eccomerce-website/internals/core/model/request"
	"Eccomerce-website/internals/core/port/repository"
	"context"
	"time"

	"go.uber.org/zap"
)

type productImageRepository struct {
	db       repository.Database
	dbLogger *zap.Logger
}

func NewProductImageRepository(db repository.Database, dbLogger *zap.Logger) repository.ProductImageRepository {
	return &productImageRepository{
		db:       db,
		dbLogger: dbLogger,
	}
}

func (p productImageRepository) InsertProductImage(ctx context.Context, request request.ProductImageRequest, requestID string) (*int, string, error) {
	DB := p.db.GetDB()
	file := request.File
	productItemId := request.ProductItemId

	imageUrl, err := cloudinaryupload.UploadToCloudinary(file, p.dbLogger, requestID)
	if err != nil {
		return nil, "", err
	}

	query := `INSERT INTO product_image(product_item_id, image_url) VALUES(?, ?)`
	result, err := DB.ExecContext(ctx, query, productItemId, imageUrl)
	if err != nil {
		errorResponse := entity.UnableToSave.Wrap(err, "failed to insert product image").WithProperty(entity.StatusCode, 500)
		p.dbLogger.Error("failed to create product image",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "InsertProductImage"),
			zap.String("requestID", requestID),
			zap.String("query", query),
			zap.Any("requestData", request),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, "", errorResponse
	}

	id64, err := result.LastInsertId()
	if err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "failed to get the inserted id").WithProperty(entity.StatusCode, 500)
		p.dbLogger.Error("unable to get lastInserted id",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "InsertProductImage"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, "", err
	}

	id := int(id64)

	return &id, imageUrl, nil

}
