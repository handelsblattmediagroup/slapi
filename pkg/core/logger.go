package core

import (
	"reflect"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx/fxevent"
)

func NewLogger() fxevent.Logger {
	return &Logger{zerolog.DebugLevel}
}

type Logger struct {
	level zerolog.Level
}

func (l *Logger) LogEvent(event fxevent.Event) {
	llog := log.With().Str("system", "fx").Str("event", reflect.TypeOf(event).String()).Logger()
	switch e := event.(type) {
	case *fxevent.Provided:
		llog.WithLevel(l.level).Strs("provided_types", e.OutputTypeNames).Send()
	default:
		llog.WithLevel(l.level).Send()
	}
}
