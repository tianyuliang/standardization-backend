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

type DeleteLabelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除数据元标签
func NewDeleteLabelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLabelLogic {
	return &DeleteLabelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLabelLogic) DeleteLabel(id int64) (resp *types.EmptyResp, err error) {
	// Step 1: 校验标签ID
	if id <= 0 {
		return nil, errorx.InvalidParameter("id", "标签ID必须大于0")
	}

	// Step 2: 删除标签（批量清空所有数据元的该标签）
	updateUser := "system" // TODO: 从token获取
	err = l.svcCtx.DataElementModel.DeleteLabelIds(l.ctx, []int64{id})
	if err != nil {
		logx.Errorf("删除数据元标签失败: %v", err)
		return nil, err
	}

	// Step 3: 发送MQ消息
	_ = mock.SendMqMessage(l.ctx, "deleteLabel", map[string]interface{}{
		"labelId": id,
	}, updateUser)

	return &types.EmptyResp{}, nil
}
