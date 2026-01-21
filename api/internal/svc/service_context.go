// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"github.com/tianyuliang/standardization-backend/api/internal/config"
	"github.com/tianyuliang/standardization-backend/api/internal/middleware"
	"github.com/tianyuliang/standardization-backend/model/catalog"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config       config.Config
	Validator    rest.Middleware
	CatalogModel catalog.CatalogModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		Validator:    middleware.NewValidatorMiddleware().Handle,
		CatalogModel: nil, // TODO: 初始化 CatalogModel，在 main.go 中注入
	}
}
