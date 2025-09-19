package domain

import (
	"context"
	"fmt"

	"github.com/kukymbr/withoutmedianews/internal/db"
)

func NewNewsService(repo *db.NewsRepo) *Service {
	return &Service{
		repo: repo,
	}
}

type Service struct {
	repo *db.NewsRepo
}

func (s *Service) GetList(
	ctx context.Context,
	categoryID int,
	tagID int,
	page, perPage int,
) ([]News, error) {
	// TODO: nils
	search := db.NewsSearch{}
	if categoryID > 0 {
		search.CategoryID = &categoryID
	}

	if tagID > 0 {
		search.TagID = &tagID
	}

	items, err := s.repo.NewsByFilters(
		ctx,
		&search,
		db.NewPager(page, perPage),
		db.EnabledOnly(),
		db.AlreadyPublished(),
		// db.WithoutColumns(db.Columns.News.Content),
		db.WithColumns(db.Columns.News.Category),
	)
	if err != nil {
		return nil, fmt.Errorf("read news list: %w", err)
	}

	tagIDs := collectTagIDs(items)
	newses := NewNewses(items)

	return s.enrichNewsesWithTags(ctx, newses, tagIDs)
}

func (s *Service) GetNews(ctx context.Context, id int) (News, error) {
	dto, err := s.repo.OneNews(
		ctx,
		&db.NewsSearch{ID: &id},
		db.EnabledOnly(),
		db.AlreadyPublished(),
		db.WithColumns(db.Columns.News.Category),
	)
	if err != nil {
		return News{}, fmt.Errorf("read news item: %w", err)
	}

	if dto == nil {
		return News{}, ErrNotFound
	}

	list, err := s.enrichNewsesWithTags(ctx, []News{NewNews(*dto)}, dto.TagIDs)
	if err != nil {
		return News{}, err
	}

	return list[0], nil
}

func (s *Service) GetCount(ctx context.Context, categoryID int, tagID int) (int, error) {
	search := db.NewsSearch{}
	if categoryID > 0 {
		search.CategoryID = &categoryID
	}

	if tagID > 0 {
		search.TagID = &tagID
	}

	count, err := s.repo.CountNews(ctx, &search, db.EnabledOnly())
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *Service) GetCategories(ctx context.Context) ([]Category, error) {
	categories, err := s.repo.CategoriesByFilters(ctx, nil, db.PagerNoLimit, db.EnabledOnly())
	if err != nil {
		return nil, fmt.Errorf("read categories from repo: %w", err)
	}

	return NewCategories(categories), nil
}

func (s *Service) GetTags(ctx context.Context) ([]Tag, error) {
	tags, err := s.repo.TagsByFilters(ctx, nil, db.PagerNoLimit, db.EnabledOnly())
	if err != nil {
		return nil, fmt.Errorf("read tags from repo: %w", err)
	}

	return NewTags(tags), nil
}

func (s *Service) enrichNewsesWithTags(ctx context.Context, newses []News, tagIDs []int) ([]News, error) {
	if len(newses) == 0 {
		return nil, nil
	}

	tags, err := s.repo.TagsByFilters(ctx, &db.TagSearch{IDs: tagIDs}, db.PagerNoLimit, db.EnabledOnly())
	if err != nil {
		return nil, err
	}

	index := newTagsIndex(tags)

	for i, news := range newses {
		newses[i].Tags = index.getByIDs(news.TagIds...)
	}

	return newses, nil
}

func collectTagIDs(newses []db.News) []int {
	tagIDs := make([]int, 0, len(newses))
	for _, news := range newses {
		tagIDs = append(tagIDs, news.TagIDs...)
	}

	return tagIDs
}

type tagsIndex map[int]Tag

func (i tagsIndex) getByIDs(ids ...int) []Tag {
	tags := make([]Tag, 0, len(ids))

	for _, id := range ids {
		tag := Tag{ID: id}

		if _, ok := i[id]; ok {
			tag = i[id]
		}

		tags = append(tags, tag)
	}

	return tags
}

func newTagsIndex(tags []db.Tag) tagsIndex {
	index := make(map[int]Tag, len(tags))

	for _, tag := range tags {
		index[tag.ID] = NewTag(tag)
	}

	return index
}
