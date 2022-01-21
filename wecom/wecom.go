package wecom

import (
	"ddl/common"
	"ddl/config"
	"encoding/json"
	"fmt"
)

type GetAccessTokenRes struct {
	Access_token string
}

type GetUserIDRes struct {
	UserID string
	Errmsg string
}

//微信登录部分

func GetAccessToken(corpsecret string) string {
	data := common.HttpGet("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=" + config.CORPID + "&corpsecret=" + corpsecret)
	res := GetAccessTokenRes{}
	err := json.Unmarshal([]byte(data), &res)
	if err != nil {
		fmt.Println("GetAccessToken failed")
		return ""
	} else {
		return res.Access_token
	}
}

func GetWecomID(code string) string {
	accessToken := GetAccessToken(config.CORPSECRET)
	data := common.HttpGet("https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo?access_token=" + accessToken + "&code=" + code)
	res := GetUserIDRes{}
	err := json.Unmarshal([]byte(data), &res)
	// fmt.Println(data)
	if err != nil || res.Errmsg != "ok" {
		fmt.Printf("GetWecomID failed\n%s\n", res.Errmsg)
		return ""
	} else {
		return res.UserID
	}
}
