// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package task

import (
	"context"
	"time"

	localErrorx "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	poolmodel "github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/task/pool"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 新建标准任务
func NewCreateTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTaskLogic {
	return &CreateTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTaskLogic) CreateTask(req *types.CreateTaskReq) (resp *types.TaskBaseResp, err error) {
	// Step 1: 参数校验 - 验证taskId格式（36位UUID）
	if len(req.TaskId) != 36 {
		return nil, localErrorx.TaskInvalidParam("taskId必须为36位UUID")
	}

	// Step 2: 验证字段ID格式（36位UUID）
	for _, id := range req.Ids {
		if len(id) != 36 {
			return nil, localErrorx.TaskInvalidParam("字段ID必须为36位UUID")
		}
	}

	// Step 3: 查询这些字段是否都在待新建池中
	pools, err := l.svcCtx.BusinessTablePoolModel.FindByBusinessTableFieldIds(l.ctx, req.Ids)
	if err != nil {
		logx.Errorf("查询业务表池失败: error=%v", err)
		return nil, localErrorx.TaskInvalidParam("查询业务表池失败")
	}
	if len(pools) != len(req.Ids) {
		return nil, localErrorx.TaskInvalidParam("输入的字段ID不在待新建标准池中")
	}

	// Step 4: 批量更新taskId和状态为CREATING (1)
	now := time.Now().Format("2006-01-02 15:04:05")
	for _, pool := range pools {
		pool.TaskId = req.TaskId
		pool.Status = poolmodel.PoolStatusProcessing // CREATING = 1
		pool.CreateTime = now
		err := l.svcCtx.BusinessTablePoolModel.Update(l.ctx, pool)
		if err != nil {
			logx.Errorf("更新业务表池失败: id=%d, error=%v", pool.Id, err)
			return nil, localErrorx.TaskInvalidParam("更新业务表池失败")
		}
	}

	// Step 5: 返回成功
	logx.Infof("创建任务成功: taskId=%s, 字段数=%d", req.TaskId, len(req.Ids))
	return &types.TaskBaseResp{Code: "0", Description: "创建成功"}, nil
}
