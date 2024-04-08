package handlers

import (
	"github.com/fasthttp/router"
	"github.com/grumpycatyo-collab/turbo-pancake/app/service/handlers/sourcegrp"
	"github.com/grumpycatyo-collab/turbo-pancake/business/core/source"
	//db2 "github.com/grumpycatyo-collab/turbo-pancake/business/core/source/db"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

func Handlers(r *router.Router, db *sqlx.DB, log *zerolog.Logger) {
	const version = "v1"

	sgh := sourcegrp.Handlers{
		Core: source.NewCore(log, db),
	}

	// TODO: Add version string in routes
	r.POST("/", sourcegrp.GetSourceCampaigns(&sgh))
}
