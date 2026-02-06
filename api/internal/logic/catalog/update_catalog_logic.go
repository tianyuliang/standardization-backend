// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package catalog

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCatalogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改目录
func NewUpdateCatalogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCatalogLogic {
	return &UpdateCatalogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCatalogLogic) UpdateCatalog(req *types.UpdateCatalogReq) (resp *types.EmptyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
