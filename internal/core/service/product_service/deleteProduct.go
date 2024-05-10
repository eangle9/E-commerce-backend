package productservice

import "Eccomerce-website/internal/core/model/response"

func (p productService) DeleteProduct(id int) response.Response {
	resp, status, errType, err := p.productRepo.DeleteProductById(id)
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
