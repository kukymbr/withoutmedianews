package domain

import (
	"context"
	"fmt"

	"github.com/kukymbr/withoutmedianews/internal/db"
)

func NewNewsService(repo *db.NewsRepository) *Service {
	return &Service{
		repo: repo,
	}
}

type Service struct {
	repo *db.NewsRepository
}

func (s *Service) GetList(
	ctx context.Context,
	categoryID int,
	tagID int,
	page db.PaginationReq,
) ([]News, error) {
	items, err := s.repo.GetList(ctx, categoryID, tagID, page)
	if err != nil {
		return nil, fmt.Errorf("read news list: %w", err)
	}

	return s.newNewses(ctx, items)
}

func (s *Service) GetNews(ctx context.Context, id int) (News, error) {
	item, err := s.repo.GetNews(ctx, id)
	if err != nil {
		return News{}, fmt.Errorf("read news item: %w", err)
	}

	list, err := s.newNewses(ctx, []db.News{item})
	if err != nil {
		return News{}, err
	}

	return list[0], nil
}

func (s *Service) GetCount(ctx context.Context, categoryID int, tagID int) (int, error) {
	count, err := s.repo.CountNews(ctx, categoryID, tagID)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *Service) GetCategories(ctx context.Context) ([]Category, error) {
	categories, err := s.repo.ReadCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("read categories from repo: %w", err)
	}

	return NewCategories(categories), nil
}

func (s *Service) GetTags(ctx context.Context) ([]Tag, error) {
	tags, err := s.repo.ReadTags(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("read tags from repo: %w", err)
	}

	return NewTags(tags), nil
}

func (s *Service) newNewses(ctx context.Context, items []db.News) ([]News, error) {
	if len(items) == 0 {
		return nil, nil
	}

	tagIDs := make([]int, 0, len(items))
	for _, item := range items {
		tagIDs = append(tagIDs, item.TagIDs...)
	}

	tags, err := s.repo.ReadTags(ctx, tagIDs)
	if err != nil {
		return nil, err
	}

	newses := make([]News, 0, len(items))
	index := newTagsIndex(tags)

	for _, dto := range items {
		news := NewNews(dto)
		tags := make([]Tag, 0, len(dto.TagIDs))

		for _, tagID := range dto.TagIDs {
			tag := Tag{ID: tagID}
			if t, ok := index[tag.ID]; ok {
				tag = t
			}

			tags = append(tags, tag)
		}

		news.Tags = tags
		newses = append(newses, news)
	}

	return newses, nil
}

func newTagsIndex(tags []db.Tag) map[int]Tag {
	index := make(map[int]Tag, len(tags))

	for _, tag := range tags {
		index[tag.ID] = NewTag(tag)
	}

	return index
}
