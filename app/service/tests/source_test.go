package tests

import (
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
	"sync"
	"syscall"
	"testing"
	"time"
)

var build = "develop"

func BenchmarkSourceCampaigns(b *testing.B) {
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

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
		}
		log.Error().Msgf("parsing config: %v", err)
	}

	db, err := database.Open(database.Config{
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		Host:     cfg.DB.Host,
		Name:     cfg.DB.Name,
	})
	if err != nil {
		log.Error().Msgf("connecting to db: %v", err)
	}
	defer db.Close()

	if err := dbschema.InitDB(&log, db); err != nil {
		log.Error().Msgf("DB initialization: %v", err)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	r := router.New()

	handlers.Handlers(r, db, &log)
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

	concurrency := 5
	responseTimes := make([]time.Duration, 0, b.N*concurrency)

	// =================================================================================================================
	// Starting benchmarking
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(concurrency)

		for c := 0; c < concurrency; c++ {
			go func() {
				startTime := time.Now()

				ctx := &fasthttp.RequestCtx{}
				ctx.Request.SetRequestURI("/v1/source/1/campaigns")
				r.Handler(ctx)

				elapsed := time.Since(startTime)
				responseTimes = append(responseTimes, elapsed)

				wg.Done()
			}()
		}

		wg.Wait() // Wait for all goroutines to finish
	}

	b.StopTimer()

	var total time.Duration
	var min, max time.Duration
	for i, rt := range responseTimes {
		total += rt
		if i == 0 || rt < min {
			min = rt
		}
		if rt > max {
			max = rt
		}
	}

	avg := total / time.Duration(len(responseTimes))
	b.ReportMetric(float64(avg)/float64(time.Millisecond), "avg_response_time_ms")
	b.ReportMetric(float64(min)/float64(time.Millisecond), "min_response_time_ms")
	b.ReportMetric(float64(max)/float64(time.Millisecond), "max_response_time_ms")

	b.StopTimer()

	log.Info().Msg("Initiating server shutdown after benchmark completion.")
	if err := server.Shutdown(); err != nil {
		log.Error().Msgf("could not stop server gracefully: %v", err)
	} else {
		log.Info().Msg("Server shutdown gracefully.")
	}

	// Ending Benchmarking
	// =================================================================================================================

	select {
	case err := <-serverErrors:
		log.Error().Msgf("server error: %v", err)

	case sig := <-shutdown:
		log.Info().Msgf("Shutdown started with signal %s", sig)
		defer log.Info().Msgf("Shutdown completed with signal %s", sig)

		if err := server.Shutdown(); err != nil {
			log.Error().Msgf("could not stop server gracefully: %v", err)
		}
	}

}
