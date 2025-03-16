package foundation

import (
	"Eccomerce-website/internal/constant/errors"
	"Eccomerce-website/internal/constant/state"

	// "Eccomerce-website/platform/asset"
	"context"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type State struct {
	// AuthDomains      state.AuthDomains
	HTTPConfig state.HTTPTransport
	// TokenConfig      state.TokenKey
	// UploadParams     state.UploadParams
	RedisConfig state.RedisConfig
}

func InitState(logger Logger) State {
	httpconfig := state.HTTPTransport{
		MaxIdleConnsPerHost: viper.GetInt("http.max_idle_conns_per_host"),
		MaxIdleConns:        viper.GetInt("http.max_idle_conns"),
		MaxConnsPerHost:     viper.GetInt("http.max_conns_per_host"),
		Timeout:             viper.GetDuration("http.timeout"),
	}
	if err := httpconfig.Validate(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid http configuration")
		logger.Fatal(context.Background(), "all http fields are required", zap.Error(err))
	}

	redisConfig := state.RedisConfig{
		EventType: viper.GetString("redis.subscriber"),
		RedisURL:  viper.GetString("redis.redis_url"),
	}
	if err := redisConfig.Validate(); err != nil {
		logger.Fatal(context.Background(), "missing some redis configs",
			zap.Error(err),
		)
	}

	// assets := GetMapSlice("assets")
	// fileTypes := make([]state.FileType, 0, len(assets))

	// for _, v := range assets {
	// 	var fileType state.FileType

	// 	fileType.SetValues(v)
	// 	fileTypes = append(fileTypes, fileType)
	// }
	return State{
		HTTPConfig:  httpconfig,
		RedisConfig: redisConfig,
		// UploadParams: asset.SetParams(logger, state.UploadParams{
		// 	FileTypes: fileTypes,
		// }),
	}
}

func GetMapSlice(path string) []map[string]any {
	value := viper.Get(path)
	mapInterfaceSlice, ok := value.([]any)
	if !ok {
		return nil
	}

	var mapStringAny []map[string]any
	for _, v := range mapInterfaceSlice {
		v, ok := v.(map[string]any)
		if ok {
			mapStringAny = append(mapStringAny, v)
		}
	}

	return mapStringAny
}
