package main

import (
	"drone-message/model"
	"fmt"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func DbLog(p *Plugin) (err error) {
	// 如果开启了记录日志则连接数据库
	if p.Extra.Db.DbLog {
		if db, err = GetDb(p); err != nil {
			model.DbMsg = err.Error()
			return
		}
		_ = db.AutoMigrate(&model.PublishLog{})
		// 连接数据库
		// 初始化数据表
		// 插入记录
	}
	return
}

func GetDb(p *Plugin) (*gorm.DB, error) {
	if db == nil {
		var err error
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			p.Extra.Db.DbUsername, p.Extra.Db.DbPassword, p.Extra.Db.DbHost, p.Extra.Db.DbPort, p.Extra.Db.DbName)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}

func WriteLog(info *model.PublishLog) error {
	db.Create(&info)
	return nil
}
