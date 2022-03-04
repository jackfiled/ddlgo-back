package admin

import (
	"ddl/auth"
	"ddl/common"
	"ddl/database"
	"ddl/push"
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

	userInfo, err := auth.GetCookieUserInfo(c)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			c.String(http.StatusBadRequest, err.Error())
		} else {
			c.String(http.StatusInternalServerError, err.Error())
		}
	}

	switch data.Class {
	case "":
		c.String(http.StatusBadRequest, err.Error())
		return
	case "dddd":
		if userInfo.Permission&(1<<int64(data.NoticeType)) == 0 {
			c.String(http.StatusBadRequest, "权限不足")
			return
		}
	default:
		fmt.Println(data.Class)
		if userInfo.Class != -1 && "2021211"+data.Class != strconv.Itoa(int(userInfo.Class)) {
			c.String(http.StatusBadRequest, "权限不足")
			return
		}
		if userInfo.Permission&(1<<20<<int64(data.NoticeType)) == 0 {
			c.String(http.StatusBadRequest, "权限不足")
			return
		}
	}

	title := data.Title
	if data.Index != 0 {
		title = "【修改】" + title
	}

	err = database.DB.Save(&data).Error
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	startTime := ""
	if data.StartTime != nil {
		startTime = data.StartTime.Format("2006-01-02 15:04")
	}
	push.SendNotice(common.PartyMap[data.Class], data.NoticeType, title,
		startTime+"--"+data.DDL.Format("2006-01-02 15:04"), data.Detail,
		"http://www.squidward.top/#/home/"+data.Class+"/"+strconv.Itoa(data.Index))
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

	userInfo, err := auth.GetCookieUserInfo(c)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			c.String(http.StatusBadRequest, err.Error())
		} else {
			c.String(http.StatusInternalServerError, err.Error())
		}
		return
	}

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

	switch notice.Class {
	case "":
		c.String(http.StatusBadRequest, err.Error())
		return
	case "dddd":
		if userInfo.Permission&(1<<int64(notice.NoticeType)) == 0 {
			c.String(http.StatusBadRequest, "权限不足")
			return
		}
	default:
		fmt.Println(notice.Class)
		if userInfo.Class != -1 && "2021211"+notice.Class != strconv.Itoa(int(userInfo.Class)) {
			c.String(http.StatusBadRequest, "权限不足")
			return
		}
		if userInfo.Permission&(1<<20<<int64(notice.NoticeType)) == 0 {
			c.String(http.StatusBadRequest, "权限不足")
			return
		}
	}

	err = database.DB.Delete(&common.AdminNotice{Class: class, Index: index}).Error
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(200, true)
}
