package controller

import (
	"context"

	apihttp "github.com/kukymbr/withoutmedianews/internal/api/http"
	"github.com/kukymbr/withoutmedianews/internal/domain"
)

func NewTagsController(service *domain.DictionaryService) *TagsController {
	return &TagsController{
		service: service,
	}
}

type TagsController struct {
	service *domain.DictionaryService
}

func (c *CategoriesController) GetTags(
	ctx context.Context,
	_ apihttp.GetTagsRequestObject,
) (apihttp.GetTagsResponseObject, error) {
	tags, err := c.service.GetTags(ctx)
	if err != nil {
		return nil, err
	}

	return apihttp.GetTags200JSONResponse(apihttp.NewTags(tags)), nil
}
