package database

import (
	"ddlBackend/models"
	"ddlBackend/tool"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DDLTables DDL表列表
var DDLTables = [7]string{"dddd", "304", "305", "306", "307", "308", "309"}

// Database 数据库指针
var Database *gorm.DB

// OpenDatabase 打开数据库函数
func OpenDatabase() (err error) {
	Database, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	// 创建ddl相关的表
	for _, value := range DDLTables {
		if !Database.Migrator().HasTable(value) {
			err = Database.Table(value).AutoMigrate(&models.DDLNotice{})
			if err != nil {
				return err
			}
		}
	}

	// 创建用户表
	if !Database.Migrator().HasTable(&models.UserInformation{}) {
		err = Database.AutoMigrate(&models.UserInformation{})
		if err != nil {
			return err
		}
	}

	// 创建记录ICS信息表
	if !Database.Migrator().HasTable(&models.ICSInformation{}) {
		err = Database.AutoMigrate(&models.ICSInformation{})
		if err != nil {
			return err
		}
	}

	// 将配置文件中设置的根管理员存入数据库
	_, err = AdminLogin(tool.Setting.RootConfig.StudentID, tool.Setting.RootConfig.Password)
	if err != nil {
		// 引发没有找到的错误
		Database.Table("user_informations").Create(&tool.Setting.RootConfig)
	}

	return nil
}

// GetDDLTable 获取指定班级的DDL事件表
func GetDDLTable(className string) (*gorm.DB, error) {
	for _, value := range DDLTables {
		if className == value {
			return Database.Table(className).Order("ddl_time DESC"), nil
		}
	}
	return nil, fmt.Errorf("the table named %s not exists", className)
}

// AdminLogin 管理员登录验证函数
func AdminLogin(studentID string, password string) (*models.UserInformation, error) {
	password = tool.Sha256PasswordWithSalt(password)

	var user models.UserInformation
	result := Database.Table("user_informations").Where("student_id = ? AND password = ?", studentID, password).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

// UserLogin 用户登录函数
func UserLogin(username string, studentID string) (*models.UserInformation, error) {
	var user models.UserInformation
	result := Database.Table("user_informations").Where("username = ? AND student_id = ?", username, studentID).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

// GetICSInformation 获得ICSInformation
func GetICSInformation(studentID string, semester string) (*models.ICSInformation, error) {
	var info models.ICSInformation
	result := Database.Table("ics_informations").Where("student_id = ? AND semester = ?", studentID, semester).First(&info)

	if result.Error != nil {
		return nil, result.Error
	}

	return &info, nil
}
