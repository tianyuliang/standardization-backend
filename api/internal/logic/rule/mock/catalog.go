// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

//go:build !mock_logic_off
// +build !mock_logic_off

package mock

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/mock"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
)

// ============================================
// Catalog RPC Mock - 模块特定封装
//
// 本文件封装共享 mock 服务，保持与 rule 模块的兼容性
// 替换目标: iDeCatalogInfoService.*()
// ============================================

// CatalogCheckExist 校验目录是否存在
// 替换目标: iDeCatalogInfoService.checkCatalogIsExist(catalogId, CatalogTypeEnum.ValueRule)
func CatalogCheckExist(ctx context.Context, svcCtx *svc.ServiceContext, catalogId int64) bool {
	exists, _, _ := mock.CheckCatalogExist(ctx, catalogId)
	return exists
}

// CatalogGetCatalogName 获取目录名称
// 替换目标: iDeCatalogInfoService.getById(catalogId).getCatalogName()
func CatalogGetCatalogName(ctx context.Context, svcCtx *svc.ServiceContext, catalogId int64) string {
	name, _ := mock.GetCatalogName(ctx, catalogId)
	return name
}

// CatalogGetChildIds 获取子目录ID列表（递归）
// 替换目标: iDeCatalogInfoService.getIDList(catalogId)
func CatalogGetChildIds(ctx context.Context, svcCtx *svc.ServiceContext, catalogId int64) []int64 {
	ids, _ := mock.GetChildCatalogIds(ctx, catalogId)
	return ids
}

// CatalogIsStdFileCatalog 校验是否为标准文件目录
// 替换目标: catalog.getType().equals(CatalogTypeEnum.File)
func CatalogIsStdFileCatalog(ctx context.Context, svcCtx *svc.ServiceContext, catalogId int64) bool {
	// 对于 rule 模块，检查是否为文件目录
	exists, _ := mock.IsDataElementCatalog(ctx, catalogId)
	return exists
}

// CatalogIsRootCatalog 校验是否为根目录
// 替换目标: catalog.isRootPath()
func CatalogIsRootCatalog(ctx context.Context, svcCtx *svc.ServiceContext, catalogId int64) bool {
	// 默认不是根目录
	return false
}
