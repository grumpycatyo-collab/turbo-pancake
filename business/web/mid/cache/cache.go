package cache

import (
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
	"sync"
	"time"
)

var (
	cache     sync.Map
	semaphore = make(chan struct{}, 1)
)

// Middleware checks the cache before proceeding to the handler.
func Middleware(expiration time.Duration, next func(ctx *fasthttp.RequestCtx), log *zerolog.Logger) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		key := string(ctx.RequestURI())
		log.Info().Msg("Attempting to retrieve info from cache")
		if val, ok := cache.Load(key); ok {
			item := val.(Item)
			if time.Now().UnixNano() <= item.Expiration {
				ctx.Response.SetBodyRaw(item.Content)
				return
			}
		}

		select {
		case semaphore <- struct{}{}:
			defer func() {
				<-semaphore
			}()
			next(ctx)

			log.Info().Msg("Writing response data into cache")
			data := ctx.Response.Body()
			cache.Store(key, Item{
				Content:    data,
				Expiration: time.Now().Add(expiration).UnixNano(),
			})
		default:
			log.Info().Msg("Another user is accessing the handler, using cache instead")
			ctx.Response.SetBody([]byte("Handler is currently busy, using cache instead"))
		}
	}
}

type Item struct {
	Content    []byte
	Expiration int64
}
