package categoryservice

import "Eccomerce-website/internal/core/model/response"

func (p productCategoryService) DeleteProductCategory(id int) response.Response {
	resp, status, errType, err := p.categoryRepo.DeleteProductCategoryById(id)
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
