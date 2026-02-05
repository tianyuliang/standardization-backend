package relation_file

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type defaultRelationRuleFileModel struct {
	conn *sqlx.Conn
}

// InsertBatch 批量插入
func (m *defaultRelationRuleFileModel) InsertBatch(ctx context.Context, data []*RelationRuleFile) error {
	if len(data) == 0 {
		return nil
	}

	query := `
		INSERT INTO t_relation_rule_file (f_id, f_rule_id, f_file_id)
		VALUES `

	placeholders := make([]string, len(data))
	args := make([]interface{}, 0, len(data)*3)

	for i, item := range data {
		placeholders[i] = "(?, ?, ?)"
		args = append(args, item.Id, item.RuleId, item.FileId)
	}

	query += strings.Join(placeholders, ", ")

	_, err := m.conn.ExecContext(ctx, query, args...)
	return err
}

// DeleteByRuleId 删除规则的所有文件关联
func (m *defaultRelationRuleFileModel) DeleteByRuleId(ctx context.Context, ruleId int64) error {
	query := `DELETE FROM t_relation_rule_file WHERE f_rule_id = ?`
	_, err := m.conn.ExecContext(ctx, query, ruleId)
	return err
}

// DeleteByFileId 删除文件的所有规则关联
func (m *defaultRelationRuleFileModel) DeleteByFileId(ctx context.Context, fileId int64) error {
	query := `DELETE FROM t_relation_rule_file WHERE f_file_id = ?`
	_, err := m.conn.ExecContext(ctx, query, fileId)
	return err
}

// FindByRuleId 查询规则关联的文件
func (m *defaultRelationRuleFileModel) FindByRuleId(ctx context.Context, ruleId int64) ([]*RelationRuleFile, error) {
	query := `
		SELECT f_id, f_rule_id, f_file_id
		FROM t_relation_rule_file
		WHERE f_rule_id = ?
	`

	rows, err := m.conn.QueryContext(ctx, query, ruleId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*RelationRuleFile
	for rows.Next() {
		var item RelationRuleFile
		err := rows.Scan(&item.Id, &item.RuleId, &item.FileId)
		if err != nil {
			return nil, err
		}
		result = append(result, &item)
	}
	return result, rows.Err()
}

// DeleteByRuleIds 批量删除规则的文件关联
func (m *defaultRelationRuleFileModel) DeleteByRuleIds(ctx context.Context, ruleIds []int64) error {
	if len(ruleIds) == 0 {
		return nil
	}

	query := fmt.Sprintf(`DELETE FROM t_relation_rule_file WHERE f_rule_id IN (%s)`,
		strings.Repeat("?,", len(ruleIds)-1)+"?")

	args := make([]interface{}, len(ruleIds))
	for i, id := range ruleIds {
		args[i] = id
	}

	_, err := m.conn.ExecContext(ctx, query, args...)
	return err
}
