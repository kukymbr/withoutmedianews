package domain

import (
	"github.com/kukymbr/withoutmedianews/internal/db"
	"github.com/kukymbr/withoutmedianews/internal/pkg/ptrs"
)

func NewNewses(dtos []db.News) []News {
	newses := make([]News, len(dtos))

	for i, dto := range dtos {
		newses[i] = NewNews(dto)
	}

	return newses
}

func NewNews(dto db.News) News {
	return News{
		ID:          dto.ID,
		Title:       dto.Title,
		Author:      ptrs.PtrToValue(dto.Author),
		ShortText:   dto.ShortText,
		Content:     ptrs.PtrToValue(dto.Content),
		Category:    NewCategory(ptrs.PtrToValue(dto.Category)),
		PublishedAt: dto.PublishedAt,
		TagIds:      dto.TagIDs,
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
