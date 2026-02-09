// Code scaffolded by goctl. Safe to edit.

package dataelement

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/mock"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/dataelement/dataelement"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/dataelement/relation"
	"github.com/sony/sonyflake"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateDataElementLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建数据元
func NewCreateDataElementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDataElementLogic {
	return &CreateDataElementLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateDataElementLogic) CreateDataElement(req *types.CreateDataElementReq) (resp *types.DataElementDetailVo, err error) {
	// Step 1: 校验目录存在性和类型
	exists, catalogName, err := mock.CheckCatalogExist(l.ctx, req.CatalogId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errorx.CatalogNotExist()
	}

	// Step 2: 校验中文名称唯一性（同stdType下）
	exist, err := l.svcCtx.DataElementModel.CheckNameCnExists(l.ctx, req.NameCn, req.StdType, 0, req.DepartmentIds)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, errorx.NameCnDuplicate()
	}

	// Step 3: 校验英文名称唯一性（同部门下）
	exist, err = l.svcCtx.DataElementModel.CheckNameEnExists(l.ctx, req.NameEn, 0, req.DepartmentIds)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, errorx.NameEnDuplicate()
	}

	// Step 4: 校验关联类型和关联对象
	if req.RelationType == "codeTable" && req.DictCode == "" {
		return nil, errorx.DictIdEmpty()
	}
	if req.RelationType == "codeRule" && req.RuleId == 0 {
		return nil, errorx.RuleIdEmpty()
	}

	// Step 5: 校验码表是否存在
	if req.RelationType == "codeTable" && req.DictCode != "" {
		exist, err = mock.CheckDictExist(l.ctx, req.DictCode)
		if err != nil {
			return nil, err
		}
		if !exist {
			return nil, errorx.DictNotExist()
		}
	}

	// Step 6: 校验规则是否存在
	if req.RelationType == "codeRule" && req.RuleId > 0 {
		exist, err = mock.CheckRuleExist(l.ctx, req.RuleId)
		if err != nil {
			return nil, err
		}
		if !exist {
			return nil, errorx.RuleNotExist()
		}
	}

	// Step 7: 校验文件ID是否存在
	// TODO: 实际需要调用StdFile服务验证

	// Step 8: 生成ID和Code
	sf := sonyflake.NewSonyflake(sonyflake.Settings{})
	idUint64, err := sf.NextID()
	if err != nil {
		return nil, err
	}
	id := int64(idUint64)

	// Step 9: 构建数据元实体
	dataElement := &dataelement.DataElement{
		Id:            id,
		Code:          id, // f_de_id equals f_id
		NameEn:        req.NameEn,
		NameCn:        req.NameCn,
		Synonym:       req.Synonym,
		StdType:       req.StdType,
		DataType:      req.DataType,
		RelationType:  req.RelationType,
		CatalogId:     req.CatalogId,
		Description:   req.Description,
		Version:       1,
		State:         dataelement.StateToInt("enable"),
		DepartmentIds: req.DepartmentIds,
		ThirdDeptId:   req.ThirdDeptId,
		CreateUser:    "system", // TODO: 从token获取
		UpdateUser:    "system",
	}

	// 处理可选字段
	if req.DataLength > 0 {
		dl := int(req.DataLength)
		dataElement.DataLength = &dl
	}
	if req.DataPrecision > 0 {
		dp := int(req.DataPrecision)
		dataElement.DataPrecision = &dp
	}
	if req.DictCode != "" {
		dictCode := int64(0) // TODO: 根据实际DictCode转换
		dataElement.DictCode = &dictCode
	}
	if req.RuleId > 0 {
		dataElement.RuleId = &req.RuleId
	}
	if req.LabelId > 0 {
		dataElement.LabelId = &req.LabelId
	}

	// Step 10: 保存数据元基本信息
	_, err = l.svcCtx.DataElementModel.Insert(l.ctx, dataElement)
	if err != nil {
		logx.Errorf("保存数据元失败: %v", err)
		return nil, err
	}

	// Step 11: 保存关联文件关系
	if len(req.StdFiles) > 0 {
		relations := make([]*relation.RelationDeFile, len(req.StdFiles))
		for i, fileId := range req.StdFiles {
			relations[i] = &relation.RelationDeFile{
				DeId:   id,
				FileId: fileId,
			}
		}
		err = l.svcCtx.RelationDeFileModel.InsertBatch(l.ctx, relations)
		if err != nil {
			logx.Errorf("保存关联文件关系失败: %v", err)
		}
	}

	// Step 12: 发送MQ消息
	// TODO: 实际实现MQ消息发送
	_ = mock.SendMqMessage(l.ctx, "create", map[string]interface{}{"id": id}, "system")

	// Step 13: 构建响应
	dataLen := int32(0)
	dataPrec := int32(0)
	ruleId := int64(0)
	labelId := int64(0)
	if dataElement.DataLength != nil {
		dataLen = int32(*dataElement.DataLength)
	}
	if dataElement.DataPrecision != nil {
		dataPrec = int32(*dataElement.DataPrecision)
	}
	if dataElement.RuleId != nil {
		ruleId = *dataElement.RuleId
	}
	if dataElement.LabelId != nil {
		labelId = *dataElement.LabelId
	}

	vo := &types.DataElementDetailVo{
		Id:            dataElement.Id,
		Code:          dataElement.Code,
		NameEn:        dataElement.NameEn,
		NameCn:        dataElement.NameCn,
		Synonym:       dataElement.Synonym,
		StdType:       dataElement.StdType,
		DataType:      dataElement.DataType,
		DataLength:    dataLen,
		DataPrecision: dataPrec,
		DictCode:      "", // TODO: 转换
		RuleId:        ruleId,
		RelationType:  dataElement.RelationType,
		CatalogId:     dataElement.CatalogId,
		LabelId:       labelId,
		Description:   dataElement.Description,
		Version:       int32(dataElement.Version),
		State:         dataelement.IntToState(dataElement.State),
		DepartmentIds: dataElement.DepartmentIds,
		ThirdDeptId:   dataElement.ThirdDeptId,
		CatalogName:   catalogName,
		CreateTime:    dataElement.CreateTime,
		UpdateTime:    dataElement.UpdateTime,
		ValueRange:    CalculateValueRange(l.ctx, req.DictCode, req.DataType, int(req.DataLength), int(req.DataPrecision)),
	}

	return vo, nil
}
