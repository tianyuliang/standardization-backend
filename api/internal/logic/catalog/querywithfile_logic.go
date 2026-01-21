// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package catalog

import (
	"context"

	"github.com/tianyuliang/standardization-backend/api/internal/svc"
	"github.com/tianyuliang/standardization-backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QuerywithfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询目录与文件树列表
func NewQuerywithfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuerywithfileLogic {
	return &QuerywithfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QuerywithfileLogic) Querywithfile(req *types.QueryWithFileReq) (resp *types.QueryWithFileResp, err error) {
	// todo: add your logic here and delete this line

	return
}
