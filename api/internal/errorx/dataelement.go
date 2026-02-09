// Package errorx provides dataelement-specific error definitions
package errorx

import (
	"fmt"

	"github.com/jinguoxing/idrm-go-base/errorx"
)

// ========== dataelement-api 错误码定义 (30400-30499) ==========

const (
	// 通用错误 (30400-30409)
	ErrCodeDataElementNotExist    = 30401 // DATA_NOT_EXIST - 数据不存在
	ErrCodeDataElementParamEmpty  = 30402 // PARAMETER_EMPTY - 参数为空
	ErrCodeDataElementInvalidParam = 30403 // INVALID_PARAMETER - 参数无效
	ErrCodeDataElementDataExist   = 30404 // DATA_EXIST - 数据已存在

	// 目录相关 (30410-30419)
	ErrCodeDataElementCatalogNotExist    = 30410 // 目录不存在
	ErrCodeDataElementCatalogTypeWrong  = 30411 // 目录类型不正确

	// 关联类型相关 (30420-30429)
	ErrCodeDataElementRelationTypeEmpty = 30420 // 关联类型为空
	ErrCodeDataElementDictIdEmpty       = 30421 // 码表ID为空
	ErrCodeDataElementRuleIdEmpty       = 30422 // 规则ID为空
	ErrCodeDataElementRuleNotExist      = 30423 // 编码规则不存在
	ErrCodeDataElementDictNotExist      = 30424 // 码表不存在

	// 名称相关 (30430-30439)
	ErrCodeDataElementNameCnDuplicate   = 30430 // 中文名称重复
	ErrCodeDataElementNameEnDuplicate   = 30431 // 英文名称重复
	ErrCodeDataElementNameEnInvalid     = 30432 // 英文名称格式无效
	ErrCodeDataElementNameCnEmpty       = 30433 // 中文名称为空
	ErrCodeDataElementNameCnTooLong     = 30434 // 中文名称过长
	ErrCodeDataElementNameEnEmpty       = 30435 // 英文名称为空

	// 数据类型相关 (30440-30449)
	ErrCodeDataElementStdTypeInvalid    = 30440 // 标准分类无效
	ErrCodeDataElementDataTypeInvalid   = 30441 // 数据类型无效
	ErrCodeDataElementDataLengthInvalid = 30442 // 数据长度无效

	// 文件相关 (30450-30459)
	ErrCodeDataElementFileIdNotExist    = 30450 // 文件ID不存在

	// 状态相关 (30460-30469)
	ErrCodeDataElementDisableReasonEmpty   = 30460 // 停用原因为空
	ErrCodeDataElementDisableReasonTooLong = 30461 // 停用原因过长
)

// ========== dataelement-api 错误函数 ==========

// DataNotExist 数据不存在
func DataNotExist() error {
	return errorx.NewWithMsg(ErrCodeDataElementNotExist, "数据不存在")
}

// ParameterEmpty 参数为空
func ParameterEmpty(field string) error {
	return errorx.NewWithMsg(ErrCodeDataElementParamEmpty, fmt.Sprintf("[%s]:空", field))
}

// InvalidParameter 参数无效
func InvalidParameter(field, message string) error {
	return errorx.NewWithMsg(ErrCodeDataElementInvalidParam, fmt.Sprintf("[%s]:%s", field, message))
}

// DataElementCatalogNotExist 目录不存在
func DataElementCatalogNotExist() error {
	return errorx.NewWithMsg(ErrCodeDataElementCatalogNotExist, "数据元对应的目录不存在")
}

// CatalogTypeIncorrect 目录类型不正确
func CatalogTypeIncorrect() error {
	return errorx.NewWithMsg(ErrCodeDataElementCatalogTypeWrong, "数据元对应的目录类型不正确")
}

// RelationTypeEmpty 关联类型为空
func RelationTypeEmpty() error {
	return errorx.NewWithMsg(ErrCodeDataElementRelationTypeEmpty, "数据元关联类型不能为空")
}

// DictIdEmpty 码表ID为空
func DictIdEmpty() error {
	return errorx.NewWithMsg(ErrCodeDataElementDictIdEmpty, "数据元关联码表不能为空")
}

// RuleIdEmpty 规则ID为空
func RuleIdEmpty() error {
	return errorx.NewWithMsg(ErrCodeDataElementRuleIdEmpty, "数据元关联编码规则不能为空")
}

// RuleNotExist 编码规则不存在
func RuleNotExist() error {
	return errorx.NewWithMsg(ErrCodeDataElementRuleNotExist, "编码规则不存在或已删除")
}

// DictNotExist 码表不存在
func DictNotExist() error {
	return errorx.NewWithMsg(ErrCodeDataElementDictNotExist, "码表数据不存在或已删除")
}

// NameCnDuplicate 中文名称重复
func NameCnDuplicate() error {
	return errorx.NewWithMsg(ErrCodeDataElementNameCnDuplicate, "中文名称+标准分类不能全部重复")
}

// NameEnDuplicate 英文名称重复
func NameEnDuplicate() error {
	return errorx.NewWithMsg(ErrCodeDataElementNameEnDuplicate, "[英文名称]:同部门下不能重复")
}

// NameEnInvalid 英文名称格式无效
func NameEnInvalid() error {
	return errorx.NewWithMsg(ErrCodeDataElementNameEnInvalid, "[英文名称]:字符不符合要求或长度超过128")
}

// NameCnEmpty 中文名称为空
func NameCnEmpty() error {
	return errorx.NewWithMsg(ErrCodeDataElementNameCnEmpty, "[中文名称]:空")
}

// NameCnTooLong 中文名称过长
func NameCnTooLong() error {
	return errorx.NewWithMsg(ErrCodeDataElementNameCnTooLong, "[中文名称]:长度超过128")
}

// NameEnEmpty 英文名称为空
func NameEnEmpty() error {
	return errorx.NewWithMsg(ErrCodeDataElementNameEnEmpty, "[英文名称]:空")
}

// StdTypeInvalid 标准分类无效
func StdTypeInvalid() error {
	return errorx.NewWithMsg(ErrCodeDataElementStdTypeInvalid, "[标准分类]:输入错误")
}

// DataTypeInvalid 数据类型无效
func DataTypeInvalid() error {
	return errorx.NewWithMsg(ErrCodeDataElementDataTypeInvalid, "[数据类型]:输入错误")
}

// DataLengthInvalid 数据长度无效
func DataLengthInvalid() error {
	return errorx.NewWithMsg(ErrCodeDataElementDataLengthInvalid, "[数据长度]:输入错误")
}

// FileIdNotExist 文件ID不存在
func FileIdNotExist() error {
	return errorx.NewWithMsg(ErrCodeDataElementFileIdNotExist, "一个或多个文件id不存在")
}

// DisableReasonEmpty 停用原因为空
func DisableReasonEmpty() error {
	return errorx.NewWithMsg(ErrCodeDataElementDisableReasonEmpty, "停用必须填写停用原因")
}

// DisableReasonTooLong 停用原因过长
func DisableReasonTooLong() error {
	return errorx.NewWithMsg(ErrCodeDataElementDisableReasonTooLong, "长度超过800")
}
