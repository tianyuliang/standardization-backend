package errorx

import "github.com/jinguoxing/idrm-go-base/errorx"

// ========== rule-api 错误码定义 (30300-30399) ==========

const (
	// 通用错误 (30300-30309)
	ErrCodeRuleDataNotExist = 30301 // DATA_NOT_EXIST - 数据不存在
	ErrCodeRuleParamEmpty   = 30302 // PARAMETER_EMPTY - 参数为空
	ErrCodeRuleInvalidParam = 30303 // InvalidParameter - 参数无效

	// 名称相关 (30310-30319)
	ErrCodeRuleNameDuplicate       = 30310 // 创建时名称重复
	ErrCodeRuleNameUpdateDuplicate = 30311 // 修改时名称重复
	ErrCodeRuleCatalogNotExist     = 30312 // 目录不存在

	// REGEX 相关 (30320-30329)
	ErrCodeRuleRegexEmpty   = 30320 // 正则表达式为空
	ErrCodeRuleRegexInvalid = 30321 // 正则表达式非法

	// CUSTOM 相关 (30330-30339)
	ErrCodeRuleCustomEmpty      = 30330 // 自定义配置为空
	ErrCodeRuleSegmentLength    = 30331 // segment_length <= 0
	ErrCodeRuleDictNotExist     = 30332 // 码表不存在
	ErrCodeRuleDateFormat       = 30333 // 日期格式不支持
	ErrCodeRuleCustomValueEmpty = 30334 // value为空

	// 参数相关 (30335-30339)
	ErrCodeRuleIdsEmpty      = 30335 // 批量删除ids为空
	ErrCodeRuleReasonTooLong = 30336 // 停用原因过长
	ErrCodeRuleQueryIdsEmpty = 30337 // 批量查询ids为空
)

// ========== 辅助函数 ==========

// RuleNameDuplicate 规则名称已存在
func RuleNameDuplicate(name string) error {
	return errorx.NewWithMsg(ErrCodeRuleNameDuplicate, "规则名称已存在")
}

// RuleCatalogNotExist 目录不存在
func RuleCatalogNotExist(catalogId int64) error {
	return errorx.NewWithMsg(ErrCodeRuleCatalogNotExist, "目录id对应的目录不存在")
}

// RuleRegexEmpty 正则表达式为空
func RuleRegexEmpty() error {
	return errorx.NewWithMsg(ErrCodeRuleRegexEmpty, "正则表达式为空")
}

// RuleRegexInvalid 正则表达式非法
func RuleRegexInvalid() error {
	return errorx.NewWithMsg(ErrCodeRuleRegexInvalid, "正则表达式非法")
}

// RuleCustomEmpty 自定义配置为空
func RuleCustomEmpty() error {
	return errorx.NewWithMsg(ErrCodeRuleCustomEmpty, "不能为空")
}

// RuleSegmentLengthInvalid segment_length必须为正整数
func RuleSegmentLengthInvalid(fieldPrefix string) error {
	return errorx.NewWithMsg(ErrCodeRuleSegmentLength, fieldPrefix+"segment_length: 值必须为正整数")
}

// RuleDictNotExist 码表不存在
func RuleDictNotExist(fieldPrefix string) error {
	return errorx.NewWithMsg(ErrCodeRuleDictNotExist, fieldPrefix+"value: 码表不存在")
}

// RuleDateFormatNotSupported 日期格式不支持
func RuleDateFormatNotSupported(fieldPrefix string) error {
	return errorx.NewWithMsg(ErrCodeRuleDateFormat, fieldPrefix+"value: 不支持的日期格式")
}

// RuleCustomValueEmpty value不能为空
func RuleCustomValueEmpty(fieldPrefix string) error {
	return errorx.NewWithMsg(ErrCodeRuleCustomValueEmpty, fieldPrefix+"value: 不能为空")
}

// RuleDisableReasonEmpty 停用原因为空
func RuleDisableReasonEmpty() error {
	return errorx.NewWithMsg(ErrCodeRuleParamEmpty, "停用必须填写停用原因")
}

// RuleDisableReasonTooLong 停用原因过长
func RuleDisableReasonTooLong() error {
	return errorx.NewWithMsg(ErrCodeRuleReasonTooLong, "长度超过800")
}

// RuleIdsEmpty ids不能为空
func RuleIdsEmpty() error {
	return errorx.NewWithMsg(ErrCodeRuleIdsEmpty, "ids 不能为空")
}

// RuleQueryIdsEmpty 查询ids为空
func RuleQueryIdsEmpty() error {
	return errorx.NewWithMsg(ErrCodeRuleQueryIdsEmpty, "ids 不能为空")
}

// RuleDataNotExist 数据不存在
func RuleDataNotExist() error {
	return errorx.NewWithMsg(ErrCodeRuleDataNotExist, "数据不存在")
}

// RuleRecordNotExist 记录不存在
func RuleRecordNotExist() error {
	return errorx.NewWithMsg(ErrCodeRuleDataNotExist, "记录不存在")
}
