package tool

import (
	"bufio"
	"ddlBackend/database"
	"ddlBackend/models"
	"encoding/json"
	"os"
)

type Config struct {
	AppPort       string                 `json:"app_port"`
	JWTSecret     string                 `json:"jwt_secret"`
	PasswordSalt  string                 `json:"password_salt"`
	JWGLOutTime   int64                  `json:"jwgl_out_time"`
	JWGrpcAddress string                 `json:"jw_grpc_address"`
	RootConfig    models.UserInformation `json:"root_config"`
	UseMysql      bool                   `json:"use_mysql"`
	MysqlConfig   database.MysqlModel    `json:"mysql_config"`
}

// DefaultSetting 默认配置文件
var DefaultSetting = Config{
	AppPort:       ":8080",
	JWTSecret:     "MakeBuptGreatAgain",
	JWGLOutTime:   24,
	PasswordSalt:  "MakeBuptGreatAgain",
	JWGrpcAddress: "rrricardo.top:7000",
	RootConfig: models.UserInformation{
		Username:   "root",
		Password:   "123456",
		ClassName:  "dddd",
		StudentID:  "0000000000",
		Permission: models.Root,
	},
	UseMysql: false,
}

// Setting 配置文件对象
var Setting Config

// ReadConfig 读取配置文件
func ReadConfig() error {
	file, err := os.Open("config.json")
	// 函数执行完成后关闭文件
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			// 关闭文件中错误
			DDLLogError(err.Error())
		}
	}(file)

	if err != nil {
		// 读取配置文件错误
		// 采用默认配置
		Setting = DefaultSetting
		return err
	}

	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&Setting)
	if err != nil {
		// 解析json失败
		// 采用默认配置
		Setting = DefaultSetting
		return err
	}

	return nil
}
