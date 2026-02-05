// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rule

import (
	"context"

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
// 业务流程（参考 specs/编码规则管理接口流程说明_20260204.md 第4.4节）:
//   1. 处理目录ID（获取当前目录及所有子目录ID列表）
//   2. 构建查询条件
//   3. 分页查询
//   4. 数据处理（解析expression、查询目录/部门/引用状态）
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
	// TODO: 调用 Catalog RPC 获取子目录列表
	// catalogIds := getMockCatalogIds(req.CatalogId)
	catalogIds := []int64{req.CatalogId}

	// ====== 步骤2: 构建查询条件 ======
	opts := &rulemodel.FindOptions{
		OrgType:      optionalInt32(req.OrgType),
		State:        optionalInt32(req.State),
		RuleType:     optionalRuleType(req.RuleType),
		Keyword:      req.Keyword,
		DepartmentId: "", // TODO: 从请求中获取部门ID
		Page:         req.Offset,
		PageSize:     req.Limit,
		Sort:         "f_create_time", // 默认按创建时间排序
		Direction:    "DESC",          // 默认降序
	}

	// ====== 步骤3: 分页查询 ======
	rules, totalCount, err := l.svcCtx.RuleModel.FindByCatalogIds(l.ctx, catalogIds, opts)
	if err != nil {
		return nil, err
	}

	// ====== 步骤4: 数据处理 ======
	// TODO: 批量查询目录名称、部门信息、引用状态
	entries := make([]types.RuleResp, 0, len(rules))
	for _, r := range rules {
		item := types.RuleResp{
			Id:            r.Id,
			Name:          r.Name,
			CatalogId:     r.CatalogId,
			CatalogName:   "", // TODO: 查询目录名称
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
			StdFileIds:    nil, // TODO: 查询关联文件
			Used:          false, // TODO: 查询引用状态
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
