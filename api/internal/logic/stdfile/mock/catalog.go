// Code scaffolded by goctl. Safe to edit.
//go:build !mock_logic_off
// +build !mock_logic_off

package mock

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
)

// ============================================
// Catalog RPC Mock
//
// 替换目标: iDeCatalogInfoService.*()
// ============================================

// CatalogCheckExist 校验目录是否存在
// MOCK: 模拟目录校验，默认返回存在
// 替换目标: iDeCatalogInfoService.checkCatalogIsExist(catalogId, CatalogTypeEnum.File)
func CatalogCheckExist(ctx context.Context, svcCtx *svc.ServiceContext, catalogId int64) bool {
	// MOCK: 默认返回目录存在
	// TODO: 调用 Catalog RPC 校验
	return true
}

// CatalogGetCatalogName 获取目录名称
// MOCK: 模拟获取目录名称，返回默认格式
// 替换目标: iDeCatalogInfoService.getById(catalogId).getCatalogName()
func CatalogGetCatalogName(ctx context.Context, svcCtx *svc.ServiceContext, catalogId int64) string {
	// MOCK: 返回格式化的目录名称
	// TODO: 调用 Catalog RPC 获取目录名称
	if catalogId == 44 {
		return "全部目录"
	}
	return ""
}

// CatalogGetChildIds 获取子目录ID列表（递归）
// MOCK: 模拟获取子目录，返回自身
// 替换目标: iDeCatalogInfoService.getIDList(catalogId)
func CatalogGetChildIds(ctx context.Context, svcCtx *svc.ServiceContext, catalogId int64) []int64 {
	// MOCK: 返回自身
	// TODO: 调用 Catalog RPC 递归获取子目录列表
	return []int64{catalogId}
}

// CatalogIsStdFileCatalog 校验是否为标准文件目录
// MOCK: 模拟目录类型校验，默认返回 true
// 替换目标: catalog.getType().equals(CatalogTypeEnum.File)
func CatalogIsStdFileCatalog(ctx context.Context, svcCtx *svc.ServiceContext, catalogId int64) bool {
	// MOCK: 默认返回是标准文件目录
	// TODO: 调用 Catalog RPC 校验目录类型
	return true
}

// CatalogIsRootCatalog 校验是否为根目录
// MOCK: 模拟根目录校验，默认返回 false
// 替换目标: catalog.isRootPath()
func CatalogIsRootCatalog(ctx context.Context, svcCtx *svc.ServiceContext, catalogId int64) bool {
	// MOCK: 默认不是根目录
	// TODO: 调用 Catalog RPC 判断是否为根目录
	return false
}
