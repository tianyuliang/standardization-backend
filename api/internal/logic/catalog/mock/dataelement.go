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
// DataElement Service Mock
//
// 对应 Java: IDataElementInfoService
// ============================================

// DataElementGetPageList 检查目录或子目录下是否存在数据元
// 对应 Java: dataelementInfoService.getPageList(id, null, null, null, 1, 1, null, null, null)
// 用途: 删除目录前校验 - DataElement类型目录
// Java源码位置: DeCatalogInfoServiceImpl.checkCatalogDelete() line 502
func DataElementGetPageList(ctx context.Context, svcCtx *svc.ServiceContext, catalogId int64) bool {
	// MOCK: 返回false，表示目录下不存在数据元，允许删除
	// IPage<DataElementInfo> result = dataelementInfoService.getPageList(id, null, null, null, 1, 1, null, null, null);
	// if (!CustomUtil.isEmpty(result.getRecords())) { checkErrors.add(new CheckErrorVo(..., "目录或子目录下已存在数据元，不允许删除")); }
	// TODO: 调用 DataElement RPC
	return false
}

// DataElementGetCountMapGroupByCatalog 按目录分组统计数据元数量
// 对应 Java: dataelementInfoService.getCountMapGroupByCatalog()
// 用途: 查询目录树时统计每个目录下的数据元数量
// Java源码位置: DeCatalogInfoServiceImpl.getCatalogCountMap() line 226
func DataElementGetCountMapGroupByCatalog(ctx context.Context, svcCtx *svc.ServiceContext) map[int64]int32 {
	// MOCK: 返回空map，表示没有数据元统计
	// case DataElement: { countMap = dataelementInfoService.getCountMapGroupByCatalog(); }
	// TODO: 调用 DataElement RPC
	return make(map[int64]int32)
}
