package rule

// ============================================
// 枚举常量定义
// ============================================

// OrgTypeEnum 标准组织类型枚举
const (
	OrgTypeGroup       int32 = 0  // 团体标准
	OrgTypeEnterprise  int32 = 1  // 企业标准
	OrgTypeIndustry    int32 = 2  // 行业标准
	OrgTypeLocal       int32 = 3  // 地方标准
	OrgTypeNational    int32 = 4  // 国家标准
	OrgTypeInternational int32 = 5 // 国际标准
	OrgTypeForeign     int32 = 6  // 国外标准
	OrgTypeOther       int32 = 99 // 其他标准
)

// RuleTypeEnum 规则类型枚举
const (
	RuleTypeRegex  int32 = 0 // 正则表达式
	RuleTypeCustom int32 = 1 // 自定义配置
)

// EnableDisableStatusEnum 启用停用状态枚举
const (
	StateDisable int32 = 0 // 停用
	StateEnable  int32 = 1 // 启用
)

// CustomTypeEnum 自定义配置类型枚举
const (
	CustomTypeDict             int = 1 // 码表
	CustomTypeNumber           int = 2 // 数字
	CustomTypeEnglishLetters   int = 3 // 英文字母
	CustomTypeChineseCharacters int = 4 // 中文字符
	CustomTypeAnyCharacters    int = 5 // 任意字符
	CustomTypeDate             int = 6 // 日期
	CustomTypeSplitStr        int = 7 // 分隔字符串
)

// CustomDateFormat 自定义日期格式列表
var CustomDateFormat = []string{
	"yyyy-MM-dd",
	"yyyy/MM/dd",
	"yyyyMM",
	"yyyy/MM",
	"yyyy-MM",
	"yyyyMMdd",
	"yyyyMM",
	"yyyy/MM/dd HH:mm:ss",
	"yyyy-MM-dd HH:mm:ss",
	"yyyy/MM/dd HH:mm:ss",
}

// ============================================
// 辅助函数
// ============================================

// GetRuleTypeString 将 int32 转换为字符串
func GetRuleTypeString(ruleType int32) string {
	switch ruleType {
	case RuleTypeRegex:
		return "REGEX"
	case RuleTypeCustom:
		return "CUSTOM"
	default:
		return "REGEX"
	}
}

// GetRuleTypeInt 将字符串转换为 int32
func GetRuleTypeInt(ruleType string) int32 {
	switch ruleType {
	case "REGEX":
		return RuleTypeRegex
	case "CUSTOM":
		return RuleTypeCustom
	default:
		return RuleTypeRegex
	}
}

// GetStateString 将 int32 状态转换为字符串
func GetStateString(state int32) string {
	switch state {
	case StateDisable:
		return "DISABLE"
	case StateEnable:
		return "ENABLE"
	default:
		return "ENABLE"
	}
}

// GetStateInt 将字符串状态转换为 int32
func GetStateInt(state string) int32 {
	switch state {
	case "DISABLE":
		return StateDisable
	case "ENABLE":
		return StateEnable
	default:
		return StateEnable
	}
}

// ParseCustomTypeString 将自定义类型字符串转换为 int
func ParseCustomTypeString(customType string) int {
	switch customType {
	case "DICT":
		return CustomTypeDict
	case "NUMBER":
		return CustomTypeNumber
	case "ENGLISH_LETTERS":
		return CustomTypeEnglishLetters
	case "CHINESE_CHARACTERS":
		return CustomTypeChineseCharacters
	case "ANY_CHARACTERS":
		return CustomTypeAnyCharacters
	case "DATE":
		return CustomTypeDate
	case "SPLIT_STR":
		return CustomTypeSplitStr
	default:
		return CustomTypeDict
	}
}
