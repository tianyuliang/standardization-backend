// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rule

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/rule/mock"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRuleDetailByDataCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 内部-根据数据元编码查看规则详情
//
// 对应 Java: RuleServiceImpl.getDetailByDataCode(Long dataCode) (lines 763-769)
// 业务流程:
//  1. 根据数据元编码查询规则ID
//  2. 查询规则详情
func NewGetRuleDetailByDataCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRuleDetailByDataCodeLogic {
	return &GetRuleDetailByDataCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRuleDetailByDataCodeLogic) GetRuleDetailByDataCode(dataCode string) (resp *types.RuleResp, err error) {
	// ====== 步骤1: 根据数据元编码查询规则ID ======
	// 对应 Java: DataElementInfo dataElementInfo = iDataElementInfoService.getOneByIdOrCode(2, dataCode) (line 764)
	//            return this.queryById(dataElementInfo.getRuleId()) (line 768)
	// MOCK: mock.DataElementGetRuleIdByDataCode() - 根据数据元编码获取规则ID
	ruleId := mock.DataElementGetRuleIdByDataCode(l.ctx, l.svcCtx, dataCode)
	if ruleId == 0 {
		return nil, nil // 数据元不存在或无关联规则
	}

	// ====== 步骤2: 查询规则详情 ======
	return l.getRuleById(ruleId)
}

// getRuleById 根据规则ID查询详情
func (l *GetRuleDetailByDataCodeLogic) getRuleById(ruleId int64) (*types.RuleResp, error) {
	ruleData, err := l.svcCtx.RuleModel.FindOne(l.ctx, ruleId)
	if err != nil {
		return nil, err
	}

	// 查询关联文件
	relations, _ := l.svcCtx.RelationRuleFileModel.FindByRuleId(l.ctx, ruleId)
	stdFileIds := make([]int64, len(relations))
	for i, r := range relations {
		stdFileIds[i] = r.FileId
	}

	// 查询目录名称
	catalogName := mock.CatalogGetCatalogName(l.ctx, l.svcCtx, ruleData.CatalogId)

	// 查询是否被引用
	usedFlag := mock.DataElementRuleUsed(l.ctx, l.svcCtx, ruleId)

	resp := buildRuleResp(ruleData, catalogName, usedFlag, stdFileIds)
	return resp, nil
}
