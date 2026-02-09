// Code scaffolded by speckit. Safe to edit.

package relation

import (
	"context"
	"strings"
)

var _ RelationDeFileModel = (*defaultRelationDeFileModel)(nil)

// InsertBatch 批量插入关系
func (m *defaultRelationDeFileModel) InsertBatch(ctx context.Context, data []*RelationDeFile) error {
	if len(data) == 0 {
		return nil
	}

	query := `INSERT INTO ` + TableNameRelationDeFile + `
		(f_de_id, f_file_id, f_create_time)
		VALUES `

	placeholders := make([]string, len(data))
	args := make([]interface{}, 0, len(data)*2)

	for i, item := range data {
		placeholders[i] = "(?, ?, NOW())"
		args = append(args, item.DeId, item.FileId)
	}

	query += strings.Join(placeholders, ", ")

	_, err := m.conn.ExecContext(ctx, query, args...)
	return err
}

// DeleteByDeId 删除数据元的所有关系
func (m *defaultRelationDeFileModel) DeleteByDeId(ctx context.Context, deId int64) error {
	query := `DELETE FROM ` + TableNameRelationDeFile + ` WHERE f_de_id = ?`
	_, err := m.conn.ExecContext(ctx, query, deId)
	return err
}

// DeleteByDeIds 批量删除数据元的关系
func (m *defaultRelationDeFileModel) DeleteByDeIds(ctx context.Context, deIds []int64) error {
	if len(deIds) == 0 {
		return nil
	}

	placeholders := make([]string, len(deIds))
	args := make([]interface{}, len(deIds))
	for i, id := range deIds {
		placeholders[i] = "?"
		args[i] = id
	}

	query := `DELETE FROM ` + TableNameRelationDeFile + ` WHERE f_de_id IN (` + strings.Join(placeholders, ",") + `)`
	_, err := m.conn.ExecContext(ctx, query, args...)
	return err
}

// DeleteByFileId 删除文件的所有关系
func (m *defaultRelationDeFileModel) DeleteByFileId(ctx context.Context, fileId int64) error {
	query := `DELETE FROM ` + TableNameRelationDeFile + ` WHERE f_file_id = ?`
	_, err := m.conn.ExecContext(ctx, query, fileId)
	return err
}

// FindByDeId 查询数据元的关联文件列表
func (m *defaultRelationDeFileModel) FindByDeId(ctx context.Context, deId int64, opts *PageOptions) ([]*RelationDeFile, int64, error) {
	// 查询总数
	countQuery := `SELECT COUNT(*) FROM ` + TableNameRelationDeFile + ` WHERE f_de_id = ?`
	var totalCount int64
	err := m.conn.QueryRowContext(ctx, countQuery, deId).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	if totalCount == 0 {
		return []*RelationDeFile{}, 0, nil
	}

	// 查询数据
	offset := BuildOffset(opts.Page, opts.PageSize)
	limit := BuildLimit(opts.PageSize)

	query := `SELECT f_id, f_de_id, f_file_id, f_create_time
		FROM ` + TableNameRelationDeFile + `
		WHERE f_de_id = ?
		ORDER BY f_id
		LIMIT ? OFFSET ?`

	rows, err := m.conn.QueryContext(ctx, query, deId, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var relations []*RelationDeFile
	for rows.Next() {
		var rel RelationDeFile
		err := rows.Scan(&rel.Id, &rel.DeId, &rel.FileId, &rel.CreateTime)
		if err != nil {
			return nil, 0, err
		}
		relations = append(relations, &rel)
	}

	return relations, totalCount, nil
}

// FindByDeIds 批量查询数据元的关联文件
func (m *defaultRelationDeFileModel) FindByDeIds(ctx context.Context, deIds []int64) (map[int64][]int64, error) {
	if len(deIds) == 0 {
		return make(map[int64][]int64), nil
	}

	placeholders := make([]string, len(deIds))
	args := make([]interface{}, len(deIds))
	for i, id := range deIds {
		placeholders[i] = "?"
		args[i] = id
	}

	query := `SELECT f_de_id, f_file_id
		FROM ` + TableNameRelationDeFile + `
		WHERE f_de_id IN (` + strings.Join(placeholders, ",") + `)
		ORDER BY f_de_id, f_file_id`

	rows, err := m.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int64][]int64)
	for rows.Next() {
		var deId, fileId int64
		err := rows.Scan(&deId, &fileId)
		if err != nil {
			return nil, err
		}
		result[deId] = append(result[deId], fileId)
	}

	return result, nil
}

// CountByDeId 统计数据元的关联文件数量
func (m *defaultRelationDeFileModel) CountByDeId(ctx context.Context, deId int64) (int64, error) {
	query := `SELECT COUNT(*) FROM ` + TableNameRelationDeFile + ` WHERE f_de_id = ?`
	var count int64
	err := m.conn.QueryRowContext(ctx, query, deId).Scan(&count)
	return count, err
}

// FindFileIdsByDeId 查询数据元关联的文件ID列表
func (m *defaultRelationDeFileModel) FindFileIdsByDeId(ctx context.Context, deId int64) ([]int64, error) {
	query := `SELECT f_file_id
		FROM ` + TableNameRelationDeFile + `
		WHERE f_de_id = ?
		ORDER BY f_file_id`

	rows, err := m.conn.QueryContext(ctx, query, deId)
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
