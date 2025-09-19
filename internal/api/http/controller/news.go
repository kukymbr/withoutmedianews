package controller

import (
	"net/http"

	"github.com/kukymbr/withoutmedianews/internal/api/http"
	"github.com/kukymbr/withoutmedianews/internal/domain"
	"github.com/labstack/echo/v4"
)

func NewNewsController(service *domain.Service) *NewsController {
	return &NewsController{
		service: service,
	}
}

type NewsController struct {
	service *domain.Service
}

func (ctrl *NewsController) GetNewses(c echo.Context) error {
	var req apihttp.NewsListReq
	if err := c.Bind(&req); err != nil {
		return err
	}

	items, err := ctrl.service.GetList(
		c.Request().Context(),
		domain.NewNewsesFilter(req.CategoryID, req.TagID),
		req.Page,
		req.PerPage,
	)
	if err != nil {
		return err
	}

	// TODO: switch to summary
	resp := apihttp.NewNewsList(items)

	return c.JSON(http.StatusOK, resp)
}

func (ctrl *NewsController) GetNews(c echo.Context) error {
	var id int
	if err := echo.PathParamsBinder(c).Int("id", &id).BindError(); err != nil {
		return err
	}

	item, err := ctrl.service.GetNews(c.Request().Context(), id)
	if err != nil {
		return err
	}

	resp := apihttp.NewNews(&item)

	return c.JSON(http.StatusOK, resp)
}

func (ctrl *NewsController) GetNewsCount(c echo.Context) error {
	var req apihttp.NewsListReq
	if err := c.Bind(&req); err != nil {
		return err
	}

	count, err := ctrl.service.GetCount(
		c.Request().Context(),
		domain.NewNewsesFilter(req.CategoryID, req.TagID),
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, apihttp.NewsCountResponse{Count: count})
}
