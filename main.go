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
		tool.DDLLogError(err.Error())
		tool.DDLLogError("Read config file failed, using default setting")
	}

	// 打开数据库
	err = database.OpenDatabase()
	if err != nil {
		tool.DDLLogError(err.Error())
		return
	}

	route := gin.Default()

	// 登录
	route.POST("/login", handlers.AdminLoginHandler)
	// 获取DDL事件列表
	route.GET("/ddlNotices", handlers.ReadDDLHandler)
	route.GET("/ddlNotices/:class", handlers.ReadClassDDLHandler)
	route.GET("/ddlNotices/:class/:id", handlers.ReadClassIDDDLHandler)

	// 图片文件路径
	route.Static("/picture", "./picture")

	// 修改DDL事件列表需要身份验证
	ddlNoticesRoute := route.Group("")
	ddlNoticesRoute.Use(Middleware.JWTAuthMiddleware())
	{
		ddlNoticesRoute.POST("/ddlNotices", handlers.CreateDDLHandler)
		ddlNoticesRoute.POST("/ddlNotices/:class", handlers.CreateClassDDLHandler)
		ddlNoticesRoute.PUT("/ddlNotices/:class/:id", handlers.UpdateClassIDDDLHandler)
		ddlNoticesRoute.DELETE("/ddlNotices/:class/:id", handlers.DeleteClassIDDDLHandler)
	}

	// 其他需要身份验证的API
	adminRoute := route.Group("")
	adminRoute.Use(Middleware.JWTAuthMiddleware())
	{
		adminRoute.POST("/upload", handlers.UploadPictureHandler)
	}

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
		tool.DDLLogError(err.Error())
		return
	}
}
