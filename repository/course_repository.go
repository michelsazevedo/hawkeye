package repository

import (
	"context"

	"github.com/michelsazevedo/hawkeye/config"
	"github.com/michelsazevedo/hawkeye/domain"
	es "github.com/michelsazevedo/hawkeye/repository/elasticsearch"
	"github.com/rs/zerolog/log"
)

var index string = "courses"

type courseRepository struct {
	esClient *es.ElasticSearchRepository[*domain.Course]
}

func NewCourseRepository(conf *config.Config) domain.SearchRepository[domain.Course] {
	esClient := es.GetElasticSearchConnection[*domain.Course](conf.GetElasticsearchUrls())
	ctx := context.Background()

	if !esClient.IndexExists(ctx, index) {
		properties, _ := conf.GetElasticSearchIndex(index)

		if err := esClient.CreateIndex(ctx, index, properties); err != nil {
			log.Error().Err(err).Msgf("Error to create %s index", index)
		}
	}

	return &courseRepository{esClient: esClient}
}

func (c *courseRepository) Index(ctx context.Context, course domain.Course) error {
	return c.esClient.Index(ctx, index, &course)
}

func (c *courseRepository) Search(ctx context.Context, query string) ([]domain.Course, error) {
	esQuery := es.NewQuery(query).
		AddField("name").
		AddField("content").
		Build()

	results, err := c.esClient.Search(ctx, esQuery, index)
	if err != nil {
		return nil, err
	}

	courses := make([]domain.Course, len(results))
	for i, r := range results {
		courses[i] = *r
	}
	return courses, nil
}
