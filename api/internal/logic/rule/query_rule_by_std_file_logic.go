// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rule

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryRuleByStdFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据标准文件查询规则
//
// 业务流程:
//   1. 查询关联该文件的规则ID列表
//   2. 根据ID列表批量查询规则详情
func NewQueryRuleByStdFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryRuleByStdFileLogic {
	return &QueryRuleByStdFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryRuleByStdFileLogic) QueryRuleByStdFile(fileId int64, req *types.RuleListQuery) (resp *types.RuleListResp, err error) {
	// ====== 步骤1: 查询关联该文件的规则ID列表 ======
	relations, err := l.svcCtx.RelationRuleFileModel.FindByFileId(l.ctx, fileId)
	if err != nil {
		return nil, err
	}

	if len(relations) == 0 {
		return &types.RuleListResp{
			Entries:    []types.RuleResp{},
			TotalCount: 0,
		}, nil
	}

	// ====== 步骤2: 批量查询规则详情 ======
	ruleIds := make([]int64, 0, len(relations))
	for _, r := range relations {
		ruleIds = append(ruleIds, r.RuleId)
	}

	rules, err := l.svcCtx.RuleModel.FindByIds(l.ctx, ruleIds)
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
