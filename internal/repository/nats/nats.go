package nats

import (
	"context"

	"github.com/michelsazevedo/hawkeye/internal/config"
	"github.com/michelsazevedo/hawkeye/internal/domain"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type natsRepository struct {
	natsConn *nats.Conn
}

func NewNatsRepository(conf *config.Config) domain.BrokerRepository {
	conn, err := nats.Connect(conf.NATS.URL)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to connect to NATS")
		return nil
	}

	log.Info().Msg("Connected to NATS")
	return &natsRepository{natsConn: conn}
}

func (n *natsRepository) Subscribe(subject string, handler domain.MessageHandler) error {
	_, err := n.natsConn.Subscribe(subject, func(message *nats.Msg) {
		ctx := context.Background()
		tracer := otel.Tracer("nats.subscribe")

		ctx, span := tracer.Start(ctx, "nats.receive.message",
			trace.WithAttributes(
				attribute.String("messaging.system", "nats"),
				attribute.String("messaging.destination", subject),
				attribute.String("messaging.destination_kind", "topic"),
				attribute.String("messaging.protocol", "nats"),
				attribute.String("messaging.operation", "receive"),
			),
		)
		defer span.End()

		handler(ctx, domain.Message{
			Subject: message.Subject,
			Data:    message.Data,
		})
	})
	if err != nil {
		log.Error().Err(err).Msgf("Failed to subscribe to subject: %s", subject)
		return err
	}

	return nil
}
