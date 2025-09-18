package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-pg/pg/v10"
	"github.com/kukymbr/withoutmedianews/internal/api/http"
	"github.com/kukymbr/withoutmedianews/internal/pkg/dbkit"
	"github.com/kukymbr/withoutmedianews/internal/pkg/logkit"
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

func (r *ErrorResponder) APIError(resp http.ResponseWriter, req *http.Request, err error) {
	errData := apihttp.APIError{
		Message: err.Error(),
	}

	encoded, encodeErr := json.Marshal(errData)
	if encodeErr != nil {
		r.logger.Error("encoding error response", zap.Error(encodeErr))

		r.PlainText(resp, req, encodeErr)

		return
	}

	resp.Header().Add("Content-Type", "application/json")

	r.respond(resp, req, err, encoded)
}

func (r *ErrorResponder) PlainText(resp http.ResponseWriter, req *http.Request, err error) {
	code := getErrorCode(err)
	msg := http.StatusText(code) + ": " + err.Error()

	resp.Header().Add("Content-Type", "text/plain")

	r.respond(resp, req, err, []byte(msg))
}

func (r *ErrorResponder) respond(resp http.ResponseWriter, req *http.Request, err error, content []byte) {
	code := getErrorCode(err)

	logger := logkit.WithHTTPRequestFields(r.logger, req)
	logger = logkit.WithHTTPResponseFields(logger, code)

	logger.Warn("responding with error", zap.Error(err))

	resp.WriteHeader(code)

	_, _ = resp.Write(content)
}

func getErrorCode(err error) int {
	switch {
	case errors.Is(err, dbkit.ErrNotFound) || errors.Is(err, pg.ErrNoRows):
		return http.StatusNotFound
	case errors.Is(err, pg.ErrMultiRows):
		return http.StatusUnprocessableEntity
	}

	return http.StatusInternalServerError
}
