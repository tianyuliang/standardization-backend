// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/config"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/rule/relation_file"
	rulemodel "github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/rule/rule"
)

type ServiceContext struct {
	Config  config.Config
	DB      *sqlx.DB
	RuleModel            rulemodel.RuleModel
	RelationRuleFileModel relation_file.RelationRuleFileModel
	// TODO: 后续添加 Kafka Producer
	// KafkaProducer *kafka.Producer
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化数据库连接
	db, err := sqlx.Connect("mysql", c.DB.Default.DataSource())
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(c.DB.Default.MaxOpenConns)
	db.SetMaxIdleConns(c.DB.Default.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(c.DB.Default.ConnMaxLifetime) * time.Second)
	db.SetConnMaxIdleTime(time.Duration(c.DB.Default.ConnMaxIdleTime) * time.Second)

	// 获取底层连接用于 Model (使用 Conn 而非 DB 以支持事务)
	conn, err := db.Connx(context.Background())
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:               c,
		DB:                   db,
		RuleModel:            rulemodel.NewRuleModel(conn),
		RelationRuleFileModel: relation_file.NewRelationRuleFileModel(conn),
	}
}
