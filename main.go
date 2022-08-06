package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	route := gin.Default()

	route.GET("/ping", func(context *gin.Context) {
		context.JSONP(200, gin.H{
			"message": "pong",
		})
	})

	err := route.Run()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
