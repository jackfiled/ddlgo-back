package handlers

import (
	"ddlBackend/database"
	"ddlBackend/log"
	"ddlBackend/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// CreateDDLHandler 创建DDL事件处理函数
func CreateDDLHandler(context *gin.Context) {
	var ddlNotice models.DDLNotice

	err := context.ShouldBindJSON(&ddlNotice)
	if err != nil {
		// 绑定json数据失败
		// 返回 400 错误请求
		log.DDLLog(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result := database.Database.Create(&ddlNotice)
	if result.Error != nil {
		// 在数据库中创建失败
		// 返回 500 服务器错误
		log.DDLLog(result.Error.Error())
		context.JSONP(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
	}

	context.JSON(http.StatusCreated, ddlNotice)
	return
}

// ReadDDLHandler 读取DDL事件处理函数
func ReadDDLHandler(context *gin.Context) {
	start := context.DefaultQuery("start", "0")
	step := context.DefaultQuery("step", "20")
	noticeType := context.DefaultQuery("noticeType", "0")

	// 将过滤器参数从字符串转换为数字
	var err error
	var startNum, stepNum int
	startNum, err = strconv.Atoi(start)
	stepNum, err = strconv.Atoi(step)
	if err != nil {
		log.DDLLog(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var ddlNotices []models.DDLNotice
	database.Database.Table("ddl_notices").Order("ddl_time DESC").Where("notice_type = ?", noticeType).Offset(startNum).Limit(stepNum).Find(&ddlNotices)

	context.JSON(http.StatusOK, ddlNotices)
	return
}
