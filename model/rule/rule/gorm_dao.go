package rule

import (
	"context"
	"errors"
	"fmt"

	"github.com/jinguoxing/idrm-go-base/db"
	"gorm.io/gorm"
)

// Insert 插入规则
func (m *defaultRuleModel) Insert(ctx context.Context, data *Rule) (int64, error) {
	result := m.getDB().WithContext(ctx).Create(data)
	if result.Error != nil {
		return 0, result.Error
	}
	return data.Id, nil
}

// FindOne 查询单条规则
func (m *defaultRuleModel) FindOne(ctx context.Context, id int64) (*Rule, error) {
	var rule Rule
	err := m.getDB().WithContext(ctx).Where("f_id = ? AND f_deleted = 0", id).First(&rule).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// Update 更新规则
func (m *defaultRuleModel) Update(ctx context.Context, data *Rule) error {
	result := m.getDB().WithContext(ctx).Save(data)
	return result.Error
}

// Delete 删除规则（逻辑删除）
func (m *defaultRuleModel) Delete(ctx context.Context, id int64) error {
	result := m.getDB().WithContext(ctx).Model(&Rule{}).Where("f_id = ?", id).
		Update("f_deleted", id).Error
	return result.Error
}

// FindByIds 批量查询规则
func (m *defaultRuleModel) FindByIds(ctx context.Context, ids []int64) ([]*Rule, error) {
	var rules []*Rule
	err := m.getDB().WithContext(ctx).Where("f_id IN ?", ids).Where("f_deleted = 0").
		Order("f_id DESC").Find(&rules).Error
	return rules, err
}

// FindByNameAndOrgType 按名称和组织类型查询（唯一性校验）
func (m *defaultRuleModel) FindByNameAndOrgType(ctx context.Context, name string, orgType int32, departmentIds string) ([]*Rule, error) {
	var rules []*Rule
	query := m.getDB().WithContext(ctx).Where("f_name = ? AND f_org_type = ? AND f_deleted = 0", name, orgType)
	if departmentIds != "" {
		query = query.Where("f_department_ids = ?", departmentIds)
	}
	err := query.Find(&rules).Error
	return rules, err
}

// FindByCatalogIds 按目录ID列表查询（分页）
func (m *defaultRuleModel) FindByCatalogIds(ctx context.Context, catalogIds []int64, offset int, limit int) ([]*Rule, int64, error) {
	var rules []*Rule
	var total int64

	// 查询总数
	countQuery := m.getDB().WithContext(ctx).Model(&Rule{}).Where("f_catalog_id IN ?", catalogIds).Where("f_deleted = 0")
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询列表
	query := m.getDB().WithContext(ctx).Where("f_catalog_id IN ?", catalogIds).Where("f_deleted = 0").
		Order("f_id DESC")
	if limit > 0 {
		query = query.Offset(offset - 1).Limit(limit)
	}
	err := query.Find(&rules).Error
	return rules, total, err
}

// FindDataExists 检查数据是否存在（用于重复校验）
func (m *defaultRuleModel) FindDataExists(ctx context.Context, filterId int64, name string, departmentIds string) (bool, error) {
	var count int64
	query := m.getDB().WithContext(ctx).Model(&Rule{}).Where("f_name = ? AND f_deleted = 0", name)
	if departmentIds != "" {
		query = query.Where("f_department_ids = ?", departmentIds)
	}
	if filterId > 0 {
		query = query.Where("f_id != ?", filterId)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// UpdateState 更新状态
func (m *defaultRuleModel) UpdateState(ctx context.Context, id int64, state int32, disableReason string) error {
	updates := map[string]interface{}{
		"f_state": state,
	}
	if state == 1 {
		// 启用时清空停用原因
		updates["f_disable_reason"] = ""
	} else {
		updates["f_disable_reason"] = disableReason
	}

	result := m.getDB().WithContext(ctx).Model(&Rule{}).Where("f_id = ?", id).Updates(updates)
	return result.Error
}

// RemoveCatalog 目录移动
func (m *defaultRuleModel) RemoveCatalog(ctx context.Context, ids []int64, catalogId int64, updateUser string) error {
	// 批量更新 catalog_id，版本号+1
	result := m.getDB().WithContext(ctx).Model(&Rule{}).Where("f_id IN ?", ids).Updates(map[string]interface{}{
		"f_catalog_id":  catalogId,
		"f_version":     gorm.Expr("f_version + 1"),
		"f_update_user": updateUser,
		"f_update_time":  gorm.Expr("NOW()"),
	})
	return result.Error
}

// UpdateVersionByIds 批量更新版本号
func (m *defaultRuleModel) UpdateVersionByIds(ctx context.Context, ids []int64) error {
	result := m.getDB().WithContext(ctx).Model(&Rule{}).Where("f_id IN ?", ids).Updates(map[string]interface{}{
		"f_version":     gorm.Expr("f_version + 1"),
		"f_update_time":  gorm.Expr("NOW()"),
	})
	return result.Error
}
