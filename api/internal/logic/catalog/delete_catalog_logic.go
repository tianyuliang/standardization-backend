// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package catalog

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCatalogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除目录
func NewDeleteCatalogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCatalogLogic {
	return &DeleteCatalogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteCatalogLogic) DeleteCatalog() (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
