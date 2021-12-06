package query

import (
	"ddl/common"
	"ddl/config"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
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

func QuerySingleHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	class := values.Get("class")
	index_ := values.Get("index_")
	if class == "" || index_ == "" {
		fmt.Fprint(w, "GET 参数未指定")
		return
	}
	db, errDb := gorm.Open("mysql", config.DB_USER_PW+"@("+config.DB_HOST+")/test?charset=utf8&loc=Local&parseTime=true")
	if errDb != nil {
		fmt.Println(errDb)
	}
	defer db.Close() //用完之后关闭数据库连接

	db.LogMode(true) //开启sql debug 模式

	noticeArr := []common.DDLNotice{}

	db.Table(class).Where("index_=?", index_).Find(&noticeArr)
	data, _ := json.Marshal(noticeArr)

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json
	fmt.Fprint(w, string(data))
}

func GetListHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	class := values.Get("class")

	if class == "" {
		class = "dddd"
	}

	start := values.Get("start")
	if start == "" {
		start = "0"
	}

	step := values.Get("step")
	if step == "" {
		step = "20"
	}

	noticeType := values.Get("noticeType")
	if noticeType == "" {
		noticeType = "-1"
	}

	db, errDb := gorm.Open("mysql", config.DB_USER_PW+"@("+config.DB_HOST+")/test?charset=utf8&loc=Local&parseTime=true")
	if errDb != nil {
		fmt.Println(errDb)
	}
	defer db.Close() //用完之后关闭数据库连接

	db.LogMode(true) //开启sql debug 模式

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
		db.Table(class).Limit(step).Offset(start).Order("time DESC").Find(&noticeArr)
	} else if noticeType == "-2" {
		db.Table(class).Limit(step).Offset(start).Where("noticeType != 0").Order("state, ddl").Find(&noticeArr)
	} else {
		noticeTypeInt, _ := strconv.Atoi(noticeType)
		db.Table(class).Limit(step).Offset(start).Where("noticeType = ?", noticeTypeInt).Order("state, ddl").Find(&noticeArr)
	}

	data, _ := json.Marshal(noticeArr)

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	fmt.Fprint(w, string(data))
}
