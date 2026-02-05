package relation_file

import (
	"context"
	"gorm.io/gorm"
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
}

type defaultRelationRuleFileModel struct {
	db *gorm.DB
}

// NewRelationRuleFileModel 创建关联规则文件模型实例
func NewRelationRuleFileModel(db *gorm.DB) RelationRuleFileModel {
	return &defaultRelationRuleFileModel{db: db}
}

func (m *defaultRelationRuleFileModel) getDB() *gorm.DB {
	return m.db.Table("t_relation_rule_file")
}

// InsertBatch 批量插入
func (m *defaultRelationRuleFileModel) InsertBatch(ctx context.Context, data []*RelationRuleFile) error {
	if len(data) == 0 {
		return nil
	}
	return m.getDB().WithContext(ctx).CreateInBatches(data, 100).Error
}

// DeleteByRuleId 删除规则的所有文件关联
func (m *defaultRelationRuleFileModel) DeleteByRuleId(ctx context.Context, ruleId int64) error {
	return m.getDB().WithContext(ctx).Where("f_rule_id = ?", ruleId).Delete(&RelationRuleFile{}).Error
}

// DeleteByFileId 删除文件的所有规则关联
func (m *defaultRelationRuleFileModel) DeleteByFileId(ctx context.Context, fileId int64) error {
	return m.getDB().WithContext(ctx).Where("f_file_id = ?", fileId).Delete(&RelationRuleFile{}).Error
}

// FindByRuleId 查询规则关联的文件
func (m *defaultRelationRuleFileModel) FindByRuleId(ctx context.Context, ruleId int64) ([]*RelationRuleFile, error) {
	var result []*RelationRuleFile
	err := m.getDB().WithContext(ctx).Where("f_rule_id = ?", ruleId).Find(&result).Error
	return result, err
}
