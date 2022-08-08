package handlers

import (
	"ddlBackend/database"
	"ddlBackend/log"
	"ddlBackend/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func ReadClassIDDDLHandler(context *gin.Context) {
	className := context.Param("class")
	id := context.Param("id")

	var db *gorm.DB
	var err error
	db, err = database.GetDDLTable(className)
	if err != nil {
		// 获取指定班级的数据库失败
		// 返回 400 错误请求
		log.DDLLog(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var idNum int
	idNum, err = strconv.Atoi(id)
	if err != nil {
		// 转换id字符串失败
		// 返回 400 请求错误
		log.DDLLog(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var ddlNotice models.DDLNotice
	result := db.Where("id = ?", idNum).Find(&ddlNotice)
	if result.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"error": result.Error.Error(),
		})
	} else {
		context.JSON(http.StatusOK, ddlNotice)
	}
	return
}

func UpdateClassIDDDLHandler(context *gin.Context) {
	className := context.Param("class")
	id := context.Param("id")

	var db *gorm.DB
	var err error
	db, err = database.GetDDLTable(className)
	if err != nil {
		// 获取指定班级的数据库失败
		// 返回 400 错误请求
		log.DDLLog(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var idNum int
	idNum, err = strconv.Atoi(id)
	if err != nil {
		// 转换id字符串失败
		// 返回 400 请求错误
		log.DDLLog(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var ddlNotice models.DDLNotice
	result := db.Where("id = ?", idNum).Find(&ddlNotice)

	if result.Error != nil {
		// 读取指定的事件失败
		// 返回 404 未找到错误
		context.JSON(http.StatusNotFound, gin.H{
			"err": result.Error.Error(),
		})
		return
	}

	err = context.ShouldBindJSON(&ddlNotice)
	if err != nil {
		// 绑定请求体失败
		// 返回 400 请求错误
		log.DDLLog(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if ddlNotice.ID != uint(idNum) {
		// 请求体中的ID和url中的ID不符
		// 返回 400 请求错误
		context.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("the id %d in the url and %d in the body are not the same", ddlNotice.ID, idNum),
		})
		return
	}

	db.Save(&ddlNotice)
	context.JSON(http.StatusNoContent, gin.H{})
	return
}

func DeleteClassIDDDLHandler(context *gin.Context) {
	className := context.Param("class")
	id := context.Param("id")

	var db *gorm.DB
	var err error
	db, err = database.GetDDLTable(className)
	if err != nil {
		// 获取指定班级的数据库失败
		// 返回 400 错误请求
		log.DDLLog(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var idNum int
	idNum, err = strconv.Atoi(id)
	if err != nil {
		// 转换id字符串失败
		// 返回 400 请求错误
		log.DDLLog(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var ddlNotice models.DDLNotice
	result := db.Where("id = ?", idNum).Find(&ddlNotice)

	if result.Error != nil {
		// 读取指定的事件失败
		// 返回 404 未找到错误
		context.JSON(http.StatusNotFound, gin.H{
			"err": result.Error.Error(),
		})
		return
	}

	db.Delete(&ddlNotice)
	context.JSON(http.StatusNoContent, gin.H{})
	return
}
