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

type GetRuleDetailByDataIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 内部-根据数据元ID查看规则详情
//
// 对应 Java: RuleServiceImpl.getDetailByDataId(Long dataId) (lines 756-760)
// 业务流程:
//  1. 根据数据元ID查询规则ID
//  2. 查询规则详情
func NewGetRuleDetailByDataIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRuleDetailByDataIdLogic {
	return &GetRuleDetailByDataIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRuleDetailByDataIdLogic) GetRuleDetailByDataId(dataId int64) (resp *types.RuleResp, err error) {
	// ====== 步骤1: 根据数据元ID查询规则ID ======
	// 对应 Java: DataElementInfo dataElementInfo = iDataElementInfoService.getById(dataId) (line 757)
	//            return this.queryById(dataElementInfo.getRuleId()) (line 760)
	// MOCK: mock.DataElementGetRuleIdByDataId() - 根据数据元ID获取规则ID
	ruleId := mock.DataElementGetRuleIdByDataId(l.ctx, l.svcCtx, dataId)
	if ruleId == 0 {
		return nil, nil // 数据元不存在或无关联规则
	}

	// ====== 步骤2: 查询规则详情 ======
	// 对应 Java: return this.queryById(dataElementInfo.getRuleId()) (line 760)
	return l.getRuleById(ruleId)
}

// getRuleById 根据规则ID查询详情
func (l *GetRuleDetailByDataIdLogic) getRuleById(ruleId int64) (*types.RuleResp, error) {
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
