package handlers

import (
	"github.com/fasthttp/router"
	"github.com/grumpycatyo-collab/turbo-pancake/app/service/handlers/sourcegrp"
)

//type Config struct{}

func Routes(r *router.Router) {
	const version = "v1"
	r.POST("/", sourcegrp.GetSourceCampaigns)
}
