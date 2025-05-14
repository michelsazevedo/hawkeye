package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

type Message struct {
	Subject string
	Data    []byte
}

type MessageHandler func(ctx context.Context, msg Message)

type SubscriberService[T any] interface {
	Listen(subject string) error
}

type BrokerRepository interface {
	Subscribe(subject string, handler MessageHandler) error
}

type subscriberService[T any] struct {
	brokerRepository BrokerRepository
	esRepo           SearchRepository[T]
}

func NewSubscriberService[T any](brokerRepository BrokerRepository, esRepo SearchRepository[T]) SubscriberService[T] {
	return &subscriberService[T]{
		brokerRepository: brokerRepository,
		esRepo:           esRepo,
	}
}

func (s *subscriberService[T]) Listen(subject string) error {
	tracer := otel.Tracer("broker-handler")

	return s.brokerRepository.Subscribe(subject, func(ctx context.Context, msg Message) {
		ctx, span := tracer.Start(ctx, "subscriber.consume.course")
		defer span.End()

		var message T

		if err := json.Unmarshal(msg.Data, &message); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to Unmarshal")
			log.Error().Err(err).Msg("Failed to Unmarshal Message")
			return
		}

		ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
		defer cancel()

		if err := s.esRepo.Index(ctx, message); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to Index Document")
			log.Error().Err(err).Msg("Failed to index course")
		}
	})
}
