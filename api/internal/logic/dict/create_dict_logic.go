// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dict

import (
	"context"

	localErrorx "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	dictmodel "github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/dict/dict"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateDictLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 新增码表
func NewCreateDictLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDictLogic {
	return &CreateDictLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateDictLogic) CreateDict(req *types.CreateDictReq) (resp *types.DictVo, err error) {
	// Step 1: 参数校验
	if len(req.Enums) == 0 {
		return nil, localErrorx.DictParamEmpty("码值列表不能为空")
	}

	// Step 2: 业务校验（目录存在性、名称唯一性、码值唯一性）
	// 校验中文名称唯一性
	if err := checkChNameUnique(l.svcCtx.DictModel, req.ChName, req.OrgType); err != nil {
		return nil, err
	}

	// 校验英文名称唯一性
	if err := checkEnNameUnique(l.svcCtx.DictModel, req.EnName, req.OrgType); err != nil {
		return nil, err
	}

	// 校验码值唯一性
	if err := checkEnumCodesUnique(req.Enums); err != nil {
		return nil, err
	}

	// Step 3: 生成码表编码（雪花算法）
	code := generateSnowflakeCode()

	// 获取用户信息
	username, _ := getUserInfo(l.ctx)

	// Step 4: 保存码表
	dictData := &dictmodel.Dict{
		Code:          code,
		ChName:        req.ChName,
		EnName:        req.EnName,
		Description:   req.Description,
		CatalogId:     req.CatalogId,
		OrgType:       req.OrgType,
		Version:       1,
		State:         stateToInt("enable"),
		DepartmentIds: req.DepartmentIds,
		CreateUser:    username,
		UpdateUser:    username,
	}

	id, err := l.svcCtx.DictModel.Insert(l.ctx, dictData)
	if err != nil {
		logx.Errorf("保存码表失败: %v", err)
		return nil, localErrorx.DictInvalidParam("保存码表失败")
	}

	// Step 5: 保存码值明细
	for _, enum := range req.Enums {
		enumData := &dictmodel.DictEnum{
			DictId:     id,
			Code:       enum.Code,
			Value:      enum.Value,
			CreateUser: username,
		}
		if _, err := l.svcCtx.DictEnumModel.Insert(l.ctx, enumData); err != nil {
			logx.Errorf("保存码值失败: %v", err)
			return nil, localErrorx.DictInvalidParam("保存码值失败")
		}
	}

	// Step 6: 保存关联文件关系
	if len(req.StdFiles) > 0 {
		var relations []*dictmodel.RelationDictFile
		for _, fileId := range req.StdFiles {
			relations = append(relations, &dictmodel.RelationDictFile{
				Id:     generateSnowflakeCode(),
				DictId: id,
				FileId: fileId,
			})
		}
		if err := l.svcCtx.RelationDictFileModel.InsertBatch(l.ctx, relations); err != nil {
			logx.Errorf("保存关联文件失败: %v", err)
			return nil, localErrorx.DictInvalidParam("保存关联文件失败")
		}
	}

	// 查询完整数据返回
	dictData.Id = id
	enums, _ := l.svcCtx.DictEnumModel.FindByDictId(l.ctx, id)
	catalogName := getCatalogName(req.CatalogId)
	deptName, deptPathNames, _ := getDeptInfo(req.DepartmentIds)

	result := buildDictVo(l.ctx, dictData, enums, catalogName, deptName, deptPathNames, false)

	logx.Infof("码表新增成功: id=%d, chName=%s", id, req.ChName)
	return &result, nil
}
