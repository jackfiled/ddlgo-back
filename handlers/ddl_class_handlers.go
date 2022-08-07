package handlers

import (
	"ddlBackend/database"
	"ddlBackend/log"
	"ddlBackend/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// CreateClassDDLHandler 班级url下创建DDL事件处理函数
func CreateClassDDLHandler(context *gin.Context) {
	className := context.Param("class")

	var ddlNotice models.DDLNotice
	err := context.ShouldBindJSON(&ddlNotice)
	if err != nil {
		log.DDLLog(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if className != ddlNotice.ClassName {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "the url and the body are not the same",
		})
		return
	}

	result := database.Database.Table("ddl_notices").Create(&ddlNotice)
	if result.Error != nil {
		log.DDLLog(result.Error.Error())
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, ddlNotice)
	return
}

// ReadClassDDLHandler 班级url下读取DDL事件处理函数
func ReadClassDDLHandler(context *gin.Context) {
	className := context.Query("class")

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
	database.Database.Table("ddl_notices").Order("ddl_time DESC").Offset(startNum).Limit(stepNum).Where("class_name = ? AND notice_type = ?", className, noticeType).Find(&ddlNotices)
	context.JSON(http.StatusOK, ddlNotices)
	return
}
