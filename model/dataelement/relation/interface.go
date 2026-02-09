// Code scaffolded by speckit. Safe to edit.

package relation

import (
	"context"
)

// RelationDeFileModel 数据元-文件关系模型接口
type RelationDeFileModel interface {
	// InsertBatch 批量插入关系
	InsertBatch(ctx context.Context, data []*RelationDeFile) error

	// DeleteByDeId 删除数据元的所有关系
	DeleteByDeId(ctx context.Context, deId int64) error

	// DeleteByDeIds 批量删除数据元的关系
	DeleteByDeIds(ctx context.Context, deIds []int64) error

	// DeleteByFileId 删除文件的所有关系
	DeleteByFileId(ctx context.Context, fileId int64) error

	// FindByDeId 查询数据元的关联文件列表
	FindByDeId(ctx context.Context, deId int64, opts *PageOptions) ([]*RelationDeFile, int64, error)

	// FindByDeIds 批量查询数据元的关联文件
	FindByDeIds(ctx context.Context, deIds []int64) (map[int64][]int64, error)

	// CountByDeId 统计数据元的关联文件数量
	CountByDeId(ctx context.Context, deId int64) (int64, error)

	// FindFileIdsByDeId 查询数据元关联的文件ID列表
	FindFileIdsByDeId(ctx context.Context, deId int64) ([]int64, error)
}
