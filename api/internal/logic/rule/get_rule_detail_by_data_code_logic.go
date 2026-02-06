// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rule

import (
	"context"

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
	// TODO: 调用 DataElement RPC 获取 ruleId
	ruleId := int64(0) // getMockRuleIdByDataCode(dataCode)
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

	resp := buildRuleResp(ruleData, "", false, nil)
	return resp, nil
}
