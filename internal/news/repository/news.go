package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/kukymbr/withoutmedianews/internal/domain"
	"github.com/kukymbr/withoutmedianews/internal/pkg/dbkit"
	"go.uber.org/zap"
)

func NewNewsRepository(db *dbkit.Database, logger *zap.Logger) *NewsRepository {
	return &NewsRepository{
		db:     db,
		logger: logger,
	}
}

type NewsRepository struct {
	db     *dbkit.Database
	logger *zap.Logger
}

func (r *NewsRepository) GetNewsList(
	ctx context.Context,
	categoryID int,
	tagID int,
	page domain.PaginationReq,
) ([]domain.NewsItem, error) {
	page = page.GetNormalized()
	now := time.Now().UTC().Format("2006-01-02 15:04:05")

	ds := r.db.Goqu().
		Select().
		From(tableNameNews).
		Where(goqu.Ex{"statusId": statusPublished}).
		Where(goqu.L(`"publishedAt" AT TIME ZONE 'UTC' <= ?`, now)).
		Order(goqu.C("publishedAt").Desc()).
		Offset(page.Offset()).
		Limit(page.PerPage)

	if categoryID > 0 {
		ds = ds.Where(goqu.Ex{"categoryId": categoryID})
	}

	if tagID > 0 {
		ds = ds.Where(goqu.L("tagIds @> ?", []int{tagID}))
	}

	var dtos []newsItemDTO

	if err := dbkit.GoquScanStructs(ctx, ds, &dtos, r.logger); err != nil {
		return nil, fmt.Errorf("scan news list results: %w", err)
	}

	if err := r.enrichNewsDTOs(ctx, dtos); err != nil {
		return nil, err
	}

	items := make([]domain.NewsItem, 0, len(dtos))

	for _, dto := range dtos {
		items = append(items, dto.ToDomain())
	}

	return items, nil
}

func (r *NewsRepository) GetNewsItem(ctx context.Context, id int) (domain.NewsItem, error) {
	ds := r.db.Goqu().
		Select().
		From(tableNameNews).
		Where(goqu.Ex{"newsId": id}).
		Limit(1)

	var dto newsItemDTO

	if err := dbkit.GoquScanStruct(ctx, ds, &dto, r.logger); err != nil {
		return domain.NewsItem{}, fmt.Errorf("fetch single news item: %w", err)
	}

	if err := r.enrichNewsDTOs(ctx, []newsItemDTO{dto}); err != nil {
		return domain.NewsItem{}, err
	}

	return dto.ToDomain(), nil
}

func (r *NewsRepository) enrichNewsDTOs(ctx context.Context, dtos []newsItemDTO) error {
	if err := r.enrichWithTags(ctx, dtos); err != nil {
		return fmt.Errorf("enrich news list results with tags: %w", err)
	}

	if err := r.enrichWithCategories(ctx, dtos); err != nil {
		return fmt.Errorf("enrich news list results with categories: %w", err)
	}

	return nil
}

func (r *NewsRepository) enrichWithTags(ctx context.Context, dtos []newsItemDTO) error {
	ids := make([]int, 0, len(dtos))

	for _, dto := range dtos {
		for _, id := range dto.TagIDs {
			ids = append(ids, int(id))
		}
	}

	if len(ids) == 0 {
		return nil
	}

	tagsDTOs := make([]tagDTO, 0, len(ids))

	ds := r.db.Goqu().
		Select().
		From(tableNameTags).
		Where(goqu.Ex{"tagId": ids})

	if err := dbkit.GoquScanStructs(ctx, ds, &tagsDTOs, r.logger); err != nil {
		return fmt.Errorf("scan tags: %w", err)
	}

	tagsIndex := indexRecords(tagsDTOs)

	for i, dto := range dtos {
		dto.Tags = make([]tagDTO, 0, len(dto.TagIDs))

		for _, id := range dto.TagIDs {
			if tag, ok := tagsIndex[int(id)]; ok {
				dto.Tags = append(dto.Tags, tag)
			}
		}

		dtos[i] = dto
	}

	return nil
}

func (r *NewsRepository) enrichWithCategories(ctx context.Context, dtos []newsItemDTO) error {
	ids := make([]int, 0, len(dtos))

	for _, dto := range dtos {
		ids = append(ids, dto.CategoryID)
	}

	if len(ids) == 0 {
		return nil
	}

	categoriesDTOs := make([]categoryDTO, 0, len(ids))

	ds := r.db.Goqu().
		Select().
		From(tableNameCategories).
		Where(goqu.Ex{"categoryId": ids})

	if err := dbkit.GoquScanStructs(ctx, ds, &categoriesDTOs, r.logger); err != nil {
		return fmt.Errorf("scan categories: %w", err)
	}

	categoriesIndex := indexRecords(categoriesDTOs)

	for i, dto := range dtos {
		if category, ok := categoriesIndex[dto.CategoryID]; ok {
			dto.Category = category
		}

		dtos[i] = dto
	}

	return nil
}

func indexRecords[DTOType idProvider](dtos []DTOType) map[int]DTOType {
	index := make(map[int]DTOType, len(dtos))

	for _, dto := range dtos {
		index[dto.GetID()] = dto
	}

	return index
}
