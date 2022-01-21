package admin

import (
	"ddl/auth"
	"ddl/common"
	"ddl/database"
	"ddl/push"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func SaveHandler(c *gin.Context) {
	var err error

	userInfo, err := auth.GetCookieUserInfo(c)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			c.String(http.StatusBadRequest, err.Error())
		} else {
			c.String(http.StatusInternalServerError, err.Error())
		}
	}

	form := c.Request.Form

	noticeType, err := strconv.Atoi(form.Get("noticeType"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	switch form.Get("classes") {
	case "":
		c.String(http.StatusBadRequest, err.Error())
		return
	case "dddd":
		if userInfo.Permission&(1<<int64(noticeType)) == 0 {
			c.String(http.StatusBadRequest, "权限不足")
			return
		}
	default:
		if "2021211"+form.Get("classes") != strconv.Itoa(int(userInfo.Class)) {
			c.String(http.StatusBadRequest, "权限不足")
			return
		}
		if userInfo.Permission&(1<<16<<int64(noticeType)) == 0 {
			c.String(http.StatusBadRequest, "权限不足")
			return
		}
	}

	index := 0
	if form.Get("index_") != "" {
		index, err = strconv.Atoi(form.Get("index_"))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
	}

	ddl, err := time.Parse("2006-01-02 15:04", form.Get("ddl"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	notice := common.DDLNotice{
		Index:      index,
		Title:      form.Get("title"),
		Detail:     form.Get("detail"),
		DDL:        ddl,
		NoticeType: noticeType,
		Img:        form.Get("img"),
	}
	if form.Get("startTime") != "" {
		startTime, err := time.Parse("2006-01-02 15:04", form.Get("ddl"))
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		notice.StartTime = startTime
	}

	err = database.DB.Table(form.Get("classes")).Save(&notice).Error
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	title := notice.Title
	if form.Get("index_") != "" {
		title = "【修改】" + title
	}
	push.SendNotice(common.PartyMap[form.Get("classes")], notice.NoticeType, title,
		form.Get("startTime")+"--"+form.Get("ddl"), notice.Detail,
		"http://www.squidward.top/?cla="+form.Get("classes")+"#"+strconv.Itoa(notice.Index))

	c.JSON(200, notice)
}

func DeleteHandler(c *gin.Context) {
	form := c.Request.Form

	userInfo, err := auth.GetCookieUserInfo(c)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			c.String(http.StatusBadRequest, err.Error())
		} else {
			c.String(http.StatusInternalServerError, err.Error())
		}
	}

	var notice common.DDLNotice
	err = database.DB.Table(form.Get("classes")).Where("index_=?", form.Get("index_")).First(&notice).Error
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	switch form.Get("classes") {
	case "":
		c.String(http.StatusBadRequest, err.Error())
		return
	case "dddd":
		if userInfo.Permission&(1<<int64(notice.NoticeType)) == 0 {
			c.String(http.StatusBadRequest, "权限不足")
			return
		}
	default:
		if "2021211"+form.Get("classes") != strconv.Itoa(int(userInfo.Class)) {
			c.String(http.StatusBadRequest, "权限不足")
			return
		}
		if userInfo.Permission&(1<<16<<int64(notice.NoticeType)) == 0 {
			c.String(http.StatusBadRequest, "权限不足")
			return
		}
	}

	err = database.DB.Table(form.Get("classes")).Delete("index_=?", form.Get("index_")).Error
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(200, true)
}
