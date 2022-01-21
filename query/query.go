package query

import (
	"ddl/common"
	"ddl/database"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql" //这个一定要引入哦！！
)

func Base64Decode(sEnc string) string {

	// Base64 Standard Decoding
	sDec, err := base64.StdEncoding.DecodeString(sEnc)
	if err != nil {
		fmt.Printf("Error decoding string: %s ", err.Error())
		return ""
	}

	return string(sDec)
}

func QuerySingleHandler(c *gin.Context) {
	class := c.Request.FormValue("class")
	index_ := c.Request.FormValue("index_")
	if class == "" || index_ == "" {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}

	noticeArr := []common.DDLNotice{}

	database.DB.Table(class).Where("index_=?", index_).Find(&noticeArr)
	data, _ := json.Marshal(noticeArr)

	// w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	// w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	c.Data(200, "application/json", data)
}

func GetListHandler(c *gin.Context) {
	class := c.Request.FormValue("class")

	if class == "" {
		class = "dddd"
	}

	start := c.Request.FormValue("start")
	if start == "" {
		start = "0"
	}

	step := c.Request.FormValue("step")
	if step == "" {
		step = "20"
	}

	noticeType := c.Request.FormValue("noticeType")
	if noticeType == "" {
		noticeType = "-1"
	}

	noticeArr := []common.DDLNotice{}

	// 	db.Table(class).Find(&noticeArr).Order("ddl DESC")

	// 	for _, notice := range noticeArr {
	// 		nowStamp := time.Now().Unix()
	// 		ddlStamp := notice.DDL.Unix()
	// 		stateNew := notice.State

	// 		if nowStamp-ddlStamp > 259200 {
	// 			stateNew = 3 // 不显示 i.e. 超时 3 天以上
	// 		} else if ddlStamp <= nowStamp {
	// 			stateNew = 2 // 超时
	// 		} else if ddlStamp-nowStamp > 86400 {
	// 			stateNew = 1 // 进行中
	// 		} else {
	// 			stateNew = 0 // 紧急 i.e. 距 ddl 1 天以内
	// 		}

	// 		if stateNew != notice.State {
	// 			var noticeOp DDLNotice
	// 			db.Table(class).Where("index_ = ?", notice.Index).Take(&noticeOp)
	// 			noticeOp.State = stateNew
	// 			db.Table(class).Save(&noticeOp)//此处有BUG！！！！！！
	// 		} else if stateNew == 3 {
	// 			// ddl DESC 排序，如果已经 discarded 且不需更新则后面也不需更新
	// 			// 这个优化其实没啥用，把上面那个判完之后不交互数据库了，算几个时间戳减法也不浪费时间
	// 			// TODO) 真该优化了的话到时候在 SELECT * 的时候整个 LIMIT
	// 			break
	// 		}
	// 	}

	if noticeType == "-1" {
		database.DB.Table(class).Limit(step).Offset(start).Order("time DESC").Find(&noticeArr)
	} else if noticeType == "-2" {
		database.DB.Table(class).Limit(step).Offset(start).Where("noticeType != 0").Order("state, ddl").Find(&noticeArr)
	} else {
		noticeTypeInt, _ := strconv.Atoi(noticeType)
		database.DB.Table(class).Limit(step).Offset(start).Where("noticeType = ?", noticeTypeInt).Order("state, ddl").Find(&noticeArr)
	}

	data, _ := json.Marshal(noticeArr)

	// w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	// w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	c.Data(200, "application/json", data)
}
