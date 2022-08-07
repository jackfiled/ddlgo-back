package database

import (
	"ddlBackend/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Database 数据库指针
var Database *gorm.DB

// OpenDatabase 打开数据库函数
func OpenDatabase() (err error) {
	Database, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	// 创建记录DDL事件的表
	if !Database.Migrator().HasTable(&models.DDLNotice{}) {
		err = Database.Migrator().CreateTable(&models.DDLNotice{})
	}

	return err
}
