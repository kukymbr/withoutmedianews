package apihttp

import (
	"context"
	"fmt"

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
		domain.PaginationReq{
			Page:    req.Params.Page,
			PerPage: req.Params.PerPage,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("read news list: %w", err)
	}

	resp := make(GetNews200JSONResponse, 0, len(items))

	for _, item := range items {
		tags := make([]Tag, 0, len(item.Tags))
		for _, tag := range item.Tags {
			tags = append(tags, Tag(tag))
		}

		resp = append(resp, NewsItem{
			Author:      item.Author,
			Category:    Category(item.Category),
			Content:     item.Content,
			ID:          item.ID,
			PublishedAt: item.PublishedAt,
			ShortText:   item.ShortText,
			Tags:        tags,
			Title:       item.Title,
		})
	}

	return resp, nil
}
