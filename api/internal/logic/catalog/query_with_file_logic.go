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

type QueryWithFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询目录及文件树
func NewQueryWithFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryWithFileLogic {
	return &QueryWithFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryWithFileLogic) QueryWithFile(req *types.QueryWithFileReq) (resp *types.CatalogListByFileResp, err error) {
	// Step 1: 获取文件目录列表
	// 对应 Java: DeCatalogInfoServiceImpl.getCatalogFileTree()
	// 默认查询 type=4 (文件类型) 的目录
	var catalogs []*catalogmodel.Catalog
	keyword := ""
	if req.Keyword != "" {
		keyword = req.Keyword
	}

	// Step 1.1: 如果有关键字，模糊搜索所有类型的目录
	// 对应 Java: DeCatalogInfoServiceImpl.queryByKeyword()
	if keyword != "" {
		// 按名称模糊搜索所有类型目录
		catalogs, err = l.svcCtx.CatalogModel.FindByName(l.ctx, keyword, catalogmodel.CatalogTypeFile)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("QueryWithFile FindByName failed: %v", err)
			return nil, baseerrorx.NewWithMsg(errorx.ErrCodeCatalogError, "查询目录失败")
		}
	} else {
		// Step 1.2: 无关键字时，返回所有文件类型目录 (type=4)
		// 对应 Java: DeCatalogInfoServiceImpl.queryList() - 获取文件目录
		catalogs, err = l.svcCtx.CatalogModel.FindByType(l.ctx, catalogmodel.CatalogTypeFile)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("QueryWithFile FindByType failed: %v", err)
			return nil, baseerrorx.NewWithMsg(errorx.ErrCodeCatalogError, "查询目录失败")
		}
	}

	// Step 2: 转换为平铺列表格式（排除level=1的根目录）
	// 对应 Java: 返回 CatalogInfoVo 平铺列表
	catalogInfos := make([]*types.CatalogInfoVo, 0)
	for _, catalog := range catalogs {
		// 只返回level > 1的目录（排除根目录）
		if catalog.Level > 1 {
			catalogInfos = append(catalogInfos, &types.CatalogInfoVo{
				Id:          catalog.Id,
				CatalogName: catalog.CatalogName,
				Level:       catalog.Level,
				ParentId:    catalog.ParentId,
				Type:        catalog.Type,
			})
		}
	}

	// Step 3: 获取文件列表
	// 对应 Java: StdFileMgrService.queryList() 或 getByName()
	var files []*types.FileCountVo
	if keyword != "" {
		// Step 3.1: 按关键字搜索文件
		// 对应 Java: stdFileMgrService.getByName(keyword)
		files = getMockFiles(keyword)
	} else {
		// Step 3.2: 获取所有文件
		// 对应 Java: stdFileMgrService.queryList()
		files = getMockFiles("")
	}

	// Step 4: 组装响应
	// 对应 Java: 返回 CatalogListByFileResp
	return &types.CatalogListByFileResp{
		Catalogs: catalogInfos,
		Files:    files,
	}, nil
}
