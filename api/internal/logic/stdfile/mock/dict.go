// Code scaffolded by goctl. Safe to edit.
//go:build !mock_logic_off
// +build !mock_logic_off

package mock

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
)

// ============================================
// Dict RPC Mock
//
// 替换目标: IDictService.*()
// ============================================

// DictQueryPageByFileId 根据文件ID分页查询关联的码表
// MOCK: 模拟查询码表，返回空列表
// 替换目标: dictService.queryPageByFileId(fileId, offset, limit)
func DictQueryPageByFileId(ctx context.Context, svcCtx *svc.ServiceContext, fileId int64, offset, limit int) ([]interface{}, error) {
	// MOCK: 返回空列表
	// TODO: 调用 Dict RPC 查询关联码表
	return []interface{}{}, nil
}

// DictAddRelation 添加文件与码表的关联关系
// MOCK: 模拟添加关联关系
// 替换目标: dictService.addRelation(fileId, dictIds)
func DictAddRelation(ctx context.Context, svcCtx *svc.ServiceContext, fileId int64, dictIds []int64) error {
	// MOCK: 直接返回成功
	// TODO: 调用 Dict RPC 添加关联关系
	return nil
}

// DictQueryByFileId 根据文件ID查询关联的码表ID列表
// MOCK: 模拟查询关联ID，返回空列表
// 替换目标: dictService.queryByFileId(fileId)
func DictQueryByFileId(ctx context.Context, svcCtx *svc.ServiceContext, fileId int64) []int64 {
	// MOCK: 返回空列表
	// TODO: 调用 Dict RPC 查询关联码表ID
	return []int64{}
}
