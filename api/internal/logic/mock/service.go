// Package mock provides shared mock services for all modules
package mock

import (
	"context"
	"fmt"

	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	"github.com/sony/sonyflake"
)

// ============================================
// Catalog Service Mock
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
// Dict Service Mock
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
// Rule Service Mock
// ============================================

// CheckRuleExist 校验规则是否存在
// TODO: 替换为 Rule RPC 调用
// Java: ruleService.queryById(ruleId)
func CheckRuleExist(ctx context.Context, ruleId int64) (bool, error) {
	// MOCK: 假设规则存在
	return true, nil
}

// ============================================
// StdFile Service Mock
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
// Java: stdFileService.queryById(fileId)
func CheckStdFileExist(ctx context.Context, fileId int64) (bool, error) {
	// MOCK: 假设文件存在
	return true, nil
}

// ============================================
// DataElement Service Mock
// ============================================

// GetDataElementInfo 获取数据元信息
// TODO: 替换为 DataElement RPC 调用
// Java: dataElementInfoService.getById(dataElementId)
func GetDataElementInfo(ctx context.Context, svcCtx *svc.ServiceContext, dataElementId int64) (*types.DataElementInfo, error) {
	// MOCK: 默认返回nil
	return nil, nil
}

// GetDataElementDetailVo 获取数据元详情
// TODO: 替换为 DataElement RPC 调用
// Java: dataElementInfoService.getDetailVo(dataElementId)
func GetDataElementDetailVo(ctx context.Context, svcCtx *svc.ServiceContext, dataElementId int64) (*types.DataElementDetailVo, error) {
	// MOCK: 默认返回nil
	return nil, nil
}

// ============================================
// MQ Service Mock
// ============================================

// SendMqMessage 发送MQ消息
// TODO: 替换为 Kafka 发送
// Java: kafkaProducerService.sendMessage(MqTopic.MQ_MESSAGE_SAILOR, mqInfo)
func SendMqMessage(ctx context.Context, topic string, message interface{}, user string) error {
	// MOCK: 模拟发送成功
	return nil
}

// ============================================
// ID Generator Mock
// ============================================

var sonyFlake *sonyflake.Sonyflake

func init() {
	// 初始化 Sonyflake 分布式ID生成器
	sonyFlake = sonyflake.NewSonyflake(sonyflake.Settings{})
}

// GenerateTaskId 生成任务ID (分布式ID)
// TODO: 使用分布式ID生成器
func GenerateTaskId() string {
	id, _ := sonyFlake.NextID()
	// 转换为字符串，实际使用时可能需要特定格式
	return fmt.Sprintf("%d", id)
}

// GenerateInt64Id 生成int64类型的ID
func GenerateInt64Id() int64 {
	id, _ := sonyFlake.NextID()
	return int64(id)
}

// ============================================
// Recommendation Service Mock
// ============================================

// CallStdRecService 调用标准推荐服务
// TODO: 实现HTTP调用推荐服务
// Java: HttpClient.post(recommendationServiceUrl + "/std-rec/rec", req)
func CallStdRecService(ctx context.Context, svcCtx *svc.ServiceContext, req *types.StdRecReq) (*types.StdRecResp, error) {
	// MOCK: 默认返回空结果
	return &types.StdRecResp{
		Code:        "0",
		Description: "成功",
		Data:        []types.StdRecItem{},
	}, nil
}

// CallRuleRecService 调用规则推荐服务
// TODO: 实现HTTP调用规则推荐服务
// Java: HttpClient.post(recommendationServiceUrl + "/rule-rec/rec", req)
func CallRuleRecService(ctx context.Context, svcCtx *svc.ServiceContext, req *types.RuleRecReq) (*types.StdRecResp, error) {
	// MOCK: 默认返回空结果
	return &types.StdRecResp{
		Code:        "0",
		Description: "成功",
		Data:        []types.StdRecItem{},
	}, nil
}

// ============================================
// AfService Mock
// ============================================

// GetTaskDetailDto 获取任务详情
// TODO: 调用 AfService RPC 获取任务详情
// Java: afService.getTaskDetailDto(taskId)
func GetTaskDetailDto(ctx context.Context, svcCtx *svc.ServiceContext, taskId string) (*types.TaskDetailResp, error) {
	// MOCK: 默认返回nil
	return nil, nil
}

// ============================================
// Webhook Service Mock
// ============================================

// SendTaskCallback 发送任务完成回调
// TODO: 实现HTTP POST回调
// Java: HttpClient.post(webhook, callbackData)
func SendTaskCallback(ctx context.Context, svcCtx *svc.ServiceContext, webhook string, taskId string) error {
	// MOCK: 默认成功
	return nil
}

// ============================================
// 额外的 Catalog 函数别名 (兼容其他模块命名)
// ============================================

// CatalogCheckExist 校验目录是否存在 (rule模块使用)
// 别名指向 CheckCatalogExist
func CatalogCheckExist(ctx context.Context, svcCtx *svc.ServiceContext, catalogId int64) bool {
	exists, _, _ := CheckCatalogExist(ctx, catalogId)
	return exists
}

// CatalogGetCatalogName 获取目录名称 (rule模块使用)
// 别名指向 GetCatalogName
func CatalogGetCatalogName(ctx context.Context, svcCtx *svc.ServiceContext, catalogId int64) string {
	name, _ := GetCatalogName(ctx, catalogId)
	return name
}

// CatalogGetChildIds 获取子目录ID列表 (rule模块使用)
// 别名指向 GetChildCatalogIds
func CatalogGetChildIds(ctx context.Context, svcCtx *svc.ServiceContext, catalogId int64) []int64 {
	ids, _ := GetChildCatalogIds(ctx, catalogId)
	return ids
}

// CatalogIsStdFileCatalog 校验是否为标准文件目录 (rule模块使用)
func CatalogIsStdFileCatalog(ctx context.Context, svcCtx *svc.ServiceContext, catalogId int64) bool {
	exists, _ := IsDataElementCatalog(ctx, catalogId)
	return exists
}

// CatalogIsRootCatalog 校验是否为根目录 (rule模块使用)
func CatalogIsRootCatalog(ctx context.Context, svcCtx *svc.ServiceContext, catalogId int64) bool {
	// MOCK: 默认不是根目录
	return false
}

// ============================================
// 额外的 DataElement 函数别名 (catalog模块使用)
// ============================================

// DataElementGetPageList 检查目录下是否存在数据元 (catalog模块使用)
func DataElementGetPageList(ctx context.Context, svcCtx *svc.ServiceContext, catalogId int64) bool {
	// MOCK: 返回false，表示目录下不存在数据元
	return false
}

// DataElementGetCountMapGroupByCatalog 按目录分组统计数据元数量
func DataElementGetCountMapGroupByCatalog(ctx context.Context, svcCtx *svc.ServiceContext) map[int64]int32 {
	// MOCK: 返回空map
	return make(map[int64]int32)
}
