package asset

// import (
// 	"Eccomerce-website/initiator/foundation"
// 	"Eccomerce-website/internal/constant/state"
// 	"Eccomerce-website/platform"
// 	"context"
// 	"io"
// 	"mime/multipart"
// 	"os"
// 	"path"
// )

// type fileAsset struct {
// 	log      foundation.Logger
// 	basePath string
// }

// func Init(log foundation.Logger, basePath string) platform.Asset {
// 	return &fileAsset{
// 		log:      log,
// 		basePath: basePath,
// 	}
// }

// func SetParams(log foundation.Logger, params state.UploadParams) state.UploadParams {
// 	for k, v := range params.FileTypes {
// 		if v.Name == "" {
// 			log.Fatal(context.Background(), "name is required for an asset type")
// 		}

// 		if len(v.Types) == 0 {
// 			log.Fatal(context.Background(), "types is required for an asset type")
// 		}

// 		if v.MaxSize == 0 {
// 			params.FileTypes[k].MaxSize = 5 * 1024 * 1024
// 		}
// 	}

// 	return params
// }

// func (f *fileAsset) SaveAsset(_ context.Context, asset multipart.File, dst string) error {
// 	err := os.MkdirAll(path.Join(f.basePath, path.Dir(dst)), 0777)
// 	if err != nil {
// 		return err
// 	}

// 	out, err := os.Create(path.Join(f.basePath, dst))
// 	if err != nil {
// 		return err
// 	}

// 	_, err = io.Copy(out, asset)
// 	if err != nil {
// 		return err
// 	}

// 	err = out.Close()

// 	return err
// }
