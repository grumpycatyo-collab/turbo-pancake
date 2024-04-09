package main

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/ardanlabs/conf"
	"github.com/fasthttp/router"
	"github.com/grumpycatyo-collab/turbo-pancake/app/service/handlers"
	"github.com/grumpycatyo-collab/turbo-pancake/business/data/dbschema"
	"github.com/grumpycatyo-collab/turbo-pancake/business/sys/database"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var build = "develop"

func main() {
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	if err := run(&log); err != nil {
		fmt.Printf("\nStartup error: %w \n", err)
		os.Exit(1)
	}
}

func run(log *zerolog.Logger) error {

	cfg := struct {
		conf.Version
		Web struct {
			APIPort         string        `conf:"default::8080"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			IdleTimeout     time.Duration `conf:"default:120s"`
			ShutdownTimeout time.Duration `conf:"default:20s,mask"`
		}
		DB struct {
			User       string `conf:"default:admin"`
			Password   string `conf:"default:admin,mask"`
			Host       string `conf:"default:localhost"`
			Name       string `conf:"default:db"`
			DisableTLS bool   `conf:"default:true"`
		}
	}{
		Version: conf.Version{
			SVN: build,
		},
	}

	const prefix = "SOURCES"
	help, err := conf.ParseOSArgs(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	log.Info().Msg("Connecting to MariaDB")
	db, err := database.Open(database.Config{
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		Host:     cfg.DB.Host,
		Name:     cfg.DB.Name,
	})
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}
	defer db.Close()

	if err := dbschema.InitDB(log, db); err != nil {
		return fmt.Errorf("DB initialization: %w", err)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	r := router.New()

	handlers.Handlers(r, db, log)
	server := &fasthttp.Server{
		Handler:      r.Handler,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Info().Msgf("FastHTTP is initializing on port%s", cfg.Web.APIPort)
		serverErrors <- server.ListenAndServe(cfg.Web.APIPort)
	}()

	// Graceful shutdown in case of multiple services
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Info().Msgf("Shutdown started with signal %s", sig)
		defer log.Info().Msgf("Shutdown completed with signal %s", sig)

		if err := server.Shutdown(); err != nil {
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
