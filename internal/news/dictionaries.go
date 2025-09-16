package news

import (
	"context"
	"fmt"

	"github.com/kukymbr/withoutmedianews/internal/domain"
)

type CategoriesReaderRepository interface {
	ReadCategoriesList(ctx context.Context) ([]domain.Category, error)
}

type TagsReaderRepository interface {
	ReadTagsList(ctx context.Context) ([]domain.Tag, error)
}

func NewDictionariesService(
	categoriesRepo CategoriesReaderRepository,
	tagsRepo TagsReaderRepository,
) *Dictionaries {
	return &Dictionaries{
		categoriesRepo: categoriesRepo,
		tagsRepo:       tagsRepo,
	}
}

type Dictionaries struct {
	categoriesRepo CategoriesReaderRepository
	tagsRepo       TagsReaderRepository
}

func (d *Dictionaries) GetCategories(ctx context.Context) ([]domain.Category, error) {
	categories, err := d.categoriesRepo.ReadCategoriesList(ctx)
	if err != nil {
		return nil, fmt.Errorf("read categories from repo: %w", err)
	}

	return categories, nil
}

func (d *Dictionaries) GetTags(ctx context.Context) ([]domain.Tag, error) {
	tags, err := d.tagsRepo.ReadTagsList(ctx)
	if err != nil {
		return nil, fmt.Errorf("read tags from repo: %w", err)
	}

	return tags, nil
}
