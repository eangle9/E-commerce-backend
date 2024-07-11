package service

import (
	"Eccomerce-website/internal/core/entity"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"net/http"
)

func (u userService) GetUsers(request request.PaginationQuery) (response.Response, error) {
	page := request.Page
	perPage := request.PerPage

	if page == 0 {
		page = 1
	}

	if perPage == 0 {
		perPage = 5
	}

	if err := request.Validate(); err != nil {
		errorResponse := entity.ValidationError.Wrap(err, "failed pagination query validation").WithProperty(entity.StatusCode, 400)
		return response.Response{}, errorResponse
	}

	offset := (page - 1) * perPage

	users, err := u.userRepo.ListUsers(offset, perPage)
	if err != nil {
		return response.Response{}, err
	}

	response := response.Response{
		Data:       users,
		StatusCode: http.StatusOK,
		Message:    "you have get list of users successfully!",
	}

	return response, nil

}
