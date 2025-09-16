package controller

import (
	"context"

	apihttp "github.com/kukymbr/withoutmedianews/internal/api/http"
	"github.com/kukymbr/withoutmedianews/internal/domain"
)

func NewCategoriesController(service *domain.DictionaryService) *CategoriesController {
	return &CategoriesController{
		service: service,
	}
}

type CategoriesController struct {
	service *domain.DictionaryService
}

func (c *CategoriesController) GetCategories(
	ctx context.Context,
	_ apihttp.GetCategoriesRequestObject,
) (apihttp.GetCategoriesResponseObject, error) {
	categories, err := c.service.GetCategories(ctx)
	if err != nil {
		return nil, err
	}

	resp := make(apihttp.GetCategories200JSONResponse, 0, len(categories))

	for _, category := range categories {
		resp = append(resp, apihttp.Category{
			ID:    category.ID,
			Title: category.Title,
		})
	}

	return resp, nil
}
