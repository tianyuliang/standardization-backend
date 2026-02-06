// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package catalog

import (
	"context"

	baseerrorx "github.com/jinguoxing/idrm-go-base/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	catalogmodel "github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/catalog/catalog"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 检索目录
func NewQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryLogic {
	return &QueryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryLogic) Query(req *types.QueryReq) (resp []types.CatalogInfoVo, err error) {
	// Step 1: 参数校验
	// 对应 Java: DeCatalogInfoController.checkPost()

	// Step 1.1: 校验type有效性
	// 对应 Java: DeCatalogInfoServiceImpl.checkType()
	if req.Type == 0 {
		return nil, errorx.CatalogTypeEmpty()
	}

	if !catalogmodel.IsValidCatalogType(req.Type) {
		return nil, errorx.CatalogTypeInvalid()
	}

	// Step 2: 根据是否有关键字进行查询
	// 对应 Java: DeCatalogInfoController.query()
	keyword := req.Keyword

	// Step 2.1: 如果有关键字，进行模糊搜索
	// 对应 Java: if (StringUtils.isNotBlank(keyword))
	if keyword != "" {
		// Step 2.1.1: 按关键字模糊查询
		// 对应 Java: catalogService.findByName(keyword, type)
		catalogs, err := l.svcCtx.CatalogModel.FindByName(l.ctx, keyword, req.Type)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("Query FindByName failed: %v", err)
			return nil, baseerrorx.NewWithMsg(errorx.ErrCodeCatalogError, "查询目录失败")
		}

		// Step 2.1.2: 转换为响应格式（平铺列表）
		// 对应 Java: 转换为 CatalogInfoVo 列表
		result := make([]types.CatalogInfoVo, 0, len(catalogs))
		for _, c := range catalogs {
			result = append(result, modelToInfoVo(c))
		}

		return result, nil
	}

	// Step 2.2: 如果没有关键字，查询所有level>1的目录
	// 对应 Java: else分支，查询所有非根目录

	// Step 2.2.1: 查询该类型的所有目录
	allCatalogs, err := l.svcCtx.CatalogModel.FindByType(l.ctx, req.Type)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("Query FindByType failed: %v", err)
		return nil, baseerrorx.NewWithMsg(errorx.ErrCodeCatalogError, "查询目录失败")
	}

	// Step 2.2.2: 过滤出level>1的目录（排除根目录）
	// 对应 Java: filter level > 1
	result := make([]types.CatalogInfoVo, 0, len(allCatalogs))
	for _, c := range allCatalogs {
		if c.Level > 1 {
			result = append(result, modelToInfoVo(c))
		}
	}

	return result, nil
}

// modelToInfoVo 将Catalog模型转换为CatalogInfoVo（平铺格式）
// 对应 Java: CatalogInfoVo 的构建
func modelToInfoVo(c *catalogmodel.Catalog) types.CatalogInfoVo {
	return types.CatalogInfoVo{
		Id:          c.Id,
		CatalogName: c.CatalogName,
		Level:       c.Level,
		ParentId:    c.ParentId,
		Type:        c.Type,
	}
}
