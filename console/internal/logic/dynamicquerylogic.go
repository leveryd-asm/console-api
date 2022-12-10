package logic

import (
	"context"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"

	"console-api/console/internal/svc"
	"console-api/console/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DynamicQueryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewDynamicQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext, r *http.Request) *DynamicQueryLogic {
	return &DynamicQueryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
	}
}

func (l *DynamicQueryLogic) DynamicQuery(req *types.DynamicQueryRequest) (resp *types.DynamicQueryResponse, err error) {
	params, err := httpx.GetFormValues(l.r)
	if err != nil {
		return nil, err
	}

	total, rows, err := l.svcCtx.DynamicModel.Query(req.TableName, req.Limit, req.Offset, req.OrderBy, req.Asc,
		req.FuzzyQuery, params)

	if err != nil {
		return nil, err
	}
	return &types.DynamicQueryResponse{
		Total: int(total),
		Rows:  rows,
	}, nil
}
