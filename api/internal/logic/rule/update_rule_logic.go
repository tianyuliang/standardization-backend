// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rule

import (
	"context"
	"time"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/rule/relation_file"
	rulemodel "github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/rule/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据ID修改编码规则
//
// 业务流程（参考 specs/编码规则管理接口流程说明_20260204.md 第4.2节）:
//   1. 校验规则是否存在
//   2. 表达式校验（与创建规则相同）
//   3. 名称唯一性校验（排除自身）
//   4. 目录存在性校验
//   5. 版本变更检测
//      - 需要递增版本号的字段：name, catalog_id, department_ids, org_type, description, rule_type, expression, 关联文件
//   6. 无变更直接返回原数据
//   7. 有变更则更新数据（版本号+1、更新关联文件）
//   8. 发送MQ消息
//
// 异常处理：
//   - 30301: 记录不存在
//   - 30311: 修改时名称与其他记录重复
//   - 30312: 目录不存在
func NewUpdateRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRuleLogic {
	return &UpdateRuleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRuleLogic) UpdateRule(id int64, req *types.UpdateRuleReq) (resp *types.RuleResp, err error) {
	// ====== 步骤1: 校验规则是否存在 ======
	oldRule, err := l.svcCtx.RuleModel.FindOne(l.ctx, id)
	if err != nil || oldRule == nil {
		return nil, errorx.RuleRecordNotExist() // [错误码 30301]
	}

	// ====== 步骤2: 表达式校验 ======
	// TODO: 调用 ValidateExpression(ruleType, regex, custom)

	// ====== 步骤3: 名称唯一性校验（排除自身）=====
	// TODO: 查询同名规则，排除当前id
	// - 如果存在其他同名记录，返回 errorx.RuleNameDuplicate(name) [错误码 30311]

	// ====== 步骤4: 目录存在性校验 ======
	// TODO: 调用 Catalog RPC 校验 catalogId 是否存在
	// - 如果不存在，返回 errorx.RuleCatalogNotExist(req.CatalogId) [错误码 30312]

	// ====== 步骤5: 版本变更检测 ======
	// TODO: 查询原有关联文件
	// oldFiles, _ := l.svcCtx.RelationRuleFileModel.FindByRuleId(l.ctx, id)

	// 检测关键字段是否发生变化
	needVersionIncrement := checkVersionChange(oldRule, req, nil)

	if !needVersionIncrement {
		// ====== 步骤6: 无变更直接返回原数据 ======
		return buildRuleResp(oldRule, "", false, nil), nil
	}

	// ====== 步骤7: 有变更则更新数据 ======
	newRule := &rulemodel.Rule{
		Id:            oldRule.Id,
		Name:          req.Name,
		CatalogId:     req.CatalogId,
		OrgType:       req.OrgType,
		Description:   req.Description,
		RuleType:      ruleTypeToInt(req.RuleType),
		Version:       oldRule.Version + 1, // 版本号+1
		Expression:    getExpression(req.RuleType, req.Regex, req.Custom),
		State:         oldRule.State, // 保持原有状态
		DepartmentIds: req.DepartmentIds,
		UpdateTime:    time.Now(),
		UpdateUser:    "", // TODO: 从 Token 获取
		// 继承原有字段
		CreateTime:    oldRule.CreateTime,
		CreateUser:    oldRule.CreateUser,
		Deleted:       oldRule.Deleted,
	}

	// TODO: 开启事务
	// 7.1 更新 t_rule
	err = l.svcCtx.RuleModel.Update(l.ctx, newRule)
	if err != nil {
		return nil, err
	}

	// 7.2 更新关联文件 t_relation_rule_file
	// TODO: 删除旧关联，插入新关联
	// - l.svcCtx.RelationRuleFileModel.DeleteByRuleId(l.ctx, id)
	// - l.svcCtx.RelationRuleFileModel.InsertBatch(l.ctx, newFiles)

	// ====== 步骤8: 发送MQ消息 ======
	// TODO: 调用 SendRuleMQMessage(producer, []newRule, "update")

	// ====== 步骤9: 构建响应 ======
	return buildRuleResp(newRule, "", false, nil), nil
}

// ====== 辅助函数 ======

// checkVersionChange 检测是否需要递增版本号
// 需要递增版本号的字段：name, catalog_id, department_ids, org_type, description, rule_type, expression, 关联文件
func checkVersionChange(old *rulemodel.Rule, req *types.UpdateRuleReq, oldFiles []*relation_file.RelationRuleFile) bool {
	// TODO: 比较以下字段是否变化
	// - name, catalog_id, department_ids, org_type, description
	// - rule_type, expression (regex/custom)
	// - 关联文件
	return true
}
