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

	// Step 1.2: 校验父目录存在（parentId必填，对应Java中不能为null）
	// 对应 Java: DeCatalogInfoServiceImpl.checkPost() - parentId校验
	if req.ParentId == 0 {
		return nil, baseerrorx.NewWithMsg(errorx.ErrCodeCatalogMissingParam, "父目录ID不能为空")
	}

	parent, err := l.svcCtx.CatalogModel.FindOne(l.ctx, req.ParentId)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("CreateCatalog FindOne parent failed: %v", err)
		return nil, errorx.CatalogParentNotExist()
	}

	// Step 1.3: 校验父目录类型（如果有传入type）
	// 对应 Java: DeCatalogInfoServiceImpl.checkPost() - 类型一致性校验
	// 注意：Java中创建时type会自动继承父目录，不需要客户端传递
	// Go中如果客户端传递了type，需要校验是否与父目录一致
	if req.Type != 0 && parent.Type != req.Type {
		return nil, errorx.CatalogTypeMismatch()
	}

	// Step 1.4: 校验目录级别不超限
	// 对应 Java: DeCatalogInfoServiceImpl.checkPost() - level校验
	// Java: parent.level >= 255 时报错（新目录level会是256）
	if parent.Level >= 255 {
		return nil, errorx.CatalogLevelOutOfRange()
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
	level := parent.Level + 1
	catalogType := parent.Type // 自动继承父目录类型

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
