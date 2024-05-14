package productitemservice

import "Eccomerce-website/internal/core/model/response"

func (p productItemService) DeleteProductItem(id int) response.Response {
	resp, status, errType, err := p.itemRepo.DeleteProductItemById(id)
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
