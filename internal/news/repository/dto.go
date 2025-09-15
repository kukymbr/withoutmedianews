package repository

import (
	"time"

	"github.com/kukymbr/withoutmedianews/internal/domain"
	"github.com/lib/pq"
)

type idProvider interface {
	GetID() int
}

type newsItemDTO struct {
	ID int `db:"newsId"`

	Title     string `db:"title"`
	ShortText string `db:"shortText"`
	Content   string `db:"content"`
	Author    string `db:"author"`

	CategoryID int           `db:"categoryId"`
	TagIDs     pq.Int64Array `db:"tagIds"`

	Category categoryDTO `db:"-"`
	Tags     []tagDTO    `db:"-"`

	PublishedAt time.Time `db:"publishedAt"`
	CreatedAt   time.Time `db:"createdAt"`

	Status int `db:"statusId"`
}

func (dto newsItemDTO) GetID() int {
	return dto.ID
}

func (dto newsItemDTO) ToDomain() domain.NewsItem {
	return domain.NewsItem{
		ID: dto.ID,

		Title:     dto.Title,
		ShortText: dto.ShortText,
		Content:   dto.Content,
		Author:    dto.Author,

		PublishedAt: dto.PublishedAt,
	}
}

type tagDTO struct {
	ID   int    `db:"tagId"`
	Name string `db:"name"`
}

func (dto tagDTO) GetID() int {
	return dto.ID
}

func (dto tagDTO) ToDomain() domain.Tag {
	return domain.Tag(dto)
}

type categoryDTO struct {
	ID    int    `db:"categoryId"`
	Title string `db:"title"`
}

func (dto categoryDTO) GetID() int {
	return dto.ID
}

func (dto categoryDTO) ToDomain() domain.Category {
	return domain.Category(dto)
}
