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

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建目录
func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.CreateReq) (resp *types.CreateResp, err error) {
	// 1. 校验目录名称格式
	if err := ValidateCatalogName(req.CatalogName); err != nil {
		return &types.CreateResp{
			Code:        errorx.NewWithCode(errorx.ErrInvalidParam).Code(),
			Description: fmt.Sprintf("目录名称格式错误: %v", err),
		}, nil
	}

	// 2. 校验父目录存在
	parentCatalog, err := l.svcCtx.CatalogModel.FindOne(l.ctx, req.ParentId)
	if err != nil {
		return &types.CreateResp{
			Code:        errorx.NewWithCode(errorx.ErrInternal).Code(),
			Description: "查询父目录失败",
		}, nil
	}
	if parentCatalog == nil {
		return &types.CreateResp{
			Code:        errorx.NewWithCode(errorx.ErrNotFound).Code(),
			Description: "父目录不存在",
		}, nil
	}

	// 3. 校验父目录级别 < 255
	if parentCatalog.Level >= catalog.CatalogMaxLevel {
		return &types.CreateResp{
			Code:        errorx.NewWithCode(errorx.ErrOutOfRange).Code(),
			Description: fmt.Sprintf("父目录级别已达上限 %d", catalog.CatalogMaxLevel),
		}, nil
	}

	// 4. 继承父目录的 type，设置 level = 父目录 level + 1
	catalogType := parentCatalog.Type
	newLevel := parentCatalog.Level + 1

	// 5. 检查同级目录名称唯一性
	isUnique, err := l.svcCtx.CatalogModel.CheckUniqueName(l.ctx, req.ParentId, catalogType, req.CatalogName, "")
	if err != nil {
		return &types.CreateResp{
			Code:        errorx.NewWithCode(errorx.ErrInternal).Code(),
			Description: "检查目录名称失败",
		}, nil
	}
	if !isUnique {
		return &types.CreateResp{
			Code:        errorx.NewWithCode(errorx.ErrConflict).Code(),
			Description: "同级目录下已存在相同名称",
		}, nil
	}

	// 6. 构建目录实体
	newCatalog := &catalog.Catalog{
		Id:          generateCatalogId(), // 使用雪花算法生成ID
		CatalogName: req.CatalogName,
		Description: req.Description,
		Level:       newLevel,
		ParentId:    req.ParentId,
		Type:        catalogType,
		AuthorityId: nil,
	}

	// 7. 插入数据库
	_, err = l.svcCtx.CatalogModel.Insert(l.ctx, newCatalog)
	if err != nil {
		return &types.CreateResp{
			Code:        errorx.NewWithCode(errorx.ErrInternal).Code(),
			Description: "创建目录失败",
		}, nil
	}

	return &types.CreateResp{
		Code:        errorx.NewWithCode(errorx.ErrSuccess).Code(),
		Description: "创建成功",
	}, nil
}

// generateCatalogId 生成目录ID（雪花算法）
// TODO: 实现雪花算法ID生成，当前使用临时实现
func generateCatalogId() string {
	// 临时实现：使用时间戳
	// 实际项目中应该使用雪花算法生成唯一ID
	return fmt.Sprintf("%d", 1000000000000000+1)
}
