// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package task

import (
	"context"

	localErrorx "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	poolmodel "github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/task/pool"

	"github.com/zeromicro/go-zero/core/logx"
)

type FinishTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 完成任务
func NewFinishTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FinishTaskLogic {
	return &FinishTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FinishTaskLogic) FinishTask(taskId string) (resp *types.TaskBaseResp, err error) {
	// Step 1: 参数校验 - 验证taskId格式（36位UUID）
	if len(taskId) != 36 {
		return nil, localErrorx.TaskInvalidParam("taskId必须为36位UUID")
	}

	// Step 2: 查询该任务的所有池记录（Java源码: queryWrapper.eq(getTaskId, taskId), list(queryWrapper)）
	pools, err := l.svcCtx.BusinessTablePoolModel.FindByTaskId(l.ctx, taskId)
	if err != nil {
		logx.Errorf("查询任务池记录失败: taskId=%s, error=%v", taskId, err)
		return nil, localErrorx.TaskInvalidParam("查询任务池记录失败")
	}

	if len(pools) == 0 {
		logx.Infof("任务没有关联的池记录: taskId=%s", taskId)
		return &types.TaskBaseResp{Code: "0", Description: "完成成功"}, nil
	}

	// Step 3: 验证所有记录都有关联的数据元（Java源码: if isEmpty(dataElementId) throw exception）
	var ids []int64
	for _, pool := range pools {
		if pool.DataElementId == 0 {
			logx.Errorf("池记录没有关联数据元: id=%d, tableField=%s", pool.Id, pool.TableField)
			return nil, localErrorx.TaskInvalidParam("不能存在没有关联数据元的标准字段")
		}
		ids = append(ids, pool.Id)

		// TODO: 验证数据元标准分类是否一致（Java源码: dataElementInfoService.getById(), validate stdType）
		// dataElementInfo, err := mock.GetDataElementInfo(l.ctx, l.svcCtx, pool.DataElementId)
		// if err != nil {
		//     return nil, localErrorx.TaskInvalidParam("获取数据元信息失败")
		// }
		// taskDetailDto, err := mock.GetTaskDetailDto(l.ctx, l.svcCtx, taskId)
		// if err != nil {
		//     return nil, localErrorx.TaskInvalidParam("获取任务详情失败")
		// }
		// 验证数据元标准分类与任务组织类型是否一致...
	}

	// Step 4: 批量更新状态为CREATED（Java源码: setState(CREATED), updateBatchById(list)）
	err = l.svcCtx.BusinessTablePoolModel.UpdateBatchStatus(l.ctx, ids, poolmodel.PoolStatusCompleted)
	if err != nil {
		logx.Errorf("批量更新池状态失败: taskId=%s, count=%d, error=%v", taskId, len(ids), err)
		return nil, localErrorx.TaskInvalidParam("批量更新池状态失败")
	}

	logx.Infof("完成任务成功: taskId=%s, count=%d", taskId, len(ids))
	return &types.TaskBaseResp{Code: "0", Description: "完成成功"}, nil
}
