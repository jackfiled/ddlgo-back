package handlers

import (
	"context"
	"ddlBackend/database"
	"ddlBackend/models"
	"ddlBackend/protos"
	"ddlBackend/tool"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
)

func GetSemesterCalendarHandler(context *gin.Context) {
	var model models.GetSemesterCalendarModel
	err := context.ShouldBindJSON(&model)
	if err != nil {
		// 绑定请求结构体出错
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, err = database.GetICSInformation(model.StudentID, model.Semester)
	if err != nil {
		// 说明没有请求过
		courses, icsStream, err := grpcGetSemester(model)
		if err != nil {
			tool.DDLLogError(err.Error())
			// RPC中出错
			context.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// 在数据库创建这个记录
		newInformation := models.ICSInformation{
			StudentID: model.StudentID,
			Semester:  model.Semester,
			ICSStream: icsStream,
		}
		database.Database.Table("user_informations").Create(&newInformation)

		context.JSON(http.StatusOK, courses)
		return
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "最近才请求过，请稍后再试",
		})
		return
	}
}

// GetICSFileHandler 返回ICS文件处理函数
func GetICSFileHandler(context *gin.Context) {
	studentID := context.Param("id")
	semester := context.Param("semester")
	// 去掉最后的.ics文件后缀
	semester = semester[:len(semester)-4]

	info, err := database.GetICSInformation(studentID, semester)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"error": "请先获得课表再尝试下载ICS日历文件",
		})
		return
	}

	context.Data(http.StatusOK, "text/calendar", info.ICSStream)
	return
}

// grpcGetSemester 远程过程调用获得课表的函数
func grpcGetSemester(model models.GetSemesterCalendarModel) ([]*protos.Course, []byte, error) {
	connection, err := grpc.Dial("localhost:7000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	defer func(connection *grpc.ClientConn) {
		err := connection.Close()
		if err != nil {
			// 如果关闭连接失败
			tool.DDLLogError(err.Error())
		}
	}(connection)

	client := protos.NewJwglerClient(connection)

	response, err := client.GetSemester(context.Background(), &protos.GetSemesterRequest{
		StudentID: model.StudentID,
		Password:  model.Password,
		Semester:  model.Semester,
	})

	if err != nil {
		// 调用远程函数出错
		return nil, nil, err
	}

	return response.Courses, response.IcsStream, nil
}
