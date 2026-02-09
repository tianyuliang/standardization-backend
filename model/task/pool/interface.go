// Code scaffolded by speckit. Safe to edit.

package pool

import (
	"context"
)

// BusinessTablePoolModel 业务表标准创建池模型接口
type BusinessTablePoolModel interface {
	// Insert 插入记录
	Insert(ctx context.Context, data *BusinessTablePool) (int64, error)

	// FindOne 查询单条记录
	FindOne(ctx context.Context, id int64) (*BusinessTablePool, error)

	// Update 更新记录
	Update(ctx context.Context, data *BusinessTablePool) error

	// Delete 删除记录（逻辑删除）
	Delete(ctx context.Context, id int64) error

	// FindByStatus 按状态查询列表
	FindByStatus(ctx context.Context, status int32) ([]*BusinessTablePool, error)

	// FindByTableName 按表名查询
	FindByTableName(ctx context.Context, tableName string) (*BusinessTablePool, error)

	// FindByTableNameAndField 按表名和字段查询
	FindByTableNameAndField(ctx context.Context, tableName, tableField string) (*BusinessTablePool, error)

	// FindByCreateUserPhone 按创建人电话查询
	FindByCreateUserPhone(ctx context.Context, phone string, status int32) ([]*BusinessTablePool, error)

	// FindWithPagination 分页查询
	FindWithPagination(ctx context.Context, keyword string, status int32, page, pageSize int) ([]*BusinessTablePool, int64, error)

	// UpdateStatus 更新状态
	UpdateStatus(ctx context.Context, id int64, status int32) error

	// DeleteTaskId 删除任务关联（将f_task_id设置为NULL）
	DeleteTaskId(ctx context.Context, id int64) error

	// FindByBusinessTableFieldIds 按业务表字段ID列表查询
	FindByBusinessTableFieldIds(ctx context.Context, fieldIds []string) ([]*BusinessTablePool, error)

	// CountByTaskId 按任务ID统计记录数
	CountByTaskId(ctx context.Context, taskId string) (int64, error)

	// CountByTaskIdWithDataElementId 按任务ID统计已关联数据元的记录数
	CountByTaskIdWithDataElementId(ctx context.Context, taskId string) (int64, error)

	// UpdateDataElementId 更新数据元ID
	UpdateDataElementId(ctx context.Context, id int64, dataElementId int64) error

	// DeleteDataElementId 删除数据元ID（将f_data_element_id设置为NULL）
	DeleteDataElementId(ctx context.Context, id int64) error

	// FindOneByBusinessTableFieldId 按业务表字段ID查询
	FindOneByBusinessTableFieldId(ctx context.Context, businessTableFieldId string) (*BusinessTablePool, error)

	// FindByTableNamesAndStates 按表名列表和状态列表查询
	FindByTableNamesAndStates(ctx context.Context, tableNames []string, states []int32) ([]*BusinessTablePool, error)

	// FindByTaskId 按任务ID查询列表
	FindByTaskId(ctx context.Context, taskId string) ([]*BusinessTablePool, error)

	// UpdateBatchStatus 批量更新状态
	UpdateBatchStatus(ctx context.Context, ids []int64, status int32) error
}
