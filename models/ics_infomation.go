package models

import "time"

// ICSInformation 存储获得ICS文件以及相关的信息
type ICSInformation struct {
	ID uint `gorm:"primaryKey"`
	// StudentID 学号
	StudentID string
	// Semester 该课表的所在学期
	Semester string
	// ICSStream ICS文件字节流
	ICSStream []byte
	// UpdatedAt 更新的时间
	UpdatedAt time.Time
}

// GetSemesterCalendarModel 获得课表的请求体
type GetSemesterCalendarModel struct {
	StudentID string `json:"student_id"`
	Password  string `json:"password"`
	Semester  string `json:"semester"`
}
