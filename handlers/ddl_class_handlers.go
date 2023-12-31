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

// CreateClassDDLHandler 班级url下创建DDL事件处理函数
func CreateClassDDLHandler(context *gin.Context) {
	className := context.Param("class")

	ok, err := checkClassAdminPermission(context, className)
	if err != nil {
		// 解析令牌和验证权限中遇到问题
		// 返回 500 服务器错误
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !ok {
		context.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	var ddlNotice models.DDLNotice
	err = context.ShouldBindJSON(&ddlNotice)
	if err != nil {
		// 请求体绑定失败
		// 返回 400 错误请求
		tool.DDLLogError(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if ddlNotice.NoticeType == models.ALL {
		// 类别设置为 1 全部
		// 不允许创建
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "all notice type of ddl is prohibited",
		})
		return
	}

	if className != ddlNotice.ClassName {
		// 请求url和请求体中班级不符
		// 返回 400 错误请求
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "the url and the body are not the same",
		})
		return
	}

	var db *gorm.DB
	db, err = database.GetDDLTable(className)
	if err != nil {
		// 无法获取对应班级的数据库
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
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, ddlNotice)
}

// ReadClassDDLHandler 班级url下读取DDL事件处理函数
func ReadClassDDLHandler(context *gin.Context) {
	className := context.Param("class")

	start := context.DefaultQuery("start", "0")
	step := context.DefaultQuery("step", "20")
	noticeType := context.DefaultQuery("noticeType", "0")

	// 将过滤器参数从字符串转换为数字
	var err error
	var startNum, stepNum, noticeTypeNum int
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
	noticeTypeNum, err = strconv.Atoi(noticeType)
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
	var db *gorm.DB
	db, err = database.GetDDLTable(className)
	if err != nil {
		// 无法获取对应班级的数据库
		// 返回 400 错误请求
		tool.DDLLogError(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if noticeTypeNum == models.ALL {
		db.Where("notice_type != ?", models.ALL).Offset(startNum).Limit(stepNum).Find(&ddlNotices)
	} else {
		db.Where("notice_type = ?", noticeTypeNum).Offset(startNum).Limit(stepNum).Find(&ddlNotices)
	}
	context.JSON(http.StatusOK, ddlNotices)
}

// checkClassAdminPermission 检查当前请求令牌的持有者是否有权限修改当前班级的内容
func checkClassAdminPermission(context *gin.Context, classname string) (bool, error) {
	claims, err := tool.GetClaimsInContext(context)

	if err != nil {
		return false, err
	}

	// 权限为根管理员
	// 或者为本班的管理员
	// 即可修改该班的内容
	if claims.Permission > models.Administrator && (claims.ClassName == classname || claims.Permission > models.User) {
		return true, nil
	} else {
		return false, nil
	}
}
