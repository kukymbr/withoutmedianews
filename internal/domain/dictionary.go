package domain

import (
	"context"
	"fmt"

	"github.com/kukymbr/withoutmedianews/internal/db"
)

func NewDictionaryService(repo *db.DictionaryRepository) *DictionaryService {
	return &DictionaryService{
		repo: repo,
	}
}

type DictionaryService struct {
	repo *db.DictionaryRepository
}

func (d *DictionaryService) GetCategories(ctx context.Context) ([]Category, error) {
	categories, err := d.repo.ReadCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("read categories from repo: %w", err)
	}

	return NewCategories(categories), nil
}

func (d *DictionaryService) GetTags(ctx context.Context) ([]Tag, error) {
	tags, err := d.repo.ReadTags(ctx)
	if err != nil {
		return nil, fmt.Errorf("read tags from repo: %w", err)
	}

	return NewTags(tags), nil
}
