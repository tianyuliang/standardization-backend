// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rule

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 编码规则列表查询
func NewListRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRuleLogic {
	return &ListRuleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListRuleLogic) ListRule(req *types.RuleListQuery) (resp *types.RuleListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
