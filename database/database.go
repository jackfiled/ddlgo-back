package database

import (
	"ddl/config"

	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func Open() (err error) {
	DB, err = gorm.Open("mysql", config.DB_USER_PW+"@("+config.DB_HOST+")/test?charset=utf8&loc=Local&parseTime=true")
	DB.LogMode(true) //开启sql debug 模式
	return err
}
