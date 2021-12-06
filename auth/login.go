package auth

import (
	"ddl/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //这个一定要引入哦！！
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	// fmt.Println(values)
	code := values.Get("code")
	ref := values.Get("ref")
	if ref == "" {
		ref = "http://squidward.top/"
	}
	if code == "" {
		fmt.Fprintf(w, "参数错误")
	} else {
		userID := LoginGetUserID(code)
		// fmt.Println(userID)
		if userID != "" {
			userInfo := GetUserInfo(userID)
			SetCookieUserInfo(w, userInfo)
			request, _ := http.NewRequest("GET", ref, nil)
			http.Redirect(w, request, ref, http.StatusFound)
			// fmt.Println(userInfo)
			// fmt.Fprintf(w, "%s %d %d %d", userInfo.UserID, userInfo.StudentID, userInfo.Class, userInfo.Permission)
		}
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	// fmt.Println(values)
	ref := values.Get("ref")
	if ref == "" {
		ref = "http://squidward.top/"
	}

	DelCookieUserInfo(w)
	request, _ := http.NewRequest("GET", ref, nil)
	http.Redirect(w, request, ref, http.StatusFound)
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

func LoginGetAccessToken() string {
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

func LoginGetUserID(code string) string {
	accessToken := LoginGetAccessToken()
	data := HttpGet("https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo?access_token=" + accessToken + "&code=" + code)
	res := GetUserIDRes{}
	err := json.Unmarshal([]byte(data), &res)
	// fmt.Println(data)
	if err != nil || res.Errmsg != "ok" {
		fmt.Printf("GetUserID failed\n%s\n", res.Errmsg)
		return ""
	} else {
		return res.UserID
	}
}

func GetUserInfo(userID string) UserInfo {
	db, errDb := gorm.Open("mysql", config.DB_USER_PW+"@("+config.DB_HOST+")/test?charset=utf8&loc=Local")
	if errDb != nil {
		fmt.Println(errDb)
	}
	defer db.Close() //用完之后关闭数据库连接

	db.LogMode(true) //开启sql debug 模式

	var userInfo UserInfo

	db.Where("userID=?", userID).First(&userInfo)
	// fmt.Println(userInfo)

	return userInfo
}

//密码登录部分

func LoginPassword(w http.ResponseWriter, key string) {

	permission, exist := config.AdminKey[key]
	if exist {
		userInfo := UserInfo{Permission: permission}
		SetCookieUserInfo(w, userInfo)
	}
}
