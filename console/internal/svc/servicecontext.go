package svc

import (
	"console-api/console/internal/config"
	"console-api/console/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config       config.Config
	DynamicModel model.DynamicModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	open, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
	if c.Mysql.Debug {
		open = open.Debug()
	}

	if err != nil {
		return nil
	}
	return &ServiceContext{
		Config:       c,
		DynamicModel: model.NewDynamicModel(open),
	}
}
