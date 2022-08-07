package main

import (
	"ddlBackend/database"
	"ddlBackend/handlers"
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

	route.GET("/ddlNotices", handlers.ReadDDLHandler)
	route.POST("/ddlNotices", handlers.CreateDDLHandler)

	err = route.Run()
	if err != nil {
		log.DDLLog(err.Error())
		return
	}
}
