package handlers

import (
	"ddlBackend/database"
	"ddlBackend/models"
	"ddlBackend/tool"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateDDLHandler 创建DDL事件处理函数
func CreateDDLHandler(context *gin.Context) {
	ok, err := tool.CheckPermission(context, models.Root)
	if err != nil {
		// 验证身份中出错
		// 返回 500 服务器错误
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !ok {
		// 验证未通过
		// 返回 401 未授权错误
		context.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	var ddlNotice models.DDLNotice

	err = context.ShouldBindJSON(&ddlNotice)
	if err != nil {
		// 绑定json数据失败
		// 返回 400 错误请求
		tool.DDLLogError(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var db *gorm.DB
	db, err = database.GetDDLTable(ddlNotice.ClassName)
	if err != nil {
		// 获取对应班级数据库失败
		// 返回 400 错误请求
		tool.DDLLogError(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result := db.Create(&ddlNotice)
	if result.Error != nil {
		// 在数据库中创建失败
		// 返回 500 服务器错误
		tool.DDLLogError(result.Error.Error())
		context.JSONP(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, ddlNotice)
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
	if err != nil {
		// 请求参数转换失败
		// 返回 400 错误请求
		tool.DDLLogError(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	stepNum, err = strconv.Atoi(step)
	if err != nil {
		// 请求参数转换失败
		// 返回 400 错误请求
		tool.DDLLogError(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var ddlNotices []models.DDLNotice
	for _, value := range database.DDLTables {
		db, _ := database.GetDDLTable(value)
		var list []models.DDLNotice
		db.Where("notice_type = ?", noticeType).Offset(startNum).Limit(stepNum).Find(&list)
		ddlNotices = append(ddlNotices, list...)
	}

	context.JSON(http.StatusOK, ddlNotices)
}
