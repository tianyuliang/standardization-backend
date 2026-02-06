package catalog

import "time"

// Catalog 目录实体
type Catalog struct {
	Id          int64     `db:"f_id" json:"id"`
	CatalogName string    `db:"f_catalog_name" json:"catalogName"`
	Description string    `db:"f_description" json:"description"`
	Level       int32     `db:"f_level" json:"level"`
	ParentId    int64     `db:"f_parent_id" json:"parentId"`
	Type        int32     `db:"f_type" json:"type"`
	AuthorityId string    `db:"f_authority_id" json:"authorityId"`
	CreateTime  time.Time `db:"f_create_time" json:"createTime"`
	CreateUser  string    `db:"f_create_user" json:"createUser"`
	UpdateTime  time.Time `db:"f_update_time" json:"updateTime"`
	UpdateUser  string    `db:"f_update_user" json:"updateUser"`
	Deleted     int64     `db:"f_deleted" json:"deleted"`
}

// CatalogTreeNodeVo 目录树节点（用于查询响应）
type CatalogTreeNodeVo struct {
	Id           int64                `json:"id"`
	CatalogName  string               `json:"catalogName"`
	Level        int32                `json:"level"`
	ParentId     int64                `json:"parentId"`
	Type         int32                `json:"type"`
	Children     []*CatalogTreeNodeVo `json:"children,omitempty"`
	Count        int32                `json:"count,omitempty"`
	HaveChildren bool                 `json:"haveChildren"`
}

// CatalogWithFileVo 目录及文件树节点
type CatalogWithFileVo struct {
	Id           int64                `json:"id"`
	CatalogName  string               `json:"catalogName"`
	Level        int32                `json:"level"`
	ParentId     int64                `json:"parentId"`
	Type         int32                `json:"type"`
	Children     []*CatalogWithFileVo `json:"children,omitempty"`
	Files        []*FileCountVo       `json:"files,omitempty"`
	HaveChildren bool                 `json:"haveChildren"`
}

// FileCountVo 文件计数
type FileCountVo struct {
	FileId    int64  `json:"fileId"`
	FileName  string `json:"fileName"`
	CatalogId int64  `json:"catalogId"`
}

// CatalogInfoVo 目录信息（平铺列表）
type CatalogInfoVo struct {
	Id          int64  `json:"id"`
	CatalogName string `json:"catalogName"`
	Level       int32  `json:"level"`
	ParentId    int64  `json:"parentId"`
	Type        int32  `json:"type"`
}
