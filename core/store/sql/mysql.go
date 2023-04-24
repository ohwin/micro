package sql

import (
	"github.com/ohwin/micro/core/config"
	. "github.com/ohwin/micro/core/constant"
	"github.com/ohwin/micro/core/log"
	"github.com/ohwin/micro/core/store"

	"gorm.io/driver/mysql"
	"gorm.io/gorm/logger"

	"gorm.io/gorm"
)

func Init() {
	c := config.App

	if c.Mysql.IsNil() {
		return
	}

	dsn := c.Mysql.DSN()
	log.Debug(dsn)
	var level logger.LogLevel
	switch c.Env {
	case Release:
		level = logger.Error
	default:
		level = logger.Info
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(level),
	})
	if err != nil {
		panic(err)
	}
	store.DB = db
}

func Migrate(tables ...interface{}) {
	err := store.DB.AutoMigrate(tables)
	if err != nil {
		panic(err)
	}
}
