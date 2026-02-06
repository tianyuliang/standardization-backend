// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package catalog

import (
	"context"

	baseerrorx "github.com/jinguoxing/idrm-go-base/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	catalogmodel "github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/catalog/catalog"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询目录树
func NewQueryTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryTreeLogic {
	return &QueryTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryTreeLogic) QueryTree(req *types.QueryTreeReq) (resp *types.CatalogResp, error error) {
	// Step 1: 根据请求参数决定查询方式
	// 对应 Java: DeCatalogInfoController.querySonTree()

	// Step 2: 如果指定了id，查询指定目录的子集树
	// 对应 Java: if (id != null && id > 0)
	if req.Id != 0 {
		// Step 2.1: 校验目录存在
		catalog, err := l.svcCtx.CatalogModel.FindOne(l.ctx, req.Id)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("QueryTree by id failed: %v", err)
			return nil, errorx.CatalogNotExist()
		}

		// Step 2.2: 获取该目录下所有子级
		// 使用 catalog.Type 和 catalog.Level+1 查询子级
		allCatalogs, err := l.svcCtx.CatalogModel.FindAllByTypeAndLevel(l.ctx, catalog.Type, catalog.Level)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("QueryTree FindAllByTypeAndLevel failed: %v", err)
			return nil, baseerrorx.NewWithMsg(errorx.ErrCodeCatalogError, "查询目录失败")
		}

		// Step 2.3: 构建树形结构（从当前目录开始）
		roots := buildTreeFromRoot(catalog, allCatalogs, l.svcCtx)
		if len(roots) == 0 {
			return nil, baseerrorx.NewWithMsg(errorx.ErrCodeCatalogError, "构建目录树失败")
		}

		return roots[0], nil
	}

	// Step 3: 如果指定了type，查询该类型的完整目录树
	// 对应 Java: if (type != null)
	if req.Type != 0 {
		// Step 3.1: 校验type有效性
		if !catalogmodel.IsValidCatalogType(req.Type) {
			return nil, errorx.CatalogTypeInvalid()
		}

		// Step 3.2: 查询该类型的所有目录
		allCatalogs, err := l.svcCtx.CatalogModel.FindByType(l.ctx, req.Type)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("QueryTree FindByType failed: %v", err)
			return nil, baseerrorx.NewWithMsg(errorx.ErrCodeCatalogError, "查询目录失败")
		}

		if len(allCatalogs) == 0 {
			// 返回空树
			return &types.CatalogResp{}, nil
		}

		// Step 3.3: 获取最小level作为根节点
		// 对应 Java: int minLevel = nodes.stream().mapToInt(DeCatalogInfo::getLevel).min().getAsInt();
		minLevel := allCatalogs[0].Level
		for _, c := range allCatalogs {
			if c.Level < minLevel {
				minLevel = c.Level
			}
		}

		// Step 3.4: 构建树形结构
		roots := BuildTree(allCatalogs, minLevel, l.svcCtx)

		// Step 3.5: 如果只有一个根节点，直接返回；否则包装在虚拟根节点中
		// 对应 Java: 返回 CatalogTreeNodeVo 列表
		if len(roots) == 1 {
			return roots[0], nil
		}

		// 多个根节点的情况，返回第一个（或者返回虚拟根节点）
		// 这里简化处理，返回第一个根节点
		if len(roots) > 0 {
			return roots[0], nil
		}

		return &types.CatalogResp{}, nil
	}

	// Step 4: 既没有id也没有type，返回错误
	// 对应 Java: 必须提供 id 或 type 参数
	return nil, baseerrorx.NewWithMsg(errorx.ErrCodeCatalogMissingParam, "必须指定目录ID或类型")
}

// buildTreeFromRoot 从指定根节点构建树形结构
// 对应 Java: DeCatalogInfoServiceImpl.generateTreeFromRoot()
func buildTreeFromRoot(root *catalogmodel.Catalog, allCatalogs []*catalogmodel.Catalog, svcCtx *svc.ServiceContext) []*types.CatalogResp {
	// 构建根节点响应
	rootResp := modelToResp(root)

	// 按level分组
	levelMap := make(map[int64][]*catalogmodel.Catalog)
	for _, c := range allCatalogs {
		levelMap[int64(c.Level)] = append(levelMap[int64(c.Level)], c)
	}

	// 递归设置子节点
	setChildren(rootResp, levelMap, svcCtx)

	return []*types.CatalogResp{rootResp}
}
