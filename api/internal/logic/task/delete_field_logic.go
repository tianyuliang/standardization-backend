// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package task

import (
	"context"

	localErrorx "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFieldLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 移除字段
func NewDeleteFieldLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFieldLogic {
	return &DeleteFieldLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteFieldLogic) DeleteField(id int64) (resp *types.TaskBaseResp, err error) {
	// Step 1: 参数校验
	if id <= 0 {
		return nil, localErrorx.TaskParamEmpty("id")
	}

	// Step 2: 删除字段（逻辑删除）
	err = l.svcCtx.BusinessTablePoolModel.Delete(l.ctx, id)
	if err != nil {
		logx.Errorf("删除字段失败: id=%d, error=%v", id, err)
		return nil, localErrorx.TaskInvalidParam("删除字段失败")
	}

	// Step 3: 返回成功
	logx.Infof("删除字段成功: id=%d", id)
	return &types.TaskBaseResp{Code: "0", Description: "删除成功"}, nil
}
