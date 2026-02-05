package relation_file

import (
	"context"
)

// RelationRuleFileModel 关联规则文件模型接口
type RelationRuleFileModel interface {
	// InsertBatch 批量插入
	InsertBatch(ctx context.Context, data []*RelationRuleFile) error

	// DeleteByRuleId 删除规则的所有文件关联
	DeleteByRuleId(ctx context.Context, ruleId int64) error

	// DeleteByFileId 删除文件的所有规则关联
	DeleteByFileId(ctx context.Context, fileId int64) error

	// FindByRuleId 查询规则关联的文件
	FindByRuleId(ctx context.Context, ruleId int64) ([]*RelationRuleFile, error)

	// DeleteByRuleIds 批量删除规则的文件关联
	DeleteByRuleIds(ctx context.Context, ruleIds []int64) error
}

// RelationRuleFile 规则-文件关联表
type RelationRuleFile struct {
	Id     int64 `db:"f_id" json:"id"`
	RuleId int64 `db:"f_rule_id" json:"ruleId"`
	FileId int64 `db:"f_file_id" json:"fileId"`
}
