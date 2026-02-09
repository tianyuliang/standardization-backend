// Code scaffolded by speckit. Safe to edit.

package task

import (
	"context"
)

// TaskStdCreateModel 标准创建任务模型接口
type TaskStdCreateModel interface {
	// Insert 插入任务
	Insert(ctx context.Context, data *TaskStdCreate) (int64, error)

	// FindOne 查询单个任务
	FindOne(ctx context.Context, id int64) (*TaskStdCreate, error)

	// Update 更新任务
	Update(ctx context.Context, data *TaskStdCreate) error

	// Delete 删除任务
	Delete(ctx context.Context, id int64) error

	// FindUncompleted 查询未处理任务列表
	FindUncompleted(ctx context.Context, keyword string, page, pageSize int) ([]*TaskStdCreate, int64, error)

	// FindCompleted 查询已完成任务列表
	FindCompleted(ctx context.Context, keyword string, page, pageSize int) ([]*TaskStdCreate, int64, error)

	// FindByStatus 按状态查询任务
	FindByStatus(ctx context.Context, status int32) ([]*TaskStdCreate, error)

	// FindByTaskNo 按任务编号查询
	FindByTaskNo(ctx context.Context, taskNo string) (*TaskStdCreate, error)

	// CountByStatus 统计各状态任务数量
	CountByStatus(ctx context.Context, status int32) (int64, error)
}

// TaskStdCreateResultModel 任务结果模型接口
type TaskStdCreateResultModel interface {
	// Insert 插入结果
	Insert(ctx context.Context, data *TaskStdCreateResult) (int64, error)

	// FindByTaskId 按任务ID查询结果列表
	FindByTaskId(ctx context.Context, taskId int64) ([]*TaskStdCreateResult, error)

	// DeleteByTaskId 删除任务的所有结果
	DeleteByTaskId(ctx context.Context, taskId int64) error

	// InsertBatch 批量插入结果
	InsertBatch(ctx context.Context, data []*TaskStdCreateResult) error
}

// FindTaskOptions 查询选项
type FindTaskOptions struct {
	Keyword  string
	Page     int
	PageSize int
}
