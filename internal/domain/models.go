package domain

import (
	"github.com/kukymbr/withoutmedianews/internal/db"
	"github.com/kukymbr/withoutmedianews/internal/pkg/sqlvalues"
)

func NewNews(dto db.News) News {
	category := Category{}

	if dto.Category != nil {
		category = NewCategory(*dto.Category)
	}

	return News{
		ID:          dto.ID,
		Title:       dto.Title,
		Author:      sqlvalues.PtrToValue(dto.Author),
		ShortText:   dto.ShortText,
		Content:     sqlvalues.PtrToValue(dto.Content),
		Category:    category,
		PublishedAt: dto.PublishedAt,
	}
}

func NewTag(dto db.Tag) Tag {
	return Tag{
		ID:   dto.ID,
		Name: dto.Name,
	}
}

func NewTags(dtos []db.Tag) []Tag {
	tags := make([]Tag, 0, len(dtos))

	for _, dto := range dtos {
		tags = append(tags, NewTag(dto))
	}

	return tags
}

func NewCategory(dto db.Category) Category {
	return Category{
		ID:    dto.ID,
		Title: dto.Title,
	}
}

func NewCategories(dtos []db.Category) []Category {
	categories := make([]Category, 0, len(dtos))

	for _, dto := range dtos {
		categories = append(categories, NewCategory(dto))
	}

	return categories
}
