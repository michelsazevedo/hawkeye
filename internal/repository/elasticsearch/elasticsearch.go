package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/michelsazevedo/hawkeye/internal/domain"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var (
	es   *elasticsearch.Client
	once sync.Once
	err  error
)

type esResponse[T domain.Searchable] struct {
	Hits struct {
		Hits []struct {
			Source T `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

type ElasticSearchRepository[T domain.Searchable] struct {
	client *elasticsearch.Client
}

func GetElasticSearchConnection[T domain.Searchable](URL []string) *ElasticSearchRepository[T] {
	once.Do(func() {
		es, err = elasticsearch.NewClient(elasticsearch.Config{Addresses: URL})
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create Elasticsearch client")
		}

		res, err := es.Info()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to connect to Elasticsearch")
		}
		defer res.Body.Close()

		log.Info().Msg("Connected to Elasticsearch")
	})

	return &ElasticSearchRepository[T]{client: es}
}

func (es *ElasticSearchRepository[T]) IndexExists(ctx context.Context, name string) bool {
	res, err := es.client.Indices.Exists([]string{name}, es.client.Indices.Exists.WithContext(ctx))
	if err != nil {
		log.Error().Err(err).Msgf("Failed to check if index %s exists", name)
		return false
	}
	defer res.Body.Close()

	return res.StatusCode == 200
}

func (es *ElasticSearchRepository[T]) CreateIndex(ctx context.Context, name string, mapping []byte) error {
	res, err := es.client.Indices.Create(
		name,
		es.client.Indices.Create.WithBody(bytes.NewReader(mapping)),
		es.client.Indices.Create.WithContext(ctx),
	)

	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	if res.IsError() {
		return fmt.Errorf("error creating index: %s", body)
	}

	log.Info().Str("index", name).Msg("Index created successfully")
	return nil
}

func (es *ElasticSearchRepository[T]) Index(ctx context.Context, index string, doc T) error {
	tracer := otel.Tracer("elasticsearch.repository")
	ctx, span := tracer.Start(ctx, "elasticsearch.index.course",
		trace.WithAttributes(
			attribute.String("es.index", index),
			attribute.String("document.id", doc.GetID()),
		),
	)
	defer span.End()

	body, err := json.Marshal(doc)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to marshal object")
		log.Error().Err(err).Msg("Failed to marshal object for indexing")
		return err
	}

	res, err := es.client.Index(
		index,
		bytes.NewReader(body),
		es.client.Index.WithDocumentID(doc.GetID()),
		es.client.Index.WithContext(ctx),
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to index Document")
		log.Error().Err(err).Msg("Failed to index Document")
		return err
	}

	defer res.Body.Close()
	log.Info().Msgf("Document indexed: %s", doc.GetID())

	return nil
}

func (es *ElasticSearchRepository[T]) Search(ctx context.Context, esQuery any, index string) ([]T, error) {
	body, err := json.Marshal(esQuery)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal Elasticsearch query")
		return nil, err
	}

	res, err := es.client.Search(
		es.client.Search.WithIndex(index),
		es.client.Search.WithBody(bytes.NewReader(body)),
		es.client.Search.WithContext(ctx),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to search Elasticsearch")
		return nil, err
	}
	defer res.Body.Close()

	return decodeResponse[T](res)
}

func decodeResponse[T domain.Searchable](res *esapi.Response) ([]T, error) {
	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("elasticsearch returned error on search: %s", body)
	}

	var response esResponse[T]

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %s", err)
	}

	var results []T

	for _, hit := range response.Hits.Hits {
		results = append(results, hit.Source)
	}

	return results, nil
}
