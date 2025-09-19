package controller

import (
	"net/http"

	apihttp "github.com/kukymbr/withoutmedianews/internal/api/http"
	"github.com/kukymbr/withoutmedianews/internal/domain"
	"github.com/labstack/echo/v4"
)

func NewCategoriesController(service *domain.Service) *CategoriesController {
	return &CategoriesController{
		service: service,
	}
}

type CategoriesController struct {
	service *domain.Service
}

func (ctrl *CategoriesController) GetCategories(c echo.Context) error {
	categories, err := ctrl.service.GetCategories(c.Request().Context())
	if err != nil {
		return err
	}

	resp := apihttp.NewCategories(categories)

	return c.JSON(http.StatusOK, resp)
}
