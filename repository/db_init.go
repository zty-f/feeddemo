package repository

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm/logger"

	"gorm.io/gorm"
)

var db *gorm.DB

// Init 初始化数据库连接操作
func Init() error {
	var err error
	// 本地环境
	//dsn := "root:123456@tcp(127.0.0.1:3306)/feed?charset=utf8mb4&parseTime=True&loc=Local"
	// 测试环境
	dsn := "platform_root:Bu^Wd8k7JU7FS7MxBNruMmS2BQ8qjJ@tcp(rm-2ze58oz13gl4089i37o.mysql.rds.aliyuncs.com:3306)/feed_test?charset=utf8&loc=Local&parseTime=True"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	return err

}
