package handlers

import (
	"context"
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

	courses, _, err := grpcGetSemester(model)
	if err != nil {
		tool.DDLLogError(err.Error())
		// RPC中出错
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, courses)
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
