// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package catalog

import (
	"context"
	"fmt"
	"strings"

	"github.com/jinguoxing/idrm-go-base/errorx"
	"github.com/tianyuliang/standardization-backend/api/internal/svc"
	"github.com/tianyuliang/standardization-backend/api/internal/types"
	"github.com/tianyuliang/standardization-backend/model/catalog/catalog"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 按名称检索目录
func NewQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryLogic {
	return &QueryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryLogic) Query(req *types.QueryReq) (resp *types.QueryResp, err error) {
	// 1. 校验 type 参数
	if req.Type < catalog.CatalogTypeDataElement || req.Type > catalog.CatalogTypeFile {
		return &types.QueryResp{
			Code:        errorx.NewWithCode(errorx.ErrInvalidParam).Code(),
			Description: fmt.Sprintf("目录类型无效，有效值为: %d=%s, %d=%s, %d=%s, %d=%s",
				catalog.CatalogTypeDataElement, catalog.GetCatalogTypeText(catalog.CatalogTypeDataElement),
				catalog.CatalogTypeDict, catalog.GetCatalogTypeText(catalog.CatalogTypeDict),
				catalog.CatalogTypeEncodingRule, catalog.GetCatalogTypeText(catalog.CatalogTypeEncodingRule),
				catalog.CatalogTypeFile, catalog.GetCatalogTypeText(catalog.CatalogTypeFile)),
		}, nil
	}

	// 2. 处理目录名称参数
	var catalogName string
	if req.CatalogName != nil {
		// SQL 特殊字符转义
		catalogName = escapeSqlLike(*req.CatalogName)
	}

	// 3. 执行查询
	catalogs, err := l.svcCtx.CatalogModel.FindByName(l.ctx, catalogName, req.Type)
	if err != nil {
		return &types.QueryResp{
			Code:        errorx.NewWithCode(errorx.ErrInternal).Code(),
			Description: "查询目录失败",
		}, nil
	}

	// 4. 转换为 VO 格式（FindByName 已经过滤 level > 1）
	var result []*types.CatalogInfoVo
	for _, cat := range catalogs {
		vo := &types.CatalogInfoVo{
			Id:          cat.Id,
			CatalogName: cat.CatalogName,
			Description: cat.Description,
			Level:       cat.Level,
			ParentId:    cat.ParentId,
			Type:        cat.Type,
			AuthorityId: "",
			Count:       0,
		}
		if cat.AuthorityId != nil {
			vo.AuthorityId = *cat.AuthorityId
		}
		result = append(result, vo)
	}

	return &types.QueryResp{
		Code:        errorx.NewWithCode(errorx.ErrSuccess).Code(),
		Description: "查询成功",
		Data:        result,
	}, nil
}

// escapeSqlLike 转义 SQL LIKE 查询的特殊字符
func escapeSqlLike(input string) string {
	// 转义 SQL LIKE 特殊字符: % _ \
	// % -> \%
	// _ -> \_
	// \ -> \\
	// 然后在查询时使用 LIKE '%value%'
	input = strings.ReplaceAll(input, `\`, `\\`)
	input = strings.ReplaceAll(input, `%`, `\%`)
	input = strings.ReplaceAll(input, `_`, `\_`)
	return input
}
