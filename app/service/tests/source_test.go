package tests

import (
	"errors"
	"fmt"
	"github.com/ardanlabs/conf"
	"github.com/grumpycatyo-collab/turbo-pancake/app/service/handlers/sourcegrp"
	"github.com/grumpycatyo-collab/turbo-pancake/business/core/source"
	"github.com/grumpycatyo-collab/turbo-pancake/business/sys/database"
	"github.com/grumpycatyo-collab/turbo-pancake/business/web/mid/cache"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
	"os"
	"sync"
	"testing"
	"time"
)

func InitDependencies() (*sqlx.DB, zerolog.Logger) {
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
	}{}

	const prefix = "SOURCES"
	help, err := conf.ParseOSArgs(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
		}
		fmt.Printf("parsing config: %v", err)
	}

	log.Info().Msg("Connecting to MariaDB")
	db, err := database.Open(database.Config{
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		Host:     cfg.DB.Host,
		Name:     cfg.DB.Name,
	})
	if err != nil {
		fmt.Printf("connecting to db: %v", err)
	}
	defer db.Close()

	return db, log
}

func BenchmarkSourceCampaigns(b *testing.B) {
	db, log := InitDependencies()

	sgh := sourcegrp.Handlers{
		Core: source.NewCore(&log, db),
	}
	var wg sync.WaitGroup

	stop := make(chan struct{})

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-time.Tick(1 * time.Second):
					ctx := &fasthttp.RequestCtx{}
					ctx.SetUserValue("id", "1")
					const cacheDuration = 5 * time.Minute
					cache.Middleware(cacheDuration, sourcegrp.GetSourceCampaigns(&sgh), &log)(ctx)
				case <-stop:
					return
				}
			}
		}()
	}

	go func() {
		time.Sleep(10 * time.Second)
		close(stop)
	}()

	wg.Wait()
}
