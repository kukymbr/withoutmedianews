package controller

import (
	"context"
	"fmt"

	"github.com/kukymbr/withoutmedianews/internal/api/http"
	"github.com/kukymbr/withoutmedianews/internal/db"
	"github.com/kukymbr/withoutmedianews/internal/pkg/dbkit"
	"github.com/kukymbr/withoutmedianews/internal/pkg/sqlvalues"
)

func NewNewsController(repo *db.WithoutmedianewsRepo) *NewsController {
	return &NewsController{
		repo: repo,
	}
}

type NewsController struct {
	repo *db.WithoutmedianewsRepo
}

func (c *NewsController) GetNewses(
	ctx context.Context,
	req apihttp.GetNewsesRequestObject,
) (apihttp.GetNewsesResponseObject, error) {
	search := db.NewsSearch{}

	if req.Params.CategoryID > 0 {
		search.CategoryID = &req.Params.CategoryID
	}

	if req.Params.TagID > 0 {
		// TODO: add wrapper
		search.With("? = ANY(t.?)", req.Params.TagID, db.Columns.News.TagIDs)
	}

	items, err := c.repo.NewsByFilters(
		ctx,
		&search,
		db.NewPager(req.Params.Page, req.Params.PerPage),
		db.EnabledOnly(),
	)
	if err != nil {
		return nil, fmt.Errorf("get news list: %w", err)
	}

	resp := make(apihttp.GetNewses200JSONResponse, 0, len(items))

	for _, item := range items {
		resp = append(resp, newNews(item))
	}

	return resp, nil
}

func (c *NewsController) GetNews(
	ctx context.Context,
	request apihttp.GetNewsRequestObject,
) (apihttp.GetNewsResponseObject, error) {
	item, err := c.repo.OneNews(ctx, &db.NewsSearch{ID: &request.ID}, db.EnabledOnly())
	if err != nil {
		return nil, fmt.Errorf("get news: %w", err)
	}

	if item == nil {
		return apihttp.GetNews200JSONResponse{}, fmt.Errorf("news not found: %w", dbkit.ErrNotFound)
	}

	return apihttp.GetNews200JSONResponse(newNews(*item)), nil
}

func (c *NewsController) GetNewsCount(
	ctx context.Context,
	request apihttp.GetNewsCountRequestObject,
) (apihttp.GetNewsCountResponseObject, error) {
	search := db.NewsSearch{}

	if request.Params.CategoryID > 0 {
		search.CategoryID = &request.Params.CategoryID
	}

	if request.Params.TagID > 0 {
		// TODO: add wrapper
		search.With("? = ANY(t.?)", request.Params.TagID, db.Columns.News.TagIDs)
	}

	count, err := c.repo.CountNews(ctx, &search, db.EnabledOnly())
	if err != nil {
		return nil, fmt.Errorf("count news: %w", err)
	}

	return apihttp.GetNewsCount200JSONResponse{Count: count}, nil
}

func newNews(news db.News) apihttp.News {
	return apihttp.News{
		Author:      sqlvalues.PtrToValue(news.Author),
		Category:    apihttp.Category{},
		Content:     sqlvalues.PtrToValue(news.Content),
		ID:          news.ID,
		PublishedAt: news.PublishedAt,
		ShortText:   news.ShortText,
		Tags:        apihttp.NewTags(nil), // TODO: enrich with tags
		Title:       news.Title,
	}
}
