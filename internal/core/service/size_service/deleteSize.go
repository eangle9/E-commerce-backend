package sizeservice

import "Eccomerce-website/internal/core/model/response"

func (s sizeService) DeleteSize(id int) response.Response {
	resp, status, errType, err := s.sizeRepo.DeleteSizeById(id)
	if err != nil {
		response := response.Response{
			Status:       status,
			ErrorType:    errType,
			ErrorMessage: err.Error(),
		}
		return response
	}

	response := response.Response{
		Status:       status,
		ErrorType:    errType,
		ErrorMessage: resp,
	}

	return response
}
