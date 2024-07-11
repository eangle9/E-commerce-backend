package reviewservice

import (
	"Eccomerce-website/internal/core/dto"
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
			StatusCode: http.StatusConflict,
			Message:    err.Error(),
		}
		return response
	}

	review.ID = *id

	response := response.Response{
		Data:       review,
		StatusCode: http.StatusCreated,
		Message:    "Congratulation, review created successfully",
	}

	return response
}
