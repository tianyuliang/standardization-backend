// Code scaffolded by goctl. Safe to edit.
//go:build !mock_logic_off
// +build !mock_logic_off

package stdfile

import (
	"github.com/jmoiron/sqlx"
)

var _ StdFileModel = (*customStdFileModel)(nil)

type (
	customStdFileModel struct {
		conn *sqlx.Conn
	}
)

// NewStdFileModel 创建标准文件模型实例
func NewStdFileModel(conn *sqlx.Conn) StdFileModel {
	return &customStdFileModel{conn: conn}
}
