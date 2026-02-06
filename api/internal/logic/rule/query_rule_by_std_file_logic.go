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

type QueryRuleByStdFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据标准文件查询规则
//
// 对应 Java: RuleServiceImpl.queryByStdFile() (lines 655-674)
// 业务流程:
//  1. file_id 为空: 返回空列表
//  2. 查询关联该文件的规则ID列表
//  3. 批量查询规则详情
func NewQueryRuleByStdFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryRuleByStdFileLogic {
	return &QueryRuleByStdFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryRuleByStdFileLogic) QueryRuleByStdFile(fileId int64, req *types.RuleListQuery) (resp *types.RuleListResp, err error) {
	// ====== 步骤1: file_id 为空: 返回空列表 ======
	// 对应 Java: if (CustomUtil.isEmpty(stdFileId)) (lines 664-666)
	if fileId == 0 {
		return &types.RuleListResp{
			Entries:    []types.RuleResp{},
			TotalCount: 0,
		}, nil
	}

	// ====== 步骤2: 查询关联该文件的规则ID列表 ======
	// 对应 Java: ruleMapper.queryByStdFile(page, stdFileId, ...) (line 672)
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

	// 提取规则ID列表
	ruleIds := make([]int64, 0, len(relations))
	for _, r := range relations {
		ruleIds = append(ruleIds, r.RuleId)
	}

	// ====== 步骤3: 批量查询规则详情 ======
	// 对应 Java: return dbDataToVo(pageResult) (line 673)
	rules, err := l.svcCtx.RuleModel.FindByIds(l.ctx, ruleIds)
	if err != nil {
		return nil, err
	}

	// ====== 步骤4: 构建响应 ======
	entries := make([]types.RuleResp, 0, len(rules))
	for _, r := range rules {
		// 查询关联文件
		fileRelations, _ := l.svcCtx.RelationRuleFileModel.FindByRuleId(l.ctx, r.Id)
		stdFileIds := make([]int64, len(fileRelations))
		for i, rf := range fileRelations {
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
