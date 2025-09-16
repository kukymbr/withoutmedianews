package logkit

import (
	"net/http"

	"go.uber.org/zap"
)

const (
	// HTTP request & response
	fieldNameStatus    = "status"
	fieldNameMethod    = "method"
	fieldNamePath      = "path"
	fieldNameQuery     = "query"
	fieldNameIP        = "ip"
	fieldNameUserAgent = "user-agent"
)

func GetHTTPRequestFields(req *http.Request) []zap.Field {
	if req == nil {
		return nil
	}

	return []zap.Field{
		zap.String(fieldNameMethod, req.Method),
		zap.String(fieldNamePath, req.URL.Path),
		zap.String(fieldNameQuery, req.URL.RawQuery),
		zap.String(fieldNameIP, req.RemoteAddr),
		zap.String(fieldNameUserAgent, req.UserAgent()),
	}
}

func GetHTTPResponseFields(status int) []zap.Field {
	return []zap.Field{
		zap.Int(fieldNameStatus, status),
	}
}

func WithHTTPRequestFields(logger *zap.Logger, req *http.Request) *zap.Logger {
	fields := GetHTTPRequestFields(req)

	return logger.With(fields...)
}

func WithHTTPResponseFields(logger *zap.Logger, statusCode int) *zap.Logger {
	fields := GetHTTPResponseFields(statusCode)

	return logger.With(fields...)
}
