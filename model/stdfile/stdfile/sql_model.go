// Code scaffolded by goctl. Safe to edit.
//go:build !mock_logic_off
// +build !mock_logic_off

package stdfile

import (
	"context"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

// Ensure customStdFileModel implements the interface
var _ = sqlx.Conn{} // Blank reference to satisfy import

// Insert 插入标准文件
func (m *customStdFileModel) Insert(ctx context.Context, data *StdFile) (int64, error) {
	query := `INSERT INTO ` + TableName + ` (
		f_number, f_name, f_catalog_id, f_act_date, f_publish_date,
		f_attachment_type, f_attachment_url, f_file_name, f_org_type,
		f_description, f_state, f_department_ids, f_third_dept_id,
		f_version, f_create_time, f_create_user, f_update_time, f_update_user
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), ?, NOW(), ?)`

	result, err := m.conn.ExecContext(ctx, query,
		data.Number, data.Name, data.CatalogId, data.ActDate, data.PublishDate,
		data.AttachmentType, data.AttachmentUrl, data.FileName, data.OrgType,
		data.Description, data.State, data.DepartmentIds, data.ThirdDeptId,
		data.Version, data.CreateUser, data.UpdateUser,
	)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// FindOne 根据ID查询
func (m *customStdFileModel) FindOne(ctx context.Context, id int64) (*StdFile, error) {
	query := `SELECT ` + getAllColumns() + ` FROM ` + TableName + ` WHERE f_id = ? AND f_deleted = 0 LIMIT 1`

	var resp StdFile
	err := m.conn.QueryRowContext(ctx, query, id).Scan(
		&resp.Id, &resp.Number, &resp.Name, &resp.CatalogId,
		&resp.ActDate, &resp.PublishDate, &resp.DisableDate,
		&resp.AttachmentType, &resp.AttachmentUrl, &resp.FileName,
		&resp.OrgType, &resp.Description, &resp.State,
		&resp.DisableReason, &resp.AuthorityId, &resp.DepartmentIds,
		&resp.ThirdDeptId, &resp.Version, &resp.CreateTime,
		&resp.CreateUser, &resp.UpdateTime, &resp.UpdateUser, &resp.Deleted,
	)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Update 更新标准文件
func (m *customStdFileModel) Update(ctx context.Context, data *StdFile) error {
	query := `UPDATE ` + TableName + ` SET
		f_number = ?, f_name = ?, f_catalog_id = ?, f_act_date = ?, f_publish_date = ?,
		f_attachment_type = ?, f_attachment_url = ?, f_file_name = ?, f_org_type = ?,
		f_description = ?, f_state = ?, f_department_ids = ?, f_third_dept_id = ?,
		f_version = ?, f_update_user = ?, f_update_time = NOW()
	WHERE f_id = ? AND f_deleted = 0`

	_, err := m.conn.ExecContext(ctx, query,
		data.Number, data.Name, data.CatalogId, data.ActDate, data.PublishDate,
		data.AttachmentType, data.AttachmentUrl, data.FileName, data.OrgType,
		data.Description, data.State, data.DepartmentIds, data.ThirdDeptId,
		data.Version, data.UpdateUser, data.Id,
	)

	return err
}

// Delete 删除标准文件（软删除）
func (m *customStdFileModel) Delete(ctx context.Context, id int64) error {
	query := `UPDATE ` + TableName + ` SET f_deleted = 1, f_update_time = NOW() WHERE f_id = ?`

	_, err := m.conn.ExecContext(ctx, query, id)
	return err
}

// FindByIds 根据ID列表查询
func (m *customStdFileModel) FindByIds(ctx context.Context, ids []int64) ([]*StdFile, error) {
	if len(ids) == 0 {
		return []*StdFile{}, nil
	}

	placeholders := strings.Repeat("?,", len(ids))
	query := `SELECT ` + getAllColumns() + ` FROM ` + TableName + ` WHERE f_id IN (` + placeholders[:len(placeholders)-1] + `) AND f_deleted = 0`

	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	rows, err := m.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []*StdFile{}
	for rows.Next() {
		var item StdFile
		err := rows.Scan(
			&item.Id, &item.Number, &item.Name, &item.CatalogId,
			&item.ActDate, &item.PublishDate, &item.DisableDate,
			&item.AttachmentType, &item.AttachmentUrl, &item.FileName,
			&item.OrgType, &item.Description, &item.State,
			&item.DisableReason, &item.AuthorityId, &item.DepartmentIds,
			&item.ThirdDeptId, &item.Version, &item.CreateTime,
			&item.CreateUser, &item.UpdateTime, &item.UpdateUser, &item.Deleted,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// FindByNumber 根据标准编号查询
func (m *customStdFileModel) FindByNumber(ctx context.Context, number string) ([]*StdFile, error) {
	query := `SELECT ` + getAllColumns() + ` FROM ` + TableName + ` WHERE f_number = ? AND f_deleted = 0`

	rows, err := m.conn.QueryContext(ctx, query, number)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []*StdFile{}
	for rows.Next() {
		var item StdFile
		err := rows.Scan(
			&item.Id, &item.Number, &item.Name, &item.CatalogId,
			&item.ActDate, &item.PublishDate, &item.DisableDate,
			&item.AttachmentType, &item.AttachmentUrl, &item.FileName,
			&item.OrgType, &item.Description, &item.State,
			&item.DisableReason, &item.AuthorityId, &item.DepartmentIds,
			&item.ThirdDeptId, &item.Version, &item.CreateTime,
			&item.CreateUser, &item.UpdateTime, &item.UpdateUser, &item.Deleted,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// FindByNameAndOrgType 根据名称和组织类型查询
func (m *customStdFileModel) FindByNameAndOrgType(ctx context.Context, name string, orgType int) ([]*StdFile, error) {
	query := `SELECT ` + getAllColumns() + ` FROM ` + TableName + ` WHERE f_name = ? AND f_org_type = ? AND f_deleted = 0`

	rows, err := m.conn.QueryContext(ctx, query, name, orgType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []*StdFile{}
	for rows.Next() {
		var item StdFile
		err := rows.Scan(
			&item.Id, &item.Number, &item.Name, &item.CatalogId,
			&item.ActDate, &item.PublishDate, &item.DisableDate,
			&item.AttachmentType, &item.AttachmentUrl, &item.FileName,
			&item.OrgType, &item.Description, &item.State,
			&item.DisableReason, &item.AuthorityId, &item.DepartmentIds,
			&item.ThirdDeptId, &item.Version, &item.CreateTime,
			&item.CreateUser, &item.UpdateTime, &item.UpdateUser, &item.Deleted,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// FindByCatalogIds 根据目录ID列表查询
func (m *customStdFileModel) FindByCatalogIds(ctx context.Context, opts *FindOptions) ([]*StdFile, int64, error) {
	whereClause := ` WHERE f_deleted = 0`
	args := []interface{}{}

	if opts != nil {
		if opts.CatalogId != nil {
			whereClause += ` AND f_catalog_id = ?`
			args = append(args, *opts.CatalogId)
		}
		if opts.Keyword != "" {
			whereClause += ` AND (f_number LIKE ? OR f_name LIKE ?)`
			keyword := `%` + opts.Keyword + `%`
			args = append(args, keyword, keyword)
		}
		if opts.OrgType != nil {
			whereClause += ` AND f_org_type = ?`
			args = append(args, *opts.OrgType)
		}
		if opts.State != nil {
			whereClause += ` AND f_state = ?`
			args = append(args, *opts.State)
		}
		if opts.DepartmentId != "" {
			whereClause += ` AND f_department_ids LIKE ?`
			args = append(args, `%`+opts.DepartmentId+`%`)
		}
	}

	// 先查询总数
	countQuery := `SELECT COUNT(*) FROM ` + TableName + whereClause
	var totalCount int64
	err := m.conn.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// 添加排序和分页
	if opts != nil && opts.Sort != "" {
		whereClause += ` ORDER BY ` + opts.Sort
		if opts.Direction == "desc" {
			whereClause += ` DESC`
		} else {
			whereClause += ` ASC`
		}
	} else {
		whereClause += ` ORDER BY f_update_time DESC`
	}

	if opts != nil && opts.PageSize > 0 {
		whereClause += ` LIMIT ? OFFSET ?`
		offset := (opts.Page - 1) * opts.PageSize
		args = append(args, opts.PageSize, offset)
	}

	query := `SELECT ` + getAllColumns() + ` FROM ` + TableName + whereClause

	rows, err := m.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := []*StdFile{}
	for rows.Next() {
		var item StdFile
		err := rows.Scan(
			&item.Id, &item.Number, &item.Name, &item.CatalogId,
			&item.ActDate, &item.PublishDate, &item.DisableDate,
			&item.AttachmentType, &item.AttachmentUrl, &item.FileName,
			&item.OrgType, &item.Description, &item.State,
			&item.DisableReason, &item.AuthorityId, &item.DepartmentIds,
			&item.ThirdDeptId, &item.Version, &item.CreateTime,
			&item.CreateUser, &item.UpdateTime, &item.UpdateUser, &item.Deleted,
		)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, &item)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return items, totalCount, nil
}

// FindDataExists 检查数据是否存在
func (m *customStdFileModel) FindDataExists(ctx context.Context, filterId int64, number string, orgType int, name string, deptIds string) (*StdFile, error) {
	whereClause := ` WHERE f_deleted = 0`
	args := []interface{}{}

	if filterId > 0 {
		whereClause += ` AND f_id != ?`
		args = append(args, filterId)
	}
	if number != "" {
		whereClause += ` AND f_number = ?`
		args = append(args, number)
	}
	if name != "" {
		whereClause += ` AND f_name = ?`
		args = append(args, name)
	}
	if orgType >= 0 {
		whereClause += ` AND f_org_type = ?`
		args = append(args, orgType)
	}
	if deptIds != "" {
		whereClause += ` AND f_department_ids = ?`
		args = append(args, deptIds)
	}

	query := `SELECT ` + getAllColumns() + ` FROM ` + TableName + whereClause + ` LIMIT 1`

	var resp StdFile
	err := m.conn.QueryRowContext(ctx, query, args...).Scan(
		&resp.Id, &resp.Number, &resp.Name, &resp.CatalogId,
		&resp.ActDate, &resp.PublishDate, &resp.DisableDate,
		&resp.AttachmentType, &resp.AttachmentUrl, &resp.FileName,
		&resp.OrgType, &resp.Description, &resp.State,
		&resp.DisableReason, &resp.AuthorityId, &resp.DepartmentIds,
		&resp.ThirdDeptId, &resp.Version, &resp.CreateTime,
		&resp.CreateUser, &resp.UpdateTime, &resp.UpdateTime, &resp.Deleted,
	)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// UpdateState 更新状态
func (m *customStdFileModel) UpdateState(ctx context.Context, id int64, state int, disableReason string) error {
	var disableDate *time.Time
	if state == StateDisable {
		now := time.Now()
		disableDate = &now
	}

	query := `UPDATE ` + TableName + ` SET f_state = ?, f_disable_date = ?, f_disable_reason = ?, f_update_time = NOW() WHERE f_id = ? AND f_deleted = 0`

	_, err := m.conn.ExecContext(ctx, query, state, disableDate, disableReason, id)
	return err
}

// RemoveCatalog 移除目录
func (m *customStdFileModel) RemoveCatalog(ctx context.Context, ids []int64, catalogId int64, updateUser string) error {
	if len(ids) == 0 {
		return nil
	}

	placeholders := strings.Repeat("?,", len(ids))
	query := `UPDATE ` + TableName + ` SET f_catalog_id = ?, f_update_user = ?, f_update_time = NOW() WHERE f_id IN (` + placeholders[:len(placeholders)-1] + `) AND f_deleted = 0`

	args := make([]interface{}, 0, len(ids)+2)
	args = append(args, catalogId, updateUser)
	for _, id := range ids {
		args = append(args, id)
	}

	_, err := m.conn.ExecContext(ctx, query, args...)
	return err
}

// BatchUpdateState 批量更新状态
func (m *customStdFileModel) BatchUpdateState(ctx context.Context, ids []int64, state int, disableReason string) error {
	if len(ids) == 0 {
		return nil
	}

	var disableDate *time.Time
	if state == StateDisable {
		now := time.Now()
		disableDate = &now
	}

	placeholders := strings.Repeat("?,", len(ids))
	query := `UPDATE ` + TableName + ` SET
		f_state = ?, f_disable_date = ?, f_disable_reason = ?, f_update_time = NOW()
	WHERE f_id IN (` + placeholders[:len(placeholders)-1] + `) AND f_deleted = 0`

	args := make([]interface{}, 0, len(ids)+3)
	args = append(args, state, disableDate, disableReason)
	for _, id := range ids {
		args = append(args, id)
	}

	_, err := m.conn.ExecContext(ctx, query, args...)
	return err
}

// DeleteByIds 批量删除
func (m *customStdFileModel) DeleteByIds(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}

	placeholders := strings.Repeat("?,", len(ids))
	query := `UPDATE ` + TableName + ` SET f_deleted = 1, f_update_time = NOW() WHERE f_id IN (` + placeholders[:len(placeholders)-1] + `)`

	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	_, err := m.conn.ExecContext(ctx, query, args...)
	return err
}

// getAllColumns 获取所有列名
func getAllColumns() string {
	return `f_id, f_number, f_name, f_catalog_id, f_act_date, f_publish_date, f_disable_date,
		f_attachment_type, f_attachment_url, f_file_name, f_org_type, f_description,
		f_state, f_disable_reason, f_authority_id, f_department_ids, f_third_dept_id,
		f_version, f_create_time, f_create_user, f_update_time, f_update_user, f_deleted`
}
