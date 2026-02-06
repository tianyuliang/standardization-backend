// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package catalog

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryWithFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询目录及文件树
func NewQueryWithFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryWithFileLogic {
	return &QueryWithFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryWithFileLogic) QueryWithFile(req *types.QueryWithFileReq) (resp *types.CatalogListByFileResp, err error) {
	// todo: add your logic here and delete this line

	return
}
