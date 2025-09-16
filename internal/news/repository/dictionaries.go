package repository

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/kukymbr/withoutmedianews/internal/domain"
	"github.com/kukymbr/withoutmedianews/internal/pkg/dbkit"
	"go.uber.org/zap"
)

func NewDictionariesRepository(db *dbkit.Database, logger *zap.Logger) *DictionariesRepository {
	return &DictionariesRepository{
		db:     db,
		logger: logger,
	}
}

type DictionariesRepository struct {
	db     *dbkit.Database
	logger *zap.Logger
}

func (r *DictionariesRepository) ReadCategoriesList(ctx context.Context) ([]domain.Category, error) {
	ds := r.db.Goqu().Select().
		From(tableNameCategories).
		Order(goqu.C("sort").Asc().NullsLast())

	var dtos []categoryDTO

	if err := dbkit.GoquScanStructs(ctx, ds, &dtos, r.logger); err != nil {
		return nil, err
	}

	categories := make([]domain.Category, 0, len(dtos))

	for _, dto := range dtos {
		categories = append(categories, dto.ToDomain())
	}

	return categories, nil
}

func (r *DictionariesRepository) ReadTagsList(ctx context.Context) ([]domain.Tag, error) {
	ds := r.db.Goqu().Select().
		From(tableNameTags).
		Order(goqu.C("name").Asc())

	var dtos []tagDTO

	if err := dbkit.GoquScanStructs(ctx, ds, &dtos, r.logger); err != nil {
		return nil, err
	}

	tags := make([]domain.Tag, 0, len(dtos))

	for _, dto := range dtos {
		tags = append(tags, dto.ToDomain())
	}

	return tags, nil
}
