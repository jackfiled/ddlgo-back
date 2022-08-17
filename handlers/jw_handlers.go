package handlers

import (
	"context"
	"ddlBackend/database"
	"ddlBackend/models"
	"ddlBackend/protos"
	"ddlBackend/tool"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	info, err := database.GetICSInformation(model.StudentID, model.Semester)
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
		database.Database.Table("ics_informations").Create(&newInformation)

		context.JSON(http.StatusOK, courses)
		return
	} else {
		duration := time.Since(info.UpdatedAt)
		// 如果上次请求到现在的时间小于超时时间
		if int64(duration) <= tool.Setting.JWGLOutTime*int64(time.Hour) {
			targetTime := info.UpdatedAt.Add(time.Duration(tool.Setting.JWGLOutTime * int64(time.Hour)))
			context.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Too frequent requests, please try after %s", targetTime.Format("06-01-02 15:04")),
			})
		} else {
			// 已经超过超时时间
			courses, icsStream, err := grpcGetSemester(model)
			if err != nil {
				tool.DDLLogError(err.Error())
				// RPC中出错
				context.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			info.ICSStream = icsStream
			database.Database.Table("ics_informations").Save(info)

			context.JSON(http.StatusOK, courses)
			return
		}
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
}

// grpcGetSemester 远程过程调用获得课表的函数
func grpcGetSemester(model models.GetSemesterCalendarModel) ([]*protos.Course, []byte, error) {
	connection, err := grpc.Dial(tool.Setting.JWGrpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
