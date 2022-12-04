package cmd

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	pgx "github.com/jackc/pgx/v4"
	"github.com/nikprim/otus-homework/hw12_13_14_15_calendar/cmd/config"
	"github.com/nikprim/otus-homework/hw12_13_14_15_calendar/internal/app"
	internalhttp "github.com/nikprim/otus-homework/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/nikprim/otus-homework/hw12_13_14_15_calendar/internal/storage/memory"
	psqlstorage "github.com/nikprim/otus-homework/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func serveHTTPCommand(ctx context.Context) *cobra.Command {
	command := &cobra.Command{
		Use:   "serve-http",
		Short: "serves http api",
		RunE:  serveHTTPCommandRunE(ctx),
	}

	command.Flags().StringVar(&cfgFile, "config", "", "Path to configuration file")

	err := command.MarkFlagRequired("config")
	if err != nil {
		return nil
	}

	return command
}

func serveHTTPCommandRunE(ctx context.Context) func(cmd *cobra.Command, args []string) (err error) {
	return func(cmd *cobra.Command, args []string) (err error) {
		configFile := cmd.Flag("config").Value.String()

		cfg, err := config.ParseConfig(configFile)
		if err != nil {
			log.Error().Err(err).Msg("failed to parse config")

			return err
		}

		logLevel, err := zerolog.ParseLevel(cfg.Logger.Level)
		if err != nil {
			log.Error().Err(err).Msg("failed to install log level")

			return err
		}

		zerolog.SetGlobalLevel(logLevel)

		var store app.Storage

		switch cfg.DB.Type {
		case "psql":
			conn, err := pgx.Connect(ctx, cfg.DB.PSQL.URL)
			if err != nil {
				log.Error().Err(err).Msg("unable to connect to database")

				return err
			}

			defer func() {
				err := conn.Close(ctx)
				if err != nil {
					log.Error().Err(err).Msg("unable to close connect to database")
				}
			}()

			err = conn.Ping(ctx)
			if err != nil {
				log.Error().Err(err).Msg("unable to connect to database")

				return err
			}

			store = psqlstorage.New(conn)
		case "memory":
			store = memorystorage.New()
		default:
			err := errors.New("unknown db type")
			log.Error().Err(err).Send()

			return err
		}

		application := app.New(store)
		server := internalhttp.NewServer(cfg.HTTP.Host, cfg.HTTP.Port, application)

		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		defer cancel()

		go func() {
			<-ctx.Done()

			log.Info().Msg("stopping an http server...")

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
			defer cancel()

			if err := server.Stop(ctx); err != nil {
				log.Error().Err(err).Msg("failed to stop http server")
			}
		}()

		log.Info().Msg("calendar is running...")

		if err := server.Start(); err != nil {
			cancel()

			if !errors.Is(err, http.ErrServerClosed) {
				log.Error().Err(err).Msg("failed to start http server")

				return err
			}
		}

		return nil
	}
}
