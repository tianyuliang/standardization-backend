package rule

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type defaultRuleModel struct {
	conn *sqlx.Conn
}

// Insert 插入规则
func (m *defaultRuleModel) Insert(ctx context.Context, data *Rule) (int64, error) {
	query := `
		INSERT INTO t_rule (
			f_name, f_catalog_id, f_org_type, f_description, f_rule_type,
			f_version, f_expression, f_state, f_disable_reason, f_authority_id,
			f_department_ids, f_third_dept_id, f_create_time, f_create_user,
			f_update_time, f_update_user, f_deleted
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := m.conn.ExecContext(
		ctx,
		query,
		data.Name, data.CatalogId, data.OrgType, data.Description, data.RuleType,
		data.Version, data.Expression, data.State, data.DisableReason, data.AuthorityId,
		data.DepartmentIds, data.ThirdDeptId, data.CreateTime, data.CreateUser,
		data.UpdateTime, data.UpdateUser, data.Deleted,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// FindOne 查询单个规则
func (m *defaultRuleModel) FindOne(ctx context.Context, id int64) (*Rule, error) {
	query := `
		SELECT f_id, f_name, f_catalog_id, f_org_type, f_description, f_rule_type,
			f_version, f_expression, f_state, f_disable_reason, f_authority_id,
			f_department_ids, f_third_dept_id, f_create_time, f_create_user,
			f_update_time, f_update_user, f_deleted
		FROM t_rule
		WHERE f_id = ? AND f_deleted = 0
	`
	var rule Rule
	err := m.conn.QueryRowContext(ctx, query, id).Scan(
		&rule.Id, &rule.Name, &rule.CatalogId, &rule.OrgType, &rule.Description,
		&rule.RuleType, &rule.Version, &rule.Expression, &rule.State,
		&rule.DisableReason, &rule.AuthorityId, &rule.DepartmentIds, &rule.ThirdDeptId,
		&rule.CreateTime, &rule.CreateUser, &rule.UpdateTime, &rule.UpdateUser,
		&rule.Deleted,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &rule, err
}

// Update 更新规则
func (m *defaultRuleModel) Update(ctx context.Context, data *Rule) error {
	query := `
		UPDATE t_rule SET
			f_name = ?, f_catalog_id = ?, f_org_type = ?, f_description = ?,
			f_rule_type = ?, f_version = ?, f_expression = ?, f_state = ?,
			f_disable_reason = ?, f_authority_id = ?, f_department_ids = ?,
			f_third_dept_id = ?, f_update_time = ?, f_update_user = ?
		WHERE f_id = ?
	`
	_, err := m.conn.ExecContext(
		ctx,
		query,
		data.Name, data.CatalogId, data.OrgType, data.Description,
		data.RuleType, data.Version, data.Expression, data.State,
		data.DisableReason, data.AuthorityId, data.DepartmentIds,
		data.ThirdDeptId, data.UpdateTime, data.UpdateUser, data.Id,
	)
	return err
}

// Delete 删除规则（物理删除）
func (m *defaultRuleModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM t_rule WHERE f_id = ?`
	_, err := m.conn.ExecContext(ctx, query, id)
	return err
}

// FindByIds 批量查询规则
func (m *defaultRuleModel) FindByIds(ctx context.Context, ids []int64) ([]*Rule, error) {
	if len(ids) == 0 {
		return []*Rule{}, nil
	}
	query := fmt.Sprintf(`
		SELECT f_id, f_name, f_catalog_id, f_org_type, f_description, f_rule_type,
			f_version, f_expression, f_state, f_disable_reason, f_authority_id,
			f_department_ids, f_third_dept_id, f_create_time, f_create_user,
			f_update_time, f_update_user, f_deleted
		FROM t_rule
		WHERE f_id IN (%s) AND f_deleted = 0
	`, placeholders(len(ids)))

	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	rows, err := m.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []*Rule
	for rows.Next() {
		var rule Rule
		err := rows.Scan(
			&rule.Id, &rule.Name, &rule.CatalogId, &rule.OrgType, &rule.Description,
			&rule.RuleType, &rule.Version, &rule.Expression, &rule.State,
			&rule.DisableReason, &rule.AuthorityId, &rule.DepartmentIds, &rule.ThirdDeptId,
			&rule.CreateTime, &rule.CreateUser, &rule.UpdateTime, &rule.UpdateUser,
			&rule.Deleted,
		)
		if err != nil {
			return nil, err
		}
		rules = append(rules, &rule)
	}
	return rules, rows.Err()
}

// FindByNameAndOrgType 按名称和orgType查询
func (m *defaultRuleModel) FindByNameAndOrgType(ctx context.Context, name string, orgType int32, departmentIds string) ([]*Rule, error) {
	query := `
		SELECT f_id, f_name, f_catalog_id, f_org_type, f_description, f_rule_type,
			f_version, f_expression, f_state, f_disable_reason, f_authority_id,
			f_department_ids, f_third_dept_id, f_create_time, f_create_user,
			f_update_time, f_update_user, f_deleted
		FROM t_rule
		WHERE f_name = ? AND f_org_type = ? AND f_deleted = 0
	`
	args := []interface{}{name, orgType}

	if departmentIds != "" {
		query += ` AND f_department_ids = ?`
		args = append(args, departmentIds)
	}

	rows, err := m.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []*Rule
	for rows.Next() {
		var rule Rule
		err := rows.Scan(
			&rule.Id, &rule.Name, &rule.CatalogId, &rule.OrgType, &rule.Description,
			&rule.RuleType, &rule.Version, &rule.Expression, &rule.State,
			&rule.DisableReason, &rule.AuthorityId, &rule.DepartmentIds, &rule.ThirdDeptId,
			&rule.CreateTime, &rule.CreateUser, &rule.UpdateTime, &rule.UpdateUser,
			&rule.Deleted,
		)
		if err != nil {
			return nil, err
		}
		rules = append(rules, &rule)
	}
	return rules, rows.Err()
}

// FindByCatalogIds 按目录ID列表查询（分页）
func (m *defaultRuleModel) FindByCatalogIds(ctx context.Context, catalogIds []int64, opts *FindOptions) ([]*Rule, int64, error) {
	if len(catalogIds) == 0 {
		return []*Rule{}, 0, nil
	}

	whereClause := fmt.Sprintf("f_catalog_id IN (%s)", placeholders(len(catalogIds)))
	return m.findWithWhere(ctx, whereClause, catalogIds, opts)
}

// FindDataExists 检查数据是否存在
func (m *defaultRuleModel) FindDataExists(ctx context.Context, filterId int64, name string, departmentIds string) (*Rule, error) {
	query := `
		SELECT f_id, f_name, f_catalog_id, f_org_type, f_description, f_rule_type,
			f_version, f_expression, f_state, f_disable_reason, f_authority_id,
			f_department_ids, f_third_dept_id, f_create_time, f_create_user,
			f_update_time, f_update_user, f_deleted
		FROM t_rule
		WHERE f_name = ? AND f_deleted = 0
	`
	args := []interface{}{name}

	if filterId > 0 {
		query += ` AND f_id != ?`
		args = append(args, filterId)
	}

	if departmentIds != "" {
		query += ` AND f_department_ids = ?`
		args = append(args, departmentIds)
	}

	var rule Rule
	err := m.conn.QueryRowContext(ctx, query, args...).Scan(
		&rule.Id, &rule.Name, &rule.CatalogId, &rule.OrgType, &rule.Description,
		&rule.RuleType, &rule.Version, &rule.Expression, &rule.State,
		&rule.DisableReason, &rule.AuthorityId, &rule.DepartmentIds, &rule.ThirdDeptId,
		&rule.CreateTime, &rule.CreateUser, &rule.UpdateTime, &rule.UpdateUser,
		&rule.Deleted,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &rule, err
}

// UpdateState 更新规则状态
func (m *defaultRuleModel) UpdateState(ctx context.Context, id int64, state int32, reason string) error {
	query := `
		UPDATE t_rule
		SET f_state = ?, f_disable_reason = ?, f_update_time = ?
		WHERE f_id = ?
	`
	_, err := m.conn.ExecContext(ctx, query, state, reason, time.Now(), id)
	return err
}

// RemoveCatalog 批量移动目录
func (m *defaultRuleModel) RemoveCatalog(ctx context.Context, ids []int64, catalogId int64, updateUser string) error {
	if len(ids) == 0 {
		return nil
	}
	query := fmt.Sprintf(`
		UPDATE t_rule
		SET f_catalog_id = ?, f_version = f_version + 1, f_update_time = ?, f_update_user = ?
		WHERE f_id IN (%s)
	`, placeholders(len(ids)))
	args := []interface{}{catalogId, time.Now(), updateUser}
	for _, id := range ids {
		args = append(args, id)
	}
	_, err := m.conn.ExecContext(ctx, query, args...)
	return err
}

// UpdateVersionByIds 批量更新版本号
func (m *defaultRuleModel) UpdateVersionByIds(ctx context.Context, ids []int64, updateUser string) error {
	if len(ids) == 0 {
		return nil
	}
	query := fmt.Sprintf(`
		UPDATE t_rule
		SET f_version = f_version + 1, f_update_time = ?, f_update_user = ?
		WHERE f_id IN (%s)
	`, placeholders(len(ids)))
	args := []interface{}{time.Now(), updateUser}
	for _, id := range ids {
		args = append(args, id)
	}
	_, err := m.conn.ExecContext(ctx, query, args...)
	return err
}

// DeleteByIds 批量删除规则
func (m *defaultRuleModel) DeleteByIds(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	query := fmt.Sprintf(`DELETE FROM t_rule WHERE f_id IN (%s)`, placeholders(len(ids)))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	_, err := m.conn.ExecContext(ctx, query, args...)
	return err
}

// findWithWhere 通用查询方法
func (m *defaultRuleModel) findWithWhere(ctx context.Context, whereClause string, whereArgs []int64, opts *FindOptions) ([]*Rule, int64, error) {
	// 构建查询条件
	conditions := []string{whereClause}
	args := make([]interface{}, 0, len(whereArgs))

	for _, arg := range whereArgs {
		args = append(args, arg)
	}

	if opts.OrgType != nil {
		conditions = append(conditions, "f_org_type = ?")
		args = append(args, *opts.OrgType)
	}
	if opts.State != nil {
		conditions = append(conditions, "f_state = ?")
		args = append(args, *opts.State)
	}
	if opts.RuleType != nil {
		conditions = append(conditions, "f_rule_type = ?")
		args = append(args, *opts.RuleType)
	}
	if opts.Keyword != "" {
		conditions = append(conditions, "f_name LIKE ?")
		args = append(args, "%"+opts.Keyword+"%")
	}
	if opts.DepartmentId != "" {
		conditions = append(conditions, "f_department_ids = ?")
		args = append(args, opts.DepartmentId)
	}
	conditions = append(conditions, "f_deleted = 0")

	whereSQL := strings.Join(conditions, " AND ")

	// 查询总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM t_rule WHERE %s", whereSQL)
	var totalCount int64
	err := m.conn.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	sort := "f_create_time"
	if opts.Sort != "" {
		sort = "f_" + strings.TrimPrefix(opts.Sort, "f_")
	}
	direction := strings.ToUpper(opts.Direction)
	if direction != "ASC" && direction != "DESC" {
		direction = "DESC"
	}

	offset := (opts.Page - 1) * opts.PageSize
	query := fmt.Sprintf(`
		SELECT f_id, f_name, f_catalog_id, f_org_type, f_description, f_rule_type,
			f_version, f_expression, f_state, f_disable_reason, f_authority_id,
			f_department_ids, f_third_dept_id, f_create_time, f_create_user,
			f_update_time, f_update_user, f_deleted
		FROM t_rule
		WHERE %s
		ORDER BY %s %s
		LIMIT ? OFFSET ?
	`, whereSQL, sort, direction)

	args = append(args, opts.PageSize, offset)

	rows, err := m.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var rules []*Rule
	for rows.Next() {
		var rule Rule
		err := rows.Scan(
			&rule.Id, &rule.Name, &rule.CatalogId, &rule.OrgType, &rule.Description,
			&rule.RuleType, &rule.Version, &rule.Expression, &rule.State,
			&rule.DisableReason, &rule.AuthorityId, &rule.DepartmentIds, &rule.ThirdDeptId,
			&rule.CreateTime, &rule.CreateUser, &rule.UpdateTime, &rule.UpdateUser,
			&rule.Deleted,
		)
		if err != nil {
			return nil, 0, err
		}
		rules = append(rules, &rule)
	}
	return rules, totalCount, rows.Err()
}

// placeholders 生成占位符
func placeholders(n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat("?,", n-1) + "?"
}
