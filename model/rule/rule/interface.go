package rule

import (
	"context"
	"time"
)

// RuleModel 规则模型接口
type RuleModel interface {
	// Insert 插入规则
	Insert(ctx context.Context, data *Rule) (int64, error)

	// FindOne 查询单条规则
	FindOne(ctx context.Context, id int64) (*Rule, error)

	// Update 更新规则
	Update(ctx context.Context, data *Rule) error

	// Delete 删除规则（物理删除）
	Delete(ctx context.Context, id int64) error

	// FindByIds 批量查询规则
	FindByIds(ctx context.Context, ids []int64) ([]*Rule, error)

	// FindByNameAndOrgType 按名称和组织类型查询（唯一性校验）
	FindByNameAndOrgType(ctx context.Context, name string, orgType int32, departmentIds string) ([]*Rule, error)

	// FindByCatalogIds 按目录ID列表查询（分页）
	FindByCatalogIds(ctx context.Context, catalogIds []int64, opts *FindOptions) ([]*Rule, int64, error)

	// FindDataExists 检查数据是否存在
	FindDataExists(ctx context.Context, filterId int64, name string, departmentIds string) (*Rule, error)

	// UpdateState 更新状态
	UpdateState(ctx context.Context, id int64, state int32, disableReason string) error

	// RemoveCatalog 目录移动
	RemoveCatalog(ctx context.Context, ids []int64, catalogId int64, updateUser string) error

	// UpdateVersionByIds 批量更新版本号
	UpdateVersionByIds(ctx context.Context, ids []int64, updateUser string) error

	// DeleteByIds 批量删除规则
	DeleteByIds(ctx context.Context, ids []int64) error
}

// FindOptions 查询选项
type FindOptions struct {
	OrgType      *int32
	State        *int32
	RuleType     *int32
	Keyword      string
	DepartmentId string
	Page         int
	PageSize     int
	Sort         string
	Direction    string
}

// Rule 规则实体
type Rule struct {
	Id            int64     `db:"f_id" json:"id"`
	Name          string    `db:"f_name" json:"name"`
	CatalogId     int64     `db:"f_catalog_id" json:"catalogId"`
	OrgType       int32     `db:"f_org_type" json:"orgType"`
	Description   string    `db:"f_description" json:"description"`
	RuleType      int32     `db:"f_rule_type" json:"ruleType"`
	Version       int32     `db:"f_version" json:"version"`
	Expression    string    `db:"f_expression" json:"-"`
	State         int32     `db:"f_state" json:"state"`
	DisableReason string    `db:"f_disable_reason" json:"disableReason"`
	AuthorityId   string    `db:"f_authority_id" json:"authorityId"`
	DepartmentIds string    `db:"f_department_ids" json:"departmentIds"`
	ThirdDeptId   string    `db:"f_third_dept_id" json:"thirdDeptId"`
	CreateTime    time.Time `db:"f_create_time" json:"createTime"`
	CreateUser    string    `db:"f_create_user" json:"createUser"`
	UpdateTime    time.Time `db:"f_update_time" json:"updateTime"`
	UpdateUser    string    `db:"f_update_user" json:"updateUser"`
	Deleted       int64     `db:"f_deleted" json:"deleted"`
}

// RuleVo 编码规则视图对象
type RuleVo struct {
	Id                 int64       `json:"id"`
	Name               string      `json:"name"`
	CatalogId          int64       `json:"catalogId"`
	CatalogName        string      `json:"catalogName,omitempty"`
	FullCatalogName    string      `json:"fullCatalogName,omitempty"`
	OrgType            int32       `json:"orgType"`
	Description        string      `json:"description,omitempty"`
	RuleType           string      `json:"ruleType"`
	Version            int32       `json:"version"`
	Regex              string      `json:"regex,omitempty"`
	Custom             []RuleCustom `json:"custom,omitempty"`
	State              string      `json:"state"`
	DisableReason      string      `json:"disableReason,omitempty"`
	StdFiles           []int64     `json:"stdFiles,omitempty"`
	UsedFlag           bool        `json:"usedFlag"`
	DepartmentId       string      `json:"departmentId,omitempty"`
	DepartmentName     string      `json:"departmentName,omitempty"`
	DepartmentPathNames string     `json:"departmentPathNames,omitempty"`
	CreateTime         string      `json:"createTime,omitempty"`
	CreateUser         string      `json:"createUser,omitempty"`
	UpdateTime         string      `json:"updateTime,omitempty"`
	UpdateUser         string      `json:"updateUser,omitempty"`
}

// RuleCustom 自定义规则配置
type RuleCustom struct {
	SegmentLength int    `json:"segment_length"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	Value         string `json:"value"`
}
