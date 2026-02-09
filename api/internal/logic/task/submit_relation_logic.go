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

type SubmitRelationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 标准关联提交
func NewSubmitRelationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitRelationLogic {
	return &SubmitRelationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubmitRelationLogic) SubmitRelation(req *types.StagingRelationReq) (resp *types.TaskBaseResp, err error) {
	// Step 1: 参数校验
	if req.TaskId <= 0 {
		return nil, localErrorx.TaskParamEmpty("taskId")
	}
	if req.BusinessTable == "" {
		return nil, localErrorx.TaskParamEmpty("businessTable")
	}

	// Step 2: 保存任务数据
	task := &taskmodel.TaskStdCreate{
		TaskNo:  GenerateTaskNo(),
		Table:   req.BusinessTable,
		Status:  taskmodel.TaskStatusProcessing,
		Webhook: "",
		Deleted: 0,
	}

	id, err := l.svcCtx.TaskStdCreateModel.Insert(l.ctx, task)
	if err != nil {
		logx.Errorf("保存任务失败: %v", err)
		return nil, localErrorx.TaskInvalidParam("保存任务失败")
	}

	// Step 3: 异步调用推荐算法（TODO: 后续实现 HTTP 调用）
	// TODO: 调用推荐服务
	logx.Infof("TODO: 调用推荐算法服务: taskId=%d", id)

	// Step 4: 返回成功
	logx.Infof("提交任务成功: id=%d, taskNo=%s", id, task.TaskNo)
	return &types.TaskBaseResp{Code: "0", Description: "提交成功"}, nil
}
