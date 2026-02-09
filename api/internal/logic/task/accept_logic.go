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
	poolmodel "github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/task/pool"

	"github.com/zeromicro/go-zero/core/logx"
)

type AcceptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 采纳
func NewAcceptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AcceptLogic {
	return &AcceptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AcceptLogic) Accept(req *types.AcceptReq) (resp *types.TaskBaseResp, err error) {
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

	// Step 3: 批量更新状态为ADOPTED (3)
	for _, id := range idsLong {
		pool, _ := l.svcCtx.BusinessTablePoolModel.FindOne(l.ctx, id)
		pool.Status = poolmodel.PoolStatusAdopted // ADOPTED = 3
		err := l.svcCtx.BusinessTablePoolModel.Update(l.ctx, pool)
		if err != nil {
			logx.Errorf("更新业务表池状态失败: id=%d, error=%v", id, err)
			return nil, localErrorx.TaskInvalidParam("更新状态失败")
		}
	}

	// Step 4: 返回成功
	logx.Infof("批量采纳成功: count=%d", len(idsLong))
	return &types.TaskBaseResp{Code: "0", Description: "采纳成功"}, nil
}
