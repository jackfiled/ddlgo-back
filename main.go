package main

import (
	"ddlBackend/Middleware"
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

	// 用户管理相关API需要验证
	route.GET("/users", handlers.ReadUsersHandler).Use(Middleware.JWTAuthMiddleware())
	route.POST("/users", handlers.CreateUserHandler).Use(Middleware.JWTAuthMiddleware())
	route.GET("/users/:id", handlers.ReadSingleUserHandler).Use(Middleware.JWTAuthMiddleware())
	route.PUT("/users/:id", handlers.UpdateUserHandler).Use(Middleware.JWTAuthMiddleware())
	route.DELETE("/users/:id", handlers.DeleteUserHandler).Use(Middleware.JWTAuthMiddleware())
	route.POST("/users/login", handlers.AdminLoginHandler)

	err = route.Run(tool.Setting.AppPort)
	if err != nil {
		tool.DDLLog(err.Error())
		return
	}
}
