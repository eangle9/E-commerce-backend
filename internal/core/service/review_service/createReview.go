package reviewservice

import (
	"Eccomerce-website/internal/core/dto"
	"Eccomerce-website/internal/core/entity"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"Eccomerce-website/internal/core/port/repository"
	"Eccomerce-website/internal/core/port/service"
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type reviewService struct {
	reviewRepo    repository.ReviewRepository
	serviceLogger *zap.Logger
}

func NewReviewService(reviewRepo repository.ReviewRepository, serviceLogger *zap.Logger) service.ReviewService {
	return &reviewService{
		reviewRepo:    reviewRepo,
		serviceLogger: serviceLogger,
	}
}

func (r reviewService) CreateReview(ctx context.Context, request request.ReviewRequest, requestID string) (response.Response, error) {
	if err := request.Validate(); err != nil {
		errorResponse := entity.ValidationError.Wrap(err, "createReview validation error").WithProperty(entity.StatusCode, 400)
		r.serviceLogger.Error("validation error",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "serviceLayer"),
			zap.String("function", "CreateReview"),
			zap.String("requestID", requestID),
			zap.Any("requestData", request),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return response.Response{}, errorResponse
	}

	review := dto.Review{
		UserID:    request.UserID,
		ProductID: request.ProductID,
		Rating:    request.Rating,
		Comment:   request.Comment,
	}

	id, err := r.reviewRepo.InsertReview(ctx, review, requestID)
	if err != nil {
		return response.Response{}, err
	}

	review.ID = *id

	response := response.Response{
		Data:       review,
		StatusCode: http.StatusCreated,
		Message:    "Congratulation, review created successfully",
	}

	return response, nil
}
