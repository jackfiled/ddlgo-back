package common

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type DDLNotice struct {
	Index      int        `gorm:"column:index_;primary_key" json:"index_"`
	Time       time.Time  `gorm:"column:time" json:"time"`
	DDL        time.Time  `gorm:"column:ddl" json:"ddl"`
	StartTime  *time.Time `gorm:"column:startTime" json:"startTime"`
	Title      string     `gorm:"column:title" json:"title"`
	Detail     string     `gorm:"column:detail" json:"detail"`
	NoticeType int        `gorm:"column:noticeType" json:"noticeType"`
	Img        string     `gorm:"column:img" json:"img"`
}

//用于接收管理端数据的结构
type AdminNotice struct {
	Class      string     `gorm:"-" json:"class"` //不写入数据库
	Index      int        `gorm:"column:index_;primary_key" json:"index_"`
	DDL        time.Time  `gorm:"column:ddl" json:"ddl"`
	StartTime  *time.Time `gorm:"column:startTime" json:"startTime"`
	Title      string     `gorm:"column:title" json:"title"`
	Detail     string     `gorm:"column:detail" json:"detail"`
	NoticeType int        `gorm:"column:noticeType" json:"noticeType"`
	Img        string     `gorm:"column:img" json:"img"`
}

func (an *AdminNotice) TableName() string {
	return an.Class
}

var PartyMap = map[string]string{
	"dddd": "1",
	"304":  "2",
	"305":  "7",
	"306":  "8",
	"307":  "9",
	"308":  "10",
	"309":  "11",
	"test": "6",
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

func Cors() gin.HandlerFunc {
	// 	const originList = `http://localhost:3000
	// http://127.0.0.1:3000`
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.GetHeader("Origin")
		if c.GetHeader("Origin") != "" && (strings.Contains(origin, "http://localhost:") || strings.Contains(origin, "http://127.0.0.1:")) {
			c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"))
			// c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
		}
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
