// Code scaffolded by goctl. Safe to edit.
//go:build !mock_logic_off
// +build !mock_logic_off

package mock

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
)

// ============================================
// DataElement RPC Mock
//
// 替换目标: IDataElementInfoService.*()
// ============================================

// DataElementQueryPageByFileId 根据文件ID分页查询关联的数据元
// MOCK: 模拟查询数据元，返回空列表
// 替换目标: dataelementInfoService.queryPageByFileId(fileId, offset, limit)
func DataElementQueryPageByFileId(ctx context.Context, svcCtx *svc.ServiceContext, fileId int64, offset, limit int) ([]interface{}, error) {
	// MOCK: 返回空列表
	// TODO: 调用 DataElement RPC 查询关联数据元
	return []interface{}{}, nil
}

// DataElementAddRelation 添加文件与数据元的关联关系
// MOCK: 模拟添加关联关系
// 替换目标: dataelementInfoService.addRelation(fileId, deIds)
func DataElementAddRelation(ctx context.Context, svcCtx *svc.ServiceContext, fileId int64, deIds []int64) error {
	// MOCK: 直接返回成功
	// TODO: 调用 DataElement RPC 添加关联关系
	return nil
}

// DataElementQueryByFileId 根据文件ID查询关联的数据元ID列表
// MOCK: 模拟查询关联ID，返回空列表
// 替换目标: dataelementInfoService.queryByFileId(fileId)
func DataElementQueryByFileId(ctx context.Context, svcCtx *svc.ServiceContext, fileId int64) []int64 {
	// MOCK: 返回空列表
	// TODO: 调用 DataElement RPC 查询关联数据元ID
	return []int64{}
}
