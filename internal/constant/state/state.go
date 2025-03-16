package state

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type RedisConfig struct {
	EventType string `json:"event_type,omitempty"`
	RedisURL  string `json:"redis_url,omitempty"`
}

func (r RedisConfig) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.EventType, validation.Required.Error("event_type is required")),
		validation.Field(&r.RedisURL, validation.Required.Error("reddis_url is required")),
	)
}

type UploadParams struct {
	FileTypes []FileType
}

type FileType struct {
	Name            string
	Types           []string
	MaxSize         int64
	AllowCustomName bool
}

func (f *FileType) SetValues(values map[string]any) {
	name, ok := values["name"].(string)
	if ok {
		f.Name = name
	}
	types, ok := values["types"].([]any)
	if ok {
		for _, v := range types {
			typeString, ok := v.(string)
			if ok {
				f.Types = append(f.Types, typeString)
			}
		}
	}

	size, ok := values["max_size"].(int)
	if ok {
		f.MaxSize = int64(size)
	}

	allowCustomName, ok := values["allow_custom_name"].(bool)
	if ok {
		f.AllowCustomName = allowCustomName
	}
}

type HTTPTransport struct {
	MaxIdleConnsPerHost int
	MaxIdleConns        int
	MaxConnsPerHost     int
	Timeout             time.Duration
}

func (a HTTPTransport) Validate() error {
	return validation.Validate([]int{
		a.MaxConnsPerHost, a.MaxIdleConns, a.MaxIdleConnsPerHost, int(a.Timeout),
	}, validation.Each(validation.Required))
}
