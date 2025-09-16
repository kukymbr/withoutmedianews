//go:build api_test

package apihttptests

import (
	"net/http"
	"testing"
)

func (s *HTTPTestSuite) TestGetNews() {
	client := s.getClient()

	tests := []struct {
		Name   string
		Params *GetNewsParams
		Assert func(resp *GetNewsResponse, err error)
	}{
		{
			Name:   "No filters",
			Params: &GetNewsParams{},
			Assert: func(resp *GetNewsResponse, err error) {
				s.Require().NoError(err)

				s.Require().Equal(http.StatusOK, resp.StatusCode())
				s.Require().Len(*resp.JSON200, 1)
				s.assertNewsItem(givenNewsItem1(), (*resp.JSON200)[0])
			},
		},
		{
			Name:   "With category filter",
			Params: &GetNewsParams{CategoryID: 1},
			Assert: func(resp *GetNewsResponse, err error) {
				s.Require().NoError(err)

				s.Require().Equal(http.StatusOK, resp.StatusCode())
				s.Require().Len(*resp.JSON200, 1)
				s.assertNewsItem(givenNewsItem1(), (*resp.JSON200)[0])
			},
		},
		{
			Name:   "With tag filter",
			Params: &GetNewsParams{TagID: 1},
			Assert: func(resp *GetNewsResponse, err error) {
				s.Require().NoError(err)

				s.Require().Equal(http.StatusOK, resp.StatusCode())
				s.Require().Len(*resp.JSON200, 1)
				s.assertNewsItem(givenNewsItem1(), (*resp.JSON200)[0])
			},
		},
		{
			Name:   "With category & tag filters both",
			Params: &GetNewsParams{CategoryID: 1, TagID: 1},
			Assert: func(resp *GetNewsResponse, err error) {
				s.Require().NoError(err)

				s.Require().Equal(http.StatusOK, resp.StatusCode())
				s.Require().Len(*resp.JSON200, 1)
				s.assertNewsItem(givenNewsItem1(), (*resp.JSON200)[0])
			},
		},
		{
			Name:   "With unknown category filter",
			Params: &GetNewsParams{CategoryID: 100},
			Assert: func(resp *GetNewsResponse, err error) {
				s.Require().NoError(err)

				s.Require().Equal(http.StatusOK, resp.StatusCode())
				s.Require().Len(*resp.JSON200, 0)
			},
		},
		{
			Name:   "With unknown tag filter",
			Params: &GetNewsParams{TagID: 100},
			Assert: func(resp *GetNewsResponse, err error) {
				s.Require().NoError(err)

				s.Require().Equal(http.StatusOK, resp.StatusCode())
				s.Require().Len(*resp.JSON200, 0)
			},
		},
		{
			Name:   "With empty page",
			Params: &GetNewsParams{Page: 100},
			Assert: func(resp *GetNewsResponse, err error) {
				s.Require().NoError(err)

				s.Require().Equal(http.StatusOK, resp.StatusCode())
				s.Require().Len(*resp.JSON200, 0)
			},
		},
	}

	for _, test := range tests {
		s.T().Run(test.Name, func(t *testing.T) {
			resp, err := client.GetNewsWithResponse(s.T().Context(), test.Params)

			test.Assert(resp, err)
		})
	}
}

func (s *HTTPTestSuite) TestGetNewsItem() {
	client := s.getClient()

	tests := []struct {
		Name   string
		ID     NumericID
		Assert func(resp *GetNewsItemResponse, err error)
	}{
		{
			Name: "Known item",
			ID:   givenNewsItem1().ID,
			Assert: func(resp *GetNewsItemResponse, err error) {
				s.Require().NoError(err)

				s.Require().Equal(http.StatusOK, resp.StatusCode())
				s.Require().NotNil(resp.JSON200)
				s.assertNewsItem(givenNewsItem1(), *resp.JSON200)
			},
		},
		{
			Name: "Unknown item",
			ID:   100,
			Assert: func(resp *GetNewsItemResponse, err error) {
				s.Require().NoError(err)

				s.Require().Equal(http.StatusNotFound, resp.StatusCode())
			},
		},
	}

	for _, test := range tests {
		s.T().Run(test.Name, func(t *testing.T) {
			resp, err := client.GetNewsItemWithResponse(s.T().Context(), test.ID)

			test.Assert(resp, err)
		})
	}
}
