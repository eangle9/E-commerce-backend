package validationdata

import (
	"Eccomerce-website/internal/core/entity"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
)

func ImageFileValidation(file *multipart.FileHeader) error {
	maxUploadSize := 8 * 1024 * 1024
	fileSize := file.Size

	if fileSize > int64(maxUploadSize) {
		err := errors.New("the uploaded product image is too large.Please upload a size less than 8MB")
		errorResponse := entity.FileTooLarge.Wrap(err, "file too large").WithProperty(entity.StatusCode, 413)
		return errorResponse
	}

	validExt := map[string]bool{
		".jpeg": true,
		".png":  true,
		".jpg":  true,
		".gif":  true,
		".webp": true,
		".svg":  true,
	}

	ext := filepath.Ext(file.Filename)
	if !validExt[ext] {
		err := fmt.Errorf("image with extension '%s' is invalid", ext)
		errorResponse := entity.InvalidExtension.Wrap(err, "invalid file extension").WithProperty(entity.StatusCode, 415)
		return errorResponse
	}

	return nil
}
