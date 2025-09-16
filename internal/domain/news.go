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

func (n *Service) GetList(
	ctx context.Context,
	categoryID int,
	tagID int,
	page db.PaginationReq,
) ([]News, error) {
	items, err := n.repo.GetList(ctx, categoryID, tagID, page)
	if err != nil {
		return nil, fmt.Errorf("read news list: %w", err)
	}

	return NewNewses(items), nil
}

func (n *Service) GetNews(ctx context.Context, id int) (News, error) {
	item, err := n.repo.GetNews(ctx, id)
	if err != nil {
		return News{}, fmt.Errorf("read news item: %w", err)
	}

	return NewNews(item), nil
}

func (n *Service) GetCount(ctx context.Context, categoryID int, tagID int) (int, error) {
	count, err := n.repo.CountNews(ctx, categoryID, tagID)
	if err != nil {
		return 0, err
	}

	return count, nil
}
