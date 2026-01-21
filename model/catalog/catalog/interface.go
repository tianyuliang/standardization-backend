package catalog

import "context"

// CatalogModel 目录数据访问接口
type CatalogModel interface {
	// Insert 创建目录
	Insert(ctx context.Context, data *Catalog) (*Catalog, error)

	// FindOne 根据 ID 查询单个目录
	FindOne(ctx context.Context, id string) (*Catalog, error)

	// FindByParent 查询父目录下的所有子目录
	FindByParent(ctx context.Context, parentId string) ([]*Catalog, error)

	// FindByType 根据类型查询所有目录
	FindByType(ctx context.Context, catalogType int32) ([]*Catalog, error)

	// FindByName 模糊查询目录名称
	FindByName(ctx context.Context, name string, catalogType int32) ([]*Catalog, error)

	// FindByTypeAndName 根据类型和精确名称查询
	FindByTypeAndName(ctx context.Context, catalogType int32, name string) (*Catalog, error)

	// FindByTypeAndParent 查询指定类型和父目录下的子目录
	FindByTypeAndParent(ctx context.Context, catalogType int32, parentId string) ([]*Catalog, error)

	// Update 更新目录
	Update(ctx context.Context, data *Catalog) error

	// Delete 删除目录（单条）
	Delete(ctx context.Context, id string) error

	// DeleteBatch 批量删除
	DeleteBatch(ctx context.Context, ids []string) error

	// WithTx 创建事务
	WithTx(tx interface{}) CatalogModel

	// Trans 事务执行
	Trans(ctx context.Context, fn func(ctx context.Context, model CatalogModel) error) error

	// FindAllDescendants 递归查找所有子孙目录（辅助方法）
	FindAllDescendants(ctx context.Context, id string) ([]*Catalog, error)

	// IsDescendant 判断 targetId 是否是 ancestorId 的子孙节点（辅助方法）
	IsDescendant(ctx context.Context, ancestorId, targetId string) (bool, error)

	// CheckUniqueName 检查同级目录下名称是否唯一（辅助方法）
	CheckUniqueName(ctx context.Context, parentId string, catalogType int32, name string, excludeId string) (bool, error)
}
