package modules

import (
	"go.uber.org/fx"

	"github.com/michelsazevedo/hawkeye/api"
	"github.com/michelsazevedo/hawkeye/domain"
	"github.com/michelsazevedo/hawkeye/repository"
)

func Modules() fx.Option {
	return fx.Options(
		fx.Provide(
			repository.NewCourseRepository,
			fx.Annotate(
				domain.NewSearchService[domain.Course],
				fx.As(new(domain.SearchService[domain.Course])),
			),
			fx.Annotate(
				api.NewSearchHandler[domain.Course],
				fx.As(new(api.SearchHandler[domain.Course])),
			),
		),
	)
}
