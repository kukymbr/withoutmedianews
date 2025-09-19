package controller

import (
	"net/http"

	apihttp "github.com/kukymbr/withoutmedianews/internal/api/http"
	"github.com/kukymbr/withoutmedianews/internal/domain"
	"github.com/labstack/echo/v4"
)

func NewTagsController(service *domain.Service) *TagsController {
	return &TagsController{
		service: service,
	}
}

type TagsController struct {
	service *domain.Service
}

func (ctrl *TagsController) GetTags(c echo.Context) error {
	tags, err := ctrl.service.GetTags(c.Request().Context())
	if err != nil {
		return err
	}

	resp := apihttp.NewTags(tags)

	return c.JSON(http.StatusOK, resp)
}
