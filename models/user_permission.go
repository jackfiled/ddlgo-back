package models

const (
	// User 用户 只可读取信息
	User = iota
	// Administrator 管理员 可控制小班的内容
	Administrator
	// Root 根管理员 可控制所有的内容
	Root
)
