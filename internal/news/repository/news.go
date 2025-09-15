package repository

import (
	"context"
	"database/sql"

	"github.com/kukymbr/withoutmedianews/internal/domain"
)

func NewNewsRepository(db *sql.DB) *NewsRepository {
	return &NewsRepository{
		db: db,
	}
}

type NewsRepository struct {
	db *sql.DB
}

func (r *NewsRepository) GetNewsList(
	ctx context.Context,
	categoryID int,
	tagID int,
	page domain.PaginationReq,
) ([]domain.NewsItem, error) {
	panic("not implemented")
}
