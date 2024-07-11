package sizeservice

import "Eccomerce-website/internal/core/model/response"

func (s sizeService) DeleteSize(id int) response.Response {
	resp, status, _, err := s.sizeRepo.DeleteSizeById(id)
	if err != nil {
		response := response.Response{
			StatusCode: status,
			Message:    err.Error(),
		}
		return response
	}

	response := response.Response{
		StatusCode: status,
		Message:    resp,
	}

	return response
}
