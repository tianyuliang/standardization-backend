// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rule

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryRuleByIdsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量查询规则
//
// 业务流程:
//   1. 参数校验（Ids非空）
//   2. 批量查询规则详情
//   3. 构建响应
func NewQueryRuleByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryRuleByIdsLogic {
	return &QueryRuleByIdsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryRuleByIdsLogic) QueryRuleByIds(req *types.QueryByIdsReq) (resp *types.RuleListResp, err error) {
	// ====== 步骤1: 参数校验 ======
	if len(req.Ids) == 0 {
		// TODO: 返回 errorx.RuleQueryIdsEmpty()
		return &types.RuleListResp{
			Entries:    []types.RuleResp{},
			TotalCount: 0,
		}, nil
	}

	// ====== 步骤2: 批量查询规则详情 ======
	rules, err := l.svcCtx.RuleModel.FindByIds(l.ctx, req.Ids)
	if err != nil {
		return nil, err
	}

	// ====== 步骤3: 构建响应 ======
	entries := make([]types.RuleResp, 0, len(rules))
	for _, r := range rules {
		item := types.RuleResp{
			Id:             r.Id,
			Name:           r.Name,
			CatalogId:      r.CatalogId,
			CatalogName:    "", // TODO: 查询目录名称
			OrgType:        r.OrgType,
			Description:    r.Description,
			RuleType:       intToRuleType(r.RuleType),
			Version:        r.Version,
			Expression:     r.Expression,
			State:          r.State,
			DisableReason:  r.DisableReason,
			AuthorityId:    r.AuthorityId,
			DepartmentIds:  r.DepartmentIds,
			ThirdDeptId:    r.ThirdDeptId,
			CreateTime:     timeToStr(r.CreateTime),
			CreateUser:     r.CreateUser,
			UpdateTime:     timeToStr(r.UpdateTime),
			UpdateUser:     r.UpdateUser,
			StdFileIds:     nil,
			Used:           false, // TODO: 查询引用状态
		}
		entries = append(entries, item)
	}

	return &types.RuleListResp{
		Entries:    entries,
		TotalCount: int64(len(entries)),
	}, nil
}
