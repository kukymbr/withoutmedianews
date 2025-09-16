package server

import "github.com/kukymbr/withoutmedianews/internal/api/http/controllers"

type Server struct {
	*controllers.NewsController
	*controllers.CategoriesController
	*controllers.TagsController
}
