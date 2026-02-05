// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rule

import (
	"context"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 编码规则详情查看
//
// 业务流程（参考 specs/编码规则管理接口流程说明_20260204.md 第4.3节）:
//   1. 查询规则基本信息
//   2. 查询目录信息（目录名称）
//   3. 根据ruleType解析expression（返回regex或custom）
//   4. 查询关联文件列表
//   5. 查询部门信息（部门名称、路径名称）
//   6. 查询是否被引用（usedFlag）
//
// 特殊说明：本接口无异常抛出，记录不存在时返回null（code=0）
func NewGetRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRuleLogic {
	return &GetRuleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRuleLogic) GetRule(id int64) (resp *types.RuleResp, err error) {
	// ====== 步骤1: 查询规则基本信息 ======
	ruleData, err := l.svcCtx.RuleModel.FindOne(l.ctx, id)
	if err != nil || ruleData == nil {
		// 记录不存在时返回null，不抛异常
		return nil, nil
	}

	// ====== 步骤2: 查询目录信息 ======
	// TODO: 调用 Catalog RPC 获取目录名称
	// catalogName := getMockCatalogName(ruleData.CatalogId)
	catalogName := ""

	// ====== 步骤3: 根据ruleType解析expression ======
	// TODO: 解析 expression
	// - RuleType = 0 (REGEX): 返回 expression 作为 expression
	// - RuleType = 1 (CUSTOM): 返回 expression 反序列化后的 custom
	expression := ruleData.Expression

	// ====== 步骤4: 查询关联文件列表 ======
	// TODO: 查询 t_relation_rule_file
	// files, _ := l.svcCtx.RelationRuleFileModel.FindByRuleId(l.ctx, id)
	// stdFileIds := extractFileIds(files)
	var stdFileIds []int64

	// ====== 步骤5: 查询部门信息 ======
	// TODO: 调用部门服务获取部门名称、路径名称
	// departmentName := getMockDeptName(ruleData.DepartmentIds)
	// departmentPathNames := getMockDeptPathNames(ruleData.DepartmentIds)

	// ====== 步骤6: 查询是否被引用 ======
	// TODO: 调用 DataElement RPC 检查是否被引用
	// usedFlag := getMockUsedFlag(id)
	used := false

	// ====== 步骤7: 构建响应 ======
	return &types.RuleResp{
		Id:            ruleData.Id,
		Name:          ruleData.Name,
		CatalogId:     ruleData.CatalogId,
		CatalogName:   catalogName,
		OrgType:       ruleData.OrgType,
		Description:   ruleData.Description,
		RuleType:      intToRuleType(ruleData.RuleType),
		Version:       ruleData.Version,
		Expression:    expression,
		State:         ruleData.State,
		DisableReason: ruleData.DisableReason,
		AuthorityId:   ruleData.AuthorityId,
		DepartmentIds: ruleData.DepartmentIds,
		ThirdDeptId:   ruleData.ThirdDeptId,
		StdFileIds:    stdFileIds,
		Used:          used,
		CreateTime:    timeToStr(ruleData.CreateTime),
		CreateUser:    ruleData.CreateUser,
		UpdateTime:    timeToStr(ruleData.UpdateTime),
		UpdateUser:    ruleData.UpdateUser,
	}, nil
}
