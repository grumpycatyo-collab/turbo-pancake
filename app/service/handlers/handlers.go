package handlers

import (
	"fmt"
	"github.com/fasthttp/router"
	"github.com/grumpycatyo-collab/turbo-pancake/app/service/handlers/sourcegrp"
	"github.com/grumpycatyo-collab/turbo-pancake/business/core/source"
	"github.com/grumpycatyo-collab/turbo-pancake/business/web/mid/cache"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"net/http/pprof"
	"time"
)

func Handlers(r *router.Router, db *sqlx.DB, log *zerolog.Logger) {
	const version = "v1"
	const cacheDuration = 5 * time.Minute

	sgh := sourcegrp.Handlers{
		Core: source.NewCore(log, db),
	}

	r.GET(fmt.Sprintf("/%s/source/{id}/campaigns", version), cache.Middleware(cacheDuration, sourcegrp.GetSourceCampaigns(&sgh), log))

	r.GET("/debug/pprof/", fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Index))
	r.GET("/debug/pprof/cmdline", fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Cmdline))
	r.GET("/debug/pprof/profile", fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Profile))
	r.GET("/debug/pprof/symbol", fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Symbol))
	r.GET("/debug/pprof/trace", fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Trace))
}
