// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package catalog

import (
	"context"
	"time"

	baseerrorx "github.com/jinguoxing/idrm-go-base/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCatalogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改目录
func NewUpdateCatalogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCatalogLogic {
	return &UpdateCatalogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCatalogLogic) UpdateCatalog(id int64, req *types.UpdateCatalogReq) (resp *types.EmptyResp, err error) {
	// Step 1: 校验目录存在
	// 对应 Java: DeCatalogInfoServiceImpl.update() - 校验目录存在
	catalog, err := l.svcCtx.CatalogModel.FindOne(l.ctx, id)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("UpdateCatalog FindOne failed: %v", err)
		return nil, errorx.CatalogNotExist()
	}

	// Step 2: 不允许修改根目录
	// 对应 Java: DeCatalogInfoServiceImpl.update() - 根目录校验
	if catalog.Level <= 1 {
		return nil, errorx.CatalogCannotModifyRoot()
	}

	// Step 3: 目录名称格式校验
	// 对应 Java: DeCatalogInfoServiceImpl.checkPost() 中的目录名称校验
	if err := ValidateCatalogName(req.CatalogName); err != nil {
		return nil, err
	}

	// Step 4: 同级目录名称唯一性校验（排除自身）
	// 对应 Java: DeCatalogInfoServiceImpl.checkNameDuplicate() - 排除自身
	siblings, err := l.svcCtx.CatalogModel.FindByParentId(l.ctx, catalog.ParentId)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("UpdateCatalog FindByParentId failed: %v", err)
		return nil, baseerrorx.NewWithMsg(errorx.ErrCodeCatalogError, "查询目录失败")
	}

	for _, sibling := range siblings {
		// 排除自身，检查同名
		if sibling.Id != id && sibling.CatalogName == req.CatalogName {
			return nil, errorx.CatalogNameDuplicate()
		}
	}

	// Step 5: 更新目录数据
	// 对应 Java: DeCatalogInfoServiceImpl.update()
	catalog.CatalogName = req.CatalogName
	catalog.Description = req.Description
	catalog.UpdateTime = time.Now()

	err = l.svcCtx.CatalogModel.Update(l.ctx, catalog)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("UpdateCatalog Update failed: %v", err)
		return nil, baseerrorx.NewWithMsg(errorx.ErrCodeCatalogError, "修改目录失败")
	}

	// Step 6: 返回成功响应
	return &types.EmptyResp{}, nil
}
