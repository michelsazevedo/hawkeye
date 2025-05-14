package domain

import (
	"context"
)

type SearchService[T any] interface {
	Index(ctx context.Context, doc T) error
	Search(ctx context.Context, query string) ([]T, error)
}

type SearchRepository[T any] interface {
	Index(ctx context.Context, doc T) error
	Search(ctx context.Context, query string) ([]T, error)
}

type searchService[T any] struct {
	searchRepository SearchRepository[T]
}

func NewSearchService[T any](searchRepository SearchRepository[T]) SearchService[T] {
	return &searchService[T]{searchRepository: searchRepository}
}

func (s *searchService[T]) Index(ctx context.Context, doc T) error {
	return s.searchRepository.Index(ctx, doc)
}

func (s *searchService[T]) Search(ctx context.Context, query string) ([]T, error) {
	return s.searchRepository.Search(ctx, query)
}
