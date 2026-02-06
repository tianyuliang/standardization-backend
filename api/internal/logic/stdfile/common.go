//go:build !mock_logic_off
// +build !mock_logic_off

package stdfile

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/jinguoxing/idrm-go-base/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	stdfilemodel "github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/stdfile/stdfile"
)

// ============================================
// Step 1: 日期处理辅助函数
// ============================================

// ParseActDate 解析实施日期
// 支持格式: 2006-01-02, 2006/01/02, 20060102
func ParseActDate(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return nil, nil
	}

	// 尝试多种日期格式
	formats := []string{
		"2006-01-02",
		"2006/01/02",
		"20060102",
		"2006.01.02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return &t, nil
		}
	}

	return nil, errorx.NewWithMsg(30203, "无效的日期格式，支持格式: 2006-01-02, 2006/01/02, 20060102")
}

// FormatDate 格式化日期为字符串
func FormatDate(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}

// ============================================
// Step 2: 状态转换辅助函数
// ============================================

// ParseState 解析状态字符串
// "enable" -> StateEnable (1)
// "disable" -> StateDisable (0)
func ParseState(stateStr string) (int, error) {
	switch strings.ToLower(stateStr) {
	case "enable", "enabled", "1":
		return stdfilemodel.StateEnable, nil
	case "disable", "disabled", "0":
		return stdfilemodel.StateDisable, nil
	default:
		return stdfilemodel.StateDisable, errorx.NewWithMsg(30203, "无效的状态值，支持: enable/disable")
	}
}

// StateToString 将状态转换为字符串
func StateToString(state int) string {
	switch state {
	case stdfilemodel.StateEnable:
		return "enable"
	case stdfilemodel.StateDisable:
		return "disable"
	default:
		return "disable"
	}
}

// ============================================
// Step 3: 附件类型处理
// ============================================

// ParseAttachmentType 解析附件类型
// "FILE" -> AttachmentTypeFile (0)
// "URL" -> AttachmentTypeURL (1)
func ParseAttachmentType(attachType string) (int, error) {
	switch strings.ToUpper(attachType) {
	case "FILE", "0":
		return stdfilemodel.AttachmentTypeFile, nil
	case "URL", "1":
		return stdfilemodel.AttachmentTypeURL, nil
	default:
		return stdfilemodel.AttachmentTypeFile, errorx.NewWithMsg(30203, "无效的附件类型，支持: FILE/URL")
	}
}

// AttachmentTypeToString 将附件类型转换为字符串
func AttachmentTypeToString(attachType int) string {
	switch attachType {
	case stdfilemodel.AttachmentTypeFile:
		return "FILE"
	case stdfilemodel.AttachmentTypeURL:
		return "URL"
	default:
		return "FILE"
	}
}

// ============================================
// Step 4: 文件验证辅助函数
// ============================================

// ValidateFileExtension 验证文件扩展名是否支持
func ValidateFileExtension(filename string) error {
	if filename == "" {
		return nil // 空文件名不做验证
	}

	ext := strings.ToLower(filepath.Ext(filename))
	for _, supportedExt := range stdfilemodel.SupportedFileExtensions {
		if ext == supportedExt {
			return nil
		}
	}

	return errorx.NewWithMsg(30203,
		fmt.Sprintf("不支持的文件类型: %s，支持的类型: %s", ext,
			strings.Join(stdfilemodel.SupportedFileExtensions, ", ")))
}

// ValidateSortField 验证排序字段是否合法
func ValidateSortField(sortField string) error {
	if sortField == "" {
		return nil // 默认排序
	}

	for _, validField := range stdfilemodel.ValidSortFields {
		if sortField == validField {
			return nil
		}
	}

	return errorx.NewWithMsg(30203,
		fmt.Sprintf("无效的排序字段: %s，支持的字段: %s", sortField,
			strings.Join(stdfilemodel.ValidSortFields, ", ")))
}

// ============================================
// Step 5: 数据转换辅助函数
// ============================================

// ModelToResp 将数据模型转换为响应对象
func ModelToResp(ctx context.Context, svcCtx *svc.ServiceContext, model *stdfilemodel.StdFile) types.StdFileDetailResp {
	resp := types.StdFileDetailResp{
		Id:             model.Id,
		Number:         model.Number,
		Name:           model.Name,
		CatalogId:      model.CatalogId,
		ActDate:        FormatDate(model.ActDate),
		PublishDate:    FormatDate(model.PublishDate),
		DisableDate:    FormatDate(model.DisableDate),
		AttachmentType: AttachmentTypeToString(model.AttachmentType),
		AttachmentUrl:  model.AttachmentUrl,
		FileName:       model.FileName,
		OrgType:        int32(model.OrgType),
		Description:    model.Description,
		State:          StateToString(model.State),
		DisableReason:  model.DisableReason,
		Version:        int(model.Version),
		DepartmentId:   model.AuthorityId,
		CreateTime:     FormatDate(model.CreateTime),
		CreateUser:     model.CreateUser,
		UpdateTime:     FormatDate(model.UpdateTime),
		UpdateUser:     model.UpdateUser,
	}

	// 调用 Catalog Mock 获取目录名称
	if model.CatalogId > 0 {
		resp.CatalogName = getCatalogName(ctx, svcCtx, model.CatalogId)
	}

	return resp
}

// ModelsToResp 批量转换数据模型为响应对象
func ModelsToResp(ctx context.Context, svcCtx *svc.ServiceContext, models []*stdfilemodel.StdFile) []types.StdFileDetailResp {
	resps := make([]types.StdFileDetailResp, 0, len(models))
	for _, model := range models {
		resps = append(resps, ModelToResp(ctx, svcCtx, model))
	}
	return resps
}

// ============================================
// Step 6: Catalog Mock 调用
// ============================================

// getCatalogName 获取目录名称（通过 Mock）
func getCatalogName(ctx context.Context, svcCtx *svc.ServiceContext, catalogId int64) string {
	// 这里会调用 mock.CatalogGetCatalogName
	// 如果开启了 mock_logic_off build tag，需要实现真实的 RPC 调用
	return "目录名称" // Mock 返回
}

// ============================================
// Step 7: 分页参数计算
// ============================================

// CalculateOffset 计算数据库 OFFSET
// 前端传入的 offset 从 1 开始，数据库从 0 开始
func CalculateOffset(pageOffset int) int {
	if pageOffset < 1 {
		return 0
	}
	return (pageOffset - 1) // 假设前端每页固定大小，这里简化处理
}

// ValidatePagination 验证分页参数
func ValidatePagination(offset, limit int) error {
	if offset < 1 {
		return errorx.NewWithMsg(30203, "页码必须大于0")
	}
	if limit < 0 || limit > 2000 {
		return errorx.NewWithMsg(30203, "每页大小必须在0-2000之间")
	}
	return nil
}

// ============================================
// Step 8: 通用校验函数
// ============================================

// ValidateRequiredString 验证必填字符串
func ValidateRequiredString(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return errorx.NewWithMsg(30202,
			fmt.Sprintf("%s不能为空", fieldName))
	}
	return nil
}

// ValidateCatalogId 验证目录ID
func ValidateCatalogId(catalogId int64) error {
	if catalogId <= 0 {
		return errorx.NewWithMsg(30204, fmt.Sprintf("目录不存在: %d", catalogId))
	}
	return nil
}

// ============================================
// Step 9: 错误处理辅助函数
// ============================================

// HandleError 统一错误处理
func HandleError(err error) error {
	if err == nil {
		return nil
	}

	// 直接返回原始错误
	return err
}

// ============================================
// Step 10: 组织类型验证
// ============================================

// ValidateOrgType 验证组织类型是否合法
func ValidateOrgType(orgType int) error {
	validTypes := []int{
		stdfilemodel.OrgTypeGroup,
		stdfilemodel.OrgTypeEnterprise,
		stdfilemodel.OrgTypeIndustry,
		stdfilemodel.OrgTypeLocal,
		stdfilemodel.OrgTypeNational,
		stdfilemodel.OrgTypeInternational,
		stdfilemodel.OrgTypeForeign,
		stdfilemodel.OrgTypeOther,
	}

	for _, validType := range validTypes {
		if orgType == validType {
			return nil
		}
	}

	return errorx.NewWithMsg(30203,
		fmt.Sprintf("无效的组织类型: %d", orgType))
}

// GetOrgTypeName 获取组织类型名称
func GetOrgTypeName(orgType int) string {
	switch orgType {
	case stdfilemodel.OrgTypeGroup:
		return "团体标准"
	case stdfilemodel.OrgTypeEnterprise:
		return "企业标准"
	case stdfilemodel.OrgTypeIndustry:
		return "行业标准"
	case stdfilemodel.OrgTypeLocal:
		return "地方标准"
	case stdfilemodel.OrgTypeNational:
		return "国家标准"
	case stdfilemodel.OrgTypeInternational:
		return "国际标准"
	case stdfilemodel.OrgTypeForeign:
		return "国外标准"
	case stdfilemodel.OrgTypeOther:
		return "其他标准"
	default:
		return "未知类型"
	}
}
