package dict

// 状态常量
const (
	StateEnable  int32 = 1 // 启用
	StateDisable int32 = 0 // 停用
)

// 状态字符串常量
const (
	StateStringEnable  = "enable"
	StateStringDisable = "disable"
)

// 表名
const (
	TableNameDict             = "t_dict"
	TableNameDictEnum         = "t_dict_enum"
	TableNameRelationDictFile = "t_relation_dict_file"
)

// 标准分类枚举值 (OrgType)
const (
	OrgTypeGroup         int32 = 0  // 团体标准
	OrgTypeEnterprise    int32 = 1  // 企业标准
	OrgTypeIndustry      int32 = 2  // 行业标准
	OrgTypeLocal         int32 = 3  // 地方标准
	OrgTypeNational      int32 = 4  // 国家标准
	OrgTypeInternational int32 = 5  // 国际标准
	OrgTypeForeign       int32 = 6  // 国外标准
	OrgTypeOther         int32 = 99 // 其他标准
)

// 字段长度限制
const (
	MaxChNameLength        = 128
	MaxEnNameLength        = 128
	MaxDescriptionLength   = 300
	MaxDisableReasonLength = 800
	MaxDepartmentIdsLength = 350
	MaxEnumCodeLength      = 50
	MaxEnumValueLength     = 128
)

// StateToInt 状态字符串转整数
func StateToInt(state string) int32 {
	if state == StateStringEnable {
		return StateEnable
	}
	return StateDisable
}

// IntToState 状态整数转字符串
func IntToState(state int32) string {
	if state == StateEnable {
		return StateStringEnable
	}
	return StateStringDisable
}
