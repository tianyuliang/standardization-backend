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

type QuerywithfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询目录与文件树列表
func NewQuerywithfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuerywithfileLogic {
	return &QuerywithfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QuerywithfileLogic) Querywithfile(req *types.QueryWithFileReq) (resp *types.QueryWithFileResp, err error) {
	// 1. 处理目录名称参数
	var catalogName string
	if req.CatalogName != nil {
		catalogName = *req.CatalogName
	}

	// 2. 查询目录列表（模糊查询所有类型）
	var allCatalogs []*catalog.Catalog
	if catalogName != "" {
		// 模糊查询每种类型
		for _, catalogType := range []int32{
			catalog.CatalogTypeDataElement,
			catalog.CatalogTypeDict,
			catalog.CatalogTypeEncodingRule,
			catalog.CatalogTypeFile,
		} {
			catalogs, err := l.svcCtx.CatalogModel.FindByName(l.ctx, catalogName, catalogType)
			if err != nil {
				return &types.QueryWithFileResp{
					Code:        errorx.NewWithCode(errorx.ErrInternal).Code(),
					Description: "查询目录失败",
				}, nil
			}
			allCatalogs = append(allCatalogs, catalogs...)
		}
	} else {
		// 如果没有名称过滤，查询所有目录（不包括根目录）
		for _, catalogType := range []int32{
			catalog.CatalogTypeDataElement,
			catalog.CatalogTypeDict,
			catalog.CatalogTypeEncodingRule,
			catalog.CatalogTypeFile,
		} {
			catalogs, err := l.svcCtx.CatalogModel.FindByName(l.ctx, "", catalogType)
			if err != nil {
				return &types.QueryWithfileResp{
					Code:        errorx.NewWithCode(errorx.ErrInternal).Code(),
					Description: "查询目录失败",
				}, nil
			}
			allCatalogs = append(allCatalogs, catalogs...)
		}
	}

	// 3. 转换为 VO 格式
	var catalogVos []*types.Catalog
	for _, cat := range allCatalogs {
		catalogVo := &types.Catalog{
			Id:          cat.Id,
			CatalogName: cat.CatalogName,
			Description: cat.Description,
			Level:       cat.Level,
			ParentId:    cat.ParentId,
			Type:        cat.Type,
			AuthorityId: "",
		}
		if cat.AuthorityId != nil {
			catalogVo.AuthorityId = *cat.AuthorityId
		}
		catalogVos = append(catalogVos, catalogVo)
	}

	// 4. 查询文件列表（当前阶段返回空列表，待文件模块实现）
	var files []*types.FileCountVo

	// 5. 组装响应
	result := &types.CatalogListByFileVo{
		Catalogs: catalogVos,
		Files:    files,
	}

	return &types.QueryWithFileResp{
		Code:        errorx.NewWithCode(errorx.ErrSuccess).Code(),
		Description: "查询成功",
		Data:        result,
	}, nil
}
