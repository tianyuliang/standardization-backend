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

	// Delete 删除规则（逻辑删除）
	Delete(ctx context.Context, id int64) error

	// FindByIds 批量查询规则
	FindByIds(ctx context.Context, ids []int64) ([]*Rule, error)

	// FindByNameAndOrgType 按名称和组织类型查询（唯一性校验）
	FindByNameAndOrgType(ctx context.Context, name string, orgType int32, departmentIds string) ([]*Rule, error)

	// FindByCatalogIds 按目录ID列表查询（分页）
	FindByCatalogIds(ctx context.Context, catalogIds []int64, offset int, limit int) ([]*Rule, int64, error)

	// FindDataExists 检查数据是否存在
	FindDataExists(ctx context.Context, filterId int64, name string, departmentIds string) (bool, error)

	// UpdateState 更新状态
	UpdateState(ctx context.Context, id int64, state int32, disableReason string) error

	// RemoveCatalog 目录移动
	RemoveCatalog(ctx context.Context, ids []int64, catalogId int64, updateUser string) error

	// UpdateVersionByIds 批量更新版本号
	UpdateVersionByIds(ctx context.Context, ids []int64) error
}

// Rule 规则实体
type Rule struct {
	Id            int64      `gorm:"column:f_id;primaryKey" json:"id"`
	Name          string     `gorm:"column:f_name;size:128;notNull" json:"name"`
	CatalogId     int64      `gorm:"column:f_catalog_id;notNull" json:"catalogId"`
	OrgType       int32      `gorm:"column:f_org_type;notNull" json:"orgType"`
	Description   string     `gorm:"column:f_description;size:300" json:"description"`
	RuleType      int32      `gorm:"column:f_rule_type;notNull;default:0" json:"ruleType"`
	Version       int32      `gorm:"column:f_version;notNull;default:1" json:"version"`
	Expression    string     `gorm:"column:f_expression;size:1024;notNull" json:"-"`
	State         int32      `gorm:"column:f_state;notNull;default:1" json:"state"`
	DisableReason string     `gorm:"column:f_disable_reason;size:1024" json:"disableReason"`
	AuthorityId   string     `gorm:"column:f_authority_id;size:100" json:"authorityId"`
	DepartmentIds string     `gorm:"column:f_department_ids;size:350" json:"departmentIds"`
	ThirdDeptId   string     `gorm:"column:f_third_dept_id;size:36" json:"thirdDeptId"`
	CreateTime    time.Time  `gorm:"column:f_create_time" json:"createTime"`
	CreateUser    string     `gorm:"column:f_create_user;size:128" json:"createUser"`
	UpdateTime    time.Time  `gorm:"column:f_update_time" json:"updateTime"`
	UpdateUser    string     `gorm:"column:f_update_user;size:128" json:"updateUser"`
	Deleted       int64      `gorm:"column:f_deleted;notNull;default:0" json:"deleted"`
}

func (Rule) TableName() string {
	return "t_rule"
}
