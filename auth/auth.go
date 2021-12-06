package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"ddl/config"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
)

type UserInfo struct {
	UserID     string `gorm:"column:userID"`
	StudentID  int32  `gorm:"column:studentID;primary_key"`
	Permission int64  `gorm:"column:permission"`
	Class      int32  `gorm:"column:class"`
}

func Encrypt(text string, key []byte) (string, error) {
	var iv = key[:aes.BlockSize]
	encrypted := make([]byte, len(text))
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	encrypter := cipher.NewCFBEncrypter(block, iv)
	encrypter.XORKeyStream(encrypted, []byte(text))
	return hex.EncodeToString(encrypted), nil
}

func Decrypt(encrypted string, key []byte) (string, error) {
	var err error
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	src, err := hex.DecodeString(encrypted)
	if err != nil {
		return "", err
	}
	var iv = key[:aes.BlockSize]
	decrypted := make([]byte, len(src))
	var block cipher.Block
	block, err = aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	decrypter := cipher.NewCFBDecrypter(block, iv)
	decrypter.XORKeyStream(decrypted, src)
	return string(decrypted), nil
}

//保存用户信息
func SetCookieUserInfo(w http.ResponseWriter, userInfo UserInfo) {
	expiration := time.Now().AddDate(0, 0, 14)
	data, _ := json.Marshal(userInfo)
	fmt.Println(string(data))
	value, _ := Encrypt(string(data), []byte(config.ENCRYPT_KEY))
	cookie := http.Cookie{Name: "UserInfo", Value: value, Expires: expiration}
	http.SetCookie(w, &cookie)
}

//读取用户信息
func GetCookieUserInfo(r *http.Request) (UserInfo, error) {
	var userInfo UserInfo

	cookie, err := r.Cookie("UserInfo")
	if err != nil {
		fmt.Println("Get cookie failed")
		return userInfo, err
	}
	data, err := Decrypt(cookie.Value, []byte(config.ENCRYPT_KEY))
	if err != nil {
		fmt.Println("Decrypt failed")
		return userInfo, err
	}
	err = json.Unmarshal([]byte(data), &userInfo)
	if err != nil {
		fmt.Println("Json unmarshal failed")
		return userInfo, err
	}
	return userInfo, nil
}

//更新用户信息
func UpdateCookieUserInfo(w http.ResponseWriter, studentID int32) {
	db, errDb := gorm.Open("mysql", config.DB_USER_PW+"@("+config.DB_HOST+")/test?charset=utf8&loc=Local&parseTime=true")
	if errDb != nil {
		fmt.Println(errDb)
	}
	defer db.Close() //用完之后关闭数据库连接
	db.LogMode(true) //开启sql debug 模式

	var userInfo UserInfo
	db.Table("user").Where("studentID=?", studentID).Take(&userInfo)

	expiration := time.Now().AddDate(0, 0, 14)
	data, _ := json.Marshal(userInfo)
	fmt.Println(string(data))
	value, _ := Encrypt(string(data), []byte(config.ENCRYPT_KEY))
	cookie := http.Cookie{Name: "UserInfo", Value: value, Expires: expiration}
	http.SetCookie(w, &cookie)
}

//清除登陆状态
func DelCookieUserInfo(w http.ResponseWriter) {
	expiration := time.Now().AddDate(0, 0, -1)
	var userInfo UserInfo
	data, _ := json.Marshal(userInfo)
	fmt.Println(string(data))
	value, _ := Encrypt(string(data), []byte(config.ENCRYPT_KEY))
	cookie := http.Cookie{Name: "UserInfo", Value: value, Expires: expiration}
	http.SetCookie(w, &cookie)
}

//设置用户权限，保存在数据库中
func SetUserPermission(studentID int32, permission int64) UserInfo {
	db, errDb := gorm.Open("mysql", config.DB_USER_PW+"@("+config.DB_HOST+")/test?charset=utf8&loc=Local&parseTime=true")
	if errDb != nil {
		fmt.Println(errDb)
	}
	defer db.Close() //用完之后关闭数据库连接
	db.LogMode(true) //开启sql debug 模式

	var userInfo UserInfo
	db.Table("user").Where("studentID=?", studentID).Update("Permission", permission).First(&userInfo)
	return userInfo
}

func AuthDemoHandler(w http.ResponseWriter, r *http.Request) {
	userInfo, err := GetCookieUserInfo(r)
	if err != nil {
		fmt.Println(err)
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
		<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
</head>
<body>
		未登录 
		<a href="https://open.weixin.qq.com/connect/oauth2/authorize?
		appid=ww8a5308483ff283cc
		&redirect_uri=%s
		&response_type=code
		&scope=snsapi_base
		&state=STATE
		#wechat_redirect">
		微信一键登录</a>

		<a href="https://open.work.weixin.qq.com/wwopen/sso/qrConnect?
		appid=ww8a5308483ff283cc
		&agentid=1000003
		&redirect_uri=%s
		">
		微信扫码登录</a>
</body>
</html>
		
		`,
			"http%3A%2F%2Fsquidward.top%3A8000%2Fapi%2Flogin%3Fref%3Dhttp%3A%2F%2Fsquidward.top%3A8000%2Fapi%2Fauth_demo",
			"http%3A%2F%2Fsquidward.top%3A8000%2Fapi%2Flogin%3Fref%3Dhttp%3A%2F%2Fsquidward.top%3A8000%2Fapi%2Fauth_demo")
		return
	}

	//更新一下
	UpdateCookieUserInfo(w, userInfo.StudentID)
	userInfo, _ = GetCookieUserInfo(r)
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
</head>
<body>
已登录 %s %d %d %d<br/> 
<a href="/api/logout?ref=%s">
点击退出登录</a>
</body>
</html>
		
		`, userInfo.UserID, userInfo.StudentID, userInfo.Class, userInfo.Permission, "http%3A%2F%2Fsquidward.top%3A8000%2Fapi%2Fauth_demo")

}
