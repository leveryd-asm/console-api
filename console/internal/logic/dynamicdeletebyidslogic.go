package logic

import (
	"context"
	"fmt"

	"console-api/console/internal/svc"
	"console-api/console/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DynamicDeleteByIdsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDynamicDeleteByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DynamicDeleteByIdsLogic {
	return &DynamicDeleteByIdsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DynamicDeleteByIdsLogic) DynamicDeleteByIds(req *types.DynamicDeleteByQueryIdsRequest) (resp *types.DynamicDeleteResponse, err error) {
	result, err := l.svcCtx.DynamicModel.DeleteByIds(req.TableName, req.Ids)
	if err != nil {
		return nil, err
	}
	return &types.DynamicDeleteResponse{
		Message:      fmt.Sprintf("delete %d rows", result),
		AffectedRows: int(result),
	}, nil
}
