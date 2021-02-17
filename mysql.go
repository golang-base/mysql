package mysql

//package main

import (
	"database/sql"
	"errors"
	_mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type Config struct {
	DataSourceName  string // 数据库连接字符串
	MaxOpenConn     int    // 最大打开的连接数
	MaxIdelConn     int    // 最大空闲的连接数
	ConnMaxLifetime int    // 连接最大的生命时间
	ConnMaxIdleTime int    // 连接最大的空闲时间
}

var gormDB *gorm.DB

func Init(config *Config) (err error) {
	gormConfig := getGromConfig(config)
	gormDB, err = gorm.Open(_mysql.Open(config.DataSourceName), gormConfig)
	if err != nil {
		return err
	}
	sqlDB, err := gormDB.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxOpenConns(config.MaxOpenConn)
	sqlDB.SetMaxIdleConns(config.MaxIdelConn)
	sqlDB.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(config.ConnMaxIdleTime) * time.Second)
	return
}

func GetGromDB() (*gorm.DB, error) {
	if gormDB == nil {
		return nil, errors.New("db 未初始化·")
	}
	if gormDB.Error != nil {
		return nil, gormDB.Error
	}
	return gormDB, nil
}

func GetSqlDB() (*sql.DB, error) {
	if gormDB == nil {
		return nil, errors.New("db 未初始化·")
	}
	return gormDB.DB()
}

func Release() (err error) {
	sqlDB, err := GetSqlDB()
	if err != nil {
		return
	}
	return sqlDB.Close()
}

func getGromConfig(config *Config) (gormConfig *gorm.Config) {
	gormConfig = &gorm.Config{}
	gormConfig.Logger = logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // 彩色打印
		})
	gormConfig.Logger.LogMode(logger.Info)
	gormConfig.DryRun = false // 是否只生成sql 而不执行
	return
}

func main() {
	config := Config{
		DataSourceName:  "",
		MaxOpenConn:     20,
		MaxIdelConn:     10,
		ConnMaxIdleTime: 100,
		ConnMaxLifetime: 100,
	}
	err := Init(&config)
	if err == nil {
		println("数据库连接成功")
	} else {
		println(err.Error())
	}
}
