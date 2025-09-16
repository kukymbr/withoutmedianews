package db

import (
	"time"

	"github.com/lib/pq"
)

type idProvider interface {
	GetID() int
}

type News struct {
	ID int `db:"newsId"`

	Title     string `db:"title"`
	ShortText string `db:"shortText"`
	Content   string `db:"content"`
	Author    string `db:"author"`

	CategoryID int           `db:"categoryId"`
	TagIDs     pq.Int64Array `db:"tagIds"`

	Category Category `db:"-"`
	Tags     []Tag    `db:"-"`

	PublishedAt time.Time `db:"publishedAt"`
	CreatedAt   time.Time `db:"createdAt"`

	Status int `db:"statusId"`
}

func (dto News) GetID() int {
	return dto.ID
}

type Tag struct {
	ID   int    `db:"tagId"`
	Name string `db:"name"`
}

func (dto Tag) GetID() int {
	return dto.ID
}

type Category struct {
	ID    int    `db:"categoryId"`
	Title string `db:"title"`
}

func (dto Category) GetID() int {
	return dto.ID
}
