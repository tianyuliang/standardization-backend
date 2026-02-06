// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rule

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/rule/mock"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	rulemodel "github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/rule/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 编码规则列表查询
//
// 对应 Java: RuleServiceImpl.queryByCatalog() (lines 140-160)
// 业务流程:
//  1. 处理目录ID（获取当前目录及所有子目录ID列表）
//  2. 构建查询条件
//  3. 分页查询
//  4. 数据处理（解析expression、查询目录/部门/引用状态）
//
// 特殊说明：本接口无异常抛出
func NewListRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRuleLogic {
	return &ListRuleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListRuleLogic) ListRule(req *types.RuleListQuery) (resp *types.RuleListResp, err error) {
	// ====== 步骤1: 处理目录ID（获取当前目录及所有子目录ID列表） ======
	// 对应 Java: iDeCatalogInfoService.getIDList(catalogId) (line 147)
	// MOCK: mock.CatalogGetChildIds() - 获取子目录列表
	catalogIds := mock.CatalogGetChildIds(l.ctx, l.svcCtx, req.CatalogId)

	// ====== 步骤2: 构建查询条件 ======
	opts := &rulemodel.FindOptions{
		OrgType:      optionalInt32(req.OrgType),
		State:        optionalInt32(req.State),
		RuleType:     optionalRuleType(req.RuleType),
		Keyword:      req.Keyword,
		DepartmentId: "", // TODO: 添加到请求参数
		Page:         req.Offset,
		PageSize:     req.Limit,
		Sort:         "", // TODO: 添加到请求参数
		Direction:    "", // TODO: 添加到请求参数
	}

	// ====== 步骤3: 分页查询 ======
	// 对应 Java: ruleMapper.queryByCatalog(page, catalogIds, keyword, orgType, state, departmentId, ruleType) (line 152)
	rules, totalCount, err := l.svcCtx.RuleModel.FindByCatalogIds(l.ctx, catalogIds, opts)
	if err != nil {
		return nil, err
	}

	// ====== 步骤4: 数据处理 ======
	// 对应 Java: dbDataToVo(pageResult) (line 155)
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
			Expression:    r.Expression, // TODO: 根据 ruleType 解析
			State:         r.State,
			DisableReason: r.DisableReason,
			AuthorityId:   r.AuthorityId,
			DepartmentIds: r.DepartmentIds,
			ThirdDeptId:   r.ThirdDeptId,
			StdFileIds:    stdFileIds,
			Used:          usedFlag,
			CreateTime:    timeToStr(r.CreateTime),
			CreateUser:    r.CreateUser,
			UpdateTime:    timeToStr(r.UpdateTime),
			UpdateUser:    r.UpdateUser,
		}
		entries = append(entries, item)
	}

	return &types.RuleListResp{
		TotalCount: totalCount,
		Entries:    entries,
	}, nil
}

// ====== 辅助函数 ======

// optionalInt32 将 int32 转换为指针，-1 表示不过滤
func optionalInt32(v int32) *int32 {
	if v == -1 {
		return nil
	}
	return &v
}

// optionalRuleType 将字符串转换为指针，空表示不过滤
func optionalRuleType(ruleType string) *int32 {
	if ruleType == "" {
		return nil
	}
	v := rulemodel.GetRuleTypeInt(ruleType)
	return &v
}
