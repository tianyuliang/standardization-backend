// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package task

import (
	"context"

	localErrorx "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	taskmodel "github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/task/task"

	"github.com/zeromicro/go-zero/core/logx"
)

type StagingRelationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 标准关联暂存
func NewStagingRelationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StagingRelationLogic {
	return &StagingRelationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StagingRelationLogic) StagingRelation(req *types.StagingRelationReq) (resp *types.TaskBaseResp, err error) {
	// Step 1: 参数校验
	if req.TaskId <= 0 {
		return nil, localErrorx.TaskParamEmpty("taskId")
	}
	if req.BusinessTable == "" {
		return nil, localErrorx.TaskParamEmpty("businessTable")
	}

	// Step 2: 暂存任务数据（不调用推荐服务）
	task := &taskmodel.TaskStdCreate{
		TaskNo:  GenerateTaskNo(),
		Table:   req.BusinessTable,
		Status:  taskmodel.TaskStatusProcessing,
		Webhook: "",
		Deleted: 0,
	}

	id, err := l.svcCtx.TaskStdCreateModel.Insert(l.ctx, task)
	if err != nil {
		logx.Errorf("暂存任务失败: %v", err)
		return nil, localErrorx.TaskInvalidParam("暂存失败")
	}

	// Step 3: 返回成功
	logx.Infof("暂存任务成功: id=%d, taskNo=%s", id, task.TaskNo)
	return &types.TaskBaseResp{Code: "0", Description: "暂存成功"}, nil
}
