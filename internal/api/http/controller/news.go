package controller

import (
	"context"

	"github.com/kukymbr/withoutmedianews/internal/api/http"
	"github.com/kukymbr/withoutmedianews/internal/domain"
)

func NewNewsController(service *domain.Service) *NewsController {
	return &NewsController{
		service: service,
	}
}

type NewsController struct {
	service *domain.Service
}

func (c *NewsController) GetNewses(
	ctx context.Context,
	req apihttp.GetNewsesRequestObject,
) (apihttp.GetNewsesResponseObject, error) {
	items, err := c.service.GetList(
		ctx,
		domain.NewNewsesFilter(req.Params.CategoryID, req.Params.TagID),
		req.Params.Page,
		req.Params.PerPage,
	)
	if err != nil {
		return nil, err
	}

	resp := make(apihttp.GetNewses200JSONResponse, 0, len(items))

	for _, item := range items {
		resp = append(resp, newNewsListable(item))
	}

	return resp, nil
}

func (c *NewsController) GetNews(
	ctx context.Context,
	request apihttp.GetNewsRequestObject,
) (apihttp.GetNewsResponseObject, error) {
	item, err := c.service.GetNews(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	return apihttp.GetNews200JSONResponse(newNews(item)), nil
}

func (c *NewsController) GetNewsCount(
	ctx context.Context,
	req apihttp.GetNewsCountRequestObject,
) (apihttp.GetNewsCountResponseObject, error) {
	count, err := c.service.GetCount(
		ctx, domain.NewNewsesFilter(req.Params.CategoryID, req.Params.TagID),
	)
	if err != nil {
		return nil, err
	}

	return apihttp.GetNewsCount200JSONResponse{Count: count}, nil
}

func newCategory(dto domain.Category) apihttp.Category {
	return apihttp.Category{
		ID:    dto.ID,
		Title: dto.Title,
	}
}

func newTag(dto domain.Tag) apihttp.Tag {
	return apihttp.Tag{
		ID:   dto.ID,
		Name: dto.Name,
	}
}

func newTags(dtos []domain.Tag) []apihttp.Tag {
	tags := make([]apihttp.Tag, 0, len(dtos))

	for _, dto := range dtos {
		tags = append(tags, newTag(dto))
	}

	return tags
}

func newNews(dto domain.News) apihttp.News {
	return apihttp.News{
		Author:      dto.Author,
		Category:    newCategory(dto.Category),
		Content:     dto.Content,
		ID:          dto.ID,
		PublishedAt: dto.PublishedAt,
		ShortText:   dto.ShortText,
		TagIds:      dto.TagIds,
		Tags:        newTags(dto.Tags),
		Title:       dto.Title,
	}
}

func newNewsListable(dto domain.News) apihttp.NewsListable {
	return apihttp.NewsListable{
		Author:      dto.Author,
		Category:    newCategory(dto.Category),
		ID:          dto.ID,
		PublishedAt: dto.PublishedAt,
		ShortText:   dto.ShortText,
		TagIds:      dto.TagIds,
		Tags:        newTags(dto.Tags),
		Title:       dto.Title,
	}
}
