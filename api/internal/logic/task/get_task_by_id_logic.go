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

type GetTaskByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 任务详情
func NewGetTaskByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTaskByIdLogic {
	return &GetTaskByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTaskByIdLogic) GetTaskById(id int64) (resp *types.TaskDetailResp, err error) {
	// Step 1: 查询任务
	task, err := l.svcCtx.TaskStdCreateModel.FindOne(l.ctx, id)
	if err != nil {
		logx.Errorf("查询任务失败: id=%d, error=%v", id, err)
		return nil, localErrorx.TaskNotExist()
	}
	if task == nil {
		return nil, localErrorx.TaskDataNotExist()
	}

	// Step 2: 查询任务结果
	results, err := l.svcCtx.TaskStdCreateResultModel.FindByTaskId(l.ctx, id)
	if err != nil {
		logx.Errorf("查询任务结果失败: taskId=%d, error=%v", id, err)
		results = []*taskmodel.TaskStdCreateResult{}
	}

	// Step 3: 构建响应
	resultList := make([]types.TaskStdCreateResult, 0, len(results))
	for _, r := range results {
		resultList = append(resultList, buildTaskResultResp(r))
	}

	logx.Infof("查询任务详情成功: id=%d, resultCount=%d", id, len(resultList))
	return &types.TaskDetailResp{
		Task:    buildTaskResp(task),
		Results: resultList,
	}, nil
}
