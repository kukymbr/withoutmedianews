package db

import (
	"context"
	"errors"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type WithoutmedianewsRepo struct {
	db      orm.DB
	filters map[string][]Filter
	sort    map[string][]SortField
	join    map[string][]string
}

// NewWithoutmedianewsRepo returns new repository
func NewWithoutmedianewsRepo(db orm.DB) WithoutmedianewsRepo {
	return WithoutmedianewsRepo{
		db: db,
		filters: map[string][]Filter{
			Tables.Category.Name: {StatusFilter},
			Tables.News.Name:     {StatusFilter},
			Tables.Tag.Name:      {StatusFilter},
		},
		sort: map[string][]SortField{
			Tables.Category.Name: {{Column: Columns.Category.Title, Direction: SortAsc}},
			Tables.News.Name:     {{Column: Columns.News.CreatedAt, Direction: SortDesc}},
			Tables.Tag.Name:      {{Column: Columns.Tag.ID, Direction: SortDesc}},
		},
		join: map[string][]string{
			Tables.Category.Name: {TableColumns},
			Tables.News.Name:     {TableColumns, Columns.News.Category},
			Tables.Tag.Name:      {TableColumns},
		},
	}
}

// WithTransaction is a function that wraps WithoutmedianewsRepo with pg.Tx transaction.
func (wr WithoutmedianewsRepo) WithTransaction(tx *pg.Tx) WithoutmedianewsRepo {
	wr.db = tx
	return wr
}

// WithEnabledOnly is a function that adds "statusId"=1 as base filter.
func (wr WithoutmedianewsRepo) WithEnabledOnly() WithoutmedianewsRepo {
	f := make(map[string][]Filter, len(wr.filters))
	for i := range wr.filters {
		f[i] = make([]Filter, len(wr.filters[i]))
		copy(f[i], wr.filters[i])
		f[i] = append(f[i], StatusEnabledFilter)
	}
	wr.filters = f

	return wr
}

/*** Category ***/

// FullCategory returns full joins with all columns
func (wr WithoutmedianewsRepo) FullCategory() OpFunc {
	return WithColumns(wr.join[Tables.Category.Name]...)
}

// DefaultCategorySort returns default sort.
func (wr WithoutmedianewsRepo) DefaultCategorySort() OpFunc {
	return WithSort(wr.sort[Tables.Category.Name]...)
}

// CategoryByID is a function that returns Category by ID(s) or nil.
func (wr WithoutmedianewsRepo) CategoryByID(ctx context.Context, id int, ops ...OpFunc) (*Category, error) {
	return wr.OneCategory(ctx, &CategorySearch{ID: &id}, ops...)
}

// OneCategory is a function that returns one Category by filters. It could return pg.ErrMultiRows.
func (wr WithoutmedianewsRepo) OneCategory(ctx context.Context, search *CategorySearch, ops ...OpFunc) (*Category, error) {
	obj := &Category{}
	err := buildQuery(ctx, wr.db, obj, search, wr.filters[Tables.Category.Name], PagerTwo, ops...).Select()

	if errors.Is(err, pg.ErrMultiRows) {
		return nil, err
	} else if errors.Is(err, pg.ErrNoRows) {
		return nil, nil
	}

	return obj, err
}

// CategoriesByFilters returns Category list.
func (wr WithoutmedianewsRepo) CategoriesByFilters(ctx context.Context, search *CategorySearch, pager Pager, ops ...OpFunc) (categories []Category, err error) {
	err = buildQuery(ctx, wr.db, &categories, search, wr.filters[Tables.Category.Name], pager, ops...).Select()
	return
}

// CountCategories returns count
func (wr WithoutmedianewsRepo) CountCategories(ctx context.Context, search *CategorySearch, ops ...OpFunc) (int, error) {
	return buildQuery(ctx, wr.db, &Category{}, search, wr.filters[Tables.Category.Name], PagerOne, ops...).Count()
}

// AddCategory adds Category to DB.
func (wr WithoutmedianewsRepo) AddCategory(ctx context.Context, category *Category, ops ...OpFunc) (*Category, error) {
	q := wr.db.ModelContext(ctx, category)
	applyOps(q, ops...)
	_, err := q.Insert()

	return category, err
}

// UpdateCategory updates Category in DB.
func (wr WithoutmedianewsRepo) UpdateCategory(ctx context.Context, category *Category, ops ...OpFunc) (bool, error) {
	q := wr.db.ModelContext(ctx, category).WherePK()
	if len(ops) == 0 {
		q = q.ExcludeColumn(Columns.Category.ID)
	}
	applyOps(q, ops...)
	res, err := q.Update()
	if err != nil {
		return false, err
	}

	return res.RowsAffected() > 0, err
}

// DeleteCategory set statusId to deleted in DB.
func (wr WithoutmedianewsRepo) DeleteCategory(ctx context.Context, id int) (deleted bool, err error) {
	category := &Category{ID: id, StatusID: StatusDeleted}

	return wr.UpdateCategory(ctx, category, WithColumns(Columns.Category.StatusID))
}

/*** News ***/

// FullNews returns full joins with all columns
func (wr WithoutmedianewsRepo) FullNews() OpFunc {
	return WithColumns(wr.join[Tables.News.Name]...)
}

// DefaultNewsSort returns default sort.
func (wr WithoutmedianewsRepo) DefaultNewsSort() OpFunc {
	return WithSort(wr.sort[Tables.News.Name]...)
}

// NewsByID is a function that returns News by ID(s) or nil.
func (wr WithoutmedianewsRepo) NewsByID(ctx context.Context, id int, ops ...OpFunc) (*News, error) {
	return wr.OneNews(ctx, &NewsSearch{ID: &id}, ops...)
}

// OneNews is a function that returns one News by filters. It could return pg.ErrMultiRows.
func (wr WithoutmedianewsRepo) OneNews(ctx context.Context, search *NewsSearch, ops ...OpFunc) (*News, error) {
	obj := &News{}
	err := buildQuery(ctx, wr.db, obj, search, wr.filters[Tables.News.Name], PagerTwo, ops...).Select()

	if errors.Is(err, pg.ErrMultiRows) {
		return nil, err
	} else if errors.Is(err, pg.ErrNoRows) {
		return nil, nil
	}

	return obj, err
}

// NewsByFilters returns News list.
func (wr WithoutmedianewsRepo) NewsByFilters(ctx context.Context, search *NewsSearch, pager Pager, ops ...OpFunc) (newsList []News, err error) {
	err = buildQuery(ctx, wr.db, &newsList, search, wr.filters[Tables.News.Name], pager, ops...).Select()
	return
}

// CountNews returns count
func (wr WithoutmedianewsRepo) CountNews(ctx context.Context, search *NewsSearch, ops ...OpFunc) (int, error) {
	return buildQuery(ctx, wr.db, &News{}, search, wr.filters[Tables.News.Name], PagerOne, ops...).Count()
}

// AddNews adds News to DB.
func (wr WithoutmedianewsRepo) AddNews(ctx context.Context, news *News, ops ...OpFunc) (*News, error) {
	q := wr.db.ModelContext(ctx, news)
	if len(ops) == 0 {
		q = q.ExcludeColumn(Columns.News.CreatedAt)
	}
	applyOps(q, ops...)
	_, err := q.Insert()

	return news, err
}

// UpdateNews updates News in DB.
func (wr WithoutmedianewsRepo) UpdateNews(ctx context.Context, news *News, ops ...OpFunc) (bool, error) {
	q := wr.db.ModelContext(ctx, news).WherePK()
	if len(ops) == 0 {
		q = q.ExcludeColumn(Columns.News.ID, Columns.News.CreatedAt)
	}
	applyOps(q, ops...)
	res, err := q.Update()
	if err != nil {
		return false, err
	}

	return res.RowsAffected() > 0, err
}

// DeleteNews set statusId to deleted in DB.
func (wr WithoutmedianewsRepo) DeleteNews(ctx context.Context, id int) (deleted bool, err error) {
	news := &News{ID: id, StatusID: StatusDeleted}

	return wr.UpdateNews(ctx, news, WithColumns(Columns.News.StatusID))
}

/*** Tag ***/

// FullTag returns full joins with all columns
func (wr WithoutmedianewsRepo) FullTag() OpFunc {
	return WithColumns(wr.join[Tables.Tag.Name]...)
}

// DefaultTagSort returns default sort.
func (wr WithoutmedianewsRepo) DefaultTagSort() OpFunc {
	return WithSort(wr.sort[Tables.Tag.Name]...)
}

// TagByID is a function that returns Tag by ID(s) or nil.
func (wr WithoutmedianewsRepo) TagByID(ctx context.Context, id int, ops ...OpFunc) (*Tag, error) {
	return wr.OneTag(ctx, &TagSearch{ID: &id}, ops...)
}

// OneTag is a function that returns one Tag by filters. It could return pg.ErrMultiRows.
func (wr WithoutmedianewsRepo) OneTag(ctx context.Context, search *TagSearch, ops ...OpFunc) (*Tag, error) {
	obj := &Tag{}
	err := buildQuery(ctx, wr.db, obj, search, wr.filters[Tables.Tag.Name], PagerTwo, ops...).Select()

	if errors.Is(err, pg.ErrMultiRows) {
		return nil, err
	} else if errors.Is(err, pg.ErrNoRows) {
		return nil, nil
	}

	return obj, err
}

// TagsByFilters returns Tag list.
func (wr WithoutmedianewsRepo) TagsByFilters(ctx context.Context, search *TagSearch, pager Pager, ops ...OpFunc) (tags []Tag, err error) {
	err = buildQuery(ctx, wr.db, &tags, search, wr.filters[Tables.Tag.Name], pager, ops...).Select()
	return
}

// CountTags returns count
func (wr WithoutmedianewsRepo) CountTags(ctx context.Context, search *TagSearch, ops ...OpFunc) (int, error) {
	return buildQuery(ctx, wr.db, &Tag{}, search, wr.filters[Tables.Tag.Name], PagerOne, ops...).Count()
}

// AddTag adds Tag to DB.
func (wr WithoutmedianewsRepo) AddTag(ctx context.Context, tag *Tag, ops ...OpFunc) (*Tag, error) {
	q := wr.db.ModelContext(ctx, tag)
	applyOps(q, ops...)
	_, err := q.Insert()

	return tag, err
}

// UpdateTag updates Tag in DB.
func (wr WithoutmedianewsRepo) UpdateTag(ctx context.Context, tag *Tag, ops ...OpFunc) (bool, error) {
	q := wr.db.ModelContext(ctx, tag).WherePK()
	if len(ops) == 0 {
		q = q.ExcludeColumn(Columns.Tag.ID)
	}
	applyOps(q, ops...)
	res, err := q.Update()
	if err != nil {
		return false, err
	}

	return res.RowsAffected() > 0, err
}

// DeleteTag set statusId to deleted in DB.
func (wr WithoutmedianewsRepo) DeleteTag(ctx context.Context, id int) (deleted bool, err error) {
	tag := &Tag{ID: id, StatusID: StatusDeleted}

	return wr.UpdateTag(ctx, tag, WithColumns(Columns.Tag.StatusID))
}
