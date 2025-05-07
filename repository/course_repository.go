package repository

import (
	"context"

	"github.com/michelsazevedo/hawkeye/config"
	"github.com/michelsazevedo/hawkeye/domain"
)

type courseRepository struct {
}

func NewCourseRepository(conf *config.Config) domain.SearchRepository[domain.Course] {
	return &courseRepository{}
}

func (c *courseRepository) Search(ctx context.Context, query string) ([]domain.Course, error) {
	courses := make([]domain.Course, 0)

	return courses, nil
}
