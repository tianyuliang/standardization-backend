// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/config"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/middleware"
	catalogmodel "github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/catalog/catalog"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/rule/relation_file"
	rulemodel "github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/rule/rule"
	stdfilemodel "github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/stdfile/stdfile"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config                config.Config
	db                    *sqlx.DB
	RuleModel             rulemodel.RuleModel
	CatalogModel          catalogmodel.CatalogModel
	RelationRuleFileModel relation_file.RelationRuleFileModel
	StdFileModel          stdfilemodel.StdFileModel
	TokenCheck            rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 使用MySQL驱动
	driver := "mysql"
	dataSource := c.DB.Default.DataSource()

	db, err := sqlx.Connect(driver, dataSource)
	if err != nil {
		panic(err)
	}

	conn, err := db.Connx(context.Background())
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:                c,
		db:                    db,
		RuleModel:             rulemodel.NewRuleModel(conn),
		CatalogModel:          catalogmodel.NewCatalogModel(conn),
		RelationRuleFileModel: relation_file.NewRelationRuleFileModel(conn),
		StdFileModel:          stdfilemodel.NewStdFileModel(conn),
		TokenCheck:            middleware.NewTokenCheckMiddleware().Handle,
	}
}
