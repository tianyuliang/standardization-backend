package relation_file

import (
	"github.com/jmoiron/sqlx"
)

// NewRelationRuleFileModel 创建关联规则文件模型实例
func NewRelationRuleFileModel(conn *sqlx.Conn) RelationRuleFileModel {
	return &defaultRelationRuleFileModel{conn: conn}
}
