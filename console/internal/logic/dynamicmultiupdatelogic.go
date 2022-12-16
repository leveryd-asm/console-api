package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/jsonx"
	"io"
	"net/http"
	"strings"

	"console-api/console/internal/svc"
	"console-api/console/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DynamicMultiUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewDynamicMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext, r *http.Request) *DynamicMultiUpdateLogic {
	return &DynamicMultiUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
	}
}

func (l *DynamicMultiUpdateLogic) DynamicMultiUpdate(req *types.DynamicMultiUpdateRequest) (resp *types.DynamicUpdateResponse, err error) {
	r := l.r
	if r.ContentLength > 0 && strings.Contains(r.Header.Get(ContentType), ApplicationJson) {
		// copy from "httpx.ParseJsonBody"
		m := make(map[string]interface{})
		err := jsonx.UnmarshalFromReader(io.LimitReader(r.Body, maxBodyLen), &m)
		if err != nil {
			return nil, err
		}
		err = l.svcCtx.DynamicModel.MultiUpdate(req.TableName, m)
		if err != nil {
			return nil, err
		}
		return &types.DynamicUpdateResponse{
			Message: "success",
		}, nil
	}
	return &types.DynamicUpdateResponse{
		Message: "fail",
	}, nil
}
