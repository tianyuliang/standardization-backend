// Code scaffolded by goctl. Safe to edit.
//go:build !mock_logic_off
// +build !mock_logic_off

package stdfile

import (
	"context"
)

// StdFileModel 标准文件模型接口
type StdFileModel interface {
	// 基础CRUD
	Insert(ctx context.Context, data *StdFile) (int64, error)
	FindOne(ctx context.Context, id int64) (*StdFile, error)
	Update(ctx context.Context, data *StdFile) error
	Delete(ctx context.Context, id int64) error

	// 查询方法
	FindByIds(ctx context.Context, ids []int64) ([]*StdFile, error)
	FindByNumber(ctx context.Context, number string) ([]*StdFile, error)
	FindByNameAndOrgType(ctx context.Context, name string, orgType int) ([]*StdFile, error)
	FindByCatalogIds(ctx context.Context, opts *FindOptions) ([]*StdFile, int64, error)
	FindDataExists(ctx context.Context, filterId int64, number string, orgType int, name string, deptIds string) (*StdFile, error)

	// 更新方法
	UpdateState(ctx context.Context, id int64, state int, disableReason string) error
	RemoveCatalog(ctx context.Context, ids []int64, catalogId int64, updateUser string) error
	BatchUpdateState(ctx context.Context, ids []int64, state int, disableReason string) error

	// 批量操作
	DeleteByIds(ctx context.Context, ids []int64) error
}

// FindOptions 查询选项
type FindOptions struct {
	CatalogId    *int64
	Keyword      string
	OrgType      *int
	State        *int
	DepartmentId string
	Page         int
	PageSize     int
	Sort         string
	Direction    string
}
