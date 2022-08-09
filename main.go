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

	route.POST("/login", handlers.AdminLoginHandler)

	route.GET("/ddlNotices", handlers.ReadDDLHandler)
	route.POST("/ddlNotices", handlers.CreateDDLHandler)

	route.GET("/ddlNotices/:class", handlers.ReadClassDDLHandler)
	route.POST("/ddlNotices/:class", handlers.CreateClassDDLHandler)

	route.GET("/ddlNotices/:class/:id", handlers.ReadClassIDDDLHandler)
	route.PUT("/ddlNotices/:class/:id", handlers.UpdateClassIDDDLHandler)
	route.DELETE("/ddlNotices/:class/:id", handlers.DeleteClassIDDDLHandler)

	// 用户管理相关API需要验证
	userRoute := route.Group("/users")
	userRoute.Use(Middleware.JWTAuthMiddleware())
	{
		userRoute.GET("/", handlers.ReadUsersHandler)
		userRoute.POST("/", handlers.CreateUserHandler)
		userRoute.GET("/:id", handlers.ReadSingleUserHandler)
		userRoute.PUT("/:id", handlers.UpdateUserHandler)
		userRoute.DELETE("/:id", handlers.DeleteUserHandler)
	}

	err = route.Run(tool.Setting.AppPort)
	if err != nil {
		tool.DDLLog(err.Error())
		return
	}
}
