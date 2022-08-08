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

	route.GET("/ddlNotices/:class", handlers.ReadClassDDLHandler)
	route.POST("/ddlNotices/:class", handlers.CreateClassDDLHandler)

	route.GET("/ddlNotices/:class/:id", handlers.ReadClassIDDDLHandler)
	route.PUT("/ddlNotices/:class/:id", handlers.UpdateClassIDDDLHandler)
	route.DELETE("/ddlNotices/:class/:id", handlers.DeleteClassIDDDLHandler)

	err = route.Run()
	if err != nil {
		log.DDLLog(err.Error())
		return
	}
}
