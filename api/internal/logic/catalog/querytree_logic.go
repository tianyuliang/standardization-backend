// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package catalog

import (
	"context"

	"github.com/tianyuliang/standardization-backend/api/internal/svc"
	"github.com/tianyuliang/standardization-backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QuerytreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询目录树
func NewQuerytreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuerytreeLogic {
	return &QuerytreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QuerytreeLogic) Querytree(req *types.QueryTreeReq) (resp *types.QueryTreeResp, err error) {
	// todo: add your logic here and delete this line

	return
}
