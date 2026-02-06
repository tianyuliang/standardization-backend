package catalog

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

// NewCatalogModel 创建目录模型实例
func NewCatalogModel(conn *sqlx.Conn) CatalogModel {
	return &customCatalogModel{conn: conn}
}

// customCatalogModel 自定义目录模型实现
type customCatalogModel struct {
	conn *sqlx.Conn
}

// tableName 表名
const tableName = "`t_de_catalog_info`"

var _ CatalogModel = (*customCatalogModel)(nil)

// Insert 插入目录
func (m *customCatalogModel) Insert(ctx context.Context, data *Catalog) (int64, error) {
	query := fmt.Sprintf("INSERT INTO %s (f_catalog_name, f_description, f_level, f_parent_id, f_type, f_authority_id, f_create_time, f_create_user, f_update_time, f_update_user, f_deleted) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", tableName)
	result, err := m.conn.ExecContext(ctx, query, data.CatalogName, data.Description, data.Level, data.ParentId, data.Type, data.AuthorityId, data.CreateTime, data.CreateUser, data.UpdateTime, data.UpdateUser, data.Deleted)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// FindOne 查询单个目录
func (m *customCatalogModel) FindOne(ctx context.Context, id int64) (*Catalog, error) {
	query := fmt.Sprintf("SELECT f_id, f_catalog_name, f_description, f_level, f_parent_id, f_type, f_authority_id, f_create_time, f_create_user, f_update_time, f_update_user, f_deleted FROM %s WHERE f_id = ? AND f_deleted = 0 LIMIT 1", tableName)
	var resp Catalog
	err := m.conn.QueryRowContext(ctx, query, id).Scan(&resp.Id, &resp.CatalogName, &resp.Description, &resp.Level, &resp.ParentId, &resp.Type, &resp.AuthorityId, &resp.CreateTime, &resp.CreateUser, &resp.UpdateTime, &resp.UpdateUser, &resp.Deleted)
	switch err {
	case nil:
		return &resp, nil
	case sql.ErrNoRows:
		return nil, sql.ErrNoRows
	default:
		return nil, err
	}
}

// Update 更新目录
func (m *customCatalogModel) Update(ctx context.Context, data *Catalog) error {
	query := fmt.Sprintf("UPDATE %s SET f_catalog_name = ?, f_description = ?, f_level = ?, f_parent_id = ?, f_type = ?, f_authority_id = ?, f_update_time = ?, f_update_user = ? WHERE f_id = ?", tableName)
	_, err := m.conn.ExecContext(ctx, query, data.CatalogName, data.Description, data.Level, data.ParentId, data.Type, data.AuthorityId, data.UpdateTime, data.UpdateUser, data.Id)
	return err
}

// Delete 删除目录（逻辑删除）
func (m *customCatalogModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("UPDATE %s SET f_deleted = f_id + 1 WHERE f_id = ?", tableName)
	_, err := m.conn.ExecContext(ctx, query, id)
	return err
}

// FindByType 按类型查询目录
func (m *customCatalogModel) FindByType(ctx context.Context, catalogType int32) ([]*Catalog, error) {
	query := fmt.Sprintf("SELECT f_id, f_catalog_name, f_description, f_level, f_parent_id, f_type, f_authority_id, f_create_time, f_create_user, f_update_time, f_update_user, f_deleted FROM %s WHERE f_type = ? AND f_deleted = 0 ORDER BY f_level, f_id", tableName)
	rows, err := m.conn.QueryContext(ctx, query, catalogType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Catalog
	for rows.Next() {
		var item Catalog
		err = rows.Scan(&item.Id, &item.CatalogName, &item.Description, &item.Level, &item.ParentId, &item.Type, &item.AuthorityId, &item.CreateTime, &item.CreateUser, &item.UpdateTime, &item.UpdateUser, &item.Deleted)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}

// FindByTypeAndLevel 按类型和最小级别查询目录
func (m *customCatalogModel) FindByTypeAndLevel(ctx context.Context, catalogType int32, minLevel int32) ([]*Catalog, error) {
	query := fmt.Sprintf("SELECT f_id, f_catalog_name, f_description, f_level, f_parent_id, f_type, f_authority_id, f_create_time, f_create_user, f_update_time, f_update_user, f_deleted FROM %s WHERE f_type = ? AND f_level >= ? AND f_deleted = 0 ORDER BY f_level, f_id", tableName)
	rows, err := m.conn.QueryContext(ctx, query, catalogType, minLevel)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Catalog
	for rows.Next() {
		var item Catalog
		err = rows.Scan(&item.Id, &item.CatalogName, &item.Description, &item.Level, &item.ParentId, &item.Type, &item.AuthorityId, &item.CreateTime, &item.CreateUser, &item.UpdateTime, &item.UpdateUser, &item.Deleted)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}

// FindByName 按名称模糊查询目录
func (m *customCatalogModel) FindByName(ctx context.Context, name string, catalogType int32) ([]*Catalog, error) {
	keyword := strings.ToLower(strings.TrimSpace(name))
	if keyword == "" {
		return []*Catalog{}, nil
	}
	if !strings.HasPrefix(keyword, "%") {
		keyword = "%" + keyword
	}
	if !strings.HasSuffix(keyword, "%") {
		keyword = keyword + "%"
	}

	query := fmt.Sprintf("SELECT f_id, f_catalog_name, f_description, f_level, f_parent_id, f_type, f_authority_id, f_create_time, f_create_user, f_update_time, f_update_user, f_deleted FROM %s WHERE f_type = ? AND f_level > 1 AND LOWER(f_catalog_name) LIKE ? AND f_deleted = 0 ORDER BY f_level, f_id", tableName)
	rows, err := m.conn.QueryContext(ctx, query, catalogType, keyword)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Catalog
	for rows.Next() {
		var item Catalog
		err = rows.Scan(&item.Id, &item.CatalogName, &item.Description, &item.Level, &item.ParentId, &item.Type, &item.AuthorityId, &item.CreateTime, &item.CreateUser, &item.UpdateTime, &item.UpdateUser, &item.Deleted)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}

// FindByParentId 按父目录ID查询
func (m *customCatalogModel) FindByParentId(ctx context.Context, parentId int64) ([]*Catalog, error) {
	query := fmt.Sprintf("SELECT f_id, f_catalog_name, f_description, f_level, f_parent_id, f_type, f_authority_id, f_create_time, f_create_user, f_update_time, f_update_user, f_deleted FROM %s WHERE f_parent_id = ? AND f_deleted = 0 ORDER BY f_catalog_name", tableName)
	rows, err := m.conn.QueryContext(ctx, query, parentId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Catalog
	for rows.Next() {
		var item Catalog
		err = rows.Scan(&item.Id, &item.CatalogName, &item.Description, &item.Level, &item.ParentId, &item.Type, &item.AuthorityId, &item.CreateTime, &item.CreateUser, &item.UpdateTime, &item.UpdateUser, &item.Deleted)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}

// FindByIds 批量查询目录
func (m *customCatalogModel) FindByIds(ctx context.Context, ids []int64) ([]*Catalog, error) {
	if len(ids) == 0 {
		return []*Catalog{}, nil
	}

	// 构建IN查询
	placeholders := strings.Repeat("?,", len(ids))
	placeholders = placeholders[:len(placeholders)-1]

	query := fmt.Sprintf("SELECT f_id, f_catalog_name, f_description, f_level, f_parent_id, f_type, f_authority_id, f_create_time, f_create_user, f_update_time, f_update_user, f_deleted FROM %s WHERE f_id IN (%s) AND f_deleted = 0 ORDER BY f_level, f_id", tableName, placeholders)

	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	rows, err := m.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Catalog
	for rows.Next() {
		var item Catalog
		err = rows.Scan(&item.Id, &item.CatalogName, &item.Description, &item.Level, &item.ParentId, &item.Type, &item.AuthorityId, &item.CreateTime, &item.CreateUser, &item.UpdateTime, &item.UpdateUser, &item.Deleted)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}

// FindAllByTypeAndLevel 按类型和最大级别查询目录
func (m *customCatalogModel) FindAllByTypeAndLevel(ctx context.Context, catalogType int32, maxLevel int32) ([]*Catalog, error) {
	query := fmt.Sprintf("SELECT f_id, f_catalog_name, f_description, f_level, f_parent_id, f_type, f_authority_id, f_create_time, f_create_user, f_update_time, f_update_user, f_deleted FROM %s WHERE f_type = ? AND f_level <= ? AND f_deleted = 0 ORDER BY f_level, f_id", tableName)
	rows, err := m.conn.QueryContext(ctx, query, catalogType, maxLevel)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Catalog
	for rows.Next() {
		var item Catalog
		err = rows.Scan(&item.Id, &item.CatalogName, &item.Description, &item.Level, &item.ParentId, &item.Type, &item.AuthorityId, &item.CreateTime, &item.CreateUser, &item.UpdateTime, &item.UpdateUser, &item.Deleted)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}

// FindTree 查询目录树（按类型和根级别）
func (m *customCatalogModel) FindTree(ctx context.Context, catalogType int32, rootLevel int32) ([]*Catalog, error) {
	// 查询所有该类型的目录（level >= rootLevel）
	query := fmt.Sprintf("SELECT f_id, f_catalog_name, f_description, f_level, f_parent_id, f_type, f_authority_id, f_create_time, f_create_user, f_update_time, f_update_user, f_deleted FROM %s WHERE f_type = ? AND f_level >= ? AND f_deleted = 0 ORDER BY f_level, f_id", tableName)
	rows, err := m.conn.QueryContext(ctx, query, catalogType, rootLevel)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Catalog
	for rows.Next() {
		var item Catalog
		err = rows.Scan(&item.Id, &item.CatalogName, &item.Description, &item.Level, &item.ParentId, &item.Type, &item.AuthorityId, &item.CreateTime, &item.CreateUser, &item.UpdateTime, &item.UpdateUser, &item.Deleted)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}

// FindChildren 查询子目录
func (m *customCatalogModel) FindChildren(ctx context.Context, parentId int64, catalogType int32) ([]*Catalog, error) {
	query := fmt.Sprintf("SELECT f_id, f_catalog_name, f_description, f_level, f_parent_id, f_type, f_authority_id, f_create_time, f_create_user, f_update_time, f_update_user, f_deleted FROM %s WHERE f_parent_id = ? AND f_type = ? AND f_deleted = 0 ORDER BY f_catalog_name", tableName)
	rows, err := m.conn.QueryContext(ctx, query, parentId, catalogType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Catalog
	for rows.Next() {
		var item Catalog
		err = rows.Scan(&item.Id, &item.CatalogName, &item.Description, &item.Level, &item.ParentId, &item.Type, &item.AuthorityId, &item.CreateTime, &item.CreateUser, &item.UpdateTime, &item.UpdateUser, &item.Deleted)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}

// DeleteByIds 批量删除目录（逻辑删除）
func (m *customCatalogModel) DeleteByIds(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}

	// 构建IN查询
	placeholders := strings.Repeat("?,", len(ids))
	placeholders = placeholders[:len(placeholders)-1]

	query := fmt.Sprintf("UPDATE %s SET f_deleted = f_id + 1 WHERE f_id IN (%s)", tableName, placeholders)

	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	_, err := m.conn.ExecContext(ctx, query, args...)
	return err
}
