package main

import (
	"context"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
	"log"
	"net/http"
	"os"
	"serenitylabs.cloud/slapi"
	"serenitylabs.cloud/slapi/pkg/api"
	"serenitylabs.cloud/slapi/pkg/fxutil"
	"serenitylabs.cloud/slapi/pkg/ghinternal"
)

func main() {
	log.SetOutput(zerolog.ConsoleWriter{Out: os.Stderr})

	app := fx.New(
		fx.Provide(
			slapi.NewRouter,
			api.GetDefaultConfig,
			ghinternal.Provider,
		),
		fx.Invoke(SetupServer),
		fx.WithLogger(fxutil.NewLogger),
	)

	if app.Err() != nil {
		panic(app.Err())
	}

	app.Run()
}

func SetupServer(lc fx.Lifecycle, router *slapi.Core, config *api.Config) *http.Server {
	server := &http.Server{
		Handler: router,
		Addr:    config.ListenAddr,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				err := server.ListenAndServe()
				if err != nil && err != http.ErrServerClosed {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})

	return server
}
