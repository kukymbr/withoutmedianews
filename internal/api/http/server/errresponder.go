package server

import (
	"errors"
	"net/http"

	"github.com/go-pg/pg/v10"
	"github.com/kukymbr/withoutmedianews/internal/api/http"
	"github.com/kukymbr/withoutmedianews/internal/domain"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func NewErrorResponder(logger *zap.Logger) *ErrorResponder {
	return &ErrorResponder{
		logger: logger,
	}
}

type ErrorResponder struct {
	logger *zap.Logger
}

func (r *ErrorResponder) APIError(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	code := getErrorCode(err)
	errData := apihttp.APIError{
		Message: err.Error(),
	}

	if err := c.JSON(code, &errData); err != nil {
		r.PlainText(err, c)
	}
}

func (r *ErrorResponder) PlainText(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	code := getErrorCode(err)
	msg := http.StatusText(code) + ": " + err.Error()

	_ = c.String(code, msg)
}

func getErrorCode(err error) int {
	switch {
	case errors.Is(err, domain.ErrNotFound) || errors.Is(err, pg.ErrNoRows):
		return http.StatusNotFound
	case errors.Is(err, pg.ErrMultiRows):
		return http.StatusUnprocessableEntity
	}

	return http.StatusInternalServerError
}
