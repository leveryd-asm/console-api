package handler

import (
	"net/http"

	"console-api/console/internal/logic"
	"console-api/console/internal/svc"
	"console-api/console/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func dynamicInsertOneHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DynamicInsertOneRequest
		if err := httpx.ParsePath(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewDynamicInsertOneLogic(r.Context(), svcCtx, r)
		resp, err := l.DynamicInsertOne(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
