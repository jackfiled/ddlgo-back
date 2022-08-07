package models

import (
	"gorm.io/gorm"
	"time"
)

// DDLNotice DDL事件模型
type DDLNotice struct {
	// gorm约定
	gorm.Model
	// Title DDL事件的标题
	Title string
	// Detail DDL事件的详情
	Detail string
	// DDLTime DDL时间
	DDLTime time.Time
	// NoticeType DDL的分类
	NoticeType int
}
