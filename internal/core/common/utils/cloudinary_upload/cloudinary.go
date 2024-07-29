package cloudinaryupload

import (
	"Eccomerce-website/internal/core/entity"
	"context"
	"mime/multipart"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"go.uber.org/zap"
)

func UploadToCloudinary(file *multipart.FileHeader, dbLogger *zap.Logger, requestID string) (string, error) {
	defer func() {
		os.Remove("./internal/core/common/upload/" + file.Filename)
	}()

	cloudinary_url := "cloudinary://828868553985341:GmQAncUhIfAWjOqQWFyXRHyFCCc@note-zipper"
	cld, err := cloudinary.NewFromURL(cloudinary_url)
	if err != nil {
		errorResponse := entity.AppInternalError.Wrap(err, "failed to get a new cloudinary instance").WithProperty(entity.StatusCode, 500)
		dbLogger.Error("unable to get a new instance by using cloudinary url",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "UploadToCloudinary"),
			zap.String("requestID", requestID),
			zap.String("cloudinaryUrl", cloudinary_url),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return "", errorResponse
	}

	var ctx = context.Background()
	resp, err := cld.Upload.Upload(ctx, "./internal/core/common/upload/"+file.Filename, uploader.UploadParams{PublicID: "my_ecommerce" + "-" + file.Filename})
	if err != nil {
		errorResponse := entity.UnableToSaveFile.Wrap(err, "failed to store image file to cloudinary").WithProperty(entity.StatusCode, 500)
		dbLogger.Error("unable to save image to cloudinary",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "UploadToCloudinary"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return "", errorResponse
	}

	return resp.SecureURL, nil
}

// func generateUid() string {
// 	uniqueId := uuid.New().String()
// 	return uniqueId
// }
