//go:build api_test

// Package apihttptests contains API integration tests.
//
// Initialized database with test data inserted is a mandatory requirement
// (see the `db/testdata/testdata.sql` file).
//
// Environment variables are required to establish the database connection:
//
// ```env
//
//	DB_HOST=localhost
//	DB_PORT=5432
//	DB_USERNAME=postgres
//	DB_PASSWORD=postgres
//	DB_DATABASE=postgres
//
// ```
//
// Tests execution example:
//
// ```shell
// DB_PORT=5432 go test -tags api_test ./internal/api/httptests/...
// ```
package apihttptests

import (
	"net/http/httptest"
	"testing"

	"github.com/kukymbr/withoutmedianews/internal/app"
	"github.com/kukymbr/withoutmedianews/internal/config"
	"github.com/kukymbr/withoutmedianews/internal/pkg/logkit"
	"github.com/stretchr/testify/suite"
)

func TestHTTPClientSuite(t *testing.T) {
	suite.Run(t, new(HTTPTestSuite))
}

type HTTPTestSuite struct {
	suite.Suite

	ctn *app.Container
}

func (s *HTTPTestSuite) SetupSuite() {
	conf, err := config.New()
	if err != nil {
		s.T().Fatal(err)
	}

	logger := logkit.NewTestLogger(s.T())

	s.ctn = app.NewContainer(s.T().Context(), conf, logger)
}

func (s *HTTPTestSuite) TearDownSuite() {
	_ = s.ctn.Close()
}

func (s *HTTPTestSuite) getClient() *ClientWithResponses {
	s.T().Helper()

	handler := s.ctn.GetRouter()
	server := httptest.NewServer(handler)

	s.T().Cleanup(func() {
		server.Close()
	})

	client, err := NewClientWithResponses(server.URL + "/api/v1")
	s.Require().NoError(err)

	return client
}

func (s *HTTPTestSuite) assertNews(expected News, item News) {
	s.T().Helper()

	s.Equal(expected.ID, item.ID)
	s.Equal(expected.Author, item.Author)
	s.Equal(expected.Title, item.Title)
	s.Equal(expected.Content, item.Content)
	s.Equal(expected.ShortText, item.ShortText)
	s.Equal(expected.Category, item.Category)
	s.ElementsMatch(expected.Tags, item.Tags)
}
