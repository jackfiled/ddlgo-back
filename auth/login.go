package auth

import (
	"ddl/config"
	"ddl/database"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetAccessTokenRes struct {
	Access_token string
}

type GetUserIDRes struct {
	UserID string
	Errmsg string
}

func (UserInfo) TableName() string {
	return "user"
}

func LoginHandler(c *gin.Context) {
	// fmt.Println(values)
	code := c.Request.FormValue("code")
	ref := c.Request.FormValue("ref")
	if ref == "" {
		ref = "http://squidward.top/"
	}
	if code == "" {
		c.String(http.StatusBadRequest, "参数错误")
	} else {
		userID := getWecomID(code)
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

func HttpGet(url string) string {
	res, err := http.Get(url)
	if err != nil {
		return ""
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return ""
	}
	return string(robots)
}

//微信登录部分

func getAccessToken() string {
	data := HttpGet("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=" + config.CORPID + "&corpsecret=" + config.CORPSECRET)
	res := GetAccessTokenRes{}
	err := json.Unmarshal([]byte(data), &res)
	if err != nil {
		fmt.Println("GetAccessToken failed")
		return ""
	} else {
		return res.Access_token
	}
}

func getWecomID(code string) string {
	accessToken := getAccessToken()
	data := HttpGet("https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo?access_token=" + accessToken + "&code=" + code)
	res := GetUserIDRes{}
	err := json.Unmarshal([]byte(data), &res)
	// fmt.Println(data)
	if err != nil || res.Errmsg != "ok" {
		fmt.Printf("GetWecomID failed\n%s\n", res.Errmsg)
		return ""
	} else {
		return res.UserID
	}
}

func GetUserInfo(userID string) UserInfo {
	var userInfo UserInfo
	database.DB.Where("userID=?", userID).First(&userInfo)
	// fmt.Println(userInfo)

	return userInfo
}

//密码登录部分

func LoginPassword(c *gin.Context, key string) {

	permission, exist := config.AdminKey[key]
	if exist {
		userInfo := UserInfo{Permission: permission}
		SetCookieUserInfo(c, userInfo)
	}
}
