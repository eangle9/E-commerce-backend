package reviewservice

import (
	"Eccomerce-website/internal/core/dto"
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"Eccomerce-website/internal/core/port/repository"
	"Eccomerce-website/internal/core/port/service"
	"net/http"
)

type reviewService struct {
	reviewRepo repository.ReviewRepository
}

func NewReviewService(reviewRepo repository.ReviewRepository) service.ReviewService {
	return &reviewService{
		reviewRepo: reviewRepo,
	}
}

func (r reviewService) CreateReview(request request.ReviewRequest) response.Response {
	review := dto.Review{
		UserID:    request.UserID,
		ProductID: request.ProductID,
		Rating:    request.Rating,
		Comment:   request.Comment,
	}

	id, err := r.reviewRepo.InsertReview(review)
	if err != nil {
		response := response.Response{
			Status:       http.StatusConflict,
			ErrorType:    "DUPLICATE_ENTRY",
			ErrorMessage: err.Error(),
		}
		return response
	}

	review.ID = *id

	response := response.Response{
		Data:         review,
		Status:       http.StatusCreated,
		ErrorType:    errorcode.Success,
		ErrorMessage: "Congratulation, review created successfully",
	}

	return response
}
