package handlers

import (
	"ddlBackend/models"
	"ddlBackend/tool"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"net/http"
	"regexp"
	"strings"
)

// UploadPictureHandler 上传文件处理函数
func UploadPictureHandler(context *gin.Context) {
	ok, err := tool.CheckPermission(context, models.Administrator)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !ok {
		context.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	file, err := context.FormFile("picture")

	if err != nil {
		tool.DDLLog(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 利用正则表达式判断文件名是否符合规则
	pattern, _ := regexp.Compile("^.*\\.((png)|(jpg))$")
	if !pattern.MatchString(file.Filename) {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Only png and jpg picture allowed",
		})
		return
	}

	// 获得图片的后缀名
	fileNames := strings.Split(file.Filename, ".")
	fileEndName := fileNames[len(fileNames)-1]

	// 利用uuid生成图片文件的新名称
	fileName := "./picture/" + uuid.NewV4().String() + "." + fileEndName
	err = context.SaveUploadedFile(file, fileName)
	if err != nil {
		tool.DDLLog(err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		// 这里只取文件名称第一位之后的值
		// 第一位是. 是没有必要的
		"address": fileName[1:],
	})
	return
}
