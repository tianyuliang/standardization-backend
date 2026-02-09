// Code scaffolded by goctl. Safe to edit.

package dataelement

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/mock"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveCatalogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量移动数据元到指定目录
func NewRemoveCatalogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveCatalogLogic {
	return &RemoveCatalogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveCatalogLogic) RemoveCatalog(req *types.RemoveCatalogReq) (resp *types.EmptyResp, err error) {
	// Step 1: 校验ID列表
	if len(req.Ids) == 0 {
		return nil, errorx.ParameterEmpty("ids")
	}

	// Step 2: 校验目标目录存在性
	if req.CatalogId <= 0 {
		return nil, errorx.InvalidParameter("catalogId", "目录ID必须大于0")
	}

	exists, _, err := mock.CheckCatalogExist(l.ctx, req.CatalogId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errorx.CatalogNotExist()
	}

	// Step 3: 移动目录（版本号+1）
	updateUser := "system" // TODO: 从token获取
	err = l.svcCtx.DataElementModel.MoveCatalog(l.ctx, req.Ids, req.CatalogId, updateUser)
	if err != nil {
		logx.Errorf("移动数据元目录失败: %v", err)
		return nil, err
	}

	// Step 4: 发送MQ消息
	for _, id := range req.Ids {
		_ = mock.SendMqMessage(l.ctx, "moveCatalog", map[string]interface{}{
			"id":        id,
			"catalogId": req.CatalogId,
		}, updateUser)
	}

	return &types.EmptyResp{}, nil
}
