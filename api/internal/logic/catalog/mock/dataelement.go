// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

//go:build !mock_logic_off
// +build !mock_logic_off

package mock

import (
	"context"

	sharedmock "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/mock"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
)

// ============================================
// DataElement RPC Mock - 模块特定封装
//
// 本文件封装共享 mock 服务，保持与 catalog 模块的兼容性
// ============================================

// DataElementGetPageList 检查目录下是否存在数据元
func DataElementGetPageList(ctx context.Context, svcCtx *svc.ServiceContext, catalogId int64) bool {
	return sharedmock.DataElementGetPageList(ctx, svcCtx, catalogId)
}

// DataElementGetCountMapGroupByCatalog 按目录分组统计数据元数量
func DataElementGetCountMapGroupByCatalog(ctx context.Context, svcCtx *svc.ServiceContext) map[int64]int32 {
	return sharedmock.DataElementGetCountMapGroupByCatalog(ctx, svcCtx)
}
