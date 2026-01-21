package catalog

import "regexp"

const (
	// CatalogNamePattern 目录名称正则表达式
	// 规则：1-20字符，中文/英文/数字/符号_-, 且_-不能作为首字符
	// 对应 Java: Constants.getRegexENOrCNVarL(1, 20)
	// Java 正则: ^(?!_)(?!-)[\u4E00-\u9FA5\uF900-\uFA2D\w-]{1,20}$
	CatalogNamePattern = `^(?!_)(?!-)[\x{4e00}-\x{9fa5}\x{f900}-\x{fa2d}\w-]{1,20}$`

	// CatalogNameMaxLength 目录名称最大长度
	CatalogNameMaxLength = 20

	// CatalogNameMinLength 目录名称最小长度
	CatalogNameMinLength = 1

	// CatalogMaxLevel 目录最大层级
	CatalogMaxLevel = 255

	// CatalogRootLevel 根目录层级
	CatalogRootLevel = 1
)

// CatalogType 目录类型枚举
const (
	CatalogTypeDataElement  = 1  // 数据元
	CatalogTypeDict         = 2  // 码表
	CatalogTypeEncodingRule = 3  // 编码规则
	CatalogTypeFile         = 4  // 文件
	CatalogTypeRoot         = 0  // 根目录(保留)
	CatalogTypeOther        = 99 // 其他(保留)
)

// CatalogTypeText 目录类型文本映射
var CatalogTypeText = map[int32]string{
	CatalogTypeDataElement:  "数据元",
	CatalogTypeDict:         "码表",
	CatalogTypeEncodingRule: "编码规则",
	CatalogTypeFile:         "文件",
	CatalogTypeRoot:         "根目录",
	CatalogTypeOther:        "其他",
}

// GetCatalogTypeText 获取目录类型文本
func GetCatalogTypeText(catalogType int32) string {
	if text, ok := CatalogTypeText[catalogType]; ok {
		return text
	}
	return "未知"
}

// CatalogNameRegexp 编译后的目录名称正则表达式（包级别变量，避免重复编译）
var CatalogNameRegexp = regexp.MustCompile(CatalogNamePattern)
