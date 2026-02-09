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

type UpdateDescriptionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改字段说明
func NewUpdateDescriptionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDescriptionLogic {
	return &UpdateDescriptionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateDescriptionLogic) UpdateDescription(req *types.UpdateDescriptionReq) (resp *types.TaskBaseResp, err error) {
	// Step 1: 参数校验
	if req.FieldId <= 0 {
		return nil, localErrorx.TaskParamEmpty("fieldId")
	}

	// Step 2: 查询字段是否存在
	field, err := l.svcCtx.BusinessTablePoolModel.FindOne(l.ctx, req.FieldId)
	if err != nil {
		logx.Errorf("查询字段失败: fieldId=%d, error=%v", req.FieldId, err)
		return nil, localErrorx.TaskNotExist()
	}
	if field == nil {
		return nil, localErrorx.TaskDataNotExist()
	}

	// Step 3: 更新字段说明
	field.FieldDescription = req.FieldDescription
	err = l.svcCtx.BusinessTablePoolModel.Update(l.ctx, field)
	if err != nil {
		logx.Errorf("更新字段说明失败: fieldId=%d, error=%v", req.FieldId, err)
		return nil, localErrorx.TaskInvalidParam("更新字段说明失败")
	}

	// Step 4: 返回成功
	logx.Infof("更新字段说明成功: fieldId=%d", req.FieldId)
	return &types.TaskBaseResp{Code: "0", Description: "更新成功"}, nil
}
