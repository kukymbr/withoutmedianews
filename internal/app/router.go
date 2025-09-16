package app

import (
	"net/http"

	apihttp "github.com/kukymbr/withoutmedianews/internal/api/http"
	"github.com/kukymbr/withoutmedianews/internal/api/http/server"
)

func initRouter(server apihttp.StrictServerInterface, errResponder *server.ErrorResponder) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /openapi.yaml", func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Content-Type", "text/yaml")
		resp.WriteHeader(http.StatusOK)

		_, _ = resp.Write(apihttp.GetOpenAPIContent())
	})

	return apihttp.HandlerFromMuxWithBaseURL(
		apihttp.NewStrictHandlerWithOptions(server, nil, apihttp.StrictHTTPServerOptions{
			ResponseErrorHandlerFunc: errResponder.APIError,
		}),
		mux,
		"/api/v1",
	)
}
