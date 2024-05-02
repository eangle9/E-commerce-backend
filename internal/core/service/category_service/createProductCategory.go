package categoryservice

import (
	"Eccomerce-website/internal/core/dto"
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/request"
	"Eccomerce-website/internal/core/model/response"
	"Eccomerce-website/internal/core/port/repository"
	"Eccomerce-website/internal/core/port/service"
	"net/http"
)

type productCategoryService struct {
	categoryRepo repository.ProductCategoryRepository
}

func NewProductCategoryRepository(repo repository.ProductCategoryRepository) service.ProductCategoryService {
	return &productCategoryService{
		categoryRepo: repo,
	}
}

func (p productCategoryService) CreateProductCategory(request request.ProductCategoryRequest) response.Response {
	parentId := request.ParentID
	name := request.Name

	category := dto.ProductCategory{
		ParentID: parentId,
		Name:     name,
	}

	id, err := p.categoryRepo.InsertProductCategory(category)
	if err != nil {
		response := response.Response{
			Status:       http.StatusConflict,
			ErrorType:    "DUPLICATE_ENTRY",
			ErrorMessage: err.Error(),
		}
		return response
	}

	category.ID = *id

	response := response.Response{
		Data:         category,
		Status:       http.StatusCreated,
		ErrorType:    errorcode.Success,
		ErrorMessage: "Congratulation, product category created successfully!",
	}

	return response
}
