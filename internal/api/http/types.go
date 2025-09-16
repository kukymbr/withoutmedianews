package apihttp

import "github.com/kukymbr/withoutmedianews/internal/domain"

func tagsFromDomain(dt []domain.Tag) []Tag {
	tags := make([]Tag, 0, len(dt))

	for _, tag := range dt {
		tags = append(tags, Tag(tag))
	}

	return tags
}
