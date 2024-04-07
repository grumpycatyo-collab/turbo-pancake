package main

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/ardanlabs/conf"
	"github.com/grumpycatyo-collab/turbo-pancake/business/sys/database"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
	"os"
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

func requestHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hello, world!\n")
}

func run(log *zerolog.Logger) error {

	cfg := struct {
		conf.Version
		Web struct {
			APIPort      string        `conf:"default::8080"`
			ReadTimeout  time.Duration `conf:"default:5s"`
			WriteTimeout time.Duration `conf:"default:10s"`
			IdleTimeout  time.Duration `conf:"default:120s"`
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

	server := &fasthttp.Server{
		Handler:      requestHandler,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
	}

	log.Info().Msgf("FastHTTP is initializing on port%s", cfg.Web.APIPort)
	if err := server.ListenAndServe(cfg.Web.APIPort); err != nil {
		return fmt.Errorf("Error in ListenAndServe: %w\n", err)
	}
	return nil
}
