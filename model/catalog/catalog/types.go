package catalog

import "time"

// Catalog 目录实体
type Catalog struct {
	Id          string    `gorm:"column:f_id;primaryKey;size:20" json:"id"`                                 // 目录唯一标识（string 存储以兼容 JSON）
	CatalogName string    `gorm:"column:f_catalog_name;size:20;not null" json:"catalogName"`                  // 目录名称
	Description string    `gorm:"column:f_description;size:255" json:"description"`                           // 目录说明
	Level       int32     `gorm:"column:f_level;unsigned;not null" json:"level"`                              // 目录级别
	ParentId    string    `gorm:"column:f_parent_id;size:20;not null;index:idx_parent_id" json:"parentId"`   // 父级标识
	Type        int32     `gorm:"column:f_type;unsigned;not null;index:idx_type_level" json:"type"`           // 目录类型
	AuthorityId *string   `gorm:"column:f_authority_id;size:20" json:"authorityId,omitempty"`                // 权限域（可选）
	CreatedAt   time.Time `gorm:"column:f_created_at;not null" json:"createdAt,omitempty"`                   // 创建时间
	UpdatedAt   time.Time `gorm:"column:f_updated_at;not null" json:"updatedAt,omitempty"`                   // 更新时间
}

// TableName 指定表名
func (Catalog) TableName() string {
	return "t_de_catalog_info"
}
