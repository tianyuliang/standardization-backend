package rule

import "time"

// RelationRuleFile 规则-文件关联实体
type RelationRuleFile struct {
	Id     int64 `gorm:"column:f_id;primaryKey" json:"id"`
	RuleId int64 `gorm:"column:f_rule_id;notNull" json:"ruleId"`
	FileId int64 `gorm:"column:f_file_id;notNull" json:"fileId"`
}

func (RelationRuleFile) TableName() string {
	return "t_relation_rule_file"
}

// RuleVo 规则视图对象
type RuleVo struct {
	Id            int64    `json:"id"`
	Name          string   `json:"name"`
	CatalogId     int64    `json:"catalogId"`
	CatalogName   string   `json:"catalogName"`     // 从 RPC 获取
	OrgType       int32    `json:"orgType"`
	Description   string   `json:"description"`
	RuleType      string   `json:"ruleType"`         // "REGEX" or "CUSTOM"
	Version       int32    `json:"version"`
	Expression    string   `json:"expression"`       // 隐藏原始表达式
	State         int32    `json:"state"`
	DisableReason string   `json:"disableReason"`
	AuthorityId   string   `json:"authorityId"`
	DepartmentIds  string   `json:"departmentIds"`    // 返回最后一段
	ThirdDeptId   string   `json:"thirdDeptId"`
	CreateTime    string   `json:"createTime"`
	CreateUser    string   `json:"createUser"`
	UpdateTime    string   `json:"updateTime"`
	UpdateUser    string   `json:"updateUser"`
	StdFileIds    []int64  `json:"stdFileIds"`
	Used          bool     `json:"used"`             // 是否被数据元引用
}

// RuleCustom 自定义配置
type RuleCustom struct {
	Type          string `json:"type"`           // DICT/NUMBER/ENGLISH_LETTERS/CHINESE_CHARACTERS/ANY_CHARACTERS/DATE/SPLIT_STR
	Value         string `json:"value"`          // 值
	SegmentLength int    `json:"segmentLength"`   // 段长度
	Count         int    `json:"count"`          // 数量
}

// FindOptions 查询选项
type FindOptions struct {
	Keyword      string
	CatalogIds   []int64
	OrgType      int32
	State        int32
	RuleType     string
	Offset       int
	Limit        int
}
