package admin

import (
	"ddl/common"
	"ddl/config"
	"fmt"

	"github.com/jinzhu/gorm"
)

func Insert(class string, notice common.DDLNotice) int16 {
	db, errDb := gorm.Open("mysql", config.DB_USER_PW+"@("+config.DB_HOST+")/test?charset=utf8&loc=Local&parseTime=true")
	if errDb != nil {
		fmt.Println(errDb)
	}
	defer db.Close() //用完之后关闭数据库连接
	db.LogMode(true) //开启sql debug 模式

	db.Table(class).Create(&notice)

	return int16(notice.Index)
}

func Update(class string, notice common.DDLNotice) {
	db, errDb := gorm.Open("mysql", config.DB_USER_PW+"@("+config.DB_HOST+")/test?charset=utf8&loc=Local&parseTime=true")
	if errDb != nil {
		fmt.Println(errDb)
	}
	defer db.Close() //用完之后关闭数据库连接
	db.LogMode(true) //开启sql debug 模式

	db.Save(&notice)
}

func Delete(class string, index_ int16) {
	db, errDb := gorm.Open("mysql", config.DB_USER_PW+"@("+config.DB_HOST+")/test?charset=utf8&loc=Local&parseTime=true")
	if errDb != nil {
		fmt.Println(errDb)
	}
	defer db.Close() //用完之后关闭数据库连接
	db.LogMode(true) //开启sql debug 模式

	db.Table(class).Delete("index_=?", index_)

}

func InputHandler() {

}

func ModifyHandler() {

}

func DeleteHandler() {

}
