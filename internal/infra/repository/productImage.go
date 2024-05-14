package repository

import (
	cloudinaryupload "Eccomerce-website/internal/core/common/utils/cloudinary_upload"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/port/repository"
)

type productImageRepository struct {
	db repository.Database
}

func NewProductImageRepository(db repository.Database) repository.ProductImageRepository {
	return &productImageRepository{
		db: db,
	}
}

func (p productImageRepository) InsertProductImage(request request.ProductImageRequest) (*int, string, error) {
	DB := p.db.GetDB()
	file := request.File
	productItemId := request.ProductItemId

	imageUrl, err := cloudinaryupload.UploadToCloudinary(file)
	if err != nil {
		return nil, "", err
	}

	query := `INSERT INTO product_image(product_item_id, image_url) VALUES(?, ?)`
	result, err := DB.Exec(query, productItemId, imageUrl)
	if err != nil {
		return nil, "", err
	}

	id64, err := result.LastInsertId()
	if err != nil {
		return nil, "", err
	}

	id := int(id64)

	return &id, imageUrl, nil

}
