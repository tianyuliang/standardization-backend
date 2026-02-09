// Code scaffolded by speckit. Safe to edit.

package pool

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

var (
	// DefBusinessTablePoolModel 默认业务表池模型
	DefBusinessTablePoolModel BusinessTablePoolModel
)

// NewBusinessTablePoolModel 创建业务表池模型
func NewBusinessTablePoolModel(conn *sqlx.Conn) BusinessTablePoolModel {
	if DefBusinessTablePoolModel != nil {
		return DefBusinessTablePoolModel
	}
	return &defaultBusinessTablePoolModel{
		conn: conn,
	}
}

// SetBusinessTablePoolModel 设置业务表池模型
func SetBusinessTablePoolModel(model BusinessTablePoolModel) {
	DefBusinessTablePoolModel = model
}

type defaultBusinessTablePoolModel struct {
	conn *sqlx.Conn
}

// Transaction 事务辅助函数
func (m *defaultBusinessTablePoolModel) Transact(ctx context.Context, fn func(session *sql.Tx) error) error {
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
