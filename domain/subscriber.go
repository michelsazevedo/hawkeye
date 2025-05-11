package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rs/zerolog/log"
)

type Message struct {
	Subject string
	Data    []byte
}

type MessageHandler func(msg Message)

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
	return s.brokerRepository.Subscribe(subject, func(msg Message) {
		var message T

		if err := json.Unmarshal(msg.Data, &message); err != nil {
			log.Error().Err(err).Msg("Failed to unmarshal message")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := s.esRepo.Index(ctx, message); err != nil {
			log.Error().Err(err).Msg("Failed to index course")
		}
	})
}
