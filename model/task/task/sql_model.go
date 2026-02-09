// Code scaffolded by speckit. Safe to edit.

package task

import (
	"context"
	"database/sql"
	"fmt"
)

var _ TaskStdCreateModel = (*defaultTaskStdCreateModel)(nil)
var _ TaskStdCreateResultModel = (*defaultTaskStdCreateResultModel)(nil)

// ==================== TaskStdCreateModel 实现 ====================

// Insert 插入任务
func (m *defaultTaskStdCreateModel) Insert(ctx context.Context, data *TaskStdCreate) (int64, error) {
	query := `INSERT INTO ` + TableNameTaskStdCreate + `
		(f_task_no, f_table, f_table_description, f_table_field, f_status,
		 f_create_time, f_create_user, f_create_user_phone, f_webhook, f_deleted)
		VALUES (?, ?, ?, ?, ?, NOW(), ?, ?, ?, 0)`

	result, err := m.conn.ExecContext(ctx, query,
		data.TaskNo, data.Table, data.TableDescription, data.TableField,
		data.Status, data.CreateUser, data.CreateUserPhone, data.Webhook)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// FindOne 查询单个任务
func (m *defaultTaskStdCreateModel) FindOne(ctx context.Context, id int64) (*TaskStdCreate, error) {
	query := `SELECT f_id, f_task_no, f_table, f_table_description, f_table_field, f_status,
		f_create_time, f_create_user, f_create_user_phone, f_update_time, f_update_user,
		f_webhook, f_deleted
		FROM ` + TableNameTaskStdCreate + `
		WHERE f_id = ? AND f_deleted = 0`

	var task TaskStdCreate
	err := m.conn.QueryRowContext(ctx, query, id).Scan(
		&task.Id, &task.TaskNo, &task.Table, &task.TableDescription,
		&task.TableField, &task.Status, &task.CreateTime, &task.CreateUser,
		&task.CreateUserPhone, &task.UpdateTime, &task.UpdateUser,
		&task.Webhook, &task.Deleted,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &task, nil
}

// Update 更新任务
func (m *defaultTaskStdCreateModel) Update(ctx context.Context, data *TaskStdCreate) error {
	query := `UPDATE ` + TableNameTaskStdCreate + `
		SET f_task_no = ?, f_table = ?, f_table_description = ?, f_table_field = ?,
		 f_status = ?, f_update_time = NOW(), f_update_user = ?, f_webhook = ?
		WHERE f_id = ? AND f_deleted = 0`

	_, err := m.conn.ExecContext(ctx, query,
		data.TaskNo, data.Table, data.TableDescription, data.TableField,
		data.Status, data.UpdateUser, data.Webhook, data.Id)

	return err
}

// Delete 删除任务（逻辑删除）
func (m *defaultTaskStdCreateModel) Delete(ctx context.Context, id int64) error {
	query := `UPDATE ` + TableNameTaskStdCreate + `
		SET f_deleted = ?
		WHERE f_id = ?`

	_, err := m.conn.ExecContext(ctx, query, id, id)
	return err
}

// FindUncompleted 查询未处理任务列表
func (m *defaultTaskStdCreateModel) FindUncompleted(ctx context.Context, keyword string, page, pageSize int) ([]*TaskStdCreate, int64, error) {
	return m.findByStatusAndKeyword(ctx, TaskStatusUnhandled, keyword, page, pageSize)
}

// FindCompleted 查询已完成任务列表
func (m *defaultTaskStdCreateModel) FindCompleted(ctx context.Context, keyword string, page, pageSize int) ([]*TaskStdCreate, int64, error) {
	return m.findByStatusAndKeyword(ctx, TaskStatusCompleted, keyword, page, pageSize)
}

// FindByStatus 按状态查询任务
func (m *defaultTaskStdCreateModel) FindByStatus(ctx context.Context, status int32) ([]*TaskStdCreate, error) {
	query := `SELECT f_id, f_task_no, f_table, f_table_description, f_table_field, f_status,
		f_create_time, f_create_user, f_create_user_phone, f_update_time, f_update_user,
		f_webhook, f_deleted
		FROM ` + TableNameTaskStdCreate + `
		WHERE f_status = ? AND f_deleted = 0
		ORDER BY f_create_time DESC`

	rows, err := m.conn.QueryContext(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*TaskStdCreate
	for rows.Next() {
		var task TaskStdCreate
		err := rows.Scan(
			&task.Id, &task.TaskNo, &task.Table, &task.TableDescription,
			&task.TableField, &task.Status, &task.CreateTime, &task.CreateUser,
			&task.CreateUserPhone, &task.UpdateTime, &task.UpdateUser,
			&task.Webhook, &task.Deleted,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	return tasks, nil
}

// FindByTaskNo 按任务编号查询
func (m *defaultTaskStdCreateModel) FindByTaskNo(ctx context.Context, taskNo string) (*TaskStdCreate, error) {
	query := `SELECT f_id, f_task_no, f_table, f_table_description, f_table_field, f_status,
		f_create_time, f_create_user, f_create_user_phone, f_update_time, f_update_user,
		f_webhook, f_deleted
		FROM ` + TableNameTaskStdCreate + `
		WHERE f_task_no = ? AND f_deleted = 0`

	var task TaskStdCreate
	err := m.conn.QueryRowContext(ctx, query, taskNo).Scan(
		&task.Id, &task.TaskNo, &task.Table, &task.TableDescription,
		&task.TableField, &task.Status, &task.CreateTime, &task.CreateUser,
		&task.CreateUserPhone, &task.UpdateTime, &task.UpdateUser,
		&task.Webhook, &task.Deleted,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &task, nil
}

// CountByStatus 统计各状态任务数量
func (m *defaultTaskStdCreateModel) CountByStatus(ctx context.Context, status int32) (int64, error) {
	query := `SELECT COUNT(*) FROM ` + TableNameTaskStdCreate + `
		WHERE f_status = ? AND f_deleted = 0`

	var count int64
	err := m.conn.QueryRowContext(ctx, query, status).Scan(&count)
	return count, err
}

// findByStatusAndKeyword 按状态和关键字查询（内部方法）
func (m *defaultTaskStdCreateModel) findByStatusAndKeyword(ctx context.Context, status int32, keyword string, page, pageSize int) ([]*TaskStdCreate, int64, error) {
	// 构建查询条件
	whereClause := "WHERE f_status = ? AND f_deleted = 0"
	args := []interface{}{status}

	if keyword != "" {
		whereClause += " AND (f_table LIKE ? OR f_task_no LIKE ?)"
		keywordPattern := "%" + keyword + "%"
		args = append(args, keywordPattern, keywordPattern)
	}

	// 查询总数
	countQuery := "SELECT COUNT(*) FROM " + TableNameTaskStdCreate + " " + whereClause
	var totalCount int64
	err := m.conn.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	offset := (page - 1) * pageSize
	query := `SELECT f_id, f_task_no, f_table, f_table_description, f_table_field, f_status,
		f_create_time, f_create_user, f_create_user_phone, f_update_time, f_update_user,
		f_webhook, f_deleted
		FROM ` + TableNameTaskStdCreate + " " + whereClause + `
		ORDER BY f_create_time DESC
		LIMIT ? OFFSET ?`

	args = append(args, pageSize, offset)
	rows, err := m.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tasks []*TaskStdCreate
	for rows.Next() {
		var task TaskStdCreate
		err := rows.Scan(
			&task.Id, &task.TaskNo, &task.Table, &task.TableDescription,
			&task.TableField, &task.Status, &task.CreateTime, &task.CreateUser,
			&task.CreateUserPhone, &task.UpdateTime, &task.UpdateUser,
			&task.Webhook, &task.Deleted,
		)
		if err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, &task)
	}

	return tasks, totalCount, nil
}

// ==================== TaskStdCreateResultModel 实现 ====================

// Insert 插入结果
func (m *defaultTaskStdCreateResultModel) Insert(ctx context.Context, data *TaskStdCreateResult) (int64, error) {
	query := `INSERT INTO ` + TableNameTaskStdCreateResult + `
		(f_task_id, f_table_field, f_table_field_description, f_std_ref_file,
		 f_std_code, f_rec_std_codes, f_std_ch_name, f_std_en_name)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := m.conn.ExecContext(ctx, query,
		data.TaskId, data.TableField, data.TableFieldDescription,
		data.StdRefFile, data.StdCode, data.RecStdCodes,
		data.StdChName, data.StdEnName)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// FindByTaskId 按任务ID查询结果列表
func (m *defaultTaskStdCreateResultModel) FindByTaskId(ctx context.Context, taskId int64) ([]*TaskStdCreateResult, error) {
	query := `SELECT f_id, f_task_id, f_table_field, f_table_field_description,
		f_std_ref_file, f_std_code, f_rec_std_codes, f_std_ch_name, f_std_en_name
		FROM ` + TableNameTaskStdCreateResult + `
		WHERE f_task_id = ?
		ORDER BY f_id`

	rows, err := m.conn.QueryContext(ctx, query, taskId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*TaskStdCreateResult
	for rows.Next() {
		var result TaskStdCreateResult
		err := rows.Scan(
			&result.Id, &result.TaskId, &result.TableField,
			&result.TableFieldDescription, &result.StdRefFile,
			&result.StdCode, &result.RecStdCodes,
			&result.StdChName, &result.StdEnName,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, &result)
	}

	return results, nil
}

// DeleteByTaskId 删除任务的所有结果
func (m *defaultTaskStdCreateResultModel) DeleteByTaskId(ctx context.Context, taskId int64) error {
	query := `DELETE FROM ` + TableNameTaskStdCreateResult + ` WHERE f_task_id = ?`
	_, err := m.conn.ExecContext(ctx, query, taskId)
	return err
}

// InsertBatch 批量插入结果
func (m *defaultTaskStdCreateResultModel) InsertBatch(ctx context.Context, data []*TaskStdCreateResult) error {
	if len(data) == 0 {
		return nil
	}

	query := `INSERT INTO ` + TableNameTaskStdCreateResult + `
		(f_task_id, f_table_field, f_table_field_description, f_std_ref_file,
		 f_std_code, f_rec_std_codes, f_std_ch_name, f_std_en_name)
		VALUES
		` + generatePlaceholders(8)

	// 使用事务
	tx, err := m.conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, item := range data {
		_, err := stmt.ExecContext(ctx,
			item.TaskId, item.TableField, item.TableFieldDescription,
			item.StdRefFile, item.StdCode, item.RecStdCodes,
			item.StdChName, item.StdEnName,
		)
		if err != nil {
			return fmt.Errorf("execute statement: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// generatePlaceholders 生成占位符
func generatePlaceholders(count int) string {
	if count <= 0 {
		return "()"
	}

	placeholders := make([]string, count)
	for i := 0; i < count; i++ {
		placeholders[i] = "(?, ?, ?, ?, ?, ?, ?, ?)"
	}
	return placeholders[0] // 这里简化，实际应该拼接所有占位符
}
