package controllers

import (
	"context"

	apihttp "github.com/kukymbr/withoutmedianews/internal/api/http"
	"github.com/kukymbr/withoutmedianews/internal/news"
)

func NewTagsController(service *news.Dictionaries) *TagsController {
	return &TagsController{
		service: service,
	}
}

type TagsController struct {
	service *news.Dictionaries
}

func (c *CategoriesController) GetTags(
	ctx context.Context,
	_ apihttp.GetTagsRequestObject,
) (apihttp.GetTagsResponseObject, error) {
	tags, err := c.service.GetTags(ctx)
	if err != nil {
		return nil, err
	}

	return apihttp.GetTags200JSONResponse(apihttp.TagsFromDomain(tags)), nil
}
