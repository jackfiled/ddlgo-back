package tool

import (
	"fmt"
	"time"
)

// DDLLogError 记录错误日志
func DDLLogError(message string) {
	timeString := time.Now().Format("01-02 15:04:05")
	fmt.Printf("%s DDL-Error: %s", timeString, message)
}
