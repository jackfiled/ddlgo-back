package admin

import (
	"ddl/common"
	"ddl/database"
)

func Insert(class string, notice common.DDLNotice) int16 {
	database.DB.Table(class).Create(&notice)

	return int16(notice.Index)
}

func Update(class string, notice common.DDLNotice) {
	database.DB.Save(&notice)
}

func Delete(class string, index_ int16) {
	database.DB.Table(class).Delete("index_=?", index_)

}

func InputHandler() {

}

func ModifyHandler() {

}

func DeleteHandler() {

}
