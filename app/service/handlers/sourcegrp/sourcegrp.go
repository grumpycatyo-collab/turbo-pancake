package sourcegrp

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"net/http"
)

func GetSourceCampaigns(ctx *fasthttp.RequestCtx) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.Unmarshal(ctx.PostBody(), &request); err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)

	}

	ctx.SetStatusCode(http.StatusOK)
	responseBytes, _ := json.Marshal(request)
	ctx.SetBody(responseBytes)
}
