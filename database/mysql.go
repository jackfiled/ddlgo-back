package database

import "fmt"

// MysqlModel Mysql数据库连接信息类
type MysqlModel struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Address      string `json:"address"`
	DatabaseName string `json:"database_name"`
}

func (model *MysqlModel) GenerateConnectionString() string {
	result := fmt.Sprintf("%s:%s@tcp(%s)/%s", model.Username, model.Password, model.Address, model.DatabaseName)
	// 配置字符串
	result = result + "?charset=utf8mb4&parseTime=True&loc=Local"
	return result
}
