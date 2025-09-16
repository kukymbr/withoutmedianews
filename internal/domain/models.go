package domain

import "github.com/kukymbr/withoutmedianews/internal/db"

func NewNews(dto db.News) News {
	return News{
		ID:          dto.ID,
		Title:       dto.Title,
		Author:      dto.Author,
		ShortText:   dto.ShortText,
		Content:     dto.Content,
		Category:    Category(dto.Category),
		Tags:        NewTags(dto.Tags),
		PublishedAt: dto.PublishedAt,
	}
}

func NewNewses(dtos []db.News) []News {
	newses := make([]News, 0, len(dtos))

	for _, dto := range dtos {
		newses = append(newses, NewNews(dto))
	}

	return newses
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
