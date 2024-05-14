package cloudinaryupload

import (
	"context"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
)

func UploadToCloudinary(file *multipart.FileHeader) (string, error) {
	defer func() {
		os.Remove("./internal/core/common/upload/" + file.Filename)
	}()

	cloudinary_url := "cloudinary://828868553985341:GmQAncUhIfAWjOqQWFyXRHyFCCc@note-zipper"
	cld, err := cloudinary.NewFromURL(cloudinary_url)
	if err != nil {
		return "", err
	}

	var ctx = context.Background()
	resp, err := cld.Upload.Upload(ctx, "./internal/core/common/upload/"+file.Filename, uploader.UploadParams{PublicID: "my_ecommerce" + "-" + file.Filename + generateUid()})
	if err != nil {
		return "", err
	}

	return resp.SecureURL, nil
}

func generateUid() string {
	uniqueId := uuid.New().String()
	return uniqueId
}
