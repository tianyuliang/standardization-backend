// Package dataelement provides common utilities for dataelement logic
package dataelement

import (
	"context"
	"fmt"
	"strings"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/mock"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/dataelement/dataelement"
)

// ============================================
// 值域计算
// ============================================

// CalculateValueRange 计算值域
// 对应 Java: DataElementServiceImpl.calculateValueRange()
// 关联码表时返回"[码值1,码值2,...]"，否则按数据类型计算
func CalculateValueRange(ctx context.Context, dictCode string, dataType int32, dataLength, dataPrecision int) string {
	if dictCode != "" {
		return calculateValueRangeFromDict(ctx, dictCode)
	}
	return calculateValueRangeFromDataType(dataType, dataLength, dataPrecision)
}

// calculateValueRangeFromDict 从码表计算值域
func calculateValueRangeFromDict(ctx context.Context, dictCode string) string {
	values, err := mock.GetDictValues(ctx, dictCode)
	if err != nil || len(values) == 0 {
		return ""
	}
	return fmt.Sprintf("[%s]", strings.Join(values, ","))
}

// calculateValueRangeFromDataType 根据数据类型计算值域
func calculateValueRangeFromDataType(dataType int32, dataLength, dataPrecision int) string {
	switch dataType {
	case 0: // Number
		if dataPrecision > 0 {
			return fmt.Sprintf("-%s ~ %d",
				power10(dataLength-dataPrecision),
				power10Int(dataLength-dataPrecision)-1)
		}
		return fmt.Sprintf("-%s ~ %d", power10(dataLength), power10Int(dataLength)-1)

	case 1: // Decimal
		if dataPrecision > 0 {
			return fmt.Sprintf("-%s.%s ~ %d.%s",
				power10(dataLength-dataPrecision),
				strings.Repeat("9", dataPrecision),
				power10Int(dataLength-dataPrecision)-1,
				strings.Repeat("9", dataPrecision))
		}
		return fmt.Sprintf("-%s ~ %d", power10(dataLength), power10Int(dataLength)-1)

	case 3: // Char
		if dataLength > 0 {
			return fmt.Sprintf("字符数 ≤ %d", dataLength)
		}
		return "字符数 ≤ 65535"

	case 4: // Date
		return "YYYY-MM-DD"

	case 5: // DateTime
		return "YYYY-MM-DD HH:mm:ss"

	case 7: // Boolean
		return "true/false"

	default:
		return ""
	}
}

// power10 计算10的n次方（返回字符串格式）
func power10(n int) string {
	if n <= 0 {
		return "1"
	}
	return "1" + strings.Repeat("0", n)
}

// power10Int 计算10的n次方（返回整数，用于算术运算）
func power10Int(n int) int {
	result := 1
	for i := 0; i < n; i++ {
		result *= 10
	}
	return result
}

// ============================================
// 版本变更检测
// ============================================

// CheckVersionChange 检查是否需要递增版本号
// 对应 Java: DataElementServiceImpl.checkVersionChange()
// 修改以下字段时版本号+1：relationType、dictCode、ruleId、nameCn、nameEn
func CheckVersionChange(oldData *dataelement.DataElement, newData *types.UpdateDataElementReq) bool {
	if oldData == nil || newData == nil {
		return false
	}

	// 关联类型变更
	if oldData.RelationType != newData.RelationType {
		return true
	}

	// 名称变更
	if oldData.NameCn != newData.NameCn || oldData.NameEn != newData.NameEn {
		return true
	}

	// 关联对象变更
	if oldData.RelationType == "codeTable" {
		if oldData.DictCode != nil && newData.DictCode != "" {
			// 比较非空情况
			return true
		}
	}

	if oldData.RelationType == "codeRule" {
		if oldData.RuleId != nil && newData.RuleId > 0 {
			// 比较非空情况
			return true
		}
	}

	return false
}
