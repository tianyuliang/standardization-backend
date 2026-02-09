// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package task

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	poolmodel "github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/task/pool"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTableFromTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 任务关联业务表
func NewGetTableFromTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTableFromTaskLogic {
	return &GetTableFromTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTableFromTaskLogic) GetTableFromTask(req *types.PageInfo) (resp *types.BusinessTableListResp, err error) {
	// Step 1: 参数处理
	page := req.Offset
	limit := req.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	// Step 2: 查询任务关联的业务表（查询处理中状态的记录）
	pools, totalCount, err := l.svcCtx.BusinessTablePoolModel.FindWithPagination(l.ctx, "", poolmodel.PoolStatusProcessing, page, limit)
	if err != nil {
		logx.Errorf("查询任务关联业务表失败: %v", err)
		return &types.BusinessTableListResp{
			TotalCount: 0,
			Data:       []types.BusinessTableResp{},
		}, nil
	}

	// Step 3: 构建响应
	data := make([]types.BusinessTableResp, 0, len(pools))
	for _, pool := range pools {
		data = append(data, buildBusinessTableResp(pool))
	}

	logx.Infof("查询任务关联业务表成功: count=%d", len(data))
	return &types.BusinessTableListResp{
		TotalCount: totalCount,
		Data:       data,
	}, nil
}
