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
	catalogmodel "github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/catalog/catalog"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCatalogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建目录
func NewCreateCatalogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCatalogLogic {
	return &CreateCatalogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCatalogLogic) CreateCatalog(req *types.CreateCatalogReq) (resp *types.CatalogResp, err error) {
	// Step 1: 参数校验
	// 对应 Java: DeCatalogInfoController.checkPost()

	// Step 1.1: 校验目录名称
	// 对应 Java: DeCatalogInfoServiceImpl.checkPost() - 名称校验
	if err := ValidateCatalogName(req.CatalogName); err != nil {
		return nil, err
	}

	// Step 1.2: 校验目录类型
	// 对应 Java: DeCatalogInfoServiceImpl.checkType()
	if err := CheckType(req.Type); err != nil {
		return nil, err
	}

	// Step 1.3: 校验父目录存在
	// 对应 Java: DeCatalogInfoServiceImpl.checkParent()
	var parent *catalogmodel.Catalog
	if req.ParentId != 0 {
		parent, err = l.svcCtx.CatalogModel.FindOne(l.ctx, req.ParentId)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("CreateCatalog FindOne parent failed: %v", err)
			return nil, errorx.CatalogParentNotExist()
		}

		// Step 1.4: 校验父目录类型与当前类型一致
		// 对应 Java: DeCatalogInfoServiceImpl.checkTypeMatch()
		if parent.Type != req.Type {
			return nil, errorx.CatalogTypeMismatch()
		}

		// Step 1.5: 校验目录级别不超限
		// 对应 Java: DeCatalogInfoServiceImpl.checkLevel()
		if parent.Level >= 255 {
			return nil, errorx.CatalogLevelOutOfRange()
		}
	}

	// Step 1.6: 校验同级目录名称不重复
	// 对应 Java: DeCatalogInfoServiceImpl.checkNameDuplicate()
	siblings, err := l.svcCtx.CatalogModel.FindByParentId(l.ctx, req.ParentId)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("CreateCatalog FindByParentId failed: %v", err)
		return nil, baseerrorx.NewWithMsg(errorx.ErrCodeCatalogError, "查询目录失败")
	}

	for _, sibling := range siblings {
		if sibling.CatalogName == req.CatalogName {
			return nil, errorx.CatalogNameDuplicate()
		}
	}

	// Step 2: 创建目录
	// 对应 Java: DeCatalogInfoServiceImpl.insert()

	// Step 2.1: 计算目录级别和类型
	// 对应 Java: 继承父目录的type，level=父目录level+1
	level := int32(1) // 默认为根目录级别
	catalogType := req.Type

	if parent != nil {
		level = parent.Level + 1
		catalogType = parent.Type // 继承父目录类型
	}

	// Step 2.2: 构建目录模型
	now := time.Now()
	catalog := &catalogmodel.Catalog{
		CatalogName: req.CatalogName,
		Description: req.Description,
		Level:       level,
		ParentId:    req.ParentId,
		Type:        catalogType,
		CreateTime:  now,
		UpdateTime:  now,
		Deleted:     0,
	}

	// Step 2.3: 插入数据库
	id, err := l.svcCtx.CatalogModel.Insert(l.ctx, catalog)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("CreateCatalog Insert failed: %v", err)
		return nil, baseerrorx.NewWithMsg(errorx.ErrCodeCatalogError, "创建目录失败")
	}

	catalog.Id = id

	// Step 3: 返回结果
	// 对应 Java: 返回 CatalogResp
	return modelToResp(catalog), nil
}
