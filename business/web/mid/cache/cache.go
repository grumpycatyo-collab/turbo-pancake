package cache

import (
	"github.com/rs/zerolog"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

type Item struct {
	Content    []byte
	Expiration int64
}

var cache sync.Map

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

		next(ctx)

		log.Info().Msg("Writing response data into cache")
		data := ctx.Response.Body()
		cache.Store(key, Item{
			Content:    data,
			Expiration: time.Now().Add(expiration).UnixNano(),
		})

	}
}
