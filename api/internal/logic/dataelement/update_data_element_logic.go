// Code scaffolded by goctl. Safe to edit.

package dataelement

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/dataelement/mock"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/dataelement/dataelement"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/dataelement/relation"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateDataElementLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 编辑数据元
func NewUpdateDataElementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDataElementLogic {
	return &UpdateDataElementLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateDataElementLogic) UpdateDataElement(req *types.UpdateDataElementReq) (resp *types.DataElementDetailVo, err error) {
	// Step 1: 校验数据元存在性
	oldData, err := l.svcCtx.DataElementModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if oldData == nil {
		return nil, errorx.DataNotExist()
	}

	// Step 2: 校验目录存在性
	if req.CatalogId > 0 {
		exists, _, err := mock.CheckCatalogExist(l.ctx, req.CatalogId)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, errorx.CatalogNotExist()
		}
	}

	// Step 3: 校验关联类型和关联对象
	if req.RelationType == "codeTable" && req.DictCode == "" {
		return nil, errorx.DictIdEmpty()
	}
	if req.RelationType == "codeRule" && req.RuleId == 0 {
		return nil, errorx.RuleIdEmpty()
	}

	// Step 4: 校验码表是否存在
	if req.RelationType == "codeTable" && req.DictCode != "" {
		exist, err := mock.CheckDictExist(l.ctx, req.DictCode)
		if err != nil {
			return nil, err
		}
		if !exist {
			return nil, errorx.DictNotExist()
		}
	}

	// Step 5: 校验规则是否存在
	if req.RelationType == "codeRule" && req.RuleId > 0 {
		exist, err := mock.CheckRuleExist(l.ctx, req.RuleId)
		if err != nil {
			return nil, err
		}
		if !exist {
			return nil, errorx.RuleNotExist()
		}
	}

	// Step 6: 判断是否需要递增版本号
	needVersionInc := CheckVersionChange(oldData, req)

	// Step 7: 构建更新实体
	dataElement := &dataelement.DataElement{
		Id:            req.Id,
		NameEn:        req.NameEn,
		NameCn:        req.NameCn,
		Synonym:       req.Synonym,
		StdType:       req.StdType,
		DataType:      req.DataType,
		RelationType:  req.RelationType,
		CatalogId:     req.CatalogId,
		Description:   req.Description,
		DepartmentIds: req.DepartmentIds,
		ThirdDeptId:   req.ThirdDeptId,
		UpdateUser:    "system", // TODO: 从token获取
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
		dictCode := int64(0) // TODO: 转换
		dataElement.DictCode = &dictCode
	}
	if req.RuleId > 0 {
		dataElement.RuleId = &req.RuleId
	}
	if req.LabelId > 0 {
		dataElement.LabelId = &req.LabelId
	}
	if needVersionInc {
		dataElement.Version = oldData.Version + 1
	} else {
		dataElement.Version = oldData.Version
	}
	dataElement.State = oldData.State

	// Step 8: 更新数据元基本信息
	err = l.svcCtx.DataElementModel.Update(l.ctx, dataElement)
	if err != nil {
		logx.Errorf("更新数据元失败: %v", err)
		return nil, err
	}

	// Step 9: 更新关联文件关系
	if req.StdFiles != nil {
		// 先删除旧关系
		_ = l.svcCtx.RelationDeFileModel.DeleteByDeId(l.ctx, req.Id)

		// 添加新关系
		if len(req.StdFiles) > 0 {
			relations := make([]*relation.RelationDeFile, len(req.StdFiles))
			for i, fileId := range req.StdFiles {
				relations[i] = &relation.RelationDeFile{
					DeId:   req.Id,
					FileId: fileId,
				}
			}
			err = l.svcCtx.RelationDeFileModel.InsertBatch(l.ctx, relations)
			if err != nil {
				logx.Errorf("保存关联文件关系失败: %v", err)
			}
		}
	}

	// Step 10: 发送MQ消息
	_ = mock.SendMqMessage(l.ctx, "update", map[string]interface{}{"id": req.Id}, "system")

	// Step 11: 查询更新后的详情
	detailReq := &types.DataElementDetailReq{
		Type:  1,
		Value: req.Id,
	}

	detailLogic := NewGetDataElementDetailLogic(l.ctx, l.svcCtx)
	return detailLogic.GetDataElementDetail(detailReq)
}
