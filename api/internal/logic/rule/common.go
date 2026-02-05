// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rule

import (
	"time"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	rulemodel "github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/rule/rule"
)

// ============================================
// 共享辅助函数
// ============================================

// intToRuleType 将 int32 转换为字符串
func intToRuleType(ruleType int32) string {
	return rulemodel.GetRuleTypeString(ruleType)
}

// ruleTypeToInt 将字符串转换为 int32
func ruleTypeToInt(ruleType string) int32 {
	return rulemodel.GetRuleTypeInt(ruleType)
}

// timeToStr 将 time.Time 转换为字符串
func timeToStr(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// buildRuleResp 构建响应对象
// TODO: 查询目录名称、部门信息、引用状态
func buildRuleResp(rule *rulemodel.Rule, catalogName string, usedFlag bool, stdFiles []int64) *types.RuleResp {
	if stdFiles == nil {
		stdFiles = []int64{}
	}
	return &types.RuleResp{
		Id:            rule.Id,
		Name:          rule.Name,
		CatalogId:     rule.CatalogId,
		CatalogName:   catalogName,
		OrgType:       rule.OrgType,
		Description:   rule.Description,
		RuleType:      rulemodel.GetRuleTypeString(rule.RuleType),
		Version:       rule.Version,
		Expression:    rule.Expression,
		State:         rule.State,
		DisableReason: rule.DisableReason,
		AuthorityId:   rule.AuthorityId,
		DepartmentIds: rule.DepartmentIds,
		ThirdDeptId:   rule.ThirdDeptId,
		StdFileIds:    stdFiles,
		Used:          usedFlag,
		CreateTime:    timeToStr(rule.CreateTime),
		CreateUser:    rule.CreateUser,
		UpdateTime:    timeToStr(rule.UpdateTime),
		UpdateUser:    rule.UpdateUser,
	}
}

// getExpression 获取表达式
// TODO: REGEX类型返回regex，CUSTOM类型返回custom的JSON序列化
func getExpression(ruleType string, regex string, custom []types.RuleCustom) string {
	if ruleType == "REGEX" {
		return regex
	}
	// TODO: json.Marshal(custom)
	return "{}"
}
