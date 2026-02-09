// Code scaffolded by speckit. Safe to edit.

package dataelement

// StateToInt 状态字符串转整数
func StateToInt(state string) int32 {
	switch state {
	case "enable":
		return DataElementStateEnabled
	case "disable":
		return DataElementStateDisabled
	default:
		return DataElementStateEnabled
	}
}

// IntToState 整数转状态字符串
func IntToState(state int32) string {
	switch state {
	case DataElementStateEnabled:
		return "enable"
	case DataElementStateDisabled:
		return "disable"
	default:
		return "enable"
	}
}

// IsValidRelationType 检查关联类型是否有效
func IsValidRelationType(relationType string) bool {
	switch relationType {
	case RelationTypeNo, RelationTypeCodeTable, RelationTypeCodeRule:
		return true
	default:
		return false
	}
}

// BuildOffset 计算分页偏移量
func BuildOffset(page, pageSize int) int {
	if page < 1 {
		page = 1
	}
	if pageSize < 0 {
		pageSize = 10
	}
	return (page - 1) * pageSize
}

// BuildLimit 构建LIMIT子句值
func BuildLimit(pageSize int) int {
	if pageSize <= 0 {
		return 10
	}
	if pageSize > 2000 {
		return 2000
	}
	return pageSize
}
