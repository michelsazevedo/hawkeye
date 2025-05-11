package nats

import (
	"github.com/michelsazevedo/hawkeye/config"
	"github.com/michelsazevedo/hawkeye/domain"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
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
		handler(domain.Message{
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
