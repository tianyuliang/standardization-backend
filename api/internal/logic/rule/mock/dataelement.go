// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

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
// 替换目标: iDataElementInfoService.*()
// ============================================

// DataElementGetRuleIdByDataId 根据数据元ID获取关联的规则ID
// MOCK: 模拟获取规则ID，默认返回0（无关联）
// 替换目标: iDataElementInfoService.getById(dataId).getRuleId()
func DataElementGetRuleIdByDataId(ctx context.Context, svcCtx *svc.ServiceContext, dataId int64) int64 {
	// MOCK: 默认返回0（无关联规则）
	// TODO: 调用 DataElement RPC 获取规则ID
	// DataElementInfo dataElement = iDataElementInfoService.getById(dataId);
	// return dataElement == null || dataElement.getRuleId() == null ? 0 : dataElement.getRuleId();
	return 0
}

// DataElementGetRuleIdByDataCode 根据数据元编码获取关联的规则ID
// MOCK: 模拟获取规则ID，默认返回0（无关联）
// 替换目标: iDataElementInfoService.getOneByIdOrCode(2, dataCode).getRuleId()
func DataElementGetRuleIdByDataCode(ctx context.Context, svcCtx *svc.ServiceContext, dataCode string) int64 {
	// MOCK: 默认返回0（无关联规则）
	// TODO: 调用 DataElement RPC 根据编码获取规则ID
	// DataElementInfo dataElement = iDataElementInfoService.getOneByIdOrCode(2, dataCode);
	// return dataElement == null || dataElement.getRuleId() == null ? 0 : dataElement.getRuleId();
	return 0
}

// DataElementRuleUsed 检查规则是否被数据元引用
// MOCK: 模拟引用状态检查，默认返回 false（未使用）
// 替换目标: iDataElementInfoService.ruleUsed(ruleIds)
func DataElementRuleUsed(ctx context.Context, svcCtx *svc.ServiceContext, ruleId int64) bool {
	// MOCK: 默认返回未使用
	// TODO: 调用 DataElement RPC 检查规则是否被引用
	// Set<Long> usedIds = iDataElementInfoService.ruleUsed(Arrays.asList(ruleId));
	// return usedIds.contains(ruleId);
	return false
}

// DataElementQueryByRuleId 查询引用该规则的数据元列表
// MOCK: 模拟查询数据元列表，返回空列表
// 替换目标: iDataElementInfoService.queryByRuleId(ruleId, offset, limit)
func DataElementQueryByRuleId(ctx context.Context, svcCtx *svc.ServiceContext, ruleId int64, offset, limit int) ([]interface{}, int64, error) {
	// MOCK: 返回空列表
	// TODO: 调用 DataElement RPC 查询引用该规则的数据元
	// IPage<DataElementInfo> resultPageData = iDataElementInfoService.queryByRuleId(ruleEntity.getId(), offset, limit);
	// return resultPageData.getRecords(), resultPageData.getTotal();
	return []interface{}{}, 0, nil
}
