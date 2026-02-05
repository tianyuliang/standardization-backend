package rule

import (
	"github.com/jinguoxing/idrm-go-base/db"
	"gorm.io/gorm"
)

// NewRuleModel 创建规则模型实例
func NewRuleModel(db *gorm.DB) RuleModel {
	return &defaultRuleModel{db: db}
}

type defaultRuleModel struct {
	db *gorm.DB
}

func (m *defaultRuleModel) getDB() *gorm.DB {
	return m.db.Table("t_rule")
}
