package fxutil

import (
	"reflect"

	"github.com/rs/zerolog/log"
	"go.uber.org/fx/fxevent"
)

func NewLogger() fxevent.Logger {
	return &Logger{}
}

type Logger struct{}

func (l *Logger) LogEvent(event fxevent.Event) {
	llog := log.With().Str("system", "fx").Str("event", reflect.TypeOf(event).String()).Logger()
	switch e := event.(type) {
	case *fxevent.Provided:
		llog.Info().Strs("provided_types", e.OutputTypeNames).Send()
	default:
		llog.Info().Send()
	}
}
