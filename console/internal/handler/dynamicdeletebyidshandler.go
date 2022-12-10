package handler

import (
	"net/http"

	"console-api/console/internal/logic"
	"console-api/console/internal/svc"
	"console-api/console/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func dynamicDeleteByIdsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DynamicDeleteByQueryIdsRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewDynamicDeleteByIdsLogic(r.Context(), svcCtx)
		resp, err := l.DynamicDeleteByIds(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
