// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/gorm"
)

type Config struct {
	rest.RestConf
	DB *gorm.DB `json:",optional"` // 数据库连接
}
