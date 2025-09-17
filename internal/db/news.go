package db

import (
	"context"
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
	"go.uber.org/zap"
)

func NewNewsRepository(db *pg.DB, logger *zap.Logger) *NewsRepository {
	return &NewsRepository{
		db:     db,
		logger: logger,
	}
}

type NewsRepository struct {
	db     *pg.DB
	logger *zap.Logger
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
	if err != nil {
		return nil, fmt.Errorf("fetch news list: %w", err)
	}

	return dtos, nil
}

func (r *NewsRepository) GetNews(ctx context.Context, id int) (News, error) {
	var dto News

	err := r.getNewsQueryBase(ctx, &dto, 0, 0).
		Where("? = ?", pg.Ident(Columns.News.ID), id).
		Limit(1).
		Select()
	if err != nil {
		return News{}, fmt.Errorf("fetch single news item: %w", err)
	}

	return dto, nil
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
	if err := query.Select(); err != nil {
		return nil, fmt.Errorf("fetch categories: %w", err)
	}

	return items, nil
}

func (r *NewsRepository) ReadTags(ctx context.Context) ([]Tag, error) {
	var items []Tag

	query := r.withPublishedFilter(r.db.ModelContext(ctx, &items), Columns.Category.StatusID)
	if err := query.Select(); err != nil {
		return nil, fmt.Errorf("fetch tags: %w", err)
	}

	return items, nil
}

func (r *NewsRepository) getNewsQueryBase(ctx context.Context, model any, categoryID int, tagID int) *pg.Query {
	now := time.Now().UTC().Format("2006-01-02 15:04:05")

	query := r.db.ModelContext(ctx, model)
	query = r.withPublishedFilter(query, Columns.News.StatusID).
		Relation(Columns.News.Category).
		Where("? AT TIME ZONE 'UTC' <= ?", pg.Ident(Columns.News.PublishedAt), now)

	if categoryID > 0 {
		query = query.Where("? = ?", pg.Ident(Columns.News.CategoryID), categoryID)
	}

	if tagID > 0 {
		query = query.Where("? = ANY(?)", tagID, pg.Ident(Columns.News.TagIDs))
	}

	return query
}

func (r *NewsRepository) withPublishedFilter(query *pg.Query, field string) *pg.Query {
	return query.Where("? = ?", pg.Ident(field), statusPublished)
}
