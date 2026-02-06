// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package catalog

import (
	"context"

	baseerrorx "github.com/jinguoxing/idrm-go-base/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCatalogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除目录
func NewDeleteCatalogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCatalogLogic {
	return &DeleteCatalogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteCatalogLogic) DeleteCatalog(id int64) (resp *types.EmptyResp, err error) {
	// Step 1: 校验目录存在
	// 对应 Java: DeCatalogInfoServiceImpl.removeWithChildren() - 校验目录存在
	catalog, err := l.svcCtx.CatalogModel.FindOne(l.ctx, id)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("DeleteCatalog FindOne failed: %v", err)
		return nil, errorx.CatalogNotExist()
	}

	// Step 2: 根目录不允许删除
	// 对应 Java: DeCatalogInfoServiceImpl.checkCatalogDelete() - 根目录校验
	if catalog.Level <= 1 {
		return nil, errorx.CatalogCannotDeleteRoot()
	}

	// Step 3: 删除前校验（检查目录及子目录下是否存在数据）
	// 对应 Java: DeCatalogInfoServiceImpl.checkCatalogDelete()
	if err := CheckCatalogDelete(l.ctx, catalog, l.svcCtx); err != nil {
		return nil, err
	}

	// Step 4: 递归获取所有子级目录ID
	// 对应 Java: DeCatalogInfoServiceImpl.getIDList()
	childIds, err := GetAllChildIds(l.ctx, id, l.svcCtx)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("DeleteCatalog GetAllChildIds failed: %v", err)
		return nil, baseerrorx.NewWithMsg(errorx.ErrCodeCatalogError, "获取子级目录失败")
	}

	// Step 5: 构建需要删除的ID列表（包含当前目录和所有子级）
	// 对应 Java: DeCatalogInfoServiceImpl.removeWithChildren()
	allIds := append(childIds, id)

	// Step 6: 批量删除目录
	// 对应 Java: 逻辑删除 f_deleted = f_id + 1
	err = l.svcCtx.CatalogModel.DeleteByIds(l.ctx, allIds)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("DeleteCatalog DeleteByIds failed: %v", err)
		return nil, baseerrorx.NewWithMsg(errorx.ErrCodeCatalogError, "删除目录失败")
	}

	// Step 7: 返回成功响应
	return &types.EmptyResp{}, nil
}
