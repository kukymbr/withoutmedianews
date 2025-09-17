package controller

import (
	"context"

	"github.com/kukymbr/withoutmedianews/internal/api/http"
	"github.com/kukymbr/withoutmedianews/internal/db"
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
		req.Params.CategoryID,
		req.Params.TagID,
		db.PaginationReq{
			Page:    req.Params.Page,
			PerPage: req.Params.PerPage,
		},
	)
	if err != nil {
		return nil, err
	}

	resp := make(apihttp.GetNewses200JSONResponse, 0, len(items))

	for _, item := range items {
		resp = append(resp, apihttp.News{
			Author:      item.Author,
			Category:    apihttp.Category(item.Category),
			Content:     item.Content,
			ID:          item.ID,
			PublishedAt: item.PublishedAt,
			ShortText:   item.ShortText,
			Tags:        apihttp.NewTags(item.Tags),
			Title:       item.Title,
		})
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

	return apihttp.GetNews200JSONResponse{
		Author:      item.Author,
		Category:    apihttp.Category(item.Category),
		Content:     item.Content,
		ID:          item.ID,
		PublishedAt: item.PublishedAt,
		ShortText:   item.ShortText,
		Tags:        apihttp.NewTags(item.Tags),
		Title:       item.Title,
	}, nil
}

func (c *NewsController) GetNewsCount(
	ctx context.Context,
	request apihttp.GetNewsCountRequestObject,
) (apihttp.GetNewsCountResponseObject, error) {
	count, err := c.service.GetCount(ctx, request.Params.CategoryID, request.Params.TagID)
	if err != nil {
		return nil, err
	}

	return apihttp.GetNewsCount200JSONResponse{Count: count}, nil
}
