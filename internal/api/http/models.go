package apihttp

import (
	"github.com/kukymbr/withoutmedianews/internal/db"
)

func NewTags(dt []db.Tag) []Tag {
	tags := make([]Tag, 0, len(dt))

	for _, tag := range dt {
		tags = append(tags, Tag{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}

	return tags
}
