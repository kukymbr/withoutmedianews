package apihttp

import (
	"context"

	"github.com/kukymbr/withoutmedianews/internal/domain"
	"github.com/kukymbr/withoutmedianews/internal/news"
)

func NewNewsController(service *news.News) *NewsController {
	return &NewsController{
		service: service,
	}
}

type NewsController struct {
	service *news.News
}

func (c *NewsController) GetNews(
	ctx context.Context,
	req GetNewsRequestObject,
) (GetNewsResponseObject, error) {
	items, err := c.service.GetList(
		ctx,
		req.Params.CategoryID,
		req.Params.TagID,
		//nolint:gosec // ignore uint conversion
		domain.PaginationReq{
			Page:    uint(req.Params.Page),
			PerPage: uint(req.Params.PerPage),
		},
	)
	if err != nil {
		return nil, err
	}

	resp := make(GetNews200JSONResponse, 0, len(items))

	for _, item := range items {
		resp = append(resp, NewsItem{
			Author:      item.Author,
			Category:    Category(item.Category),
			Content:     item.Content,
			ID:          item.ID,
			PublishedAt: item.PublishedAt,
			ShortText:   item.ShortText,
			Tags:        tagsFromDomain(item.Tags),
			Title:       item.Title,
		})
	}

	return resp, nil
}

func (c *NewsController) GetNewsItem(
	ctx context.Context,
	request GetNewsItemRequestObject,
) (GetNewsItemResponseObject, error) {
	item, err := c.service.GetNewsItem(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	return GetNewsItem200JSONResponse{
		Author:      item.Author,
		Category:    Category(item.Category),
		Content:     item.Content,
		ID:          item.ID,
		PublishedAt: item.PublishedAt,
		ShortText:   item.ShortText,
		Tags:        tagsFromDomain(item.Tags),
		Title:       item.Title,
	}, nil
}
