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
// Dict Service Mock
//
// 对应 Java: IDictService
// ============================================

// DictQueryList 检查目录或子目录下是否存在码表
// 对应 Java: dictService.queryList(id, null, null, null, 1, 1, null, null, null)
// 用途: 删除目录前校验 - Dict类型目录
// Java源码位置: DeCatalogInfoServiceImpl.checkCatalogDelete() line 509
func DictQueryList(ctx context.Context, svcCtx *svc.ServiceContext, catalogId int64) bool {
	// MOCK: 返回false，表示目录下不存在码表，允许删除
	// Result<List<DictVo>> dictResult = dictService.queryList(id, null, null, null, 1, 1, null, null, null);
	// if (!CustomUtil.isEmpty(dictResult.getData())) { checkErrors.add(new CheckErrorVo(..., "目录或子目录下已存在码表，不允许删除")); }
	// TODO: 调用 Dict RPC
	return false
}

// DictGetCountMapGroupByCatalog 按目录分组统计码表数量
// 对应 Java: dictService.getCountMapGroupByCatalog()
// 用途: 查询目录树时统计每个目录下的码表数量
// Java源码位置: DeCatalogInfoServiceImpl.getCatalogCountMap() line 230
func DictGetCountMapGroupByCatalog(ctx context.Context, svcCtx *svc.ServiceContext) map[int64]int32 {
	// MOCK: 返回空map，表示没有码表统计
	// case DeDict: { countMap = dictService.getCountMapGroupByCatalog(); }
	// TODO: 调用 Dict RPC
	return make(map[int64]int32)
}
