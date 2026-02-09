// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package task

import (
	"context"
	"fmt"
	"strconv"

	localErrorx "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelFieldLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 撤销
func NewCancelFieldLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelFieldLogic {
	return &CancelFieldLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelFieldLogic) CancelField(req *types.CancelFieldReq) (resp *types.TaskBaseResp, err error) {
	// Step 1: 参数校验
	if len(req.Ids) == 0 {
		return nil, localErrorx.TaskParamEmpty("ids")
	}

	// 将string IDs转换为Long IDs
	var idsLong []int64
	for _, idStr := range req.Ids {
		// Java源码检查: id.length() != 19 (BIGINT格式)
		if len(idStr) != 19 {
			return nil, localErrorx.TaskInvalidParam(fmt.Sprintf("ID格式错误: %s, 必须为19位", idStr))
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			return nil, localErrorx.TaskInvalidParam(fmt.Sprintf("ID格式错误: %s", idStr))
		}
		idsLong = append(idsLong, id)
	}

	// Step 2: 查询所有记录是否存在
	for _, id := range idsLong {
		pool, err := l.svcCtx.BusinessTablePoolModel.FindOne(l.ctx, id)
		if err != nil {
			logx.Errorf("查询业务表池失败: id=%d, error=%v", id, err)
			return nil, localErrorx.TaskNotExist()
		}
		if pool == nil {
			return nil, localErrorx.TaskDataNotExist()
		}
	}

	// Step 3: 批量删除任务关联 (将f_task_id设置为NULL)
	for _, id := range idsLong {
		err := l.svcCtx.BusinessTablePoolModel.DeleteTaskId(l.ctx, id)
		if err != nil {
			logx.Errorf("删除任务关联失败: id=%d, error=%v", id, err)
			return nil, localErrorx.TaskInvalidParam("删除任务关联失败")
		}
	}

	// Step 4: 返回成功
	logx.Infof("批量撤销成功: count=%d", len(idsLong))
	return &types.TaskBaseResp{Code: "0", Description: "撤销成功"}, nil
}
