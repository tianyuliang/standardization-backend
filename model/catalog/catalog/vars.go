package catalog

const (
	// CatalogNamePattern 目录名称正则表达式
	// 规则：1-20字符，中文/英文/数字/符号_-, 且_-不能作为首字符
	CatalogNamePattern = `^[\u4e00-\u9fa5a-zA-Z0-9][\u4e00-\u9fa5a-zA-Z0-9_-]{0,19}$`

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
	CatalogTypeDataElement    = 1 // 数据元
	CatalogTypeDict           = 2 // 码表
	CatalogTypeEncodingRule   = 3 // 编码规则
	CatalogTypeFile           = 4 // 文件
	CatalogTypeRoot           = 0 // 根目录(保留)
	CatalogTypeOther          = 99 // 其他(保留)
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
