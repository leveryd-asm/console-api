package logic

import (
	"context"
	"fmt"
	"net/http"

	"console-api/console/internal/svc"
	"console-api/console/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DynamicDeleteByQueryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewDynamicDeleteByQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext, r *http.Request) *DynamicDeleteByQueryLogic {
	return &DynamicDeleteByQueryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
	}
}

func (l *DynamicDeleteByQueryLogic) DynamicDeleteByQuery(req *types.DynamicDeleteByQueryRequest) (resp *types.DynamicDeleteResponse, err error) {
	count, err := l.svcCtx.DynamicModel.DeleteByQuery(req.TableName, req.FuzzyQuery, req.QueryCondition)
	if err != nil {
		return nil, err
	}
	return &types.DynamicDeleteResponse{
		Message:      fmt.Sprintf("delete %d rows", count),
		AffectedRows: int(count),
	}, nil
}
