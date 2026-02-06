// Code scaffolded by goctl. Safe to edit.
//go:build !mock_logic_off
// +build !mock_logic_off

package stdfile

// ========== 标准组织类型枚举 (OrgTypeEnum) ==========

const (
	OrgTypeGroup         int = 0  // 团体标准
	OrgTypeEnterprise    int = 1  // 企业标准
	OrgTypeIndustry      int = 2  // 行业标准
	OrgTypeLocal         int = 3  // 地方标准
	OrgTypeNational      int = 4  // 国家标准
	OrgTypeInternational int = 5  // 国际标准
	OrgTypeForeign       int = 6  // 国外标准
	OrgTypeOther         int = 99 // 其他标准
)

// ========== 附件类型枚举 (AttachmentTypeEnum) ==========

const (
	AttachmentTypeFile int = 0 // 文件附件
	AttachmentTypeURL  int = 1 // 外置链接
)

// ========== 启用停用状态枚举 (EnableDisableStatusEnum) ==========

const (
	StateDisable int = 0 // 停用
	StateEnable  int = 1 // 启用
)

// ========== 支持的文件类型 ==========

var SupportedFileExtensions = []string{
	".doc", ".docx",
	".pdf",
	".txt",
	".ppt", ".pptx",
	".xls", ".xlsx",
}

// ========== 文件大小限制 ==========

const MaxFileSize = 30 * 1024 * 1024 // 30MB

// ========== 可排序的字段 ==========

var ValidSortFields = []string{
	"f_create_time",
	"f_update_time",
	"f_act_date",
	"f_disable_date",
	"f_state",
	"f_id",
}

// ========== 表字段名 ==========

const (
	TableName = "t_std_file_mgr"

	// 字段名
	ColumnId             = "f_id"
	ColumnNumber         = "f_number"
	ColumnName           = "f_name"
	ColumnCatalogId      = "f_catalog_id"
	ColumnActDate        = "f_act_date"
	ColumnPublishDate    = "f_publish_date"
	ColumnDisableDate    = "f_disable_date"
	ColumnAttachmentType = "f_attachment_type"
	ColumnAttachmentUrl  = "f_attachment_url"
	ColumnFileName       = "f_file_name"
	ColumnOrgType        = "f_org_type"
	ColumnDescription    = "f_description"
	ColumnState          = "f_state"
	ColumnDisableReason  = "f_disable_reason"
	ColumnAuthorityId    = "f_authority_id"
	ColumnDepartmentIds  = "f_department_ids"
	ColumnThirdDeptId    = "f_third_dept_id"
	ColumnVersion        = "f_version"
	ColumnCreateTime     = "f_create_time"
	ColumnCreateUser     = "f_create_user"
	ColumnUpdateTime     = "f_update_time"
	ColumnUpdateUser     = "f_update_user"
	ColumnDeleted        = "f_deleted"
)
