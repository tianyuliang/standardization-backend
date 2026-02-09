package dict

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type defaultDictModel struct {
	conn *sqlx.Conn
}

type defaultDictEnumModel struct {
	conn *sqlx.Conn
}

type defaultRelationDictFileModel struct {
	conn *sqlx.Conn
}

// ============================================
// DictModel 实现
// ============================================

func (m *defaultDictModel) Insert(ctx context.Context, data *Dict) (int64, error) {
	query := `INSERT INTO ` + TableNameDict + `
		(f_code, f_ch_name, f_en_name, f_description, f_catalog_id, f_org_type,
		 f_version, f_state, f_disable_reason, f_authority_id, f_department_ids,
		 f_third_dept_id, f_create_time, f_create_user, f_update_time, f_update_user, f_deleted)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	data.CreateTime = now
	data.UpdateTime = now

	result, err := m.conn.ExecContext(ctx, query,
		data.Code, data.ChName, data.EnName, data.Description, data.CatalogId, data.OrgType,
		data.Version, data.State, data.DisableReason, data.AuthorityId, data.DepartmentIds,
		data.ThirdDeptId, data.CreateTime, data.CreateUser, data.UpdateTime, data.UpdateUser, 0,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	data.Id = id
	return id, nil
}

func (m *defaultDictModel) FindOne(ctx context.Context, id int64) (*Dict, error) {
	query := `SELECT f_id, f_code, f_ch_name, f_en_name, f_description, f_catalog_id, f_org_type,
		f_version, f_state, f_disable_reason, f_authority_id, f_department_ids, f_third_dept_id,
		f_create_time, f_create_user, f_update_time, f_update_user, f_deleted
		FROM ` + TableNameDict + ` WHERE f_id = ? AND f_deleted = 0 LIMIT 1`

	var data Dict
	err := m.conn.QueryRowContext(ctx, query, id).Scan(
		&data.Id, &data.Code, &data.ChName, &data.EnName, &data.Description, &data.CatalogId, &data.OrgType,
		&data.Version, &data.State, &data.DisableReason, &data.AuthorityId, &data.DepartmentIds, &data.ThirdDeptId,
		&data.CreateTime, &data.CreateUser, &data.UpdateTime, &data.UpdateUser, &data.Deleted,
	)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (m *defaultDictModel) FindByCode(ctx context.Context, code int64) (*Dict, error) {
	query := `SELECT f_id, f_code, f_ch_name, f_en_name, f_description, f_catalog_id, f_org_type,
		f_version, f_state, f_disable_reason, f_authority_id, f_department_ids, f_third_dept_id,
		f_create_time, f_create_user, f_update_time, f_update_user, f_deleted
		FROM ` + TableNameDict + ` WHERE f_code = ? AND f_deleted = 0 LIMIT 1`

	var data Dict
	err := m.conn.QueryRowContext(ctx, query, code).Scan(
		&data.Id, &data.Code, &data.ChName, &data.EnName, &data.Description, &data.CatalogId, &data.OrgType,
		&data.Version, &data.State, &data.DisableReason, &data.AuthorityId, &data.DepartmentIds, &data.ThirdDeptId,
		&data.CreateTime, &data.CreateUser, &data.UpdateTime, &data.UpdateUser, &data.Deleted,
	)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (m *defaultDictModel) Update(ctx context.Context, data *Dict) error {
	query := `UPDATE ` + TableNameDict + ` SET
		f_ch_name = ?, f_en_name = ?, f_description = ?, f_catalog_id = ?, f_org_type = ?,
		f_version = ?, f_state = ?, f_disable_reason = ?, f_department_ids = ?,
		f_update_time = ?, f_update_user = ?
		WHERE f_id = ? AND f_deleted = 0`

	data.UpdateTime = time.Now()

	_, err := m.conn.ExecContext(ctx, query,
		data.ChName, data.EnName, data.Description, data.CatalogId, data.OrgType,
		data.Version, data.State, data.DisableReason, data.DepartmentIds,
		data.UpdateTime, data.UpdateUser, data.Id,
	)

	return err
}

func (m *defaultDictModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM ` + TableNameDict + ` WHERE f_id = ?`
	_, err := m.conn.ExecContext(ctx, query, id)
	return err
}

func (m *defaultDictModel) FindByIds(ctx context.Context, ids []int64) ([]*Dict, error) {
	if len(ids) == 0 {
		return []*Dict{}, nil
	}

	query, args, err := sqlx.In(`SELECT f_id, f_code, f_ch_name, f_en_name, f_description, f_catalog_id, f_org_type,
		f_version, f_state, f_disable_reason, f_authority_id, f_department_ids, f_third_dept_id,
		f_create_time, f_create_user, f_update_time, f_update_user, f_deleted
		FROM `+TableNameDict+` WHERE f_id IN (?) AND f_deleted = 0`, ids)
	if err != nil {
		return nil, err
	}

	query = m.conn.Rebind(query)
	rows, err := m.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []*Dict
	for rows.Next() {
		var item Dict
		err := rows.Scan(
			&item.Id, &item.Code, &item.ChName, &item.EnName, &item.Description, &item.CatalogId, &item.OrgType,
			&item.Version, &item.State, &item.DisableReason, &item.AuthorityId, &item.DepartmentIds, &item.ThirdDeptId,
			&item.CreateTime, &item.CreateUser, &item.UpdateTime, &item.UpdateUser, &item.Deleted,
		)
		if err != nil {
			return nil, err
		}
		data = append(data, &item)
	}

	return data, nil
}

func (m *defaultDictModel) FindByChNameAndOrgType(ctx context.Context, chName string, orgType int32) (*Dict, error) {
	query := `SELECT f_id, f_code, f_ch_name, f_en_name, f_description, f_catalog_id, f_org_type,
		f_version, f_state, f_disable_reason, f_authority_id, f_department_ids, f_third_dept_id,
		f_create_time, f_create_user, f_update_time, f_update_user, f_deleted
		FROM ` + TableNameDict + ` WHERE f_ch_name = ? AND f_org_type = ? AND f_deleted = 0 LIMIT 1`

	var data Dict
	err := m.conn.QueryRowContext(ctx, query, chName, orgType).Scan(
		&data.Id, &data.Code, &data.ChName, &data.EnName, &data.Description, &data.CatalogId, &data.OrgType,
		&data.Version, &data.State, &data.DisableReason, &data.AuthorityId, &data.DepartmentIds, &data.ThirdDeptId,
		&data.CreateTime, &data.CreateUser, &data.UpdateTime, &data.UpdateUser, &data.Deleted,
	)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (m *defaultDictModel) FindByEnNameAndOrgType(ctx context.Context, enName string, orgType int32) (*Dict, error) {
	query := `SELECT f_id, f_code, f_ch_name, f_en_name, f_description, f_catalog_id, f_org_type,
		f_version, f_state, f_disable_reason, f_authority_id, f_department_ids, f_third_dept_id,
		f_create_time, f_create_user, f_update_time, f_update_user, f_deleted
		FROM ` + TableNameDict + ` WHERE f_en_name = ? AND f_org_type = ? AND f_deleted = 0 LIMIT 1`

	var data Dict
	err := m.conn.QueryRowContext(ctx, query, enName, orgType).Scan(
		&data.Id, &data.Code, &data.ChName, &data.EnName, &data.Description, &data.CatalogId, &data.OrgType,
		&data.Version, &data.State, &data.DisableReason, &data.AuthorityId, &data.DepartmentIds, &data.ThirdDeptId,
		&data.CreateTime, &data.CreateUser, &data.UpdateTime, &data.UpdateUser, &data.Deleted,
	)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (m *defaultDictModel) FindByCatalogIds(ctx context.Context, opts *FindOptions) ([]*Dict, int64, error) {
	whereClause := `WHERE d.f_deleted = 0`
	args := []interface{}{}

	if opts.CatalogId != nil {
		whereClause += ` AND d.f_catalog_id = ?`
		args = append(args, *opts.CatalogId)
	}

	if opts.OrgType != nil {
		whereClause += ` AND d.f_org_type = ?`
		args = append(args, *opts.OrgType)
	}

	if opts.State != nil {
		whereClause += ` AND d.f_state = ?`
		args = append(args, *opts.State)
	}

	if opts.Keyword != "" {
		whereClause += ` AND (d.f_ch_name LIKE ? OR d.f_en_name LIKE ?)`
		keyword := "%" + opts.Keyword + "%"
		args = append(args, keyword, keyword)
	}

	if opts.DepartmentId != "" {
		whereClause += ` AND d.f_department_ids LIKE ?`
		args = append(args, "%"+opts.DepartmentId+"%")
	}

	// 查询总数
	countQuery := `SELECT COUNT(*) FROM ` + TableNameDict + ` AS d ` + whereClause
	var totalCount int64
	err := m.conn.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	sortClause := ` ORDER BY d.f_update_time DESC`
	if opts.Sort != "" {
		sortClause = fmt.Sprintf(` ORDER BY d.%s %s`, opts.Sort, opts.Direction)
	}

	offset := (opts.Page - 1) * opts.PageSize
	limitClause := fmt.Sprintf(` LIMIT %d OFFSET %d`, opts.PageSize, offset)

	query := `SELECT d.f_id, d.f_code, d.f_ch_name, d.f_en_name, d.f_description, d.f_catalog_id, d.f_org_type,
		d.f_version, d.f_state, d.f_disable_reason, d.f_authority_id, d.f_department_ids, d.f_third_dept_id,
		d.f_create_time, d.f_create_user, d.f_update_time, d.f_update_user, d.f_deleted
		FROM ` + TableNameDict + ` AS d ` + whereClause + sortClause + limitClause

	rows, err := m.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var data []*Dict
	for rows.Next() {
		var item Dict
		err := rows.Scan(
			&item.Id, &item.Code, &item.ChName, &item.EnName, &item.Description, &item.CatalogId, &item.OrgType,
			&item.Version, &item.State, &item.DisableReason, &item.AuthorityId, &item.DepartmentIds, &item.ThirdDeptId,
			&item.CreateTime, &item.CreateUser, &item.UpdateTime, &item.UpdateUser, &item.Deleted,
		)
		if err != nil {
			return nil, 0, err
		}
		data = append(data, &item)
	}

	return data, totalCount, nil
}

func (m *defaultDictModel) FindByStdFileCatalog(ctx context.Context, opts *FindOptions) ([]*Dict, int64, error) {
	// JOIN t_relation_dict_file 查询关联指定目录文件的码表
	whereClause := `WHERE d.f_deleted = 0`
	args := []interface{}{}

	if opts.CatalogId != nil {
		// 子查询：查询关联该目录下所有文件的码表
		whereClause += ` AND EXISTS (
			SELECT 1 FROM ` + TableNameRelationDictFile + ` rdf
			INNER JOIN t_std_file sf ON rdf.f_file_id = sf.f_id
			WHERE rdf.f_dict_id = d.f_id
			AND (sf.f_catalog_id = ? OR sf.f_catalog_id IN (
				SELECT f_id FROM t_catalog WHERE f_path LIKE CONCAT(
					(SELECT f_path FROM t_catalog WHERE f_id = ?), '/%'
				)
			))
		)`
		args = append(args, *opts.CatalogId, *opts.CatalogId)
	} else {
		// 查询未关联任何文件的码表
		whereClause += ` AND NOT EXISTS (
			SELECT 1 FROM ` + TableNameRelationDictFile + ` rdf
			WHERE rdf.f_dict_id = d.f_id
		)`
	}

	if opts.OrgType != nil {
		whereClause += ` AND d.f_org_type = ?`
		args = append(args, *opts.OrgType)
	}

	if opts.State != nil {
		whereClause += ` AND d.f_state = ?`
		args = append(args, *opts.State)
	}

	if opts.Keyword != "" {
		whereClause += ` AND (d.f_ch_name LIKE ? OR d.f_en_name LIKE ?)`
		keyword := "%" + opts.Keyword + "%"
		args = append(args, keyword, keyword)
	}

	// 查询总数
	countQuery := `SELECT COUNT(*) FROM ` + TableNameDict + ` AS d ` + whereClause
	var totalCount int64
	err := m.conn.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	sortClause := ` ORDER BY d.f_update_time DESC`
	if opts.Sort != "" {
		sortClause = fmt.Sprintf(` ORDER BY d.%s %s`, opts.Sort, opts.Direction)
	}

	offset := (opts.Page - 1) * opts.PageSize
	limitClause := fmt.Sprintf(` LIMIT %d OFFSET %d`, opts.PageSize, offset)

	query := `SELECT d.f_id, d.f_code, d.f_ch_name, d.f_en_name, d.f_description, d.f_catalog_id, d.f_org_type,
		d.f_version, d.f_state, d.f_disable_reason, d.f_authority_id, d.f_department_ids, d.f_third_dept_id,
		d.f_create_time, d.f_create_user, d.f_update_time, d.f_update_user, d.f_deleted
		FROM ` + TableNameDict + ` AS d ` + whereClause + sortClause + limitClause

	rows, err := m.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var data []*Dict
	for rows.Next() {
		var item Dict
		err := rows.Scan(
			&item.Id, &item.Code, &item.ChName, &item.EnName, &item.Description, &item.CatalogId, &item.OrgType,
			&item.Version, &item.State, &item.DisableReason, &item.AuthorityId, &item.DepartmentIds, &item.ThirdDeptId,
			&item.CreateTime, &item.CreateUser, &item.UpdateTime, &item.UpdateUser, &item.Deleted,
		)
		if err != nil {
			return nil, 0, err
		}
		data = append(data, &item)
	}

	return data, totalCount, nil
}

func (m *defaultDictModel) FindByFileId(ctx context.Context, fileId int64) ([]*Dict, error) {
	query := `SELECT d.f_id, d.f_code, d.f_ch_name, d.f_en_name, d.f_description, d.f_catalog_id, d.f_org_type,
		d.f_version, d.f_state, d.f_disable_reason, d.f_authority_id, d.f_department_ids, d.f_third_dept_id,
		d.f_create_time, d.f_create_user, d.f_update_time, d.f_update_user, d.f_deleted
		FROM ` + TableNameDict + ` d
		INNER JOIN ` + TableNameRelationDictFile + ` rdf ON d.f_id = rdf.f_dict_id
		WHERE rdf.f_file_id = ? AND d.f_deleted = 0`

	rows, err := m.conn.QueryContext(ctx, query, fileId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []*Dict
	for rows.Next() {
		var item Dict
		err := rows.Scan(
			&item.Id, &item.Code, &item.ChName, &item.EnName, &item.Description, &item.CatalogId, &item.OrgType,
			&item.Version, &item.State, &item.DisableReason, &item.AuthorityId, &item.DepartmentIds, &item.ThirdDeptId,
			&item.CreateTime, &item.CreateUser, &item.UpdateTime, &item.UpdateUser, &item.Deleted,
		)
		if err != nil {
			return nil, err
		}
		data = append(data, &item)
	}

	return data, nil
}

func (m *defaultDictModel) FindDataExists(ctx context.Context, chName, enName string, orgType int32, filterId int64, deptIds string) (*Dict, error) {
	whereClause := `WHERE f_deleted = 0 AND f_org_type = ?`
	args := []interface{}{orgType}

	if chName != "" {
		whereClause += ` AND f_ch_name = ?`
		args = append(args, chName)
	} else if enName != "" {
		whereClause += ` AND f_en_name = ?`
		args = append(args, enName)
	}

	if filterId > 0 {
		whereClause += ` AND f_id != ?`
		args = append(args, filterId)
	}

	if deptIds != "" {
		whereClause += ` AND f_department_ids LIKE ?`
		args = append(args, "%"+deptIds+"%")
	}

	query := `SELECT f_id, f_code, f_ch_name, f_en_name, f_description, f_catalog_id, f_org_type,
		f_version, f_state, f_disable_reason, f_authority_id, f_department_ids, f_third_dept_id,
		f_create_time, f_create_user, f_update_time, f_update_user, f_deleted
		FROM ` + TableNameDict + ` ` + whereClause + ` LIMIT 1`

	var data Dict
	err := m.conn.QueryRowContext(ctx, query, args...).Scan(
		&data.Id, &data.Code, &data.ChName, &data.EnName, &data.Description, &data.CatalogId, &data.OrgType,
		&data.Version, &data.State, &data.DisableReason, &data.AuthorityId, &data.DepartmentIds, &data.ThirdDeptId,
		&data.CreateTime, &data.CreateUser, &data.UpdateTime, &data.UpdateUser, &data.Deleted,
	)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (m *defaultDictModel) UpdateState(ctx context.Context, id int64, state int32, disableReason string) error {
	query := `UPDATE ` + TableNameDict + ` SET f_state = ?, f_disable_reason = ?, f_update_time = ? WHERE f_id = ? AND f_deleted = 0`
	_, err := m.conn.ExecContext(ctx, query, state, disableReason, time.Now(), id)
	return err
}

func (m *defaultDictModel) UpdateVersionByIds(ctx context.Context, ids []int64, updateUser string) error {
	if len(ids) == 0 {
		return nil
	}

	query := `UPDATE ` + TableNameDict + ` SET f_version = f_version + 1, f_update_user = ?, f_update_time = ? WHERE f_id IN (?)`
	now := time.Now()

	query, args, err := sqlx.In(query, updateUser, now, ids)
	if err != nil {
		return err
	}

	query = m.conn.Rebind(query)
	_, err = m.conn.ExecContext(ctx, query, args...)
	return err
}

func (m *defaultDictModel) DeleteByIds(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}

	query := `DELETE FROM ` + TableNameDict + ` WHERE f_id IN (?)`
	query, args, err := sqlx.In(query, ids)
	if err != nil {
		return err
	}

	query = m.conn.Rebind(query)
	_, err = m.conn.ExecContext(ctx, query, args...)
	return err
}

// ============================================
// DictEnumModel 实现
// ============================================

func (m *defaultDictEnumModel) Insert(ctx context.Context, data *DictEnum) (int64, error) {
	query := `INSERT INTO ` + TableNameDictEnum + `
		(f_dict_id, f_code, f_value, f_create_time, f_create_user)
		VALUES (?, ?, ?, ?, ?)`

	data.CreateTime = time.Now()

	result, err := m.conn.ExecContext(ctx, query,
		data.DictId, data.Code, data.Value, data.CreateTime, data.CreateUser,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	data.Id = id
	return id, nil
}

func (m *defaultDictEnumModel) FindByDictId(ctx context.Context, dictId int64) ([]*DictEnum, error) {
	query := `SELECT f_id, f_dict_id, f_code, f_value, f_create_time, f_create_user
		FROM ` + TableNameDictEnum + ` WHERE f_dict_id = ? ORDER BY f_id`

	rows, err := m.conn.QueryContext(ctx, query, dictId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []*DictEnum
	for rows.Next() {
		var item DictEnum
		err := rows.Scan(&item.Id, &item.DictId, &item.Code, &item.Value, &item.CreateTime, &item.CreateUser)
		if err != nil {
			return nil, err
		}
		data = append(data, &item)
	}

	return data, nil
}

func (m *defaultDictEnumModel) FindPageByDictId(ctx context.Context, dictId int64, keyword string, offset, limit int) ([]*DictEnum, int64, error) {
	whereClause := `WHERE f_dict_id = ?`
	args := []interface{}{dictId}

	if keyword != "" {
		whereClause += ` AND (f_code LIKE ? OR f_value LIKE ?)`
		keywordPattern := "%" + keyword + "%"
		args = append(args, keywordPattern, keywordPattern)
	}

	// 查询总数
	countQuery := `SELECT COUNT(*) FROM ` + TableNameDictEnum + ` ` + whereClause
	var totalCount int64
	err := m.conn.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	query := `SELECT f_id, f_dict_id, f_code, f_value, f_create_time, f_create_user
		FROM ` + TableNameDictEnum + ` ` + whereClause + ` ORDER BY f_id LIMIT ? OFFSET ?`

	args = append(args, limit, (offset-1)*limit)
	rows, err := m.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var data []*DictEnum
	for rows.Next() {
		var item DictEnum
		err := rows.Scan(&item.Id, &item.DictId, &item.Code, &item.Value, &item.CreateTime, &item.CreateUser)
		if err != nil {
			return nil, 0, err
		}
		data = append(data, &item)
	}

	return data, totalCount, nil
}

func (m *defaultDictEnumModel) DeleteByDictId(ctx context.Context, dictId int64) error {
	query := `DELETE FROM ` + TableNameDictEnum + ` WHERE f_dict_id = ?`
	_, err := m.conn.ExecContext(ctx, query, dictId)
	return err
}

func (m *defaultDictEnumModel) CheckDuplicateCode(ctx context.Context, dictId int64, codes []string) (bool, error) {
	if len(codes) == 0 {
		return false, nil
	}

	query, args, err := sqlx.In(`SELECT COUNT(*) FROM `+TableNameDictEnum+` WHERE f_dict_id = ? AND f_code IN (?)`, dictId, codes)
	if err != nil {
		return false, err
	}

	query = m.conn.Rebind(query)
	var count int
	err = m.conn.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// ============================================
// RelationDictFileModel 实现
// ============================================

func (m *defaultRelationDictFileModel) InsertBatch(ctx context.Context, data []*RelationDictFile) error {
	if len(data) == 0 {
		return nil
	}

	query := `INSERT INTO ` + TableNameRelationDictFile + ` (f_id, f_dict_id, f_file_id) VALUES `
	values := make([]string, len(data))
	args := make([]interface{}, 0, len(data)*3)

	for i, item := range data {
		values[i] = `(?, ?, ?)`
		args = append(args, item.Id, item.DictId, item.FileId)
	}

	query += strings.Join(values, `, `)

	_, err := m.conn.ExecContext(ctx, query, args...)
	return err
}

func (m *defaultRelationDictFileModel) DeleteByDictId(ctx context.Context, dictId int64) error {
	query := `DELETE FROM ` + TableNameRelationDictFile + ` WHERE f_dict_id = ?`
	_, err := m.conn.ExecContext(ctx, query, dictId)
	return err
}

func (m *defaultRelationDictFileModel) FindByDictId(ctx context.Context, dictId int64) ([]int64, error) {
	query := `SELECT f_file_id FROM ` + TableNameRelationDictFile + ` WHERE f_dict_id = ?`
	rows, err := m.conn.QueryContext(ctx, query, dictId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fileIds []int64
	for rows.Next() {
		var fileId int64
		err := rows.Scan(&fileId)
		if err != nil {
			return nil, err
		}
		fileIds = append(fileIds, fileId)
	}

	return fileIds, nil
}
