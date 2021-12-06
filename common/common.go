package common

import (
	"time"
)

type DDLNotice struct {
	Index      int       `gorm:"column:index_;primary_key" json:"index_"`
	Time       time.Time `gorm:"column:time" json:"time"`
	DDL        time.Time `gorm:"column:ddl" json:"ddl"`
	StartTime  time.Time `gorm:"column:startTime" json:"startTime"`
	State      int       `gorm:"column:state" json:"state"`
	Title      string    `gorm:"column:title" json:"title"`
	Detail     string    `gorm:"column:detail" json:"detail"`
	NoticeType int       `gorm:"column:noticeType" json:"noticeType"`
}
