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
// Rule Service Mock
//
// 对应 Java: RuleService
// ============================================

// RuleQueryList 检查目录或子目录下是否存在编码规则
// 对应 Java: ruleService.queryList(id, null, null, null, 1, 1, null, null, null, null)
// 用途: 删除目录前校验 - ValueRule类型目录
// Java源码位置: DeCatalogInfoServiceImpl.checkCatalogDelete() line 516
func RuleQueryList(ctx context.Context, svcCtx *svc.ServiceContext, catalogId int64) bool {
	// MOCK: 返回false，表示目录下不存在编码规则，允许删除
	// Result<List<RuleVo>> ruleResult = ruleService.queryList(id, null, null, null, 1, 1, null, null, null, null);
	// if (!CustomUtil.isEmpty(ruleResult.getData())) { checkErrors.add(new CheckErrorVo(..., "目录或子目录下已存在编码规则，不允许删除")); }
	// TODO: 调用 Rule RPC
	return false
}

// RuleGetCountMapGroupByCatalog 按目录分组统计编码规则数量
// 对应 Java: ruleService.getCountMapGroupByCatalog()
// 用途: 查询目录树时统计每个目录下的编码规则数量
// Java源码位置: DeCatalogInfoServiceImpl.getCatalogCountMap() line 234
func RuleGetCountMapGroupByCatalog(ctx context.Context, svcCtx *svc.ServiceContext) map[int64]int32 {
	// MOCK: 返回空map，表示没有编码规则统计
	// case ValueRule: { countMap = ruleService.getCountMapGroupByCatalog(); }
	// TODO: 调用 Rule RPC
	return make(map[int64]int32)
}
