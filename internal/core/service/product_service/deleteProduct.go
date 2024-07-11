package productservice

import "Eccomerce-website/internal/core/model/response"

func (p productService) DeleteProduct(id int) response.Response {
	resp, status, _, err := p.productRepo.DeleteProductById(id)
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
