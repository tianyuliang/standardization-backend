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

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改目录
func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.UpdateReq) (resp *types.UpdateResp, err error) {
	// 1. 校验父目录 ID 不为空（对应 Java lines 631-634）
	if req.ParentId == "" {
		return &types.UpdateResp{
			Code:        errorx.NewWithCode(errorx.ErrInvalidParam).Code(),
			Description: "父目录ID不能为空",
		}, nil
	}

	// 2. 校验目录名称格式
	if err := ValidateCatalogName(req.CatalogName); err != nil {
		return &types.UpdateResp{
			Code:        errorx.NewWithCode(errorx.ErrInvalidParam).Code(),
			Description: fmt.Sprintf("目录名称格式错误: %v", err),
		}, nil
	}

	// 3. 校验目录存在
	currentCatalog, err := l.svcCtx.CatalogModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return &types.UpdateResp{
			Code:        errorx.NewWithCode(errorx.ErrInternal).Code(),
			Description: "查询目录失败",
		}, nil
	}
	if currentCatalog == nil {
		return &types.UpdateResp{
			Code:        errorx.NewWithCode(errorx.ErrNotFound).Code(),
			Description: "目录不存在",
		}, nil
	}

	// 4. 校验不是根目录（level <= 1 不允许修改）
	if currentCatalog.Level <= catalog.CatalogRootLevel {
		return &types.UpdateResp{
			Code:        errorx.NewWithCode(errorx.ErrInvalidParam).Code(),
			Description: "根目录不允许修改",
		}, nil
	}

	// 5. 校验新父目录存在
	newParentCatalog, err := l.svcCtx.CatalogModel.FindOne(l.ctx, req.ParentId)
	if err != nil {
		return &types.UpdateResp{
			Code:        errorx.NewWithCode(errorx.ErrInternal).Code(),
			Description: "查询父目录失败",
		}, nil
	}
	if newParentCatalog == nil {
		return &types.UpdateResp{
			Code:        errorx.NewWithCode(errorx.ErrNotFound).Code(),
			Description: "父目录不存在",
		}, nil
	}

	// 6. 校验新父目录级别 < 255
	if newParentCatalog.Level >= catalog.CatalogMaxLevel {
		return &types.UpdateResp{
			Code:        errorx.NewWithCode(errorx.ErrOutOfRange).Code(),
			Description: fmt.Sprintf("父目录级别已达上限 %d", catalog.CatalogMaxLevel),
		}, nil
	}

	// 7. 循环检测：新父目录不能是自身子目录
	isDescendant, err := l.svcCtx.CatalogModel.IsDescendant(l.ctx, req.Id, req.ParentId)
	if err != nil {
		return &types.UpdateResp{
			Code:        errorx.NewWithCode(errorx.ErrInternal).Code(),
			Description: "检查目录关系失败",
		}, nil
	}
	if isDescendant {
		return &types.UpdateResp{
			Code:        errorx.NewWithCode(errorx.ErrInvalidParam).Code(),
			Description: "不能将目录移动到其子目录下",
		}, nil
	}

	// 8. 类型一致性：新父目录 type 必须与当前目录一致
	if newParentCatalog.Type != currentCatalog.Type {
		return &types.UpdateResp{
			Code: errorx.NewWithCode(errorx.ErrInvalidParam).Code(),
			Description: fmt.Sprintf("目录类型不一致，当前: %s, 新父目录: %s",
				catalog.GetCatalogTypeText(currentCatalog.Type),
				catalog.GetCatalogTypeText(newParentCatalog.Type)),
		}, nil
	}

	// 9. 检查同级名称唯一性（排除自身）
	isUnique, err := l.svcCtx.CatalogModel.CheckUniqueName(l.ctx, req.ParentId, currentCatalog.Type, req.CatalogName, req.Id)
	if err != nil {
		return &types.UpdateResp{
			Code:        errorx.NewWithCode(errorx.ErrInternal).Code(),
			Description: "检查目录名称失败",
		}, nil
	}
	if !isUnique {
		return &types.UpdateResp{
			Code:        errorx.NewWithCode(errorx.ErrConflict).Code(),
			Description: "同级目录下已存在相同名称",
		}, nil
	}

	// 10. 更新目录信息
	currentCatalog.CatalogName = req.CatalogName
	currentCatalog.ParentId = req.ParentId
	currentCatalog.Description = req.Description

	err = l.svcCtx.CatalogModel.Update(l.ctx, currentCatalog)
	if err != nil {
		return &types.UpdateResp{
			Code:        errorx.NewWithCode(errorx.ErrInternal).Code(),
			Description: "更新目录失败",
		}, nil
	}

	return &types.UpdateResp{
		Code:        errorx.NewWithCode(errorx.ErrSuccess).Code(),
		Description: "更新成功",
	}, nil
}
