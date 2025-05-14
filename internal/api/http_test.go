package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/michelsazevedo/hawkeye/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	data        = domain.Course{Id: "f81d4fae-7dec-11d0-a765-00a0c91e6bf6", Name: "Go land", Headline: "Learn Go"}
	mockService = new(domain.MockSearchService[domain.Course])
	handlerApi  = NewSearchHandler(mockService)
	e           = echo.New()
)

func TestIndexHandler(t *testing.T) {
	t.Run("Returns Status Code 200", func(t *testing.T) {
		mockService.On("Search", mock.Anything, mock.AnythingOfType("string")).Return([]domain.Course{data}, nil)

		req := httptest.NewRequest(http.MethodGet, "/courses/search", nil)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, handlerApi.Search(c)) {
			if status := rec.Code; status != http.StatusOK {
				t.Errorf("Handler returned wrong status code. Expected: %d. Got: %d.",
					http.StatusOK, status)
			}
		}

		var httpResponse ApiCollectionResponse[domain.Course]

		if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &httpResponse)) {
			assert.Equal(t, len(httpResponse.Data), 1)
			assert.Nil(t, httpResponse.Errors)
		}
	})
}
