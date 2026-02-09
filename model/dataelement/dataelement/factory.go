// Code scaffolded by speckit. Safe to edit.

package dataelement

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

var (
	// DefDataElementModel 默认数据元模型
	DefDataElementModel DataElementModel
)

// NewDataElementModel 创建数据元模型
func NewDataElementModel(conn *sqlx.Conn) DataElementModel {
	if DefDataElementModel != nil {
		return DefDataElementModel
	}
	return &defaultDataElementModel{
		conn: conn,
	}
}

// SetDataElementModel 设置数据元模型
func SetDataElementModel(model DataElementModel) {
	DefDataElementModel = model
}

type defaultDataElementModel struct {
	conn *sqlx.Conn
}

// Transact 事务辅助函数
func (m *defaultDataElementModel) Transact(ctx context.Context, fn func(session *sql.Tx) error) error {
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
