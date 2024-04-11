package sourcegrp

import (
	"encoding/json"
	"errors"
	"github.com/grumpycatyo-collab/turbo-pancake/business/core/source"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
	"strings"
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

		id, err := strconv.Atoi(idStr)
		if err != nil {
			ctx.SetStatusCode(http.StatusForbidden)
		}

		domainStr := string(ctx.QueryArgs().Peek("domain"))
		domain := strings.ToLower(domainStr)

		filterStr := string(ctx.QueryArgs().Peek("filter"))
		filter := strings.ToLower(filterStr)

		campaigns, err := h.Core.QueryCampaignsBySourceID(id, domain, filter)
		if err != nil {
			switch {
			case errors.Is(err, source.ErrInvalidID):
				ctx.SetStatusCode(http.StatusBadRequest)
			case errors.Is(err, source.ErrNotFound):
				ctx.SetStatusCode(http.StatusNotFound)
			default:
				//fmt.Printf("ID[%d]: %w", id, err)
				ctx.SetStatusCode(http.StatusConflict)
			}
		}

		ctx.SetStatusCode(http.StatusOK)
		responseBytes, _ := json.Marshal(campaigns)
		ctx.Response.SetBody(responseBytes)
	}
}
