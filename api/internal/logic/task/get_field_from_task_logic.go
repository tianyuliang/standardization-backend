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

type GetFieldFromTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 任务关联字段
func NewGetFieldFromTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFieldFromTaskLogic {
	return &GetFieldFromTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFieldFromTaskLogic) GetFieldFromTask(req *types.PageInfo) (resp *types.FieldListResp, err error) {
	// Step 1: 参数处理
	page := req.Offset
	limit := req.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	// Step 2: 查询任务关联的字段（查询处理中状态的记录）
	pools, totalCount, err := l.svcCtx.BusinessTablePoolModel.FindWithPagination(l.ctx, "", poolmodel.PoolStatusProcessing, page, limit)
	if err != nil {
		logx.Errorf("查询任务关联字段失败: %v", err)
		return &types.FieldListResp{
			TotalCount: 0,
			Data:       []types.FieldResp{},
		}, nil
	}

	// Step 3: 构建响应（只返回有字段信息的记录）
	data := make([]types.FieldResp, 0)
	for _, pool := range pools {
		if pool.TableField != "" {
			data = append(data, buildFieldResp(pool))
		}
	}

	logx.Infof("查询任务关联字段成功: count=%d", len(data))
	return &types.FieldListResp{
		TotalCount: totalCount,
		Data:       data,
	}, nil
}
