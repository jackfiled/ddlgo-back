package tool

import (
	"bufio"
	"encoding/json"
	"os"
)

type Config struct {
	AppPort   string `json:"app_port"`
	JWTSecret string `json:"jwt_secret"`
}

var DefaultSetting = Config{
	AppPort:   ":8080",
	JWTSecret: "MakeBuptGreatAgain",
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
			DDLLog(err.Error())
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
