package handlers

import (
	"fmt"
	"github.com/fasthttp/router"
	"github.com/grumpycatyo-collab/turbo-pancake/app/service/handlers/sourcegrp"
	"github.com/grumpycatyo-collab/turbo-pancake/business/core/source"
	"github.com/grumpycatyo-collab/turbo-pancake/business/web/mid/cache"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp/pprofhandler"
	"time"
)

func Handlers(r *router.Router, db *sqlx.DB, log *zerolog.Logger) {
	const version = "v1"
	const cacheDuration = 5 * time.Minute

	sgh := sourcegrp.Handlers{
		Core: source.NewCore(log, db),
	}

	r.GET(fmt.Sprintf("/%s/source/{id}/campaigns", version), cache.Middleware(cacheDuration, sourcegrp.GetSourceCampaigns(&sgh), log))

	r.GET("/debug/pprof/{profile:*}", pprofhandler.PprofHandler)
}
