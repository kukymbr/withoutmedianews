package server

import "github.com/kukymbr/withoutmedianews/internal/api/http/controller"

type Server struct {
	*controller.NewsController
	*controller.CategoriesController
	*controller.TagsController
}
