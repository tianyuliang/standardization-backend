package errorx

import "github.com/jinguoxing/idrm-go-base/errorx"

const (
	// 目录模块错误码 30100-30199
	ErrCatalogNotExist     = 30100 // 目录不存在
	ErrCatalogNameInvalid  = 30101 // 目录名称格式错误
	ErrCatalogLevelInvalid = 30102 // 目录级别超出范围
	ErrCatalogHasChildren  = 30103 //目录下存在子目录
)
