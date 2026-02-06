// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rule

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/rule/mock"
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
// 对应 Java: RuleServiceImpl.queryByIds(List<Long> ids) (lines 162-176)
// 业务流程:
//  1. 参数校验（Ids非空）
//  2. 批量查询规则详情
//  3. 构建响应（解析expression为regex或custom）
func NewQueryRuleByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryRuleByIdsLogic {
	return &QueryRuleByIdsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryRuleByIdsLogic) QueryRuleByIds(req *types.QueryByIdsReq) (resp *types.RuleListResp, err error) {
	// ====== 步骤1: 参数校验 ======
	// 对应 Java: (implicitly checked)
	if len(req.Ids) == 0 {
		return nil, errorx.RuleQueryIdsEmpty()
	}

	// ====== 步骤2: 批量查询规则详情 ======
	// 对应 Java: ruleMapper.queryByIds(ids) (line 163)
	rules, err := l.svcCtx.RuleModel.FindByIds(l.ctx, req.Ids)
	if err != nil {
		return nil, err
	}

	// ====== 步骤3: 构建响应 ======
	// 对应 Java: for循环构建RuleVo (lines 165-174)
	entries := make([]types.RuleResp, 0, len(rules))
	for _, r := range rules {
		// 查询关联文件
		relations, _ := l.svcCtx.RelationRuleFileModel.FindByRuleId(l.ctx, r.Id)
		stdFileIds := make([]int64, len(relations))
		for i, rf := range relations {
			stdFileIds[i] = rf.FileId
		}

		// MOCK: mock.CatalogGetCatalogName() - 获取目录名称
		catalogName := mock.CatalogGetCatalogName(l.ctx, l.svcCtx, r.CatalogId)

		// MOCK: mock.DataElementRuleUsed() - 检查规则是否被引用
		usedFlag := mock.DataElementRuleUsed(l.ctx, l.svcCtx, r.Id)

		item := types.RuleResp{
			Id:            r.Id,
			Name:          r.Name,
			CatalogId:     r.CatalogId,
			CatalogName:   catalogName,
			OrgType:       r.OrgType,
			Description:   r.Description,
			RuleType:      intToRuleType(r.RuleType),
			Version:       r.Version,
			Expression:    r.Expression,
			State:         r.State,
			DisableReason: r.DisableReason,
			AuthorityId:   r.AuthorityId,
			DepartmentIds: r.DepartmentIds,
			ThirdDeptId:   r.ThirdDeptId,
			CreateTime:    timeToStr(r.CreateTime),
			CreateUser:    r.CreateUser,
			UpdateTime:    timeToStr(r.UpdateTime),
			UpdateUser:    r.UpdateUser,
			StdFileIds:    stdFileIds,
			Used:          usedFlag,
		}
		entries = append(entries, item)
	}

	return &types.RuleListResp{
		Entries:    entries,
		TotalCount: int64(len(entries)),
	}, nil
}
