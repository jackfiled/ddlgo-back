package main

import (
	"ddlBackend/database"
	"ddlBackend/log"
	"github.com/gin-gonic/gin"
)

func main() {
	// 打开数据库
	err := database.OpenDatabase()
	if err != nil {
		log.DDLLog(err.Error())
		return
	}

	route := gin.Default()

	err = route.Run()
	if err != nil {
		log.DDLLog(err.Error())
		return
	}
}
