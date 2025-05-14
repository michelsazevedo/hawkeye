package modules

import (
	"context"

	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"

	"github.com/michelsazevedo/hawkeye/internal/api"
	"github.com/michelsazevedo/hawkeye/internal/domain"
	"github.com/michelsazevedo/hawkeye/internal/repository"
	"github.com/michelsazevedo/hawkeye/internal/repository/nats"
	"github.com/michelsazevedo/hawkeye/pkg/observability"
	"github.com/rs/zerolog/log"
)

func Modules() fx.Option {
	return fx.Options(
		fx.Provide(
			observability.NewTracerProvider,
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
		fx.Invoke(func(lc fx.Lifecycle, tracer *sdktrace.TracerProvider) {
			otel.SetTracerProvider(tracer)

			lc.Append(fx.Hook{
				OnStop: func(ctx context.Context) error {
					log.Info().Msg("Shutting down OpenTelemetry tracer provider")
					return tracer.Shutdown(ctx)
				},
			})
		}),
	)
}
