package domain

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockSearchService[T any] struct {
	mock.Mock
}

func (m *MockSearchService[T]) Index(ctx context.Context) ([]T, error) {
	args := m.Called(ctx)
	return args.Get(0).([]T), args.Error(1)
}

func (m *MockSearchService[T]) Search(ctx context.Context, query string) ([]T, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]T), args.Error(1)
}
