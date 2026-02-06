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

// ========== catalog-api 错误码定义 (30100-30199) ==========

const (
	// 通用错误 (30100-30109)
	ErrCodeCatalogError        = 30100 // CatalogServiceError - 目录服务错误
	ErrCodeCatalogEmpty        = 30101 // Empty - 数据不存在
	ErrCodeCatalogMissingParam = 30102 // MissingParameter - 参数缺失
	ErrCodeCatalogInvalidParam = 30103 // InvalidParameter - 参数无效
	ErrCodeCatalogOutOfRange   = 30104 // OutOfRange - 超出范围
	ErrCodeCatalogConflict     = 30105 // OperationConflict - 操作冲突
	ErrCodeCatalogDataExist    = 30106 // DATA_EXIST - 数据已存在
)

// ========== catalog-api 辅助函数 ==========

// CatalogNameEmpty 目录名称为空
func CatalogNameEmpty() error {
	return errorx.NewWithMsg(ErrCodeCatalogInvalidParam, "目录名称不能为空")
}

// CatalogNameTooLong 目录名称过长
func CatalogNameTooLong() error {
	return errorx.NewWithMsg(ErrCodeCatalogInvalidParam, "目录名称长度不能超过20个字符")
}

// CatalogNameInvalidFormat 目录名称格式错误
func CatalogNameInvalidFormat() error {
	return errorx.NewWithMsg(ErrCodeCatalogInvalidParam, "目录名称只能由中英文、数字、下划线、中划线组成")
}

// CatalogNameInvalidPrefix 目录名称前缀无效
func CatalogNameInvalidPrefix() error {
	return errorx.NewWithMsg(ErrCodeCatalogInvalidParam, "目录名称不能以下划线和中划线开头")
}

// CatalogParentNotExist 父目录不存在
func CatalogParentNotExist() error {
	return errorx.NewWithMsg(ErrCodeCatalogEmpty, "无法找到对应的父目录")
}

// CatalogLevelOutOfRange 目录级别超出范围
func CatalogLevelOutOfRange() error {
	return errorx.NewWithMsg(ErrCodeCatalogOutOfRange, "目录级别取值范围(1-255)")
}

// CatalogNameDuplicate 同级目录名称重复
func CatalogNameDuplicate() error {
	return errorx.NewWithMsg(ErrCodeCatalogConflict, "同级目录名称不能重复")
}

// CatalogCannotDeleteRoot 不允许删除根目录
func CatalogCannotDeleteRoot() error {
	return errorx.NewWithMsg(ErrCodeCatalogInvalidParam, "不允许删除根目录")
}

// CatalogCannotModifyRoot 不允许修改根目录
func CatalogCannotModifyRoot() error {
	return errorx.NewWithMsg(ErrCodeCatalogInvalidParam, "不能修改根目录")
}

// CatalogHasData 目录下存在数据，不允许删除
func CatalogHasData() error {
	return errorx.NewWithMsg(ErrCodeCatalogDataExist, "目录或子目录下已存在数据，不允许删除")
}

// CatalogTypeEmpty 目录类型为空
func CatalogTypeEmpty() error {
	return errorx.NewWithMsg(ErrCodeCatalogMissingParam, "类型不能为空")
}

// CatalogTypeInvalid 目录类型无效
func CatalogTypeInvalid() error {
	return errorx.NewWithMsg(ErrCodeCatalogInvalidParam, "此类型不在有效值范围内")
}

// CatalogNewParentIsChild 新父目录不能是自身及其子目录
func CatalogNewParentIsChild() error {
	return errorx.NewWithMsg(ErrCodeCatalogInvalidParam, "新的父目录不能是自身及其子目录")
}

// CatalogTypeMismatch 目录类型不匹配
func CatalogTypeMismatch() error {
	return errorx.NewWithMsg(ErrCodeCatalogInvalidParam, "新的父目录类型不能与当前目录不一致")
}

// CatalogNotExist 目录不存在
func CatalogNotExist() error {
	return errorx.NewWithMsg(ErrCodeCatalogEmpty, "目录id对应的目录不存在")
}
