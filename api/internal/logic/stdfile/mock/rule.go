// Code scaffolded by goctl. Safe to edit.
//go:build !mock_logic_off
// +build !mock_logic_off

package mock

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
)

// ============================================
// Rule RPC Mock
//
// 替换目标: RuleService.*()
// ============================================

// RuleQueryPageByFileId 根据文件ID分页查询关联的编码规则
// MOCK: 模拟查询编码规则，返回空列表
// 替换目标: ruleService.queryPageByFileId(fileId, offset, limit)
func RuleQueryPageByFileId(ctx context.Context, svcCtx *svc.ServiceContext, fileId int64, offset, limit int) ([]interface{}, error) {
	// MOCK: 返回空列表
	// TODO: 调用 Rule RPC 查询关联编码规则
	return []interface{}{}, nil
}

// RuleAddRelation 添加文件与编码规则的关联关系
// MOCK: 模拟添加关联关系
// 替换目标: ruleService.addRelation(fileId, ruleIds)
func RuleAddRelation(ctx context.Context, svcCtx *svc.ServiceContext, fileId int64, ruleIds []int64) error {
	// MOCK: 直接返回成功
	// TODO: 调用 Rule RPC 添加关联关系
	return nil
}

// RuleQueryByFileId 根据文件ID查询关联的编码规则ID列表
// MOCK: 模拟查询关联ID，返回空列表
// 替换目标: ruleService.queryByFileId(fileId)
func RuleQueryByFileId(ctx context.Context, svcCtx *svc.ServiceContext, fileId int64) []int64 {
	// MOCK: 返回空列表
	// TODO: 调用 Rule RPC 查询关联编码规则ID
	return []int64{}
}
