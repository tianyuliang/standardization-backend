// Package mock provides mock services for dataelement logic
package mock

import (
	"context"

	"github.com/sony/sonyflake"
)

// ============================================
// Catalog Mock Service
// ============================================

// CheckCatalogExist 校验目录是否存在
// TODO: 替换为 Catalog RPC 调用
// Java: iDeCatalogInfoService.checkCatalogIsExist(catalogId)
func CheckCatalogExist(ctx context.Context, catalogId int64) (bool, string, error) {
	// MOCK: 模拟目录存在
	// 实际实现需要调用 Catalog RPC
	return true, "默认目录", nil
}

// GetCatalogName 获取目录名称
// TODO: 替换为 Catalog RPC 调用
// Java: catalog.getCatalogName()
func GetCatalogName(ctx context.Context, catalogId int64) (string, error) {
	// MOCK: 返回默认目录名
	return "默认目录", nil
}

// GetChildCatalogIds 获取子目录ID列表
// TODO: 替换为 Catalog RPC 调用
// Java: iDeCatalogInfoService.getIDList(catalogId)
func GetChildCatalogIds(ctx context.Context, catalogId int64) ([]int64, error) {
	// MOCK: 返回当前目录ID
	return []int64{catalogId}, nil
}

// IsDataElementCatalog 校验是否为数据元目录
// TODO: 替换为 Catalog RPC 调用
// Java: catalog.getType().equals(CatalogTypeEnum.DataElement)
func IsDataElementCatalog(ctx context.Context, catalogId int64) (bool, error) {
	// MOCK: 假设是数据元目录
	return true, nil
}

// ============================================
// Dict Mock Service
// ============================================

// CheckDictExist 校验码表是否存在
// TODO: 替换为 Dict RPC 调用
// Java: dictService.queryById(dictId)
func CheckDictExist(ctx context.Context, dictCode string) (bool, error) {
	// MOCK: 假设码表存在
	return true, nil
}

// GetDictValues 获取码表值列表
// TODO: 替换为 Dict RPC 调用
// Java: dictService.getDictValues(dictCode)
func GetDictValues(ctx context.Context, dictCode string) ([]string, error) {
	// MOCK: 返回默认值列表
	return []string{"值1", "值2", "值3"}, nil
}

// ============================================
// Rule Mock Service
// ============================================

// CheckRuleExist 校验规则是否存在
// TODO: 替换为 Rule RPC 调用
// Java: ruleService.queryById(ruleId)
func CheckRuleExist(ctx context.Context, ruleId int64) (bool, error) {
	// MOCK: 假设规则存在
	return true, nil
}

// ============================================
// StdFile Mock Service
// ============================================

// GetStdFileByIds 批量获取标准文件信息
// TODO: 替换为 StdFile RPC 调用
// Java: stdFileService.queryByIds(fileIds)
func GetStdFileByIds(ctx context.Context, fileIds []int64) (map[int64]string, error) {
	// MOCK: 返回文件ID到文件名的映射
	result := make(map[int64]string)
	for _, id := range fileIds {
		result[id] = "标准文件"
	}
	return result, nil
}

// CheckStdFileExist 校验标准文件是否存在
// TODO: 替换为 StdFile RPC 调用
func CheckStdFileExist(ctx context.Context, fileId int64) (bool, error) {
	// MOCK: 假设文件存在
	return true, nil
}

// ============================================
// MQ Mock Service
// ============================================

// SendMqMessage 发送MQ消息
// TODO: 替换为 Kafka 发送
// Java: kafkaProducerService.sendMessage(MqTopic.MQ_MESSAGE_SAILOR, mqInfo)
func SendMqMessage(ctx context.Context, topic string, message interface{}, user string) error {
	// MOCK: 模拟发送成功
	return nil
}

// GenerateTaskId 生成任务ID
// TODO: 使用分布式ID生成器
func GenerateTaskId() string {
	sf := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, _ := sf.NextID()
	return string(rune(id))
}
