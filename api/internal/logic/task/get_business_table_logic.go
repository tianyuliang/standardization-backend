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

type GetBusinessTableLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 业务表列表
func NewGetBusinessTableLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBusinessTableLogic {
	return &GetBusinessTableLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBusinessTableLogic) GetBusinessTable(req *types.PageInfoWithKeyword) (resp *types.BusinessTableListResp, err error) {
	// Step 1: 参数处理
	keyword := req.Keyword
	page := req.Offset
	limit := req.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	// Step 2: 查询业务表列表（查询待处理状态）
	pools, totalCount, err := l.svcCtx.BusinessTablePoolModel.FindWithPagination(l.ctx, keyword, poolmodel.PoolStatusPending, page, limit)
	if err != nil {
		logx.Errorf("查询业务表列表失败: %v", err)
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

	logx.Infof("查询业务表列表成功: count=%d", len(data))
	return &types.BusinessTableListResp{
		TotalCount: totalCount,
		Data:       data,
	}, nil
}

// buildBusinessTableResp 构建业务表响应
func buildBusinessTableResp(pool *poolmodel.BusinessTablePool) types.BusinessTableResp {
	return types.BusinessTableResp{
		Id:               pool.Id,
		TableName:        pool.TableName,
		TableDescription: pool.TableDescription,
		TableField:       pool.TableField,
		Status:           pool.Status,
		CreateUser:       pool.CreateUser,
		CreateTime:       pool.CreateTime,
	}
}
