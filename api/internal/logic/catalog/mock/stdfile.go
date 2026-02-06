// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

//go:build !mock_logic_off
// +build !mock_logic_off

package mock

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
)

// ============================================
// StdFile Service Mock
//
// 对应 Java: StdFileMgrService
// ============================================

// StdFileQueryList 检查目录或子目录下是否存在标准文件
// 对应 Java: stdFileMgrService.queryList(id, null, null, null, 1, 1, null, null, null)
// 用途: 删除目录前校验 - File类型目录
// Java源码位置: DeCatalogInfoServiceImpl.checkCatalogDelete() line 523
func StdFileQueryList(ctx context.Context, svcCtx *svc.ServiceContext, catalogId int64) bool {
	// MOCK: 返回false，表示目录下不存在标准文件，允许删除
	// Result<List<StdFileMgrVo>> stdFileResult = stdFileMgrService.queryList(id, null, null, null, 1, 1, null, null, null);
	// if (CustomUtil.isNotEmpty(stdFileResult.getData())) { checkErrors.add(new CheckErrorVo(..., "目录或子目录下已存在文件，不允许删除")); }
	// TODO: 调用 StdFile RPC
	return false
}

// StdFileGetByName 按名称模糊搜索标准文件
// 对应 Java: stdFileMgrService.getByName(catalog_name)
// 用途: 查询目录及文件树 - 搜索文件
// Java源码位置: DeCatalogInfoController.queryParentTree() line 404 (query/with_file接口)
func StdFileGetByName(ctx context.Context, svcCtx *svc.ServiceContext, keyword string) []*types.FileCountVo {
	// MOCK: 返回空列表
	// List<StdFileMgrEntity> fileList = stdFileMgrService.getByName(keyword);
	// fileList.forEach(i -> { FileCountVo fileCountVo = new FileCountVo(); ... });
	// TODO: 调用 StdFile RPC
	return []*types.FileCountVo{}
}

// StdFileGetCountMapGroupByCatalog 按目录分组统计标准文件数量
// 对应 Java: stdFileMgrService.getCountMapGroupByCatalog()
// 用途: 查询目录树时统计每个目录下的标准文件数量
// Java源码位置: DeCatalogInfoServiceImpl.getCatalogCountMap() line 238
func StdFileGetCountMapGroupByCatalog(ctx context.Context, svcCtx *svc.ServiceContext) map[int64]int32 {
	// MOCK: 返回空map，表示没有标准文件统计
	// case File: { countMap = stdFileMgrService.getCountMapGroupByCatalog(); }
	// TODO: 调用 StdFile RPC
	return make(map[int64]int32)
}
