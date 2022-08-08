package database

import (
	"ddlBackend/models"
	"errors"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DDLTables DDL表列表
var DDLTables = [7]string{"dddd", "304", "305", "306", "307", "308", "309"}

// Database 数据库指针
var Database *gorm.DB

// OpenDatabase 打开数据库函数
func OpenDatabase() (err error) {
	Database, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	// 创建ddl相关的表
	for _, value := range DDLTables {
		if !Database.Migrator().HasTable(value) {
			err = Database.Table(value).AutoMigrate(&models.DDLNotice{})
			if err != nil {
				return err
			}
		}
	}

	// 创建用户表
	if !Database.Migrator().HasTable(&models.UserInformation{}) {
		err = Database.AutoMigrate(&models.UserInformation{})
		if err != nil {
			return err
		}
	}

	return nil
}

// GetDDLTable 获取指定班级的DDL事件表
func GetDDLTable(className string) (*gorm.DB, error) {
	for _, value := range DDLTables {
		if className == value {
			return Database.Table(className).Order("ddl_time DESC"), nil
		}
	}
	return nil, errors.New(fmt.Sprintf("The table named %s not exists", className))
}
