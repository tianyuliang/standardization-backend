package catalog

import (
	"sync"

	"gorm.io/gorm"
)

var (
	// catalogModelCache 目录模型缓存
	catalogModelCache = make(map[string]*CatalogModel)
	// catalogModelLock 目录模型锁
	catalogModelLock sync.RWMutex
)

// ModelConfig 模型配置
type ModelConfig struct {
	DB *gorm.DB
}

// NewCatalogModel 创建目录模型工厂函数
func NewCatalogModel(cfg ModelConfig) CatalogModel {
	catalogModelLock.Lock()
	defer catalogModelLock.Unlock()

	// 简单缓存，实际项目中应该用更复杂的缓存策略
	cacheKey := "catalog"
	if model, ok := catalogModelCache[cacheKey]; ok && model != nil {
		return *model
	}

	model := newGormModel(cfg)
	catalogModelCache[cacheKey] = &model
	return model
}
