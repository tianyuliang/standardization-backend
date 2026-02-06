// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rule

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRuleInternalLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 内部-根据ID查看规则详情
//
// 业务流程:
//
//	与 GET /v1/rule/{id} 实现相同，仅路由不同
func NewGetRuleInternalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRuleInternalLogic {
	return &GetRuleInternalLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRuleInternalLogic) GetRuleInternal(id int64) (resp *types.RuleResp, err error) {
	// 查询规则详情
	ruleData, err := l.svcCtx.RuleModel.FindOne(l.ctx, id)
	if err != nil {
		return nil, err
	}

	// 构建响应（TODO: 查询目录名称、关联文件、引用状态）
	resp = buildRuleResp(ruleData, "", false, nil)
	return
}
