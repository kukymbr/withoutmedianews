package controller

import (
	"context"

	apihttp "github.com/kukymbr/withoutmedianews/internal/api/http"
	"github.com/kukymbr/withoutmedianews/internal/db"
)

func NewTagsController(repo *db.WithoutmedianewsRepo) *TagsController {
	return &TagsController{
		repo: repo,
	}
}

type TagsController struct {
	repo *db.WithoutmedianewsRepo
}

func (c *CategoriesController) GetTags(
	ctx context.Context,
	_ apihttp.GetTagsRequestObject,
) (apihttp.GetTagsResponseObject, error) {
	tags, err := c.repo.TagsByFilters(ctx, nil, db.PagerNoLimit, db.EnabledOnly())
	if err != nil {
		return nil, err
	}

	return apihttp.GetTags200JSONResponse(apihttp.NewTags(tags)), nil
}
