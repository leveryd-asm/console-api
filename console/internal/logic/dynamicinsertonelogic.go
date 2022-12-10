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

// 	copy from "github.com/zeromicro/go-zero/rest/internal/header"
const (
	// ApplicationJson stands for application/json.
	ApplicationJson = "application/json"
	// ContentType is the header key for Content-Type.
	ContentType = "Content-Type"
	// JsonContentType is the content type for JSON.
	JsonContentType = "application/json; charset=utf-8"

	maxBodyLen = 8 << 20 // 8MB
)

type DynamicInsertOneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewDynamicInsertOneLogic(ctx context.Context, svcCtx *svc.ServiceContext, r *http.Request) *DynamicInsertOneLogic {
	return &DynamicInsertOneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
	}
}

func (l *DynamicInsertOneLogic) DynamicInsertOne(req *types.DynamicInsertOneRequest) (resp *types.DynamicInsertOneResponse, err error) {
	r := l.r
	if r.ContentLength > 0 && strings.Contains(r.Header.Get(ContentType), ApplicationJson) {
		// copy from "httpx.ParseJsonBody"
		m := make(map[string]interface{})
		err := jsonx.UnmarshalFromReader(io.LimitReader(r.Body, maxBodyLen), &m)
		if err != nil {
			return nil, err
		}
		err = l.svcCtx.DynamicModel.InsertOne(req.TableName, m)
		if err != nil {
			return nil, err
		}
		return &types.DynamicInsertOneResponse{
			Message: "success",
		}, nil
	}

	return &types.DynamicInsertOneResponse{
		Message: "fail",
	}, nil
}
