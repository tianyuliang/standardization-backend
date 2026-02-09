// Code scaffolded by speckit. Safe to edit.

package relation

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
