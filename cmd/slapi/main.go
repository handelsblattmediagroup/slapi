package main

import (
	"os"

	"github.com/acuteaura/slapi/pkg/core"
	"github.com/acuteaura/slapi/pkg/routers/helloworld"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
	"golang.org/x/term"
)

func main() {
	if term.IsTerminal(int(os.Stderr.Fd())) {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		log.Info().Msg("detected TTY, using fancy color output")
	}

	opts := make([]fx.Option, 0)

	opts = append(opts,
		fx.WithLogger(core.NewLogger),

		fx.Provide(core.New),
		fx.Provide(core.GetConfigDefaults),
		fx.Provide(core.NewTracer),

		fx.Provide(helloworld.Provider),

		fx.Invoke(core.SetupServer),
	)

	fxApp := fx.New(
		opts...,
	)
	if fxApp.Err() != nil {
		panic(fxApp.Err())
	}

	fxApp.Run()
}
