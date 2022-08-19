package handlers

import (
	"ddlBackend/database"
	"ddlBackend/models"
	"ddlBackend/tool"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ReadUsersHandler 读取所有用户信息的处理函数
func ReadUsersHandler(context *gin.Context) {
	ok, err := tool.CheckPermission(context, models.Root)
	if err != nil {
		// 验证权限过程中出错
		// 返回 500 服务器错误
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !ok {
		// 权限不够
		// 返回 401 未授权
		context.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	var users []models.UserInformation

	result := database.Database.Table("user_informations").Find(&users)
	if result.Error != nil {
		// 如果读取中出错
		tool.DDLLogError(result.Error.Error())
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, users)
}

// ReadSingleUserHandler 读取单个用户信息处理函数
func ReadSingleUserHandler(context *gin.Context) {
	ok, err := tool.CheckPermission(context, models.Root)
	if err != nil {
		// 验证权限过程中出错
		// 返回 500 服务器错误
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !ok {
		// 权限不够
		// 返回 401 未授权
		context.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	var user models.UserInformation
	id := context.Param("id")
	idNum, err := strconv.Atoi(id)
	if err != nil {
		// url中参数读取出错
		tool.DDLLogError(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	result := database.Database.Table("user_informations").First(&user, idNum)
	if result.Error != nil {
		// 没找到
		tool.DDLLogError(result.Error.Error())
		context.JSON(http.StatusNotFound, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, user)
}

// CreateUserHandler 创建用户处理函数
func CreateUserHandler(context *gin.Context) {
	ok, err := tool.CheckPermission(context, models.Root)
	if err != nil {
		// 验证权限过程中出错
		// 返回 500 服务器错误
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !ok {
		// 权限不够
		// 返回 401 未授权
		context.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	var user models.UserInformation
	err = context.ShouldBindJSON(&user)
	if err != nil {
		// 绑定json对象出错
		tool.DDLLogError(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 将密码加盐哈希之后储存在数据库中
	user.Password = tool.Sha256PasswordWithSalt(user.Password)

	result := database.Database.Table("user_informations").Create(&user)
	if result.Error != nil {
		tool.DDLLogError(result.Error.Error())
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, user)
}

// UpdateUserHandler 更新用户信息
func UpdateUserHandler(context *gin.Context) {
	ok, err := tool.CheckPermission(context, models.Root)
	if err != nil {
		// 验证权限过程中出错
		// 返回 500 服务器错误
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !ok {
		// 权限不够
		// 返回 401 未授权
		context.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	id := context.Param("id")
	idNum, err := strconv.Atoi(id)
	if err != nil {
		tool.DDLLogError(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user models.UserInformation
	err = context.ShouldBindJSON(&user)
	if err != nil {
		tool.DDLLogError(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if user.ID != uint(idNum) {
		// 请求体和url参数不匹配
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("the id %d in the url and %d in the body are not the same", idNum, user.ID),
		})
		return
	}

	// 密码加盐哈希之后再存入数据库
	user.Password = tool.Sha256PasswordWithSalt(user.Password)

	database.Database.Table("user_informations").Save(&user)
	context.JSON(http.StatusNoContent, gin.H{})
}

// DeleteUserHandler 删除用户信息处理函数
func DeleteUserHandler(context *gin.Context) {
	ok, err := tool.CheckPermission(context, models.Root)
	if err != nil {
		// 验证权限过程中出错
		// 返回 500 服务器错误
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !ok {
		// 权限不够
		// 返回 401 未授权
		context.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	id := context.Param("id")
	idNum, err := strconv.Atoi(id)
	if err != nil {
		tool.DDLLogError(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user models.UserInformation
	result := database.Database.Table("user_informations").First(&user, idNum)
	if result.Error != nil {
		tool.DDLLogError(result.Error.Error())
		context.JSON(http.StatusNotFound, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	database.Database.Table("user_informations").Delete(&user)
	context.JSON(http.StatusNoContent, gin.H{})
}

// AdminLoginHandler 管理员登录处理函数
func AdminLoginHandler(context *gin.Context) {
	var loginModel models.AdminLoginModel

	err := context.ShouldBindJSON(&loginModel)
	if err != nil {
		// 绑定JSON出错
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	user, err := database.AdminLogin(loginModel.StudentID, loginModel.Password)
	if err != nil {
		// 数据库中查无此用户
		context.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := tool.GenerateJWTToken(*user)
	if err != nil {
		// 产生JWT令牌中错误
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"token":      token,
		"username":   user.Username,
		"class_name": user.ClassName,
		"student_id": user.StudentID,
	})
}

// UserLoginHandler 用户登录处理函数
func UserLoginHandler(context *gin.Context) {
	var loginModel models.UserLoginModel

	err := context.ShouldBindJSON(&loginModel)
	if err != nil {
		// 绑定JSON出错
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := database.UserLogin(loginModel.Username, loginModel.StudentID)
	if err != nil {
		// 数据库中查无此人
		context.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	if user.Permission > models.User {
		// 如果该用户的权限大于普通用户
		// 则不能通过学号姓名验证的方式获得令牌
		// 返回 403 Forbidden
		context.JSON(http.StatusForbidden, gin.H{
			"error": "该用户身份为管理员，请通过学号密码登录",
		})
		return
	}

	token, err := tool.GenerateJWTToken(*user)
	if err != nil {
		// 产生令牌出错
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"token":      token,
		"username":   user.Username,
		"class_name": user.ClassName,
		"student_id": user.StudentID,
	})
}

// AdminUpdatePasswordHandler 管理员修改密码处理函数
func AdminUpdatePasswordHandler(context *gin.Context) {
	claims, err := tool.GetClaimsInContext(context)
	if err != nil {
		// 解析JWT令牌信息错误
		// 返回 500 服务器错误
		tool.DDLLogError(err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 修改管理员密码的请求体保持和登录一致
	// 需要学号和新的密码
	var loginModel models.AdminLoginModel
	err = context.ShouldBindJSON(&loginModel)
	if err != nil {
		// 绑定JSON出错
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	var user models.UserInformation
	result := database.Database.Table("user_informations").Where("student_id =?", loginModel.StudentID).Find(&user)
	if result.Error != nil {
		// 数据库中未发现该用户
		context.JSON(http.StatusNotFound, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	// 修改密码的权限要求：
	// 令牌持有账号和欲修改密码的账号是同一个账号
	// 令牌的权限在管理员及以上
	if claims.StudentID == loginModel.StudentID && claims.Permission >= models.Administrator {
		user.Password = loginModel.Password
		database.Database.Table("user_informations").Save(&user)

		context.JSON(http.StatusOK, gin.H{})
		return
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "permission denied",
		})
		return
	}
}
