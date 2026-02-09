// Code scaffolded by speckit. Safe to edit.

package dataelement

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

var _ DataElementModel = (*defaultDataElementModel)(nil)

// Insert 插入数据元
func (m *defaultDataElementModel) Insert(ctx context.Context, data *DataElement) (int64, error) {
	query := `INSERT INTO ` + TableNameDataElementInfo + `
		(f_de_id, f_name_en, f_name_cn, f_synonym, f_std_type, f_data_type, f_data_length, f_data_precision,
		 f_dict_code, f_rule_id, f_relation_type, f_catalog_id, f_label_id, f_description,
		 f_version, f_state, f_department_ids, f_third_dept_id, f_create_time, f_create_user, f_deleted)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 1, ?, ?, NOW(), ?, 0)`

	result, err := m.conn.ExecContext(ctx, query,
		data.Code, data.NameEn, data.NameCn, data.Synonym, data.StdType, data.DataType,
		data.DataLength, data.DataPrecision, data.DictCode, data.RuleId, data.RelationType,
		data.CatalogId, data.LabelId, data.Description, data.State,
		data.DepartmentIds, data.ThirdDeptId, data.CreateUser)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// FindOne 查询单个数据元
func (m *defaultDataElementModel) FindOne(ctx context.Context, id int64) (*DataElement, error) {
	query := `SELECT f_id, f_de_id, f_name_en, f_name_cn, f_synonym, f_std_type, f_data_type,
		f_data_length, f_data_precision, f_dict_code, f_rule_id, f_relation_type,
		f_catalog_id, f_label_id, f_description, f_version, f_state,
		f_authority_id, f_department_ids, f_third_dept_id, f_disable_reason,
		f_create_time, f_create_user, f_update_time, f_update_user, f_deleted
		FROM ` + TableNameDataElementInfo + `
		WHERE f_id = ? AND f_deleted = 0`

	var de DataElement
	err := m.conn.QueryRowContext(ctx, query, id).Scan(
		&de.Id, &de.Code, &de.NameEn, &de.NameCn, &de.Synonym, &de.StdType, &de.DataType,
		&de.DataLength, &de.DataPrecision, &de.DictCode, &de.RuleId, &de.RelationType,
		&de.CatalogId, &de.LabelId, &de.Description, &de.Version, &de.State,
		&de.AuthorityId, &de.DepartmentIds, &de.ThirdDeptId, &de.DisableReason,
		&de.CreateTime, &de.CreateUser, &de.UpdateTime, &de.UpdateUser, &de.Deleted,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &de, nil
}

// FindOneByCode 按Code查询数据元
func (m *defaultDataElementModel) FindOneByCode(ctx context.Context, code int64) (*DataElement, error) {
	query := `SELECT f_id, f_de_id, f_name_en, f_name_cn, f_synonym, f_std_type, f_data_type,
		f_data_length, f_data_precision, f_dict_code, f_rule_id, f_relation_type,
		f_catalog_id, f_label_id, f_description, f_version, f_state,
		f_authority_id, f_department_ids, f_third_dept_id, f_disable_reason,
		f_create_time, f_create_user, f_update_time, f_update_user, f_deleted
		FROM ` + TableNameDataElementInfo + `
		WHERE f_de_id = ? AND f_deleted = 0`

	var de DataElement
	err := m.conn.QueryRowContext(ctx, query, code).Scan(
		&de.Id, &de.Code, &de.NameEn, &de.NameCn, &de.Synonym, &de.StdType, &de.DataType,
		&de.DataLength, &de.DataPrecision, &de.DictCode, &de.RuleId, &de.RelationType,
		&de.CatalogId, &de.LabelId, &de.Description, &de.Version, &de.State,
		&de.AuthorityId, &de.DepartmentIds, &de.ThirdDeptId, &de.DisableReason,
		&de.CreateTime, &de.CreateUser, &de.UpdateTime, &de.UpdateUser, &de.Deleted,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &de, nil
}

// Update 更新数据元
func (m *defaultDataElementModel) Update(ctx context.Context, data *DataElement) error {
	query := `UPDATE ` + TableNameDataElementInfo + `
		SET f_name_en = ?, f_name_cn = ?, f_synonym = ?, f_std_type = ?, f_data_type = ?,
			f_data_length = ?, f_data_precision = ?, f_dict_code = ?, f_rule_id = ?,
			f_relation_type = ?, f_catalog_id = ?, f_label_id = ?, f_description = ?,
			f_version = ?, f_state = ?, f_department_ids = ?, f_third_dept_id = ?,
			f_update_time = NOW(), f_update_user = ?
		WHERE f_id = ? AND f_deleted = 0`

	_, err := m.conn.ExecContext(ctx, query,
		data.NameEn, data.NameCn, data.Synonym, data.StdType, data.DataType,
		data.DataLength, data.DataPrecision, data.DictCode, data.RuleId,
		data.RelationType, data.CatalogId, data.LabelId, data.Description,
		data.Version, data.State, data.DepartmentIds, data.ThirdDeptId,
		data.UpdateUser, data.Id)

	return err
}

// Delete 删除数据元（逻辑删除）
func (m *defaultDataElementModel) Delete(ctx context.Context, id int64) error {
	query := `UPDATE ` + TableNameDataElementInfo + `
		SET f_deleted = ?
		WHERE f_id = ?`

	_, err := m.conn.ExecContext(ctx, query, id, id)
	return err
}

// FindByCatalogIds 按目录ID列表查询数据元
func (m *defaultDataElementModel) FindByCatalogIds(ctx context.Context, catalogIds []int64, opts *FindOptions) ([]*DataElement, int64, error) {
	if len(catalogIds) == 0 {
		return []*DataElement{}, 0, nil
	}

	// 构建IN子句
	placeholders := make([]string, len(catalogIds))
	args := make([]interface{}, len(catalogIds))
	for i, id := range catalogIds {
		placeholders[i] = "?"
		args[i] = id
	}

	// 构建WHERE条件
	where := []string{"f_deleted = 0", fmt.Sprintf("f_catalog_id IN (%s)", strings.Join(placeholders, ","))}

	if opts.State != nil {
		where = append(where, "f_state = ?")
		args = append(args, *opts.State)
	}
	if opts.StdType != nil {
		where = append(where, "f_std_type = ?")
		args = append(args, *opts.StdType)
	}
	if opts.DataType != nil {
		where = append(where, "f_data_type = ?")
		args = append(args, *opts.DataType)
	}
	if opts.RelationType != nil {
		where = append(where, "f_relation_type = ?")
		args = append(args, *opts.RelationType)
	}
	if opts.Keyword != "" {
		where = append(where, "(f_name_cn LIKE ? OR f_name_en LIKE ? OR f_synonym LIKE ?)")
		keyword := "%" + opts.Keyword + "%"
		args = append(args, keyword, keyword, keyword)
	}

	whereClause := strings.Join(where, " AND ")

	// 查询总数
	countQuery := `SELECT COUNT(*) FROM ` + TableNameDataElementInfo + ` WHERE ` + whereClause
	var totalCount int64
	countArgs := append([]interface{}{}, args...)
	err := m.conn.QueryRowContext(ctx, countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// 查询数据
	offset := BuildOffset(opts.Page, opts.PageSize)
	limit := BuildLimit(opts.PageSize)

	orderBy := "f_create_time"
	if opts.Sort != "" {
		orderBy = fmt.Sprintf("f_%s", strings.TrimPrefix(opts.Sort, "f_"))
	}
	direction := "DESC"
	if opts.Direction == "asc" {
		direction = "ASC"
	}

	dataQuery := `SELECT f_id, f_de_id, f_name_en, f_name_cn, f_synonym, f_std_type, f_data_type,
		f_data_length, f_data_precision, f_dict_code, f_rule_id, f_relation_type,
		f_catalog_id, f_label_id, f_description, f_version, f_state,
		f_authority_id, f_department_ids, f_third_dept_id, f_disable_reason,
		f_create_time, f_create_user, f_update_time, f_update_user, f_deleted
		FROM ` + TableNameDataElementInfo + `
		WHERE ` + whereClause + `
		ORDER BY ` + orderBy + ` ` + direction + `
		LIMIT ? OFFSET ?`

	args = append(args, limit, offset)

	rows, err := m.conn.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var elements []*DataElement
	for rows.Next() {
		var de DataElement
		err := rows.Scan(
			&de.Id, &de.Code, &de.NameEn, &de.NameCn, &de.Synonym, &de.StdType, &de.DataType,
			&de.DataLength, &de.DataPrecision, &de.DictCode, &de.RuleId, &de.RelationType,
			&de.CatalogId, &de.LabelId, &de.Description, &de.Version, &de.State,
			&de.AuthorityId, &de.DepartmentIds, &de.ThirdDeptId, &de.DisableReason,
			&de.CreateTime, &de.CreateUser, &de.UpdateTime, &de.UpdateUser, &de.Deleted,
		)
		if err != nil {
			return nil, 0, err
		}
		elements = append(elements, &de)
	}

	return elements, totalCount, nil
}

// FindByIds 按ID列表查询数据元
func (m *defaultDataElementModel) FindByIds(ctx context.Context, ids []int64) ([]*DataElement, error) {
	if len(ids) == 0 {
		return []*DataElement{}, nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := `SELECT f_id, f_de_id, f_name_en, f_name_cn, f_synonym, f_std_type, f_data_type,
		f_data_length, f_data_precision, f_dict_code, f_rule_id, f_relation_type,
		f_catalog_id, f_label_id, f_description, f_version, f_state,
		f_authority_id, f_department_ids, f_third_dept_id, f_disable_reason,
		f_create_time, f_create_user, f_update_time, f_update_user, f_deleted
		FROM ` + TableNameDataElementInfo + `
		WHERE f_id IN (` + strings.Join(placeholders, ",") + `) AND f_deleted = 0
		ORDER BY f_id`

	rows, err := m.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var elements []*DataElement
	for rows.Next() {
		var de DataElement
		err := rows.Scan(
			&de.Id, &de.Code, &de.NameEn, &de.NameCn, &de.Synonym, &de.StdType, &de.DataType,
			&de.DataLength, &de.DataPrecision, &de.DictCode, &de.RuleId, &de.RelationType,
			&de.CatalogId, &de.LabelId, &de.Description, &de.Version, &de.State,
			&de.AuthorityId, &de.DepartmentIds, &de.ThirdDeptId, &de.DisableReason,
			&de.CreateTime, &de.CreateUser, &de.UpdateTime, &de.UpdateUser, &de.Deleted,
		)
		if err != nil {
			return nil, err
		}
		elements = append(elements, &de)
	}

	return elements, nil
}

// FindByRuleId 按规则ID分页查询数据元
func (m *defaultDataElementModel) FindByRuleId(ctx context.Context, ruleId int64, opts *FindOptions) ([]*DataElement, int64, error) {
	if ruleId <= 0 {
		return []*DataElement{}, 0, nil
	}

	// 构建WHERE条件
	whereClause := "WHERE f_rule_id = ? AND f_deleted = 0"
	args := []interface{}{ruleId}

	// 添加关键词搜索
	if opts.Keyword != "" {
		whereClause += " AND (f_name_cn LIKE ? OR f_name_en LIKE ? OR f_de_id LIKE ?)"
		keyword := "%" + opts.Keyword + "%"
		args = append(args, keyword, keyword, keyword)
	}

	// 查询总数
	countQuery := "SELECT COUNT(*) FROM " + TableNameDataElementInfo + " " + whereClause
	var totalCount int64
	err := m.conn.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	if totalCount == 0 {
		return []*DataElement{}, 0, nil
	}

	// 构建排序
	orderBy := "ORDER BY f_id"
	if opts.Sort != "" {
		// 安全的排序字段
		validSorts := map[string]string{
			"f_id":          "f_id",
			"f_de_id":       "f_de_id",
			"f_name_cn":     "f_name_cn",
			"f_name_en":     "f_name_en",
			"f_create_time": "f_create_time",
			"f_update_time": "f_update_time",
		}
		if sortField, ok := validSorts[opts.Sort]; ok {
			direction := "DESC"
			if opts.Direction == "asc" {
				direction = "ASC"
			}
			orderBy = "ORDER BY " + sortField + " " + direction
		}
	}

	// 构建分页
	offset := BuildOffset(opts.Page, opts.PageSize)
	limit := BuildLimit(opts.PageSize)

	// 查询数据
	query := `SELECT f_id, f_de_id, f_name_en, f_name_cn, f_synonym, f_std_type, f_data_type,
		f_data_length, f_data_precision, f_dict_code, f_rule_id, f_relation_type,
		f_catalog_id, f_label_id, f_description, f_version, f_state,
		f_authority_id, f_department_ids, f_third_dept_id, f_disable_reason,
		f_create_time, f_create_user, f_update_time, f_update_user, f_deleted
		FROM ` + TableNameDataElementInfo + ` ` + whereClause + ` ` + orderBy + ` LIMIT ? OFFSET ?`

	args = append(args, limit, offset)
	rows, err := m.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var elements []*DataElement
	for rows.Next() {
		var de DataElement
		err := rows.Scan(
			&de.Id, &de.Code, &de.NameEn, &de.NameCn, &de.Synonym, &de.StdType, &de.DataType,
			&de.DataLength, &de.DataPrecision, &de.DictCode, &de.RuleId, &de.RelationType,
			&de.CatalogId, &de.LabelId, &de.Description, &de.Version, &de.State,
			&de.AuthorityId, &de.DepartmentIds, &de.ThirdDeptId, &de.DisableReason,
			&de.CreateTime, &de.CreateUser, &de.UpdateTime, &de.UpdateUser, &de.Deleted,
		)
		if err != nil {
			return nil, 0, err
		}
		elements = append(elements, &de)
	}

	return elements, totalCount, nil
}

// FindByCodes 按Code列表查询数据元
func (m *defaultDataElementModel) FindByCodes(ctx context.Context, codes []int64) ([]*DataElement, error) {
	if len(codes) == 0 {
		return []*DataElement{}, nil
	}

	placeholders := make([]string, len(codes))
	args := make([]interface{}, len(codes))
	for i, code := range codes {
		placeholders[i] = "?"
		args[i] = code
	}

	query := `SELECT f_id, f_de_id, f_name_en, f_name_cn, f_synonym, f_std_type, f_data_type,
		f_data_length, f_data_precision, f_dict_code, f_rule_id, f_relation_type,
		f_catalog_id, f_label_id, f_description, f_version, f_state,
		f_authority_id, f_department_ids, f_third_dept_id, f_disable_reason,
		f_create_time, f_create_user, f_update_time, f_update_user, f_deleted
		FROM ` + TableNameDataElementInfo + `
		WHERE f_de_id IN (` + strings.Join(placeholders, ",") + `) AND f_deleted = 0
		ORDER BY f_id`

	rows, err := m.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var elements []*DataElement
	for rows.Next() {
		var de DataElement
		err := rows.Scan(
			&de.Id, &de.Code, &de.NameEn, &de.NameCn, &de.Synonym, &de.StdType, &de.DataType,
			&de.DataLength, &de.DataPrecision, &de.DictCode, &de.RuleId, &de.RelationType,
			&de.CatalogId, &de.LabelId, &de.Description, &de.Version, &de.State,
			&de.AuthorityId, &de.DepartmentIds, &de.ThirdDeptId, &de.DisableReason,
			&de.CreateTime, &de.CreateUser, &de.UpdateTime, &de.UpdateUser, &de.Deleted,
		)
		if err != nil {
			return nil, err
		}
		elements = append(elements, &de)
	}

	return elements, nil
}

// FindByFileCatalog 按文件目录查询数据元
func (m *defaultDataElementModel) FindByFileCatalog(ctx context.Context, opts *FindOptions) ([]*DataElement, int64, error) {
	// TODO: 实现按文件目录查询（需要关联关系表）
	// 先返回空列表
	return []*DataElement{}, 0, nil
}

// FindByFileId 按文件ID查询数据元
func (m *defaultDataElementModel) FindByFileId(ctx context.Context, fileId int64, opts *FindOptions) ([]*DataElement, int64, error) {
	// TODO: 实现按文件ID查询（需要关联关系表）
	// 先返回空列表
	return []*DataElement{}, 0, nil
}

// CheckNameCnExists 检查中文名称是否存在
func (m *defaultDataElementModel) CheckNameCnExists(ctx context.Context, nameCn string, stdType int32, excludeId int64, deptIds string) (bool, error) {
	query := `SELECT COUNT(*) FROM ` + TableNameDataElementInfo + `
		WHERE f_name_cn = ? AND f_std_type = ? AND f_deleted = 0`
	args := []interface{}{nameCn, stdType}

	if excludeId > 0 {
		query += ` AND f_id != ?`
		args = append(args, excludeId)
	}

	var count int
	err := m.conn.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CheckNameEnExists 检查英文名称是否存在
func (m *defaultDataElementModel) CheckNameEnExists(ctx context.Context, nameEn string, excludeId int64, deptIds string) (bool, error) {
	query := `SELECT COUNT(*) FROM ` + TableNameDataElementInfo + `
		WHERE f_name_en = ? AND f_deleted = 0`
	args := []interface{}{nameEn}

	if excludeId > 0 {
		query += ` AND f_id != ?`
		args = append(args, excludeId)
	}

	var count int
	err := m.conn.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindDataExists 查找数据是否存在
func (m *defaultDataElementModel) FindDataExists(ctx context.Context, name string, stdType int32, excludeId int64, deptIds string) (*DataElement, error) {
	query := `SELECT f_id, f_de_id, f_name_en, f_name_cn FROM ` + TableNameDataElementInfo + `
		WHERE f_name_cn = ? AND f_std_type = ? AND f_deleted = 0`
	args := []interface{}{name, stdType}

	if excludeId > 0 {
		query += ` AND f_id != ?`
		args = append(args, excludeId)
	}

	var de DataElement
	err := m.conn.QueryRowContext(ctx, query, args...).Scan(&de.Id, &de.Code, &de.NameEn, &de.NameCn)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &de, nil
}

// UpdateState 更新状态
func (m *defaultDataElementModel) UpdateState(ctx context.Context, ids []int64, state int32, reason, updateUser string) error {
	if len(ids) == 0 {
		return nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := `UPDATE ` + TableNameDataElementInfo + `
		SET f_state = ?, f_disable_reason = ?, f_update_time = NOW(), f_update_user = ?
		WHERE f_id IN (` + strings.Join(placeholders, ",") + `) AND f_deleted = 0`

	args = append([]interface{}{state, reason, updateUser}, args...)

	_, err := m.conn.ExecContext(ctx, query, args...)
	return err
}

// MoveCatalog 移动目录
func (m *defaultDataElementModel) MoveCatalog(ctx context.Context, ids []int64, catalogId int64, updateUser string) error {
	if len(ids) == 0 {
		return nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	// 移动目录时版本号+1
	query := `UPDATE ` + TableNameDataElementInfo + `
		SET f_catalog_id = ?, f_version = f_version + 1, f_update_time = NOW(), f_update_user = ?
		WHERE f_id IN (` + strings.Join(placeholders, ",") + `) AND f_deleted = 0`

	args = append([]interface{}{catalogId, updateUser}, args...)

	_, err := m.conn.ExecContext(ctx, query, args...)
	return err
}

// DeleteLabelIds 删除标签
func (m *defaultDataElementModel) DeleteLabelIds(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := `UPDATE ` + TableNameDataElementInfo + `
		SET f_label_id = NULL, f_update_time = NOW()
		WHERE f_id IN (` + strings.Join(placeholders, ",") + `) AND f_deleted = 0`

	_, err := m.conn.ExecContext(ctx, query, args...)
	return err
}

// DeleteByIds 批量删除
func (m *defaultDataElementModel) DeleteByIds(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := `UPDATE ` + TableNameDataElementInfo + `
		SET f_deleted = ?, f_update_time = NOW()
		WHERE f_id IN (` + strings.Join(placeholders, ",") + `)`

	args = append([]interface{}{time.Now().Unix()}, args...)

	_, err := m.conn.ExecContext(ctx, query, args...)
	return err
}

// CountByCatalogId 按目录ID统计数据元数量
func (m *defaultDataElementModel) CountByCatalogId(ctx context.Context, catalogId int64) (int64, error) {
	query := `SELECT COUNT(*) FROM ` + TableNameDataElementInfo + `
		WHERE f_catalog_id = ? AND f_deleted = 0`

	var count int64
	err := m.conn.QueryRowContext(ctx, query, catalogId).Scan(&count)
	return count, err
}

// IncrementVersion 递增版本号
func (m *defaultDataElementModel) IncrementVersion(ctx context.Context, id int64) error {
	query := `UPDATE ` + TableNameDataElementInfo + `
		SET f_version = f_version + 1, f_update_time = NOW()
		WHERE f_id = ? AND f_deleted = 0`

	_, err := m.conn.ExecContext(ctx, query, id)
	return err
}
