// Code scaffolded by speckit. Safe to edit.

package task

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

var (
	// DefTaskStdCreateModel 默认任务模型
	DefTaskStdCreateModel TaskStdCreateModel
	// DefTaskStdCreateResultModel 默认任务结果模型
	DefTaskStdCreateResultModel TaskStdCreateResultModel
)

// NewTaskStdCreateModel 创建任务模型
func NewTaskStdCreateModel(conn *sqlx.Conn) TaskStdCreateModel {
	if DefTaskStdCreateModel != nil {
		return DefTaskStdCreateModel
	}
	return &defaultTaskStdCreateModel{
		conn: conn,
	}
}

// NewTaskStdCreateResultModel 创建任务结果模型
func NewTaskStdCreateResultModel(conn *sqlx.Conn) TaskStdCreateResultModel {
	if DefTaskStdCreateResultModel != nil {
		return DefTaskStdCreateResultModel
	}
	return &defaultTaskStdCreateResultModel{
		conn: conn,
	}
}

// SetTaskStdCreateModel 设置任务模型
func SetTaskStdCreateModel(model TaskStdCreateModel) {
	DefTaskStdCreateModel = model
}

// SetTaskStdCreateResultModel 设置任务结果模型
func SetTaskStdCreateResultModel(model TaskStdCreateResultModel) {
	DefTaskStdCreateResultModel = model
}

type defaultTaskStdCreateModel struct {
	conn *sqlx.Conn
}

type defaultTaskStdCreateResultModel struct {
	conn *sqlx.Conn
}

// Transaction 事务辅助函数
func (m *defaultTaskStdCreateModel) Transact(ctx context.Context, fn func(session *sql.Tx) error) error {
	tx, err := m.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}
