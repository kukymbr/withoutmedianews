package db

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/kukymbr/withoutmedianews/internal/pkg/dbkit"
	"go.uber.org/zap"
)

func NewDictionaryRepository(db *dbkit.Database, logger *zap.Logger) *DictionaryRepository {
	return &DictionaryRepository{
		db:     db,
		logger: logger,
	}
}

type DictionaryRepository struct {
	db     *dbkit.Database
	logger *zap.Logger
}

func (r *DictionaryRepository) ReadCategories(ctx context.Context) ([]Category, error) {
	ds := r.db.Goqu().Select().
		From(tableNameCategories).
		Order(goqu.C("sort").Asc().NullsLast())

	var dtos []Category

	if err := dbkit.GoquScanStructs(ctx, ds, &dtos, r.logger); err != nil {
		return nil, err
	}

	return dtos, nil
}

func (r *DictionaryRepository) ReadTags(ctx context.Context) ([]Tag, error) {
	ds := r.db.Goqu().Select().
		From(tableNameTags).
		Order(goqu.C("name").Asc())

	var dtos []Tag

	if err := dbkit.GoquScanStructs(ctx, ds, &dtos, r.logger); err != nil {
		return nil, err
	}

	return dtos, nil
}
