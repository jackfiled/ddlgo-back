package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"ddl/admin"
	"ddl/auth"
	"ddl/common"
	"ddl/config"
	"ddl/database"
	"ddl/query"

	"github.com/gin-gonic/gin"
)

func main() {
	err := database.Open()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	router := gin.Default()

	router.Use(common.Cors())

	router.GET("/", indexHandler)

	//测试用
	router.StaticFile("/login.html", "./static/login.html")

	router.GET("/api/wechatlogin", auth.WechatLoginHandler)
	router.POST("/api/login", auth.LoginHandler)
	router.GET("/api/logout", auth.LogoutHandler)

	router.GET("/WW_verify_udfdZsIBL9yNi4SN.txt", WWVerify)

	router.GET("/api/check_auth", auth.CheckAuthHandler)
	router.GET("/api/auth_demo", auth.AuthDemoHandler)
	router.GET("/api/set_permission_demo", setPermission)

	router.GET("/api/get_list", query.GetListHandler)
	router.GET("/api/query_single", query.QuerySingleHandler)

	router.POST("/api/save", admin.SaveHandler)
	router.DELETE("/api/delete", admin.DeleteHandler)
	router.POST("/api/upload_img", admin.UploadFileHandler)

	router.Run(config.WEB_ADDR)
}

func indexHandler(c *gin.Context) {
	c.String(200, "DDL系统API")
}

//微信扫码授权验证
func WWVerify(c *gin.Context) {
	c.String(200, "udfdZsIBL9yNi4SN")
}

func setPermission(c *gin.Context) {
	if c.Request.FormValue("studentID") == "" {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	studentID, _ := strconv.ParseInt(c.Request.FormValue("studentID"), 10, 32)
	permission, _ := strconv.ParseInt(c.Request.FormValue("permission"), 10, 64)
	userInfo := auth.SetUserPermission(int32(studentID), permission)
	c.String(200, strconv.Itoa(int(userInfo.StudentID))+"权限修改为"+strconv.Itoa(int(userInfo.Permission)))
}
