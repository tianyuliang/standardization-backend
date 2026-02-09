package dict

import (
	"context"
	"time"
)

// DictModel 码表模型接口
type DictModel interface {
	// Insert 插入码表
	Insert(ctx context.Context, data *Dict) (int64, error)

	// FindOne 按ID查询单条码表
	FindOne(ctx context.Context, id int64) (*Dict, error)

	// FindByCode 按Code查询单条码表
	FindByCode(ctx context.Context, code int64) (*Dict, error)

	// Update 更新码表
	Update(ctx context.Context, data *Dict) error

	// Delete 删除码表（物理删除）
	Delete(ctx context.Context, id int64) error

	// FindByIds 批量查询码表
	FindByIds(ctx context.Context, ids []int64) ([]*Dict, error)

	// FindByChNameAndOrgType 按中文名和组织类型查询（唯一性校验）
	FindByChNameAndOrgType(ctx context.Context, chName string, orgType int32) (*Dict, error)

	// FindByEnNameAndOrgType 按英文名和组织类型查询（唯一性校验）
	FindByEnNameAndOrgType(ctx context.Context, enName string, orgType int32) (*Dict, error)

	// FindByCatalogIds 按目录ID列表查询（分页）
	FindByCatalogIds(ctx context.Context, opts *FindOptions) ([]*Dict, int64, error)

	// FindByStdFileCatalog 按标准文件目录查询码表
	FindByStdFileCatalog(ctx context.Context, opts *FindOptions) ([]*Dict, int64, error)

	// FindByFileId 按文件ID查询码表
	FindByFileId(ctx context.Context, fileId int64) ([]*Dict, error)

	// FindDataExists 检查数据是否存在
	FindDataExists(ctx context.Context, chName, enName string, orgType int32, filterId int64, deptIds string) (*Dict, error)

	// UpdateState 更新状态
	UpdateState(ctx context.Context, id int64, state int32, disableReason string) error

	// UpdateVersionByIds 批量更新版本号
	UpdateVersionByIds(ctx context.Context, ids []int64, updateUser string) error

	// DeleteByIds 批量删除码表
	DeleteByIds(ctx context.Context, ids []int64) error
}

// DictEnumModel 码值模型接口
type DictEnumModel interface {
	// Insert 插入码值
	Insert(ctx context.Context, data *DictEnum) (int64, error)

	// FindByDictId 按码表ID查询所有码值
	FindByDictId(ctx context.Context, dictId int64) ([]*DictEnum, error)

	// FindPageByDictId 按码表ID分页查询码值
	FindPageByDictId(ctx context.Context, dictId int64, keyword string, offset, limit int) ([]*DictEnum, int64, error)

	// DeleteByDictId 删除指定码表的所有码值
	DeleteByDictId(ctx context.Context, dictId int64) error

	// CheckDuplicateCode 检查码值是否重复
	CheckDuplicateCode(ctx context.Context, dictId int64, codes []string) (bool, error)
}

// RelationDictFileModel 码表-文件关系模型接口
type RelationDictFileModel interface {
	// InsertBatch 批量插入关系
	InsertBatch(ctx context.Context, data []*RelationDictFile) error

	// DeleteByDictId 删除指定码表的所有关系
	DeleteByDictId(ctx context.Context, dictId int64) error

	// FindByDictId 查询指定码表关联的文件ID列表
	FindByDictId(ctx context.Context, dictId int64) ([]int64, error)
}

// FindOptions 查询选项
type FindOptions struct {
	CatalogId    *int64
	FileId       *int64
	OrgType      *int32
	State        *int32
	Keyword      string
	DepartmentId string
	Page         int
	PageSize     int
	Sort         string
	Direction    string
}

// Dict 码表实体
type Dict struct {
	Id            int64     `db:"f_id" json:"id"`
	Code          int64     `db:"f_code" json:"code"`
	ChName        string    `db:"f_ch_name" json:"chName"`
	EnName        string    `db:"f_en_name" json:"enName"`
	Description   string    `db:"f_description" json:"description"`
	CatalogId     int64     `db:"f_catalog_id" json:"catalogId"`
	OrgType       int32     `db:"f_org_type" json:"orgType"`
	Version       int32     `db:"f_version" json:"version"`
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

// DictEnum 码值实体
type DictEnum struct {
	Id         int64     `db:"f_id" json:"id"`
	DictId     int64     `db:"f_dict_id" json:"dictId"`
	Code       string    `db:"f_code" json:"code"`
	Value      string    `db:"f_value" json:"value"`
	CreateTime time.Time `db:"f_create_time" json:"createTime"`
	CreateUser string    `db:"f_create_user" json:"createUser"`
}

// RelationDictFile 码表-文件关系实体
type RelationDictFile struct {
	Id     int64 `db:"f_id" json:"id"`
	DictId int64 `db:"f_dict_id" json:"dictId"`
	FileId int64 `db:"f_file_id" json:"fileId"`
}
