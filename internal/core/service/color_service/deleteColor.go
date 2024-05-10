package colorservice

import "Eccomerce-website/internal/core/model/response"

func (c colorService) DeleteColor(id int) response.Response {
	resp, status, errType, err := c.colorRepo.DeleteColorById(id)
	if err != nil {
		response := response.Response{
			Status:       status,
			ErrorType:    errType,
			ErrorMessage: err,
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
