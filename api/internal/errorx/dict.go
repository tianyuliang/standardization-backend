package errorx

import (
	"github.com/jinguoxing/idrm-go-base/errorx"
)

// ========== dict-api 错误码定义 (30400-30499) ==========

const (
	// 通用错误 (30400-30409)
	ErrCodeDictDataNotExist = 30401 // 数据不存在
	ErrCodeDictParamEmpty   = 30402 // 参数为空
	ErrCodeDictInvalidParam = 30403 // 参数无效

	// 目录相关 (30404)
	ErrCodeDictCatalogNotExist = 30404 // 目录不存在

	// 名称相关 (30405-30406)
	ErrCodeDictChNameDuplicate = 30405 // 中文名称重复
	ErrCodeDictEnNameDuplicate = 30406 // 英文名称重复

	// 码值相关 (30407-30409)
	ErrCodeDictEnumCodeEmpty     = 30407 // 码值为空
	ErrCodeDictEnumValueEmpty    = 30408 // 码值描述为空
	ErrCodeDictEnumCodeDuplicate = 30409 // 码值重复

	// 其他 (30410)
	ErrCodeDictReasonTooLong = 30410 // 停用原因过长
)

// ========== 辅助函数 ==========

// DictDataNotExist 数据不存在
func DictDataNotExist() error {
	return errorx.NewWithMsg(ErrCodeDictDataNotExist, "数据不存在")
}

// DictParamEmpty 参数为空
func DictParamEmpty(msg string) error {
	if msg == "" {
		return errorx.NewWithMsg(ErrCodeDictParamEmpty, "参数为空")
	}
	return errorx.NewWithMsg(ErrCodeDictParamEmpty, msg)
}

// DictInvalidParam 参数无效
func DictInvalidParam(msg string) error {
	if msg == "" {
		return errorx.NewWithMsg(ErrCodeDictInvalidParam, "参数无效")
	}
	return errorx.NewWithMsg(ErrCodeDictInvalidParam, msg)
}

// DictCatalogNotExist 目录不存在
func DictCatalogNotExist(catalogId int64) error {
	return errorx.NewWithMsg(ErrCodeDictCatalogNotExist, "目录id对应的目录不存在")
}

// DictChNameDuplicate 中文名称重复
func DictChNameDuplicate() error {
	return errorx.NewWithMsg(ErrCodeDictChNameDuplicate, "码表中文名称、标准分类不能全部重复")
}

// DictEnNameDuplicate 英文名称重复
func DictEnNameDuplicate() error {
	return errorx.NewWithMsg(ErrCodeDictEnNameDuplicate, "码表英文名称、标准分类不能全部重复")
}

// DictEnumCodeEmpty 码值为空
func DictEnumCodeEmpty() error {
	return errorx.NewWithMsg(ErrCodeDictEnumCodeEmpty, "码值输入不能为空")
}

// DictEnumValueEmpty 码值描述为空
func DictEnumValueEmpty() error {
	return errorx.NewWithMsg(ErrCodeDictEnumValueEmpty, "码值描述不能为空")
}

// DictEnumCodeDuplicate 码值重复
func DictEnumCodeDuplicate() error {
	return errorx.NewWithMsg(ErrCodeDictEnumCodeDuplicate, "码值出现重复记录")
}

// DictReasonTooLong 停用原因过长
func DictReasonTooLong() error {
	return errorx.NewWithMsg(ErrCodeDictReasonTooLong, "停用原因长度不能超过800字符")
}
