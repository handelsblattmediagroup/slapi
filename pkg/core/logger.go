package core

import (
	"os"
	"reflect"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx/fxevent"
)

func SetupDefaultLoglevel() {
	logLevelRaw, ok := os.LookupEnv("SLAPI_LOG_LEVEL")
	if ok {
		logLevel, err := zerolog.ParseLevel(logLevelRaw)
		if err != nil {
			panic("bad log level")
		}
		log.Logger = log.Level(logLevel)
	} else {
		log.Logger = log.Level(zerolog.InfoLevel)
		log.Info().Msg("no loglevel defined, defaulting to INFO (set SLAPI_LOG_LEVEL)")
	}
}

func NewFxLogAdapter() fxevent.Logger {
	return &FxLogAdapter{zerolog.DebugLevel}
}

type FxLogAdapter struct {
	level zerolog.Level
}

func (l *FxLogAdapter) LogEvent(event fxevent.Event) {
	llog := log.With().Str("system", "fx").Str("event", reflect.TypeOf(event).String()).Logger()
	switch e := event.(type) {
	case *fxevent.Provided:
		llog.WithLevel(l.level).Strs("provided_types", e.OutputTypeNames).Send()
	default:
		llog.WithLevel(l.level).Send()
	}
}
