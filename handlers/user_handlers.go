package handlers

import (
	"ddlBackend/database"
	"ddlBackend/models"
	"ddlBackend/tool"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// ReadUsersHandler 读取所有用户信息的处理函数
func ReadUsersHandler(context *gin.Context) {
	var users []models.UserInformation

	result := database.Database.Table("user_informations").Find(&users)
	if result.Error != nil {
		// 如果读取中出错
		tool.DDLLog(result.Error.Error())
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, users)
	return
}

// ReadSingleUserHandler 读取单个用户信息处理函数
func ReadSingleUserHandler(context *gin.Context) {
	var user models.UserInformation
	id := context.Param("id")
	idNum, err := strconv.Atoi(id)
	if err != nil {
		// url中参数读取出错
		tool.DDLLog(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	result := database.Database.Table("user_informations").First(&user, idNum)
	if result.Error != nil {
		// 没找到
		tool.DDLLog(result.Error.Error())
		context.JSON(http.StatusNotFound, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, user)
	return
}

// CreateUserHandler 创建用户处理函数
func CreateUserHandler(context *gin.Context) {
	var user models.UserInformation
	err := context.ShouldBindJSON(&user)
	if err != nil {
		// 绑定json对象出错
		tool.DDLLog(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result := database.Database.Table("user_informations").Create(&user)
	if result.Error != nil {
		tool.DDLLog(result.Error.Error())
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, user)
	return
}

// UpdateUserHandler 更新用户信息
func UpdateUserHandler(context *gin.Context) {
	id := context.Param("id")
	idNum, err := strconv.Atoi(id)
	if err != nil {
		tool.DDLLog(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user models.UserInformation
	err = context.ShouldBindJSON(&user)
	if err != nil {
		tool.DDLLog(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if user.ID != uint(idNum) {
		// 请求体和url参数不匹配
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("the id %d in the url and %d in the body are not the same", user.ID, idNum),
		})
		return
	}

	database.Database.Table("user_informations").Save(&user)
	context.JSON(http.StatusNoContent, gin.H{})
	return
}

func DeleteUserHandler(context *gin.Context) {
	id := context.Param("id")
	idNum, err := strconv.Atoi(id)
	if err != nil {
		tool.DDLLog(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user models.UserInformation
	result := database.Database.Table("user_informations").First(&user, idNum)
	if result.Error != nil {
		tool.DDLLog(result.Error.Error())
		context.JSON(http.StatusNotFound, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	database.Database.Table("user_informations").Delete(&user)
	context.JSON(http.StatusNoContent, gin.H{})
	return
}
