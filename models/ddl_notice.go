package models

import (
	"time"
)

// DDLNotice DDL事件模型
type DDLNotice struct {
	// 数据库事件编号
	ID uint `gorm:"primaryKey" json:"id"`
	// Title DDL事件的标题
	Title string `json:"title"`
	// Detail DDL事件的详情
	Detail string `json:"detail"`
	// DDLTime DDL时间
	DDLTime time.Time `json:"ddl_time"`
	// NoticeType DDL的分类
	NoticeType int `json:"notice_type"`
}
