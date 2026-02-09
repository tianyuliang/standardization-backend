// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package task

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCompletedTasksLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 已完成任务列表
func NewGetCompletedTasksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCompletedTasksLogic {
	return &GetCompletedTasksLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCompletedTasksLogic) GetCompletedTasks(req *types.PageInfoWithKeyword) (resp *types.TaskDataListResp, err error) {
	// Step 1: 参数处理
	keyword := req.Keyword
	page := req.Offset
	limit := req.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	// Step 2: 查询已完成任务列表
	tasks, totalCount, err := l.svcCtx.TaskStdCreateModel.FindCompleted(l.ctx, keyword, page, limit)
	if err != nil {
		logx.Errorf("查询已完成任务列表失败: %v", err)
		return &types.TaskDataListResp{Data: []types.TaskResp{}, TotalCount: 0}, nil
	}

	// Step 3: 构建响应
	result := make([]types.TaskResp, 0, len(tasks))
	for _, task := range tasks {
		result = append(result, buildTaskResp(task))
	}

	logx.Infof("查询已完成任务列表成功: count=%d, total=%d", len(result), totalCount)
	return &types.TaskDataListResp{
		Data:       result,
		TotalCount: totalCount,
	}, nil
}
