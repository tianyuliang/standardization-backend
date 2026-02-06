package catalog

// CatalogTypeEnum 目录类型枚举
const (
	CatalogTypeDataElement int32 = 1 // 数据元目录
	CatalogTypeDict        int32 = 2 // 码表目录
	CatalogTypeValueRule   int32 = 3 // 编码规则目录
	CatalogTypeFile        int32 = 4 // 文件目录
)

// CatalogTypeString 目录类型字符串映射
var CatalogTypeString = map[int32]string{
	CatalogTypeDataElement: "DataElement",
	CatalogTypeDict:        "Dict",
	CatalogTypeValueRule:   "ValueRule",
	CatalogTypeFile:        "File",
}

// CatalogTypeInt 目录类型int映射
var CatalogTypeInt = map[string]int32{
	"DataElement": CatalogTypeDataElement,
	"Dict":        CatalogTypeDict,
	"ValueRule":   CatalogTypeValueRule,
	"File":        CatalogTypeFile,
}

// GetCatalogTypeString 获取目录类型字符串
func GetCatalogTypeString(catalogType int32) string {
	if s, ok := CatalogTypeString[catalogType]; ok {
		return s
	}
	return "Other"
}

// GetCatalogTypeInt 获取目录类型int
func GetCatalogTypeInt(catalogType string) int32 {
	if i, ok := CatalogTypeInt[catalogType]; ok {
		return i
	}
	return CatalogTypeDataElement
}

// IsValidCatalogType 校验目录类型是否有效
func IsValidCatalogType(catalogType int32) bool {
	_, ok := CatalogTypeString[catalogType]
	return ok
}
