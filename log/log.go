package log

import (
	"fmt"
	"time"
)

// DDLLog 日志记录函数
func DDLLog(message string) {
	now := time.Now()
	fmt.Printf("%v  %s\n", now, message)
}
