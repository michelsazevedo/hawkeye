package api

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/michelsazevedo/hawkeye/domain"
)

type ApiCollectionResponse[T any] struct {
	Data   []T      `json:"data"`
	Errors []string `json:"errors"`
}

type SearchHandler[T any] interface {
	Search(c echo.Context) error
}

type handler[T any] struct {
	searchService domain.SearchService[T]
}

func NewSearchHandler[T any](searchService domain.SearchService[T]) SearchHandler[T] {
	return &handler[T]{searchService: searchService}
}

func (h *handler[T]) Search(c echo.Context) error {
	results, err := h.searchService.Search(c.Request().Context(), c.QueryParam("q"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ApiCollectionResponse[T]{Errors: errData(err)})
	}

	return c.JSON(http.StatusOK, &ApiCollectionResponse[T]{Data: results})
}

func errData(err error) []string {
	messages := strings.Split(err.Error(), ".")
	errMsg := make([]string, 0, len(messages))

	for _, message := range messages {
		msg := strings.TrimSpace(message)
		if msg != "" {
			errMsg = append(errMsg, msg)
		}
	}
	return errMsg
}
