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
	"github.com/tianyuliang/standardization-backend/model/catalog/stub"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除目录
func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.DeleteReq) (resp *types.DeleteResp, err error) {
	// 1. 校验目录存在
	catalog, err := l.svcCtx.CatalogModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return &types.DeleteResp{
			Code:        errorx.NewWithCode(errorx.ErrInternal).Code(),
			Description: "查询目录失败",
		}, nil
	}
	if catalog == nil {
		return &types.DeleteResp{
			Code:        errorx.NewWithCode(errorx.ErrNotFound).Code(),
			Description: "目录不存在",
		}, nil
	}

	// 2. 校验不是根目录（level <= 1 不允许删除）
	if catalog.Level <= catalog.CatalogRootLevel {
		return &types.DeleteResp{
			Code:        errorx.NewWithCode(errorx.ErrInvalidParam).Code(),
			Description: "根目录不允许删除",
		}, nil
	}

	// 3. 递归获取所有子目录 ID
	descendants, err := l.svcCtx.CatalogModel.FindAllDescendants(l.ctx, req.Id)
	if err != nil {
		return &types.DeleteResp{
			Code:        errorx.NewWithCode(errorx.ErrInternal).Code(),
			Description: "查询子目录失败",
		}, nil
	}

	// 4. 调用桩模块检查关联数据（当前阶段跳过）
	// TODO: 实际项目中需要调用外部服务检查关联数据
	if err := l.checkRelatedData(catalog.Type, req.Id, descendants); err != nil {
		return &types.DeleteResp{
			Code:        errorx.NewWithCode(errorx.ErrDataExist).Code(),
			Description: err.Error(),
		}, nil
	}

	// 5. 批量删除目录及子目录
	deleteIds := append([]string{req.Id}, getCatalogIds(descendants)...)
	err = l.svcCtx.CatalogModel.DeleteBatch(l.ctx, deleteIds)
	if err != nil {
		return &types.DeleteResp{
			Code:        errorx.NewWithCode(errorx.ErrInternal).Code(),
			Description: "删除目录失败",
		}, nil
	}

	return &types.DeleteResp{
		Code:        errorx.NewWithCode(errorx.ErrSuccess).Code(),
		Description: "删除成功",
	}, nil
}

// checkRelatedData 检查关联数据
func (l *DeleteLogic) checkRelatedData(catalogType int32, catalogId string, descendants []*catalog.Catalog) error {
	// 当前阶段使用桩模块，所有检查返回 false（无关联数据）
	// 实际项目中需要通过 RPC 调用外部服务检查

	stubChecker := &stub.StubExternalChecker{}

	// 收集所有需要检查的目录ID
	idsToCheck := append([]string{catalogId}, getCatalogIds(descendants)...)

	// 根据目录类型检查不同的关联数据
	for _, id := range idsToCheck {
		var hasData bool
		var err error

		switch catalogType {
		case catalog.CatalogTypeDataElement:
			hasData, err = stubChecker.CheckDataElement(l.ctx, id)
		case catalog.CatalogTypeDict:
			hasData, err = stubChecker.CheckDict(l.ctx, id)
		case catalog.CatalogTypeEncodingRule:
			hasData, err = stubChecker.CheckRule(l.ctx, id)
		case catalog.CatalogTypeFile:
			hasData, err = stubChecker.CheckFile(l.ctx, id)
		}

		if err != nil {
			return fmt.Errorf("检查关联数据失败: %w", err)
		}

		// TODO: 当前桩模块始终返回 false，后续实现真实检查
		// if hasData {
		//     return fmt.Errorf("目录下存在关联数据，无法删除")
		// }
	}

	return nil
}

// getCatalogIds 从目录列表中提取 ID
func getCatalogIds(catalogs []*catalog.Catalog) []string {
	ids := make([]string, 0, len(catalogs))
	for _, cat := range catalogs {
		ids = append(ids, cat.Id)
	}
	return ids
}
