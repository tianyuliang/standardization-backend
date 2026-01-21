package catalog

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type gormModel struct {
	db *gorm.DB
}

// newGormModel 创建 GORM 模型
func newGormModel(cfg ModelConfig) CatalogModel {
	return &gormModel{
		db: cfg.DB,
	}
}

// Insert 创建目录
func (m *gormModel) Insert(ctx context.Context, data *Catalog) (*Catalog, error) {
	db := m.db.WithContext(ctx).Create(data)
	if db.Error != nil {
		return nil, db.Error
	}
	return data, nil
}

// FindOne 根据 ID 查询单个目录
func (m *gormModel) FindOne(ctx context.Context, id string) (*Catalog, error) {
	var result Catalog
	db := m.db.WithContext(ctx).Where("f_id = ?", id).First(&result)
	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, db.Error
	}
	return &result, nil
}

// FindByParent 查询父目录下的所有子目录
func (m *gormModel) FindByParent(ctx context.Context, parentId string) ([]*Catalog, error) {
	var results []*Catalog
	db := m.db.WithContext(ctx).Where("f_parent_id = ?", parentId).Order("f_id ASC").Find(&results)
	if db.Error != nil {
		return nil, db.Error
	}
	return results, nil
}

// FindByType 根据类型查询所有目录
func (m *gormModel) FindByType(ctx context.Context, catalogType int32) ([]*Catalog, error) {
	var results []*Catalog
	db := m.db.WithContext(ctx).Where("f_type = ?", catalogType).Order("f_level ASC, f_id ASC").Find(&results)
	if db.Error != nil {
		return nil, db.Error
	}
	return results, nil
}

// FindByName 模糊查询目录名称（小写不敏感）
func (m *gormModel) FindByName(ctx context.Context, name string, catalogType int32) ([]*Catalog, error) {
	var results []*Catalog
	query := m.db.WithContext(ctx).Where("f_type = ?", catalogType).
		Where("f_level > ?", CatalogRootLevel). // 排除根目录
		Where("f_catalog_name LIKE ?", "%"+name+"%").
		Order("f_level ASC, f_id ASC").
		Find(&results)
	if query.Error != nil {
		return nil, query.Error
	}
	return results, nil
}

// FindByTypeAndName 根据类型和精确名称查询
func (m *gormModel) FindByTypeAndName(ctx context.Context, catalogType int32, name string) (*Catalog, error) {
	var result Catalog
	db := m.db.WithContext(ctx).
		Where("f_type = ?", catalogType).
		Where("f_catalog_name = ?", name).
		First(&result)
	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, db.Error
	}
	return &result, nil
}

// FindByTypeAndParent 查询指定类型和父目录下的子目录
func (m *gormModel) FindByTypeAndParent(ctx context.Context, catalogType int32, parentId string) ([]*Catalog, error) {
	var results []*Catalog
	db := m.db.WithContext(ctx).
		Where("f_type = ?", catalogType).
		Where("f_parent_id = ?", parentId).
		Order("f_id ASC").
		Find(&results)
	if db.Error != nil {
		return nil, db.Error
	}
	return results, nil
}

// Update 更新目录
func (m *gormModel) Update(ctx context.Context, data *Catalog) error {
	db := m.db.WithContext(ctx).Save(data)
	return db.Error
}

// Delete 删除目录（单条）
func (m *gormModel) Delete(ctx context.Context, id string) error {
	db := m.db.WithContext(ctx).Where("f_id = ?", id).Delete(&Catalog{})
	return db.Error
}

// DeleteBatch 批量删除
func (m *gormModel) DeleteBatch(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	db := m.db.WithContext(ctx).Where("f_id IN ?", ids).Delete(&Catalog{})
	return db.Error
}

// WithTx 创建事务
func (m *gormModel) WithTx(tx interface{}) CatalogModel {
	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return m
	}
	return &gormModel{db: gormTx}
}

// Trans 事务执行
func (m *gormModel) Trans(ctx context.Context, fn func(ctx context.Context, model CatalogModel) error) error {
	err := m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txModel := &gormModel{db: tx}
		return fn(ctx, txModel)
	})
	return err
}

// FindAllDescendants 递归查找所有子孙目录（辅助方法）
func (m *gormModel) FindAllDescendants(ctx context.Context, id string) ([]*Catalog, error) {
	var allDescendants []*Catalog

	// 查找直接子目录
	children, err := m.FindByParent(ctx, id)
	if err != nil {
		return nil, err
	}

	allDescendants = append(allDescendants, children...)

	// 递归查找子目录的子孙目录
	for _, child := range children {
		descendants, err := m.FindAllDescendants(ctx, child.Id)
		if err != nil {
			return nil, err
		}
		allDescendants = append(allDescendants, descendants...)
	}

	return allDescendants, nil
}

// IsDescendant 判断 targetId 是否是 ancestorId 的子孙节点（辅助方法）
func (m *gormModel) IsDescendant(ctx context.Context, ancestorId, targetId string) (bool, error) {
	if ancestorId == targetId {
		return true, nil
	}

	current, err := m.FindOne(ctx, targetId)
	if err != nil {
		return false, err
	}
	if current == nil {
		return false, nil
	}

	// 递归向上查找
	for current.ParentId != ancestorId && current.Level > CatalogRootLevel {
		parent, err := m.FindOne(ctx, current.ParentId)
		if err != nil {
			return false, err
		}
		if parent == nil {
			return false, nil
		}
		current = parent
	}

	return current.ParentId == ancestorId, nil
}

// CheckUniqueName 检查同级目录下名称是否唯一（辅助方法）
func (m *gormModel) CheckUniqueName(ctx context.Context, parentId string, catalogType int32, name string, excludeId string) (bool, error) {
	var count int64
	query := m.db.WithContext(ctx).Model(&Catalog{}).
		Where("f_parent_id = ?", parentId).
		Where("f_type = ?", catalogType).
		Where("f_catalog_name = ?", name)

	if excludeId != "" {
		query = query.Where("f_id != ?", excludeId)
	}

	db := query.Count(&count)
	if db.Error != nil {
		return false, db.Error
	}

	return count == 0, nil
}

// GetMaxLevel 获取指定目录的最大子层级（辅助方法，用于循环检测）
func (m *gormModel) GetMaxLevel(ctx context.Context, id string) (int32, error) {
	descendants, err := m.FindAllDescendants(ctx, id)
	if err != nil {
		return 0, err
	}

	maxLevel := int32(0)
	for _, desc := range descendants {
		if desc.Level > maxLevel {
			maxLevel = desc.Level
		}
	}

	return maxLevel, nil
}
