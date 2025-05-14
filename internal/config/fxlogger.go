package config

import (
	"github.com/rs/zerolog"
	"go.uber.org/fx/fxevent"
)

func NewFxZerolog(logger zerolog.Logger) fxevent.Logger {
	return &fxZerolog{logger: logger}
}

type fxZerolog struct {
	logger zerolog.Logger
}

func (l *fxZerolog) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			l.logger.Error().Err(e.Err).Msg("OnStart failed")
		}
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			l.logger.Error().Err(e.Err).Msg("OnStop failed")
		}
	case *fxevent.Provided:
		if e.Err != nil {
			l.logger.Error().Err(e.Err).Msg("Provided failed")
		}
	case *fxevent.Invoked:
		if e.Err != nil {
			l.logger.Error().Err(e.Err).Msg("Invoked failed")
		}
	case *fxevent.Supplied:
		if e.Err != nil {
			l.logger.Error().Err(e.Err).Msg("Supplied failed")
		}
	}
}
