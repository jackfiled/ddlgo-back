package main

import (
	"ddlBackend/database"
	"ddlBackend/handlers"
	"ddlBackend/tool"
	"github.com/gin-gonic/gin"
)

func main() {
	// 读取配置文件
	err := tool.ReadConfig()
	if err != nil {
		tool.DDLLog(err.Error())
		tool.DDLLog("Read config file failed, using default setting")
	}

	// 打开数据库
	err = database.OpenDatabase()
	if err != nil {
		tool.DDLLog(err.Error())
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

	route.GET("/users", handlers.ReadUsersHandler)
	route.POST("/users", handlers.CreateUserHandler)
	route.GET("/users/:id", handlers.ReadSingleUserHandler)
	route.PUT("/users/:id", handlers.UpdateUserHandler)
	route.DELETE("/users/:id", handlers.DeleteUserHandler)

	err = route.Run(tool.Setting.AppPort)
	if err != nil {
		tool.DDLLog(err.Error())
		return
	}
}
