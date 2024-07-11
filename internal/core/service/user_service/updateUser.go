package service

import (
	validationdata "Eccomerce-website/internal/core/common/utils/validationData"
	"Eccomerce-website/internal/core/entity"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"fmt"
	"net/http"
)

func (u userService) UpdateUser(id int, request request.UpdateUser) (response.Response, error) {
	if err := request.Validate(); err != nil {
		errorResponse := entity.ValidationError.Wrap(err, "failed updateUser data validation").WithProperty(entity.StatusCode, 400)
		return response.Response{}, errorResponse
	}

	if request.PhoneNumber != "" {
		formatedPhone, err := validationdata.FormatPhoneNumber(request.PhoneNumber)
		if err != nil {
			return response.Response{}, err
		}
		request.PhoneNumber = formatedPhone
	}

	if request.Password != "" {
		hashPassword, err := validationdata.HasPassword(request.Password)
		if err != nil {
			return response.Response{}, err
		}

		request.Password = hashPassword
	}

	updateUser, err := u.userRepo.EditUserById(id, request)
	if err != nil {
		return response.Response{}, err
	}

	response := response.Response{
		Data:       updateUser,
		StatusCode: http.StatusOK,
		Message:    fmt.Sprintf("you have successfully updated the user with userId %d", id),
	}
	return response, nil
}
