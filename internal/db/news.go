package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/kukymbr/withoutmedianews/internal/pkg/dbkit"
)

func NewNewsRepository(db *pg.DB) *NewsRepository {
	return &NewsRepository{
		db: db,
	}
}

type NewsRepository struct {
	db *pg.DB
}

func (r *NewsRepository) GetList(
	ctx context.Context,
	categoryID int,
	tagID int,
	page PaginationReq,
) ([]News, error) {
	var dtos []News

	page = page.GetNormalized()

	err := r.getNewsQueryBase(ctx, &dtos, categoryID, tagID).
		Order(Columns.News.PublishedAt + " DESC").
		Offset(page.Offset()).
		Limit(page.PerPage).
		Select()

	return dtos, listResponseError(err)
}

func (r *NewsRepository) GetNews(ctx context.Context, id int) (News, error) {
	var dto News

	err := r.getNewsQueryBase(ctx, &dto, 0, 0).
		Where("t.? = ?", pg.Ident(Columns.News.ID), id).
		Limit(1).
		Select()

	return dto, itemResponseError(err)
}

func (r *NewsRepository) CountNews(ctx context.Context, categoryID int, tagID int) (int, error) {
	count, err := r.getNewsQueryBase(ctx, &News{}, categoryID, tagID).Count()
	if err != nil {
		return 0, fmt.Errorf("fetch news count: %w", err)
	}

	return count, nil
}

func (r *NewsRepository) ReadCategories(ctx context.Context) ([]Category, error) {
	var items []Category

	query := r.withPublishedFilter(r.db.ModelContext(ctx, &items), Columns.Category.StatusID)
	err := query.Select()

	return items, listResponseError(err)
}

func (r *NewsRepository) ReadTags(ctx context.Context, ids []int) ([]Tag, error) {
	var items []Tag

	if ids != nil && len(ids) == 0 {
		return nil, fmt.Errorf("%w: empty IDs list given", dbkit.ErrNotFound)
	}

	query := r.withPublishedFilter(r.db.ModelContext(ctx, &items), Columns.Category.StatusID)
	if len(ids) > 0 {
		query.Where("t.? IN (?)", pg.Ident(Columns.Tag.ID), pg.In(ids))
	}

	err := query.Select()

	return items, listResponseError(err)
}

func (r *NewsRepository) getNewsQueryBase(ctx context.Context, model any, categoryID int, tagID int) *pg.Query {
	now := time.Now().UTC().Format("2006-01-02 15:04:05")

	query := r.db.ModelContext(ctx, model)
	query = r.withPublishedFilter(query, Columns.News.StatusID).
		Relation(Columns.News.Category).
		Where("t.? AT TIME ZONE 'UTC' <= ?", pg.Ident(Columns.News.PublishedAt), now)

	if categoryID > 0 {
		query = query.Where("t.? = ?", pg.Ident(Columns.News.CategoryID), categoryID)
	}

	if tagID > 0 {
		query = query.Where("? = ANY(t.?)", tagID, pg.Ident(Columns.News.TagIDs))
	}

	return query
}

func (r *NewsRepository) withPublishedFilter(query *pg.Query, field string) *pg.Query {
	return query.Where("t.? = ?", pg.Ident(field), statusPublished)
}

func listResponseError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, pg.ErrNoRows):
		return nil
	}

	return err
}

func itemResponseError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, pg.ErrNoRows):
		return fmt.Errorf("%w: %s", dbkit.ErrNotFound, err.Error())
	}

	return err
}
