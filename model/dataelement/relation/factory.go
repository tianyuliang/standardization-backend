// Code scaffolded by speckit. Safe to edit.

package relation

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

var (
	// DefRelationDeFileModel 默认关系模型
	DefRelationDeFileModel RelationDeFileModel
)

// NewRelationDeFileModel 创建关系模型
func NewRelationDeFileModel(conn *sqlx.Conn) RelationDeFileModel {
	if DefRelationDeFileModel != nil {
		return DefRelationDeFileModel
	}
	return &defaultRelationDeFileModel{
		conn: conn,
	}
}

// SetRelationDeFileModel 设置关系模型
func SetRelationDeFileModel(model RelationDeFileModel) {
	DefRelationDeFileModel = model
}

type defaultRelationDeFileModel struct {
	conn *sqlx.Conn
}

// Transact 事务辅助函数
func (m *defaultRelationDeFileModel) Transact(ctx context.Context, fn func(session *sql.Tx) error) error {
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
