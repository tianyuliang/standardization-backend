// Code scaffolded by speckit. Safe to edit.

package dataelement

import (
	"context"
)

// DataElementModel 数据元模型接口
type DataElementModel interface {
	// Insert 插入数据元
	Insert(ctx context.Context, data *DataElement) (int64, error)

	// FindOne 查询单个数据元
	FindOne(ctx context.Context, id int64) (*DataElement, error)

	// FindOneByCode 按Code查询数据元
	FindOneByCode(ctx context.Context, code int64) (*DataElement, error)

	// Update 更新数据元
	Update(ctx context.Context, data *DataElement) error

	// Delete 删除数据元
	Delete(ctx context.Context, id int64) error

	// FindByCatalogIds 按目录ID列表查询数据元
	FindByCatalogIds(ctx context.Context, catalogIds []int64, opts *FindOptions) ([]*DataElement, int64, error)

	// FindByIds 按ID列表查询数据元
	FindByIds(ctx context.Context, ids []int64) ([]*DataElement, error)

	// FindByRuleId 按规则ID分页查询数据元
	FindByRuleId(ctx context.Context, ruleId int64, opts *FindOptions) ([]*DataElement, int64, error)

	// FindByCodes 按Code列表查询数据元
	FindByCodes(ctx context.Context, codes []int64) ([]*DataElement, error)

	// FindByFileCatalog 按文件目录查询数据元
	FindByFileCatalog(ctx context.Context, opts *FindOptions) ([]*DataElement, int64, error)

	// FindByFileId 按文件ID查询数据元
	FindByFileId(ctx context.Context, fileId int64, opts *FindOptions) ([]*DataElement, int64, error)

	// CheckNameCnExists 检查中文名称是否存在
	CheckNameCnExists(ctx context.Context, nameCn string, stdType int32, excludeId int64, deptIds string) (bool, error)

	// CheckNameEnExists 检查英文名称是否存在
	CheckNameEnExists(ctx context.Context, nameEn string, excludeId int64, deptIds string) (bool, error)

	// FindDataExists 查找数据是否存在
	FindDataExists(ctx context.Context, name string, stdType int32, excludeId int64, deptIds string) (*DataElement, error)

	// UpdateState 更新状态
	UpdateState(ctx context.Context, ids []int64, state int32, reason, updateUser string) error

	// MoveCatalog 移动目录
	MoveCatalog(ctx context.Context, ids []int64, catalogId int64, updateUser string) error

	// DeleteLabelIds 删除标签
	DeleteLabelIds(ctx context.Context, ids []int64) error

	// DeleteByIds 批量删除
	DeleteByIds(ctx context.Context, ids []int64) error

	// CountByCatalogId 按目录ID统计数据元数量
	CountByCatalogId(ctx context.Context, catalogId int64) (int64, error)

	// IncrementVersion 递增版本号
	IncrementVersion(ctx context.Context, id int64) error
}
