package dict

import (
	"github.com/jmoiron/sqlx"
)

var (
	// DefDictModel 默认的 DictModel 实现
	DefDictModel DictModel
	// DefDictEnumModel 默认的 DictEnumModel 实现
	DefDictEnumModel DictEnumModel
	// DefRelationDictFileModel 默认的 RelationDictFileModel 实现
	DefRelationDictFileModel RelationDictFileModel
)

// NewDictModel 创建 DictModel
func NewDictModel(conn *sqlx.Conn) DictModel {
	if DefDictModel != nil {
		return DefDictModel
	}
	return &defaultDictModel{
		conn: conn,
	}
}

// NewDictEnumModel 创建 DictEnumModel
func NewDictEnumModel(conn *sqlx.Conn) DictEnumModel {
	if DefDictEnumModel != nil {
		return DefDictEnumModel
	}
	return &defaultDictEnumModel{
		conn: conn,
	}
}

// NewRelationDictFileModel 创建 RelationDictFileModel
func NewRelationDictFileModel(conn *sqlx.Conn) RelationDictFileModel {
	if DefRelationDictFileModel != nil {
		return DefRelationDictFileModel
	}
	return &defaultRelationDictFileModel{
		conn: conn,
	}
}
