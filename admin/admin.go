package admin

import (
	"ddl/common"
	"ddl/database"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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
	// push.SendNotice(common.PartyMap[form.Get("class")], notice.NoticeType, title,
	// 	form.Get("startTime")+"--"+form.Get("ddl"), notice.Detail,
	// 	"http://www.squidward.top/?cla="+form.Get("class")+"#"+strconv.Itoa(notice.Index))
	_ = title
	c.JSON(200, data)
}

func DeleteHandler(c *gin.Context) {
	var err error
	class := c.Request.FormValue("class")
	indStr := c.Request.FormValue("index_")
	if class == "" || indStr == "" {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}

	index, err := strconv.Atoi(indStr)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
	}

	// userInfo, err := auth.GetCookieUserInfo(c)
	// if err != nil {
	// 	if errors.Is(err, http.ErrNoCookie) {
	// 		c.String(http.StatusBadRequest, err.Error())
	// 	} else {
	// 		c.String(http.StatusInternalServerError, err.Error())
	// 	}
	//  return
	// }

	notice := common.AdminNotice{Class: class, Index: index}
	err = database.DB.Model(&notice).First(&notice).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.String(http.StatusBadRequest, "通知不存在")
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// switch form.Get("class") {
	// case "":
	// 	c.String(http.StatusBadRequest, err.Error())
	// 	return
	// case "dddd":
	// 	if userInfo.Permission&(1<<int64(notice.NoticeType)) == 0 {
	// 		c.String(http.StatusBadRequest, "权限不足")
	// 		return
	// 	}
	// default:
	// 	if "2021211"+form.Get("class") != strconv.Itoa(int(userInfo.Class)) {
	// 		c.String(http.StatusBadRequest, "权限不足")
	// 		return
	// 	}
	// 	if userInfo.Permission&(1<<16<<int64(notice.NoticeType)) == 0 {
	// 		c.String(http.StatusBadRequest, "权限不足")
	// 		return
	// 	}
	// }

	err = database.DB.Delete(&common.AdminNotice{Class: class, Index: index}).Error
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(200, true)
}
