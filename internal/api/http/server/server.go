package server

import (
	"github.com/kukymbr/withoutmedianews/internal/api/http/controller"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(srv *Server) *echo.Echo {
	router := echo.New()

	router.Use(middleware.Logger())
	router.Use(middleware.Recover())

	router.HTTPErrorHandler = srv.errResponder.APIError

	v1 := router.Group("/api/v1")

	v1.GET("/newses", srv.news.GetNewses)
	v1.GET("/newses/count", srv.news.GetNewsCount)
	v1.GET("/newses/news/:id", srv.news.GetNews)

	v1.GET("/categories", srv.categories.GetCategories)
	v1.GET("/tags", srv.tags.GetTags)

	router.GET("/openapi.yaml", handleSpecRequest)

	return router
}

func New(
	news *controller.NewsController,
	categories *controller.CategoriesController,
	tags *controller.TagsController,
) *Server {
	return &Server{
		news:       news,
		categories: categories,
		tags:       tags,
	}
}

type Server struct {
	news       *controller.NewsController
	categories *controller.CategoriesController
	tags       *controller.TagsController

	errResponder ErrorResponder
}
