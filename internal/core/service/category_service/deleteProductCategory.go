package categoryservice

import "Eccomerce-website/internal/core/model/response"

func (p productCategoryService) DeleteProductCategory(id int) response.Response {
	resp, status, _, err := p.categoryRepo.DeleteProductCategoryById(id)
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
