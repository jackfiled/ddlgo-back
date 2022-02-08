package admin

import (
	"ddl/common"
	"ddl/database"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SaveHandler(c *gin.Context) {
	var err error

	data := common.AdminNotice{}
	err = c.BindJSON(&data)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println(data)

	// userInfo, err := auth.GetCookieUserInfo(c)
	// if err != nil {
	// 	if errors.Is(err, http.ErrNoCookie) {
	// 		c.String(http.StatusBadRequest, err.Error())
	// 	} else {
	// 		c.String(http.StatusInternalServerError, err.Error())
	// 	}
	// }

	// switch data.Class {
	// case "":
	// 	c.String(http.StatusBadRequest, err.Error())
	// 	return
	// case "dddd":
	// 	if userInfo.Permission&(1<<int64(data.NoticeType)) == 0 {
	// 		c.String(http.StatusBadRequest, "权限不足")
	// 		return
	// 	}
	// default:
	// 	if "2021211"+data.Class != strconv.Itoa(int(userInfo.Class)) {
	// 		c.String(http.StatusBadRequest, "权限不足")
	// 		return
	// 	}
	// 	if userInfo.Permission&(1<<16<<int64(data.NoticeType)) == 0 {
	// 		c.String(http.StatusBadRequest, "权限不足")
	// 		return
	// 	}
	// }

	// notice := common.DDLNotice{
	// 	Index:      data.Index,
	// 	Title:      data.Title,
	// 	Detail:     data.Detail,
	// 	DDL:        data.DDL,
	// 	NoticeType: data.NoticeType,
	// 	Img:        data.Img,
	// 	StartTime:  data.StartTime,
	// }

	err = database.DB.Save(&data).Error
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	title := data.Title
	if data.Index != 0 {
		title = "【修改】" + title
	}
	// push.SendNotice(common.PartyMap[form.Get("classes")], notice.NoticeType, title,
	// 	form.Get("startTime")+"--"+form.Get("ddl"), notice.Detail,
	// 	"http://www.squidward.top/?cla="+form.Get("classes")+"#"+strconv.Itoa(notice.Index))
	_ = title
	c.JSON(200, data)
}

func DeleteHandler(c *gin.Context) {
	var err error
	form := c.Request.Form

	// userInfo, err := auth.GetCookieUserInfo(c)
	// if err != nil {
	// 	if errors.Is(err, http.ErrNoCookie) {
	// 		c.String(http.StatusBadRequest, err.Error())
	// 	} else {
	// 		c.String(http.StatusInternalServerError, err.Error())
	// 	}
	// }

	// var notice common.DDLNotice
	// err = database.DB.Table(form.Get("classes")).Where("index_=?", form.Get("index_")).First(&notice).Error
	// if err != nil {
	// 	c.String(http.StatusInternalServerError, err.Error())
	// 	return
	// }

	// switch form.Get("classes") {
	// case "":
	// 	c.String(http.StatusBadRequest, err.Error())
	// 	return
	// case "dddd":
	// 	if userInfo.Permission&(1<<int64(notice.NoticeType)) == 0 {
	// 		c.String(http.StatusBadRequest, "权限不足")
	// 		return
	// 	}
	// default:
	// 	if "2021211"+form.Get("classes") != strconv.Itoa(int(userInfo.Class)) {
	// 		c.String(http.StatusBadRequest, "权限不足")
	// 		return
	// 	}
	// 	if userInfo.Permission&(1<<16<<int64(notice.NoticeType)) == 0 {
	// 		c.String(http.StatusBadRequest, "权限不足")
	// 		return
	// 	}
	// }

	err = database.DB.Table(form.Get("classes")).Delete("index_=?", form.Get("index_")).Error
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(200, true)
}
