// Code scaffolded by speckit. Safe to edit.

package pool

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
)

var _ BusinessTablePoolModel = (*defaultBusinessTablePoolModel)(nil)

// Insert 插入记录
func (m *defaultBusinessTablePoolModel) Insert(ctx context.Context, data *BusinessTablePool) (int64, error) {
	query := `INSERT INTO ` + TableNameBusinessTablePool + `
		(f_table_name, f_table_description, f_table_field, f_field_description,
		 f_data_type, f_status, f_create_user, f_create_user_phone, f_task_id, f_deleted)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 0)`

	result, err := m.conn.ExecContext(ctx, query,
		data.TableName, data.TableDescription, data.TableField, data.FieldDescription,
		data.DataType, data.Status, data.CreateUser, data.CreateUserPhone, data.TaskId)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// FindOne 查询单条记录
func (m *defaultBusinessTablePoolModel) FindOne(ctx context.Context, id int64) (*BusinessTablePool, error) {
	query := `SELECT f_id, f_table_name, f_table_description, f_table_field, f_field_description,
		f_data_type, f_status, f_create_user, f_create_user_phone, f_task_id, f_data_element_id,
		f_create_time, f_update_time, f_deleted
		FROM ` + TableNameBusinessTablePool + `
		WHERE f_id = ? AND f_deleted = 0`

	var pool BusinessTablePool
	err := m.conn.QueryRowContext(ctx, query, id).Scan(
		&pool.Id, &pool.TableName, &pool.TableDescription, &pool.TableField,
		&pool.FieldDescription, &pool.DataType, &pool.Status, &pool.CreateUser,
		&pool.CreateUserPhone, &pool.TaskId, &pool.DataElementId,
		&pool.CreateTime, &pool.UpdateTime, &pool.Deleted,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &pool, nil
}

// Update 更新记录
func (m *defaultBusinessTablePoolModel) Update(ctx context.Context, data *BusinessTablePool) error {
	query := `UPDATE ` + TableNameBusinessTablePool + `
		SET f_table_name = ?, f_table_description = ?, f_table_field = ?,
		 f_field_description = ?, f_data_type = ?, f_status = ?, f_task_id = ?, f_data_element_id = ?, f_update_time = NOW()
		WHERE f_id = ? AND f_deleted = 0`

	_, err := m.conn.ExecContext(ctx, query,
		data.TableName, data.TableDescription, data.TableField,
		data.FieldDescription, data.DataType, data.Status, data.TaskId, data.DataElementId, data.Id)

	return err
}

// Delete 删除记录（逻辑删除）
func (m *defaultBusinessTablePoolModel) Delete(ctx context.Context, id int64) error {
	query := `UPDATE ` + TableNameBusinessTablePool + `
		SET f_deleted = 1, f_update_time = NOW()
		WHERE f_id = ?`

	_, err := m.conn.ExecContext(ctx, query, id)
	return err
}

// FindByStatus 按状态查询列表
func (m *defaultBusinessTablePoolModel) FindByStatus(ctx context.Context, status int32) ([]*BusinessTablePool, error) {
	query := `SELECT f_id, f_table_name, f_table_description, f_table_field, f_field_description,
		f_data_type, f_status, f_create_user, f_create_user_phone, f_task_id, f_data_element_id,
		f_create_time, f_update_time, f_deleted
		FROM ` + TableNameBusinessTablePool + `
		WHERE f_status = ? AND f_deleted = 0
		ORDER BY f_create_time DESC`

	rows, err := m.conn.QueryContext(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pools []*BusinessTablePool
	for rows.Next() {
		var pool BusinessTablePool
		err := rows.Scan(
			&pool.Id, &pool.TableName, &pool.TableDescription, &pool.TableField,
			&pool.FieldDescription, &pool.DataType, &pool.Status, &pool.CreateUser,
			&pool.CreateUserPhone, &pool.TaskId, &pool.DataElementId,
			&pool.CreateTime, &pool.UpdateTime, &pool.Deleted,
		)
		if err != nil {
			return nil, err
		}
		pools = append(pools, &pool)
	}

	return pools, nil
}

// FindByTableName 按表名查询
func (m *defaultBusinessTablePoolModel) FindByTableName(ctx context.Context, tableName string) (*BusinessTablePool, error) {
	query := `SELECT f_id, f_table_name, f_table_description, f_table_field, f_field_description,
		f_data_type, f_status, f_create_user, f_create_user_phone, f_task_id, f_data_element_id,
		f_create_time, f_update_time, f_deleted
		FROM ` + TableNameBusinessTablePool + `
		WHERE f_table_name = ? AND f_deleted = 0
		ORDER BY f_create_time DESC
		LIMIT 1`

	var pool BusinessTablePool
	err := m.conn.QueryRowContext(ctx, query, tableName).Scan(
		&pool.Id, &pool.TableName, &pool.TableDescription, &pool.TableField,
		&pool.FieldDescription, &pool.DataType, &pool.Status, &pool.CreateUser,
		&pool.CreateUserPhone, &pool.TaskId, &pool.DataElementId,
		&pool.CreateTime, &pool.UpdateTime, &pool.Deleted,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &pool, nil
}

// FindByTableNameAndField 按表名和字段查询
func (m *defaultBusinessTablePoolModel) FindByTableNameAndField(ctx context.Context, tableName, tableField string) (*BusinessTablePool, error) {
	query := `SELECT f_id, f_table_name, f_table_description, f_table_field, f_field_description,
		f_data_type, f_status, f_create_user, f_create_user_phone, f_task_id, f_data_element_id,
		f_create_time, f_update_time, f_deleted
		FROM ` + TableNameBusinessTablePool + `
		WHERE f_table_name = ? AND f_table_field = ? AND f_deleted = 0
		LIMIT 1`

	var pool BusinessTablePool
	err := m.conn.QueryRowContext(ctx, query, tableName, tableField).Scan(
		&pool.Id, &pool.TableName, &pool.TableDescription, &pool.TableField,
		&pool.FieldDescription, &pool.DataType, &pool.Status, &pool.CreateUser,
		&pool.CreateUserPhone, &pool.TaskId, &pool.DataElementId,
		&pool.CreateTime, &pool.UpdateTime, &pool.Deleted,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &pool, nil
}

// FindByCreateUserPhone 按创建人电话查询
func (m *defaultBusinessTablePoolModel) FindByCreateUserPhone(ctx context.Context, phone string, status int32) ([]*BusinessTablePool, error) {
	query := `SELECT f_id, f_table_name, f_table_description, f_table_field, f_field_description,
		f_data_type, f_status, f_create_user, f_create_user_phone, f_task_id, f_data_element_id,
		f_create_time, f_update_time, f_deleted
		FROM ` + TableNameBusinessTablePool + `
		WHERE f_create_user_phone = ? AND f_status = ? AND f_deleted = 0
		ORDER BY f_create_time DESC`

	rows, err := m.conn.QueryContext(ctx, query, phone, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pools []*BusinessTablePool
	for rows.Next() {
		var pool BusinessTablePool
		err := rows.Scan(
			&pool.Id, &pool.TableName, &pool.TableDescription, &pool.TableField,
			&pool.FieldDescription, &pool.DataType, &pool.Status, &pool.CreateUser,
			&pool.CreateUserPhone, &pool.TaskId, &pool.DataElementId,
			&pool.CreateTime, &pool.UpdateTime, &pool.Deleted,
		)
		if err != nil {
			return nil, err
		}
		pools = append(pools, &pool)
	}

	return pools, nil
}

// FindWithPagination 分页查询
func (m *defaultBusinessTablePoolModel) FindWithPagination(ctx context.Context, keyword string, status int32, page, pageSize int) ([]*BusinessTablePool, int64, error) {
	// 构建查询条件
	whereClause := "WHERE f_deleted = 0"
	args := []interface{}{}

	if status >= 0 {
		whereClause += " AND f_status = ?"
		args = append(args, status)
	}

	if keyword != "" {
		whereClause += " AND (f_table_name LIKE ? OR f_table_field LIKE ?)"
		keywordPattern := "%" + keyword + "%"
		args = append(args, keywordPattern, keywordPattern)
	}

	// 查询总数
	countQuery := "SELECT COUNT(*) FROM " + TableNameBusinessTablePool + " " + whereClause
	var totalCount int64
	err := m.conn.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	offset := (page - 1) * pageSize
	query := `SELECT f_id, f_table_name, f_table_description, f_table_field, f_field_description,
		f_data_type, f_status, f_create_user, f_create_user_phone, f_task_id, f_data_element_id,
		f_create_time, f_update_time, f_deleted
		FROM ` + TableNameBusinessTablePool + " " + whereClause + `
		ORDER BY f_create_time DESC
		LIMIT ? OFFSET ?`

	args = append(args, pageSize, offset)
	rows, err := m.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var pools []*BusinessTablePool
	for rows.Next() {
		var pool BusinessTablePool
		err := rows.Scan(
			&pool.Id, &pool.TableName, &pool.TableDescription, &pool.TableField,
			&pool.FieldDescription, &pool.DataType, &pool.Status, &pool.CreateUser,
			&pool.CreateUserPhone, &pool.TaskId, &pool.DataElementId,
			&pool.CreateTime, &pool.UpdateTime, &pool.Deleted,
		)
		if err != nil {
			return nil, 0, err
		}
		pools = append(pools, &pool)
	}

	return pools, totalCount, nil
}

// UpdateStatus 更新状态
func (m *defaultBusinessTablePoolModel) UpdateStatus(ctx context.Context, id int64, status int32) error {
	query := `UPDATE ` + TableNameBusinessTablePool + `
		SET f_status = ?, f_update_time = NOW()
		WHERE f_id = ? AND f_deleted = 0`

	_, err := m.conn.ExecContext(ctx, query, status, id)
	return err
}

// DeleteTaskId 删除任务关联（将f_task_id设置为NULL）
func (m *defaultBusinessTablePoolModel) DeleteTaskId(ctx context.Context, id int64) error {
	query := `UPDATE ` + TableNameBusinessTablePool + `
		SET f_task_id = NULL, f_update_time = NOW()
		WHERE f_id = ? AND f_deleted = 0`

	_, err := m.conn.ExecContext(ctx, query, id)
	return err
}

// FindByBusinessTableFieldIds 按业务表字段ID列表查询
func (m *defaultBusinessTablePoolModel) FindByBusinessTableFieldIds(ctx context.Context, fieldIds []string) ([]*BusinessTablePool, error) {
	if len(fieldIds) == 0 {
		return []*BusinessTablePool{}, nil
	}

	// 构建IN查询 - 使用f_id作为业务表字段ID（因为当前DDL没有f_business_table_field_id字段）
	placeholders := make([]string, len(fieldIds))
	args := make([]interface{}, len(fieldIds))
	for i, id := range fieldIds {
		placeholders[i] = "?"
		args[i] = id
	}

	query := `SELECT f_id, f_table_name, f_table_description, f_table_field, f_field_description,
		f_data_type, f_status, f_create_user, f_create_user_phone, f_task_id, f_data_element_id,
		f_create_time, f_update_time, f_deleted
		FROM ` + TableNameBusinessTablePool + `
		WHERE f_id IN (` + strings.Join(placeholders, ",") + `) AND f_deleted = 0`

	rows, err := m.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pools []*BusinessTablePool
	for rows.Next() {
		var pool BusinessTablePool
		err := rows.Scan(
			&pool.Id, &pool.TableName, &pool.TableDescription, &pool.TableField,
			&pool.FieldDescription, &pool.DataType, &pool.Status, &pool.CreateUser,
			&pool.CreateUserPhone, &pool.TaskId, &pool.DataElementId,
			&pool.CreateTime, &pool.UpdateTime, &pool.Deleted,
		)
		if err != nil {
			return nil, err
		}
		pools = append(pools, &pool)
	}

	return pools, nil
}

// CountByTaskId 按任务ID统计记录数
func (m *defaultBusinessTablePoolModel) CountByTaskId(ctx context.Context, taskId string) (int64, error) {
	query := `SELECT COUNT(*) FROM ` + TableNameBusinessTablePool + `
		WHERE f_task_id = ? AND f_deleted = 0`

	var count int64
	err := m.conn.QueryRowContext(ctx, query, taskId).Scan(&count)
	return count, err
}

// CountByTaskIdWithDataElementId 按任务ID统计已关联数据元的记录数
func (m *defaultBusinessTablePoolModel) CountByTaskIdWithDataElementId(ctx context.Context, taskId string) (int64, error) {
	query := `SELECT COUNT(*) FROM ` + TableNameBusinessTablePool + `
		WHERE f_task_id = ? AND f_data_element_id IS NOT NULL AND f_deleted = 0`

	var count int64
	err := m.conn.QueryRowContext(ctx, query, taskId).Scan(&count)
	return count, err
}

// UpdateDataElementId 更新数据元ID
func (m *defaultBusinessTablePoolModel) UpdateDataElementId(ctx context.Context, id int64, dataElementId int64) error {
	query := `UPDATE ` + TableNameBusinessTablePool + `
		SET f_data_element_id = ?, f_update_time = NOW()
		WHERE f_id = ? AND f_deleted = 0`

	_, err := m.conn.ExecContext(ctx, query, dataElementId, id)
	return err
}

// DeleteDataElementId 删除数据元ID（将f_data_element_id设置为NULL）
func (m *defaultBusinessTablePoolModel) DeleteDataElementId(ctx context.Context, id int64) error {
	query := `UPDATE ` + TableNameBusinessTablePool + `
		SET f_data_element_id = NULL, f_update_time = NOW()
		WHERE f_id = ? AND f_deleted = 0`

	_, err := m.conn.ExecContext(ctx, query, id)
	return err
}

// FindOneByBusinessTableFieldId 按业务表字段ID查询（使用f_id）
func (m *defaultBusinessTablePoolModel) FindOneByBusinessTableFieldId(ctx context.Context, businessTableFieldId string) (*BusinessTablePool, error) {
	// 将string ID转换为int64
	id, err := strconv.ParseInt(businessTableFieldId, 10, 64)
	if err != nil {
		return nil, err
	}

	query := `SELECT f_id, f_table_name, f_table_description, f_table_field, f_field_description,
		f_data_type, f_status, f_create_user, f_create_user_phone, f_task_id, f_data_element_id,
		f_create_time, f_update_time, f_deleted
		FROM ` + TableNameBusinessTablePool + `
		WHERE f_id = ? AND f_deleted = 0`

	var pool BusinessTablePool
	err = m.conn.QueryRowContext(ctx, query, id).Scan(
		&pool.Id, &pool.TableName, &pool.TableDescription, &pool.TableField,
		&pool.FieldDescription, &pool.DataType, &pool.Status, &pool.CreateUser,
		&pool.CreateUserPhone, &pool.TaskId, &pool.DataElementId,
		&pool.CreateTime, &pool.UpdateTime, &pool.Deleted,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &pool, nil
}

// FindByTableNamesAndStates 按表名列表和状态列表查询
func (m *defaultBusinessTablePoolModel) FindByTableNamesAndStates(ctx context.Context, tableNames []string, states []int32) ([]*BusinessTablePool, error) {
	// 构建查询条件
	whereClause := "WHERE f_deleted = 0"
	args := []interface{}{}

	if len(tableNames) > 0 {
		placeholders := make([]string, len(tableNames))
		for i, name := range tableNames {
			placeholders[i] = "?"
			args = append(args, name)
		}
		whereClause += " AND f_table_name IN (" + strings.Join(placeholders, ",") + ")"
	}

	if len(states) > 0 {
		statePlaceholders := make([]string, len(states))
		for i, state := range states {
			statePlaceholders[i] = "?"
			args = append(args, state)
		}
		whereClause += " AND f_status IN (" + strings.Join(statePlaceholders, ",") + ")"
	}

	query := `SELECT f_id, f_table_name, f_table_description, f_table_field, f_field_description,
		f_data_type, f_status, f_create_user, f_create_user_phone, f_task_id, f_data_element_id,
		f_create_time, f_update_time, f_deleted
		FROM ` + TableNameBusinessTablePool + " " + whereClause + `
		ORDER BY f_create_time DESC`

	rows, err := m.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pools []*BusinessTablePool
	for rows.Next() {
		var pool BusinessTablePool
		err := rows.Scan(
			&pool.Id, &pool.TableName, &pool.TableDescription, &pool.TableField,
			&pool.FieldDescription, &pool.DataType, &pool.Status, &pool.CreateUser,
			&pool.CreateUserPhone, &pool.TaskId, &pool.DataElementId,
			&pool.CreateTime, &pool.UpdateTime, &pool.Deleted,
		)
		if err != nil {
			return nil, err
		}
		pools = append(pools, &pool)
	}

	return pools, nil
}

// FindByTaskId 按任务ID查询列表
func (m *defaultBusinessTablePoolModel) FindByTaskId(ctx context.Context, taskId string) ([]*BusinessTablePool, error) {
	query := `SELECT f_id, f_table_name, f_table_description, f_table_field, f_field_description,
		f_data_type, f_status, f_create_user, f_create_user_phone, f_task_id, f_data_element_id,
		f_create_time, f_update_time, f_deleted
		FROM ` + TableNameBusinessTablePool + `
		WHERE f_task_id = ? AND f_deleted = 0
		ORDER BY f_create_time DESC`

	rows, err := m.conn.QueryContext(ctx, query, taskId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pools []*BusinessTablePool
	for rows.Next() {
		var pool BusinessTablePool
		err := rows.Scan(
			&pool.Id, &pool.TableName, &pool.TableDescription, &pool.TableField,
			&pool.FieldDescription, &pool.DataType, &pool.Status, &pool.CreateUser,
			&pool.CreateUserPhone, &pool.TaskId, &pool.DataElementId,
			&pool.CreateTime, &pool.UpdateTime, &pool.Deleted,
		)
		if err != nil {
			return nil, err
		}
		pools = append(pools, &pool)
	}

	return pools, nil
}

// UpdateBatchStatus 批量更新状态
func (m *defaultBusinessTablePoolModel) UpdateBatchStatus(ctx context.Context, ids []int64, status int32) error {
	if len(ids) == 0 {
		return nil
	}

	// 构建IN查询
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids)+1)
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	args[len(ids)] = status

	query := `UPDATE ` + TableNameBusinessTablePool + `
		SET f_status = ?, f_update_time = NOW()
		WHERE f_id IN (` + strings.Join(placeholders, ",") + `) AND f_deleted = 0`

	_, err := m.conn.ExecContext(ctx, query, args...)
	return err
}
