package catalog

import "context"

// CatalogModel 目录模型接口
type CatalogModel interface {
	// 基础CRUD
	Insert(ctx context.Context, data *Catalog) (int64, error)
	FindOne(ctx context.Context, id int64) (*Catalog, error)
	Update(ctx context.Context, data *Catalog) error
	Delete(ctx context.Context, id int64) error

	// 查询方法
	FindByType(ctx context.Context, catalogType int32) ([]*Catalog, error)
	FindByTypeAndLevel(ctx context.Context, catalogType int32, minLevel int32) ([]*Catalog, error)
	FindByName(ctx context.Context, name string, catalogType int32) ([]*Catalog, error)
	FindByParentId(ctx context.Context, parentId int64) ([]*Catalog, error)
	FindByIds(ctx context.Context, ids []int64) ([]*Catalog, error)
	FindAllByTypeAndLevel(ctx context.Context, catalogType int32, maxLevel int32) ([]*Catalog, error)

	// 树形结构
	FindTree(ctx context.Context, catalogType int32, rootLevel int32) ([]*Catalog, error)
	FindChildren(ctx context.Context, parentId int64, catalogType int32) ([]*Catalog, error)

	// 批量操作
	DeleteByIds(ctx context.Context, ids []int64) error
}

// CatalogModelCtx Catalog模型上下文接口（带事务支持）
type CatalogModelCtx interface {
	CatalogModel
	// Trans 事务支持
	Trans(ctx context.Context, fn func(ctx context.Context, session interface{}) error) error
}
