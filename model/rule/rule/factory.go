package rule

import (
	"github.com/jmoiron/sqlx"
)

// NewRuleModel 创建规则模型实例
func NewRuleModel(conn *sqlx.Conn) RuleModel {
	return &defaultRuleModel{conn: conn}
}
