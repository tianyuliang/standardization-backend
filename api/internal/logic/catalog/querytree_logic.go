// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package catalog

import (
	"context"
	"fmt"

	"github.com/jinguoxing/idrm-go-base/errorx"
	"github.com/tianyuliang/standardization-backend/api/internal/svc"
	"github.com/tianyuliang/standardization-backend/api/internal/types"
	"github.com/tianyuliang/standardization-backend/model/catalog/catalog"

	"github.com/zeromicro/go-zero/core/logx"
)

type QuerytreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询目录树
func NewQuerytreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuerytreeLogic {
	return &QuerytreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QuerytreeLogic) Querytree(req *types.QueryTreeReq) (resp *types.QueryTreeResp, err error) {
	// 1. 校验 type 参数
	if req.Type < catalog.CatalogTypeDataElement || req.Type > catalog.CatalogTypeFile {
		return &types.QueryTreeResp{
			Code:        errorx.NewWithCode(errorx.ErrInvalidParam).Code(),
			Description: fmt.Sprintf("目录类型无效，有效值为: %d=%s, %d=%s, %d=%s, %d=%s",
				catalog.CatalogTypeDataElement, catalog.GetCatalogTypeText(catalog.CatalogTypeDataElement),
				catalog.CatalogTypeDict, catalog.GetCatalogTypeText(catalog.CatalogTypeDict),
				catalog.CatalogTypeEncodingRule, catalog.GetCatalogTypeText(catalog.CatalogTypeEncodingRule),
				catalog.CatalogTypeFile, catalog.GetCatalogTypeText(catalog.CatalogTypeFile)),
		}, nil
	}

	// 2. 查询目录数据
	var allCatalogs []*catalog.Catalog
	if req.Id != nil {
		// 如果提供了 id，查询该目录及其子树
		id := fmt.Sprintf("%d", *req.Id)
		targetCatalog, err := l.svcCtx.CatalogModel.FindOne(l.ctx, id)
		if err != nil {
			return &types.QueryTreeResp{
				Code:        errorx.NewWithCode(errorx.ErrNotFound).Code(),
				Description: "目录不存在",
			}, nil
		}
		if targetCatalog == nil {
			return &types.QueryTreeResp{
				Code:        errorx.NewWithCode(errorx.ErrNotFound).Code(),
				Description: "目录不存在",
			}, nil
		}
		if targetCatalog.Type != req.Type {
			return &types.QueryTreeResp{
				Code:        errorx.NewWithCode(errorx.ErrInvalidParam).Code(),
				Description: "目录类型不匹配",
			}, nil
		}

		// 获取该目录的所有子孙节点
		allDescendants, err := l.svcCtx.CatalogModel.FindAllDescendants(l.ctx, id)
		if err != nil {
			return &types.QueryTreeResp{
				Code:        errorx.NewWithCode(errorx.ErrInternal).Code(),
				Description: "查询目录失败",
			}, nil
		}

		// 构建树：当前节点 + 所有子孙节点
		allCatalogs = append([]*catalog.Catalog{targetCatalog}, allDescendants...)
	} else {
		// 如果未提供 id，查询该类型所有目录
		allCatalogs, err = l.svcCtx.CatalogModel.FindByType(l.ctx, req.Type)
		if err != nil {
			return &types.QueryTreeResp{
				Code:        errorx.NewWithCode(errorx.ErrInternal).Code(),
				Description: "查询目录失败",
			}, nil
		}
	}

	// 3. 构建树结构
	tree := l.buildTree(allCatalogs, req.Id)

	return &types.QueryTreeResp{
		Code:        errorx.NewWithCode(errorx.ErrSuccess).Code(),
		Description: "查询成功",
		Data:        tree,
	}, nil
}

// buildTree 构建目录树结构
func (l *QuerytreeLogic) buildTree(catalogs []*catalog.Catalog, rootId *int64) *types.CatalogTreeNodeVo {
	if len(catalogs) == 0 {
		return nil
	}

	// 创建 id -> node 映射
	nodeMap := make(map[string]*types.CatalogTreeNodeVo)
	for _, cat := range catalogs {
		nodeMap[cat.Id] = &types.CatalogTreeNodeVo{
			Id:           cat.Id,
			CatalogName:  cat.CatalogName,
			Description:  cat.Description,
			Level:        cat.Level,
			ParentId:     cat.ParentId,
			Type:         cat.Type,
			AuthorityId:  "",
			Count:        0,
			Children:     nil,
			HaveChildren: false,
		}
		if cat.AuthorityId != nil {
			nodeMap[cat.Id].AuthorityId = *cat.AuthorityId
		}
	}

	// 找出根节点
	var roots []*types.CatalogTreeNodeVo
	if rootId != nil {
		// 如果指定了 rootId，以该节点为根
		rootIdStr := fmt.Sprintf("%d", *rootId)
		if node, ok := nodeMap[rootIdStr]; ok {
			roots = append(roots, node)
		}
	} else {
		// 否则以 level=1 的节点为根
		for _, node := range nodeMap {
			if node.Level == catalog.CatalogRootLevel {
				roots = append(roots, node)
			}
		}
	}

	// 构建父子关系
	for _, cat := range catalogs {
		node, ok := nodeMap[cat.Id]
		if !ok {
			continue
		}

		// 查找父节点
		if parent, parentOk := nodeMap[cat.ParentId]; parentOk {
			if rootId == nil || (rootId != nil && parent.Id != fmt.Sprintf("%d", *rootId)) {
				parent.Children = append(parent.Children, node)
				parent.HaveChildren = true
			}
		}
	}

	// 如果只有一个根节点，直接返回；否则返回虚拟根节点
	if len(roots) == 1 {
		return roots[0]
	}

	// 多个根节点的情况（一般不会发生，因为每种类型只有一个 level=1 的根目录）
	// 返回第一个根节点
	if len(roots) > 0 {
		return roots[0]
	}

	return nil
}

// ValidateCatalogName 校验目录名称格式
func ValidateCatalogName(name string) error {
	if len(name) < catalog.CatalogNameMinLength || len(name) > catalog.CatalogNameMaxLength {
		return fmt.Errorf("目录名称长度必须在 %d-%d 之间", catalog.CatalogNameMinLength, catalog.CatalogNameMaxLength)
	}

	// 简单校验：首字符必须是中文/英文/数字
	firstChar := name[0]
	if !((firstChar >= 0x4e00 && firstChar <= 0x9fa5) || // 中文
		(firstChar >= 'a' && firstChar <= 'z') ||
		(firstChar >= 'A' && firstChar <= 'Z') ||
		(firstChar >= '0' && firstChar <= '9')) {
		return fmt.Errorf("目录名称首字符必须是中文、字母或数字")
	}

	return nil
}
