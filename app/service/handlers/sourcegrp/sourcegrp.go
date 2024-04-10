package sourcegrp

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/grumpycatyo-collab/turbo-pancake/business/core/source"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
)

type Handlers struct {
	Core source.Core
}

func GetSourceCampaigns(h *Handlers) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {

		idVal := ctx.UserValue("id")
		if idVal == nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}

		idStr, ok := idVal.(string)
		if !ok {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
		domain := string(ctx.QueryArgs().Peek("domain"))

		id, err := strconv.Atoi(idStr)
		if err != nil {
			ctx.SetStatusCode(http.StatusForbidden)
		}

		campaigns, err := h.Core.QueryCampaignsBySourceID(id, domain)
		if err != nil {
			switch {
			case errors.Is(err, source.ErrInvalidID):
				ctx.SetStatusCode(http.StatusBadRequest)
			case errors.Is(err, source.ErrNotFound):
				ctx.SetStatusCode(http.StatusNotFound)
			default:
				fmt.Printf("ID[%d]: %w", id, err)
			}
		}

		ctx.SetStatusCode(http.StatusOK)
		responseBytes, _ := json.Marshal(campaigns)
		ctx.Response.SetBody(responseBytes)
	}
}
