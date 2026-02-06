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

// intToState 将 int32 状态转换为字符串
func intToState(state int32) string {
	return rulemodel.GetStateString(state)
}

// stateToInt 将字符串状态转换为 int32
func stateToInt(state string) int32 {
	return rulemodel.GetStateInt(state)
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

// ============================================
// 业务校验函数 (TODO: 后续实现)
// ============================================

// ValidateExpression 表达式校验
// TODO: REGEX类型：校验正则表达式非空且格式正确
// TODO: CUSTOM类型：校验custom非空、segment_length>0、value有效性
func ValidateExpression(ruleType string, regex string, custom []types.RuleCustom) error {
	// TODO: 实现表达式校验逻辑
	// - REGEX: 校验 regex 非空 + regexp.Compile()
	// - CUSTOM: 校验 custom 非空 + 遍历校验每个 segment
	return nil
}

// CheckVersionChange 版本变更检测
// TODO: 检测关键字段是否发生变化
// 需要递增版本号的字段：name, catalog_id, department_ids, org_type, description, rule_type, expression, 关联文件
func CheckVersionChange(old *rulemodel.Rule, req *types.UpdateRuleReq, oldFiles []*RelationRuleFile) bool {
	// TODO: 比较以下字段是否变化
	// - name, catalog_id, department_ids, org_type, description
	// - rule_type, expression (regex/custom)
	// - 关联文件
	return true
}

// SendRuleMQMessage MQ消息发送
// TODO: 发送规则变更消息到 Kafka
func SendRuleMQMessage(producer interface{}, rules interface{}, operation string) error {
	// TODO: 实现MQ消息发送
	// - MQ Topic: MQ_MESSAGE_SAILOR
	// - 消息格式: { "header": {}, "payload": { "type": "smart-recommendation-graph", "content": { "type": "insert/update/delete", "tableName": "t_rule", "entities": [...] } } }
	return nil
}

// CheckNameUnique 名称唯一性校验
// TODO: 检查同一 orgType 下规则名称是否唯一
func CheckNameUnique(model rulemodel.RuleModel, name string, orgType int32, departmentIds string) error {
	// TODO: 查询数据库：SELECT * FROM t_rule WHERE f_name=? AND f_org_type=? AND f_deleted=0 AND f_department_ids=?
	// - 如果存在记录，返回 errorx.RuleNameDuplicate(name) [错误码 30310]
	return nil
}

// CheckCatalogIdExist 目录存在性校验
// TODO: 调用 Catalog RPC 校验目录是否存在
func CheckCatalogIdExist(catalogId int64) error {
	// TODO: 调用 Catalog RPC
	// 当前返回 mock 数据
	return nil
}

// ============================================
// 类型定义
// ============================================

// RelationRuleFile 关联规则文件 (用于版本变更检测)
type RelationRuleFile struct {
	Id     int64 `json:"id"`
	RuleId int64 `json:"ruleId"`
	FileId int64 `json:"fileId"`
}
