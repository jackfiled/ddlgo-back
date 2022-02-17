package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"ddl/config"
	"ddl/database"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var ErrExpTime = errors.New("exp time")

type UserInfo struct {
	UserID     string    `gorm:"column:userID"`
	StudentID  int32     `gorm:"column:studentID;primary_key"`
	Permission int64     `gorm:"column:permission"`
	Class      int32     `gorm:"column:class"`
	ExpTime    time.Time `gorm:"-"`
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
func SetCookieUserInfo(c *gin.Context, userInfo UserInfo) {
	userInfo.ExpTime = time.Now().AddDate(0, 0, 14)
	data, _ := json.Marshal(userInfo)
	fmt.Println(string(data))
	value, _ := Encrypt(string(data), []byte(config.ENCRYPT_KEY))
	// c.SetCookie("UserInfo", value, 14*24*3600, "/", "", false, true)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "UserInfo",
		Value:    value,
		Path:     "/",
		Domain:   "",
		MaxAge:   14 * 24 * 3600,
		Secure:   true,
		HttpOnly: false,
		SameSite: http.SameSiteNoneMode,
	})
}

//读取用户信息
func GetCookieUserInfo(c *gin.Context) (UserInfo, error) {
	var userInfo UserInfo

	cookie, err := c.Cookie("UserInfo")
	if err != nil {
		fmt.Println("Get cookie failed")
		return userInfo, err
	}
	data, err := Decrypt(cookie, []byte(config.ENCRYPT_KEY))
	if err != nil {
		fmt.Println("Decrypt failed")
		return userInfo, err
	}
	err = json.Unmarshal([]byte(data), &userInfo)
	if err != nil {
		fmt.Println("Json unmarshal failed")
		return userInfo, err
	}
	if time.Now().After(userInfo.ExpTime) {
		return userInfo, ErrExpTime
	}
	return userInfo, nil
}

//更新用户信息
func UpdateCookieUserInfo(c *gin.Context, studentID int32) {
	var userInfo UserInfo
	database.DB.Table("user").Where("studentID=?", studentID).Take(&userInfo)
	userInfo.ExpTime = time.Now().AddDate(0, 0, 14)
	data, _ := json.Marshal(userInfo)
	fmt.Println(string(data))
	value, _ := Encrypt(string(data), []byte(config.ENCRYPT_KEY))
	// c.SetCookie("UserInfo", value, 14*24*3600, "/", "", false, true)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "UserInfo",
		Value:    value,
		Path:     "/",
		Domain:   "",
		MaxAge:   14 * 24 * 3600,
		Secure:   true,
		HttpOnly: false,
		SameSite: http.SameSiteNoneMode,
	})
}

//清除登陆状态
func DelCookieUserInfo(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "UserInfo",
		Value:    "",
		Path:     "/",
		Domain:   "",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: false,
		SameSite: http.SameSiteNoneMode,
	})
}

//设置用户权限，保存在数据库中
func SetUserPermission(studentID int32, permission int64) UserInfo {
	var userInfo UserInfo
	database.DB.Table("user").Where("studentID=?", studentID).Update("Permission", permission).First(&userInfo)
	return userInfo
}

func AuthDemoHandler(c *gin.Context) {
	userInfo, err := GetCookieUserInfo(c)
	if err != nil {
		if !errors.Is(err, http.ErrNoCookie) {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		fmt.Println(err)
		c.Header("Content-Type", "text/html")
		c.String(200, "<script language='javascript'>window.location.href='/login.html'</script>")
		return
	}

	//已登陆
	//更新一下
	UpdateCookieUserInfo(c, userInfo.StudentID)
	userInfo, _ = GetCookieUserInfo(c)
	c.Header("Content-Type", "text/html")
	c.String(200, `
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

func CheckAuthHandler(c *gin.Context) {
	userInfo, err := GetCookieUserInfo(c)
	if err != nil {
		if !errors.Is(err, http.ErrNoCookie) {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(200, false)
		return
	}

	//已登陆
	//更新一下
	if userInfo.StudentID == 0 {
		userInfo.ExpTime = time.Now().AddDate(0, 0, 14)
		data, _ := json.Marshal(userInfo)
		fmt.Println(string(data))
		value, _ := Encrypt(string(data), []byte(config.ENCRYPT_KEY))
		// c.SetCookie("UserInfo", value, 14*24*3600, "/", "", false, true)
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "UserInfo",
			Value:    value,
			Path:     "/",
			Domain:   "",
			MaxAge:   14 * 24 * 3600,
			Secure:   true,
			HttpOnly: false,
			SameSite: http.SameSiteNoneMode,
		})
	} else {
		UpdateCookieUserInfo(c, userInfo.StudentID)
	}

	c.JSON(200, userInfo)
}
