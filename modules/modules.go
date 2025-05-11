package modules

import (
	"context"

	"go.uber.org/fx"

	"github.com/michelsazevedo/hawkeye/api"
	"github.com/michelsazevedo/hawkeye/domain"
	"github.com/michelsazevedo/hawkeye/repository"
	"github.com/michelsazevedo/hawkeye/repository/nats"
	"github.com/rs/zerolog/log"
)

func Modules() fx.Option {
	return fx.Options(
		fx.Provide(
			nats.NewNatsRepository,
			repository.NewCourseRepository,
			fx.Annotate(
				domain.NewSubscriberService[domain.Course],
				fx.As(new(domain.SubscriberService[domain.Course])),
			),
			fx.Annotate(
				domain.NewSearchService[domain.Course],
				fx.As(new(domain.SearchService[domain.Course])),
			),
			fx.Annotate(
				api.NewSearchHandler[domain.Course],
				fx.As(new(api.SearchHandler[domain.Course])),
			),
		),
		fx.Invoke(func(lc fx.Lifecycle, s domain.SubscriberService[domain.Course]) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						if err := s.Listen("courses.created"); err != nil {
							log.Error().Err(err).Msg("Failed to listen to NATS subject: courses.created")
						}
					}()
					return nil
				},
			})
		}),
	)
}
