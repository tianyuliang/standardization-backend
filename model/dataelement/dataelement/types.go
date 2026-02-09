// Code scaffolded by speckit. Safe to edit.

package dataelement

// TableNameDataElementInfo 数据元表名
const TableNameDataElementInfo = "t_data_element_info"

// DataElementState 数据元状态
const (
	DataElementStateDisabled int32 = 0 // 停用
	DataElementStateEnabled  int32 = 1 // 启用
)

// DataElement 数据元实体
type DataElement struct {
	Id              int64  `db:"f_id" json:"id"`
	Code            int64  `db:"f_de_id" json:"code"`
	NameEn          string `db:"f_name_en" json:"nameEn"`
	NameCn          string `db:"f_name_cn" json:"nameCn"`
	Synonym         string `db:"f_synonym" json:"synonym"`
	StdType         int32  `db:"f_std_type" json:"stdType"`
	DataType        int32  `db:"f_data_type" json:"dataType"`
	DataLength      *int   `db:"f_data_length" json:"dataLength"`
	DataPrecision   *int   `db:"f_data_precision" json:"dataPrecision"`
	DictCode        *int64 `db:"f_dict_code" json:"dictCode"`
	RuleId          *int64 `db:"f_rule_id" json:"ruleId"`
	RelationType    string `db:"f_relation_type" json:"relationType"`
	CatalogId       int64  `db:"f_catalog_id" json:"catalogId"`
	LabelId         *int64 `db:"f_label_id" json:"labelId"`
	Description     string `db:"f_description" json:"description"`
	Version         int    `db:"f_version" json:"version"`
	State           int32  `db:"f_state" json:"state"`
	AuthorityId     string `db:"f_authority_id" json:"authorityId"`
	DepartmentIds   string `db:"f_department_ids" json:"departmentIds"`
	ThirdDeptId     string `db:"f_third_dept_id" json:"thirdDeptId"`
	DisableReason   string `db:"f_disable_reason" json:"disableReason"`
	CreateTime      string `db:"f_create_time" json:"createTime"`
	CreateUser      string `db:"f_create_user" json:"createUser"`
	UpdateTime      string `db:"f_update_time" json:"updateTime"`
	UpdateUser      string `db:"f_update_user" json:"updateUser"`
	Deleted         int64  `db:"f_deleted" json:"deleted"`
}

// RelationType 关联类型枚举
const (
	RelationTypeNo         = "no"         // 无关联
	RelationTypeCodeTable  = "codeTable"  // 码表关联
	RelationTypeCodeRule   = "codeRule"   // 规则关联
)

// FindOptions 查询选项
type FindOptions struct {
	CatalogId       *int64
	StdType         *int32
	State           *int32
	DataType        *int32
	RelationType    *string
	Keyword         string
	DepartmentIds   string
	FileCatalogId   *int64
	FileId          *int64
	Page            int
	PageSize        int
	Sort            string
	Direction       string
}

// PageOptions 分页选项
type PageOptions struct {
	Page     int
	PageSize int
}
