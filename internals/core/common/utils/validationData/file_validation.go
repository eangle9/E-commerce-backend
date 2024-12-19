package validationdata

import (
	"Eccomerce-website/internals/core/entity"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

func ImageFileValidation(file *multipart.FileHeader, handlerLogger *zap.Logger, requestID string) error {
	maxUploadSize := 8 * 1024 * 1024
	fileSize := file.Size

	if fileSize > int64(maxUploadSize) {
		err := errors.New("the uploaded product image is too large.Please upload a size less than 8MB")
		errorResponse := entity.FileTooLarge.Wrap(err, "file too large").WithProperty(entity.StatusCode, 413)
		handlerLogger.Error("file too big",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "ImageFileValidation"),
			zap.String("requestID", requestID),
			zap.Int64("fileSize", fileSize),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
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
		handlerLogger.Error("incorrect image extension",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "handlerLayer"),
			zap.String("function", "ImageFileValidation"),
			zap.String("requestID", requestID),
			zap.String("extension", ext),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return errorResponse
	}

	return nil
}
