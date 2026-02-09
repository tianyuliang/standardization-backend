// Code scaffolded by goctl. Safe to edit.

package dataelement

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/dataelement/mock"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/dataelement/dataelement"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetDataElementInternalLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 内部查看数据元详情
func NewGetDataElementInternalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDataElementInternalLogic {
	return &GetDataElementInternalLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDataElementInternalLogic) GetDataElementInternal(req *types.IdReq) (resp *types.DataElementDetailVo, err error) {
	// Step 1: 根据ID查询数据元
	de, err := l.svcCtx.DataElementModel.FindOne(l.ctx, req.Id)
	if err != nil {
		logx.Errorf("内部查询数据元详情失败: %v", err)
		return nil, err
	}
	if de == nil {
		return nil, errorx.DataNotExist()
	}

	// Step 2: 获取目录名称
	catalogName, _ := mock.GetCatalogName(l.ctx, de.CatalogId)

	// Step 3: 计算值域
	valueRange := ""
	dataLen := 0
	dataPrec := 0
	if de.DataLength != nil {
		dataLen = *de.DataLength
	}
	if de.DataPrecision != nil {
		dataPrec = *de.DataPrecision
	}

	if de.DictCode != nil {
		dictCodeStr := "" // TODO: 转换
		valueRange = CalculateValueRange(l.ctx, dictCodeStr, de.DataType, dataLen, dataPrec)
	} else {
		valueRange = CalculateValueRange(l.ctx, "", de.DataType, dataLen, dataPrec)
	}

	// Step 4: 获取关联的码表名称
	dictName := ""
	if de.RelationType == "codeTable" {
		// TODO: 从Dict服务获取码表名称
		dictName = "码表名称"
	}

	// Step 5: 获取关联的规则名称
	ruleName := ""
	if de.RelationType == "codeRule" && de.RuleId != nil {
		// TODO: 从Rule服务获取规则名称
		ruleName = "规则名称"
	}

	// Step 6: 获取关联的标准文件
	var stdFiles []types.StdFileVo
	fileIds, err := l.svcCtx.RelationDeFileModel.FindFileIdsByDeId(l.ctx, de.Id)
	if err == nil && len(fileIds) > 0 {
		fileNames, _ := mock.GetStdFileByIds(l.ctx, fileIds)
		for _, fileId := range fileIds {
			stdFiles = append(stdFiles, types.StdFileVo{
				Id:       fileId,
				FileName: fileNames[fileId],
			})
		}
	}

	// Step 7: 构建响应
	dataLen32 := int32(0)
	dataPrec32 := int32(0)
	ruleId := int64(0)
	labelId := int64(0)
	if de.DataLength != nil {
		dataLen32 = int32(*de.DataLength)
	}
	if de.DataPrecision != nil {
		dataPrec32 = int32(*de.DataPrecision)
	}
	if de.RuleId != nil {
		ruleId = *de.RuleId
	}
	if de.LabelId != nil {
		labelId = *de.LabelId
	}

	return &types.DataElementDetailVo{
		Id:             de.Id,
		Code:           de.Code,
		NameEn:         de.NameEn,
		NameCn:         de.NameCn,
		Synonym:        de.Synonym,
		StdType:        de.StdType,
		DataType:       de.DataType,
		DataLength:     dataLen32,
		DataPrecision:  dataPrec32,
		DictCode:       "", // TODO: 转换
		RuleId:         ruleId,
		RelationType:   de.RelationType,
		CatalogId:      de.CatalogId,
		LabelId:        labelId,
		Description:    de.Description,
		Version:        int32(de.Version),
		State:          dataelement.IntToState(de.State),
		DepartmentIds:  de.DepartmentIds,
		ThirdDeptId:    de.ThirdDeptId,
		DisableReason:  de.DisableReason,
		CatalogName:    catalogName,
		DepartmentPath: "", // TODO: 从部门服务获取
		ValueRange:     valueRange,
		DictName:       dictName,
		RuleName:       ruleName,
		StdFiles:       stdFiles,
		CreateUser:     de.CreateUser,
		CreateTime:     de.CreateTime,
		UpdateTime:     de.UpdateTime,
	}, nil
}
