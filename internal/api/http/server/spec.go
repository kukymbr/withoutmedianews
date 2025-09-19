package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func handleSpecRequest(c echo.Context) error {
	spec, err := rawSpec()
	if err != nil {
		return err
	}

	return c.Blob(http.StatusOK, "application/json", spec)
}
