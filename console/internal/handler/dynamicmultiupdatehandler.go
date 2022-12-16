package handler

import (
	"net/http"

	"console-api/console/internal/logic"
	"console-api/console/internal/svc"
	"console-api/console/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func dynamicMultiUpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DynamicMultiUpdateRequest
		if err := httpx.ParsePath(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewDynamicMultiUpdateLogic(r.Context(), svcCtx, r)
		resp, err := l.DynamicMultiUpdate(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
