package auth

import (
	"ddl/common"
	"ddl/config"
	"ddl/database"
	"ddl/wecom"
	"net/http"

	"github.com/gin-gonic/gin"
)

func WechatLoginHandler(c *gin.Context) {
	// fmt.Println(values)
	code := c.Request.FormValue("code")
	ref := c.Request.FormValue("ref")
	if ref == "" {
		ref = "http://squidward.top/"
	}
	if code == "" {
		c.String(http.StatusBadRequest, "参数错误")
	} else {
		userID := wecom.GetWecomID(code)
		// fmt.Println(userID)
		if userID != "" {
			userInfo := GetUserInfo(userID)
			SetCookieUserInfo(c, userInfo)

			c.Header("Content-Type", "text/html")
			c.String(200, "<script language='javascript'>window.location.href='"+ref+"'</script>")

			// fmt.Println(userInfo)
			// fmt.Fprintf(w, "%s %d %d %d", userInfo.UserID, userInfo.StudentID, userInfo.Class, userInfo.Permission)
		}
	}
}

func LoginHandler(c *gin.Context) {
	// fmt.Println(values)
	key := c.Request.FormValue("key")
	pass := LoginPassword(c, key)
	if pass {
		c.JSON(200, config.AdminKey[key])
	} else {
		c.JSON(200, false)
	}

}

func LogoutHandler(c *gin.Context) {
	// fmt.Println(values)
	ref := c.Request.FormValue("ref")
	if ref == "" {
		ref = "http://squidward.top/"
	}

	DelCookieUserInfo(c)

	c.Header("Content-Type", "text/html")
	c.String(200, "<script language='javascript'>window.location.href='"+ref+"'</script>")

	// fmt.Println(userInfo)
	// fmt.Fprintf(w, "%s %d %d %d", userInfo.UserID, userInfo.StudentID, userInfo.Class, userInfo.Permission)

}

func GetUserInfo(userID string) (userInfo common.UserInfo) {
	database.DB.Where("userID=?", userID).First(&userInfo)
	// fmt.Println(userInfo)
	return
}

//密码登录部分

func LoginPassword(c *gin.Context, key string) bool {

	userInfo, exist := config.AdminKey[key]
	if exist {
		SetCookieUserInfo(c, userInfo)
	}
	return exist
}
