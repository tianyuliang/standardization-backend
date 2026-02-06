// Code scaffolded by goctl. Safe to edit.
//go:build !mock_logic_off
// +build !mock_logic_off

package stdfile

import (
	"time"
)

// StdFile 标准文件数据模型
type StdFile struct {
	Id             int64      `db:"f_id"`
	Number         string     `db:"f_number"`
	Name           string     `db:"f_name"`
	CatalogId      int64      `db:"f_catalog_id"`
	ActDate        *time.Time `db:"f_act_date"`
	PublishDate    *time.Time `db:"f_publish_date"`
	DisableDate    *time.Time `db:"f_disable_date"`
	AttachmentType int        `db:"f_attachment_type"`
	AttachmentUrl  string     `db:"f_attachment_url"`
	FileName       string     `db:"f_file_name"`
	OrgType        int        `db:"f_org_type"`
	Description    string     `db:"f_description"`
	State          int        `db:"f_state"`
	DisableReason  string     `db:"f_disable_reason"`
	AuthorityId    string     `db:"f_authority_id"`
	DepartmentIds  string     `db:"f_department_ids"`
	ThirdDeptId    string     `db:"f_third_dept_id"`
	Version        int        `db:"f_version"`
	CreateTime     *time.Time `db:"f_create_time"`
	CreateUser     string     `db:"f_create_user"`
	UpdateTime     *time.Time `db:"f_update_time"`
	UpdateUser     string     `db:"f_update_user"`
	Deleted        int        `db:"f_deleted"`
}

// TableName 返回表名
func (StdFile) TableName() string {
	return TableName
}
