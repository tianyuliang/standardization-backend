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

type GetBusinessTableFieldLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 业务表字段列表
func NewGetBusinessTableFieldLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBusinessTableFieldLogic {
	return &GetBusinessTableFieldLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBusinessTableFieldLogic) GetBusinessTableField(req *types.PageInfo) (resp *types.FieldListResp, err error) {
	// Step 1: 参数处理
	page := req.Offset
	limit := req.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	// Step 2: 查询业务表字段列表（查询待处理状态）
	// Note: 由于pool表是平铺存储表和字段的，分页查询需要特殊处理
	pools, totalCount, err := l.svcCtx.BusinessTablePoolModel.FindWithPagination(l.ctx, "", poolmodel.PoolStatusPending, page, limit)
	if err != nil {
		logx.Errorf("查询业务表字段列表失败: %v", err)
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

	logx.Infof("查询业务表字段列表成功: count=%d", len(data))
	return &types.FieldListResp{
		TotalCount: totalCount,
		Data:       data,
	}, nil
}

// buildFieldResp 构建字段响应
func buildFieldResp(pool *poolmodel.BusinessTablePool) types.FieldResp {
	return types.FieldResp{
		Id:               pool.Id,
		FieldName:        pool.TableField,
		FieldDescription: pool.FieldDescription,
		DataType:         pool.DataType,
	}
}
