package news

import (
	"context"
	"fmt"

	"github.com/kukymbr/withoutmedianews/internal/domain"
)

type ReaderRepository interface {
	GetNewsList(ctx context.Context, categoryID int, tagID int, page domain.PaginationReq) ([]domain.NewsItem, error)
	GetNewsItem(ctx context.Context, id int) (domain.NewsItem, error)
}

func NewNewsService(repo ReaderRepository) *News {
	return &News{
		repo: repo,
	}
}

type News struct {
	repo ReaderRepository
}

func (n *News) GetList(
	ctx context.Context,
	categoryID int,
	tagID int,
	page domain.PaginationReq,
) ([]domain.NewsItem, error) {
	items, err := n.repo.GetNewsList(ctx, categoryID, tagID, page)
	if err != nil {
		return nil, fmt.Errorf("read news list: %w", err)
	}

	return items, nil
}

func (n *News) GetNewsItem(ctx context.Context, id int) (domain.NewsItem, error) {
	item, err := n.repo.GetNewsItem(ctx, id)
	if err != nil {
		return domain.NewsItem{}, fmt.Errorf("read news item: %w", err)
	}

	return item, nil
}
