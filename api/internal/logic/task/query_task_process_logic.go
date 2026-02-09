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

type QueryTaskProcessLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 进度查询
func NewQueryTaskProcessLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryTaskProcessLogic {
	return &QueryTaskProcessLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryTaskProcessLogic) QueryTaskProcess(req *types.QueryProcessReq) (resp *types.ProcessResp, err error) {
	// Step 1: 参数校验 - 验证taskId格式（36位UUID）
	if len(req.TaskId) != 36 {
		return nil, localErrorx.TaskInvalidParam("taskId必须为36位UUID")
	}

	// Step 2: 查询总数（Java源码: count(queryWrapper) where taskId = ?）
	totalCount, err := l.svcCtx.BusinessTablePoolModel.CountByTaskId(l.ctx, req.TaskId)
	if err != nil {
		logx.Errorf("查询任务进度失败: taskId=%s, error=%v", req.TaskId, err)
		return nil, localErrorx.TaskInvalidParam("查询任务进度失败")
	}

	// Step 3: 查询完成数（Java源码: count(queryWrapper) where taskId = ? AND dataElementId IS NOT NULL）
	finishCount, err := l.svcCtx.BusinessTablePoolModel.CountByTaskIdWithDataElementId(l.ctx, req.TaskId)
	if err != nil {
		logx.Errorf("查询任务完成数失败: taskId=%s, error=%v", req.TaskId, err)
		return nil, localErrorx.TaskInvalidParam("查询任务完成数失败")
	}

	// Step 4: 返回进度信息（Java源码: finishNumber, totalNumber）
	logx.Infof("查询任务进度成功: taskId=%s, total=%d, finish=%d", req.TaskId, totalCount, finishCount)
	return &types.ProcessResp{
		FinishNumber: int(finishCount),
		TotalNumber:  int(totalCount),
	}, nil
}
