package chapa

import (
	"Eccomerce-website/internal/core/dto"
	"Eccomerce-website/internal/core/entity"
	"Eccomerce-website/internal/core/port/repository"
	"Eccomerce-website/platform"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	paymentIntentUrl       = "https://api.chapa.co/v1/transaction/initialize"
	verifyPaymentIntentUrl = "https://api.chapa.co/v1/transaction/verify/%v"
)

type chapa struct {
	apiKey        string
	client        *http.Client
	chapaRepo     repository.ChapaRepository
	serviceLogger *zap.Logger
}

func NewChapa(chapaRepo repository.ChapaRepository, serviceLogger *zap.Logger) platform.API {
	return &chapa{
		apiKey: viper.GetString("API_KEY"),
		client: &http.Client{
			Timeout: 1 * time.Minute,
		},
		chapaRepo:     chapaRepo,
		serviceLogger: serviceLogger,
	}
}

func (c *chapa) InitiatePayment(request *dto.PaymentRequest) (*dto.PaymentResponse, error) {
	var resp dto.PaymentResponse
	fmt.Println("secretToken", c.apiKey)
	if err := request.Validate(); err != nil {
		err = entity.ValidationError.Wrap(err, "required fields are empty").WithProperty(entity.StatusCode, 400)
		return nil, err
	}
	body, err := json.Marshal(request)
	if err != nil {
		err = entity.AppInternalError.Wrap(err, "unable to marshal payment request").WithProperty(entity.StatusCode, 500)
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, paymentIntentUrl, bytes.NewBuffer(body))
	if err != nil {
		err = entity.AppInternalError.Wrap(err, "unable to create request").WithProperty(entity.StatusCode, 500)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Close = true
	res, err := c.client.Do(req)
	if err != nil {
		err = entity.UnableToSave.Wrap(err, "unable to send initiate payment request").WithProperty(entity.StatusCode, 500)
		return nil, err
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			c.serviceLogger.Error("failed to close response body", zap.Error(err))
		}
	}()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		err = entity.UnableToRead.Wrap(err, "failed to read response body").WithProperty(entity.StatusCode, 500)
		return nil, err
	}
	if err := json.Unmarshal(resBody, &resp); err != nil {
		err = entity.UnableToRead.Wrap(err, "failed to unmarshal response body").WithProperty(entity.StatusCode, 500)
		return nil, err
	}
	return &resp, nil
}

func (c *chapa) VerifyPayment(txRef string) (*dto.VerifyResponse, error) {
	return &dto.VerifyResponse{}, nil
}
