package main

import (
	"net/http"
	"strconv"

	"ddl/auth"
	"ddl/query"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	router.GET("/", indexHandler)
	router.GET("/api/login", auth.LoginHandler)
	router.GET("/api/logout", auth.LogoutHandler)

	router.GET("/WW_verify_udfdZsIBL9yNi4SN.txt", WWVerify)

	router.GET("/api/auth_demo", auth.AuthDemoHandler)
	router.GET("/api/set_permission_demo", setPermission)

	router.GET("/api/get_list", query.GetListHandler)
	router.GET("/api/query_single", query.QuerySingleHandler)

	router.Run(":8000")
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
