package models

// UserInformation 用户信息模型
type UserInformation struct {
	// ID 数据库用户编号
	ID uint `gorm:"primaryKey" json:"id"`
	// Username 用户名
	Username string `json:"username"`
	// Password 用户密码
	Password string `json:"password"`
	// ClassName 所在班级
	ClassName string `json:"class_name"`
	// StudentID 学号
	StudentID string `json:"student_id"`
	// Permission 权限
	Permission uint `json:"permission"`
}

// AdminLoginModel 管理员登录JSON模型
type AdminLoginModel struct {
	StudentID string `json:"student_id"`
	Password  string `json:"password"`
}

// UserLoginModel 用户登录JSON模型
type UserLoginModel struct {
	Username  string `json:"username"`
	StudentID string `json:"student_id"`
}
