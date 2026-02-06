// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package catalog

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	catalogmodel "github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/catalog/catalog"
)

// ============================================
// 辅助函数
// ============================================

// timeToStr 将时间转换为字符串格式
func timeToStr(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// escapeSqlSpecialChars SQL特殊字符转义（XSS防护）
// 对应 Java: StringUtil.escapeSqlSpecialChars()
func escapeSqlSpecialChars(s string) string {
	if s == "" {
		return s
	}
	// 转义SQL特殊字符
	replacer := strings.NewReplacer(
		"\\", "\\\\",
		"'", "''",
		"\"", "\\\"",
		"\n", "\\n",
		"\r", "\\r",
		"\x00", "\\x00",
		"\x1a", "\\x1a",
	)
	return replacer.Replace(s)
}

// ============================================
// 目录名称校验
// 对应 Java: Constants.getRegexENOrCNVarL(1, 20)
// ============================================

// ValidateCatalogName 校验目录名称格式
// 对应 Java: DeCatalogInfoServiceImpl.checkPost() 中的目录名称校验
func ValidateCatalogName(name string) error {
	name = strings.TrimSpace(name)

	// 长度校验
	if name == "" {
		return errorx.CatalogNameEmpty()
	}
	if len(name) > 20 {
		return errorx.CatalogNameTooLong()
	}

	// 格式校验：中英文、数字、_、-，不能以_-开头
	// 对应正则: ^[\u4e00-\u9fa5a-zA-Z0-9][\u4e00-\u9fa5a-zA-Z0-9_-]{0,19}$
	matched, _ := regexp.MatchString(`^[\p{Han}a-zA-Z0-9][\p{Han}a-zA-Z0-9_-]*$`, name)
	if !matched {
		return errorx.CatalogNameInvalidFormat()
	}

	return nil
}

// ============================================
// 树形结构构建
// 对应 Java: DeCatalogInfoServiceImpl.generateTrees()
// ============================================

// BuildTree 构建目录树（带数据统计）
// 对应 Java: DeCatalogInfoServiceImpl.generateTrees(nodes, type, level)
func BuildTree(catalogs []*catalogmodel.Catalog, rootLevel int32, svcCtx *svc.ServiceContext) []*types.CatalogResp {
	if len(catalogs) == 0 {
		return []*types.CatalogResp{}
	}

	// 按level分组
	levelMap := make(map[int64][]*catalogmodel.Catalog)
	for _, c := range catalogs {
		levelMap[int64(c.Level)] = append(levelMap[int64(c.Level)], c)
	}

	// 获取根节点
	roots := make([]*types.CatalogResp, 0)
	for _, c := range catalogs {
		if c.Level == rootLevel {
			roots = append(roots, modelToResp(c))
		}
	}

	// 递归设置子节点
	for _, root := range roots {
		setChildren(root, levelMap, svcCtx)
	}

	return roots
}

// BuildTreeNoCount 构建目录树（不带数据统计）
// 对应 Java: DeCatalogInfoServiceImpl.generateTreesNoCount()
func BuildTreeNoCount(catalogs []*catalogmodel.Catalog, rootLevel int32) []*types.CatalogResp {
	if len(catalogs) == 0 {
		return []*types.CatalogResp{}
	}

	// 按level分组
	levelMap := make(map[int64][]*catalogmodel.Catalog)
	for _, c := range catalogs {
		levelMap[int64(c.Level)] = append(levelMap[int64(c.Level)], c)
	}

	// 获取根节点
	roots := make([]*types.CatalogResp, 0)
	for _, c := range catalogs {
		if c.Level == rootLevel {
			roots = append(roots, modelToResp(c))
		}
	}

	// 递归设置子节点
	for _, root := range roots {
		setChildrenNoCount(root, levelMap)
	}

	return roots
}

// setChildren 递归设置子节点（带数据统计）
func setChildren(parent *types.CatalogResp, levelMap map[int64][]*catalogmodel.Catalog, svcCtx *svc.ServiceContext) {
	children := levelMap[int64(parent.Level)+1]
	if len(children) == 0 {
		return
	}

	for _, child := range children {
		if child.ParentId == parent.Id {
			childResp := modelToResp(child)
			// TODO: 调用对应服务获取数据统计
			// childResp.Count = getCountMap(child.Id, child.Type)
			parent.Children = append(parent.Children, childResp)
			parent.HaveChildren = true
			setChildren(childResp, levelMap, svcCtx)
		}
	}
}

// setChildrenNoCount 递归设置子节点（不带数据统计）
func setChildrenNoCount(parent *types.CatalogResp, levelMap map[int64][]*catalogmodel.Catalog) {
	children := levelMap[int64(parent.Level)+1]
	if len(children) == 0 {
		return
	}

	for _, child := range children {
		if child.ParentId == parent.Id {
			childResp := modelToResp(child)
			parent.Children = append(parent.Children, childResp)
			parent.HaveChildren = true
			setChildrenNoCount(childResp, levelMap)
		}
	}
}

// ============================================
// 递归获取子级ID列表
// 对应 Java: DeCatalogInfoServiceImpl.getIDList(DeCatalogInfo)
// ============================================

// GetAllChildIds 递归获取所有子级ID列表
// 对应 Java: DeCatalogInfoServiceImpl.getIDList()
func GetAllChildIds(ctx context.Context, catalogId int64, svcCtx *svc.ServiceContext) ([]int64, error) {
	// 获取当前目录
	catalog, err := svcCtx.CatalogModel.FindOne(ctx, catalogId)
	if err != nil {
		return nil, err
	}

	// 获取所有level >= 当前level的目录
	allCatalogs, err := svcCtx.CatalogModel.FindByTypeAndLevel(ctx, catalog.Type, catalog.Level)
	if err != nil {
		return nil, err
	}

	// 递归收集子级ID
	childIds := make([]int64, 0)
	for _, cat := range allCatalogs {
		if isDescendant(cat, catalog) && cat.Id != catalogId {
			childIds = append(childIds, cat.Id)
		}
	}

	return childIds, nil
}

// isDescendant 判断catalog是否是ancestor的后代
func isDescendant(catalog, ancestor *catalogmodel.Catalog) bool {
	// TODO: 简化实现，完整实现需要传入父目录映射
	return catalog.ParentId == ancestor.Id || catalog.Level > ancestor.Level
}

// GetIDList 获取目录及其所有子级的ID列表（包含自身）
// 对应 Java: DeCatalogInfoServiceImpl.getIDList()
func GetIDList(ctx context.Context, catalogId int64, svcCtx *svc.ServiceContext) ([]int64, error) {
	childIds, err := GetAllChildIds(ctx, catalogId, svcCtx)
	if err != nil {
		return nil, err
	}

	// 添加自身ID
	result := append(childIds, catalogId)
	return result, nil
}

// ============================================
// 模型转换
// ============================================

// modelToResp 将Catalog模型转换为CatalogResp
func modelToResp(c *catalogmodel.Catalog) *types.CatalogResp {
	return &types.CatalogResp{
		Id:           c.Id,
		CatalogName:  c.CatalogName,
		Description:  c.Description,
		Level:        c.Level,
		ParentId:     c.ParentId,
		Type:         c.Type,
		Children:     nil,
		HaveChildren: false,
	}
}

// ============================================
// 类型校验
// 对应 Java: DeCatalogInfoServiceImpl.checkType()
// ============================================

// CheckType 校验目录类型是否有效
// 对应 Java: DeCatalogInfoServiceImpl.checkType()
func CheckType(catalogType int32) error {
	if catalogType == 0 {
		return errorx.CatalogTypeEmpty()
	}

	if !catalogmodel.IsValidCatalogType(catalogType) {
		return errorx.CatalogTypeInvalid()
	}

	return nil
}
