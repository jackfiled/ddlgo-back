package push

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

//touser := ""                            //企业号中的用户帐号，在zabbix用户Media中配置，如果配置不正常，将按部门发送。
//toparty := "1"                                              //企业号中的部门id。
var corpid = "ww8a5308483ff283cc" //企业号的标识

// var agentid = 1000002                                          //企业号中的应用id。
// var corpsecret = "EPQstC4qi51TcvtVQRzQ1HowUdJ4jrOG_cFgcIA160E" //企业号中的应用的Secret
var corpsecretMap = map[int]string{
	1000002: "EPQstC4qi51TcvtVQRzQ1HowUdJ4jrOG_cFgcIA160E",
	1000005: "tFqgBhMGuktPfuEuXmvkXktU-W6Oq7cLFAg_n-WbwYQ",
}

// var partyMap = map[string]string{
// 	"dddd": "1", //大班
// 	"304":  "2",
// 	"305":  "7",
// 	"306":  "8",
// 	"307":  "9",
// 	"308":  "10",
// 	"309":  "11",
// 	"test": "6", //测试
// }

var noticeTypeMap = map[int]string{
	0: "DDL",
	1: "活动",
	2: "思政活动",
	3: "文体活动",
	4: "志愿活动",
	5: "讲座",
	6: "竞赛",
}

type JSON struct {
	Access_token string `json:"access_token"`
}

type MESSAGES struct {
	Touser  string `json:"touser"`
	Toparty string `json:"toparty"`
	Msgtype string `json:"msgtype"`
	Agentid int    `json:"agentid"`
	Text    struct {
		//Subject string `json:"subject"`
		Content string `json:"content"`
	} `json:"text"`
	Safe int `json:"safe"`
}

type MESSAGESCRAD struct {
	Touser   string `json:"touser"`
	Toparty  string `json:"toparty"`
	Msgtype  string `json:"msgtype"`
	Agentid  int    `json:"agentid"`
	Textcard struct {
		//Subject string `json:"subject"`
		Title       string `json:"title"`
		Description string `json:"description"`
		URL         string `json:"url"`
		Btntxt      string `json:"btntxt"`
	} `json:"textcard"`
	Safe int `json:"safe"`
}

func Get_AccessToken(corpid, corpsecret string) string {
	gettoken_url := "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=" + corpid + "&corpsecret=" + corpsecret
	//print(gettoken_url)
	client := &http.Client{}
	req, _ := client.Get(gettoken_url)
	defer req.Body.Close()
	body, _ := ioutil.ReadAll(req.Body)
	//fmt.Printf("\n%q",string(body))
	var json_str JSON
	json.Unmarshal([]byte(body), &json_str)
	//fmt.Printf("\n%q",json_str.Access_token)
	return json_str.Access_token
}

func SendData(access_token, msg string) {
	send_url := "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=" + access_token
	//print(send_url)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", send_url, bytes.NewBuffer([]byte(msg)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("charset", "UTF-8")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

//构造信息结构
// func messages(touser string, toparty string, agentid int, content string) string {
// 	msg := MESSAGES{
// 		Touser:  touser,
// 		Toparty: toparty,
// 		Msgtype: "text",
// 		Agentid: agentid,
// 		Safe:    0,
// 		Text: struct {
// 			//Subject string `json:"subject"`
// 			Content string `json:"content"`
// 		}{Content: content},
// 	}
// 	sed_msg, _ := json.Marshal(msg)
// 	//fmt.Printf("%s",string(sed_msg))
// 	return string(sed_msg)
// }

// func messagesGroup(toparty string, agentid int, content string) string {
// 	msg := MESSAGES{
// 		Toparty: toparty,
// 		Msgtype: "text",
// 		Agentid: agentid,
// 		Safe:    0,
// 		Text: struct {
// 			//Subject string `json:"subject"`
// 			Content string `json:"content"`
// 		}{Content: content},
// 	}
// 	sed_msg, _ := json.Marshal(msg)
// 	//fmt.Printf("%s",string(sed_msg))
// 	return string(sed_msg)
// }

func messagesGroupCard(toparty string, agentid int, title string, description string, url string) string {
	msg := MESSAGESCRAD{
		Toparty: toparty,
		Msgtype: "textcard",
		Agentid: agentid,
		Safe:    0,
		Textcard: struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			URL         string `json:"url"`
			Btntxt      string `json:"btntxt"`
		}{Title: title, Description: description, URL: url},
	}
	sed_msg, _ := json.Marshal(msg)
	//fmt.Printf("%s",string(sed_msg))
	return string(sed_msg)
}

//发送消息 party:部门名 noticeType:通知类型 title:标题 dateTime:时间 detail:详情 url:卡片跳转链接
// func sendNotice(parties []string, noticeType int, title string, dateTime string, detail string, url string) {
// 	accessToken := Get_AccessToken(corpid, corpsecret)
// 	var partyIDs string
// 	for _, party := range parties {
// 		partyIDs += partyMap[party] + "|"
// 	}
// 	partyIDs = partyIDs[:len(partyIDs)-1]
// 	noticeTypeString := typeMap[noticeType]

// 	fmt.Println(partyIDs)
// 	// 序列化成json之后，\n会被转义，也就是变成了\\n，使用str替换，替换掉转义
// 	msg := messagesGroupCard(partyIDs, agentid, title, "<div class=\"gray\">"+dateTime+" | "+noticeTypeString+"</div><br>"+detail, url)
// 	msg = strings.Replace(msg, "\\\\", "\\", -1)

// 	//  fmt.Println(strings.Replace(msg,"\\\\","\\",-1))
// 	SendData(accessToken, msg)
// }

//partyID:部门ID noticeType:通知类型 title:标题 dateTime:时间 detail:详情 url:卡片跳转链接
func SendNotice(partyIDs string, noticeType string, title string, dateTime string, detail string, url string) {
	var agentid = 1000002
	accessToken := Get_AccessToken(corpid, corpsecretMap[agentid])

	fmt.Println(partyIDs)
	// 序列化成json之后，\n会被转义，也就是变成了\\n，使用str替换，替换掉转义
	msg := messagesGroupCard(partyIDs, agentid, title, "<div class=\"gray\">"+dateTime+" | "+noticeType+"</div><br>"+detail, url)
	msg = strings.Replace(msg, "\\\\", "\\", -1)

	//  fmt.Println(strings.Replace(msg,"\\\\","\\",-1))
	SendData(accessToken, msg)
}

//partyID:部门ID noticeType:通知类型 title:标题 dateTime:时间 detail:详情 url:卡片跳转链接
func SendNoticeActivity(partyIDs string, noticeType string, title string, dateTime string, detail string, url string) {
	var agentid = 1000005
	accessToken := Get_AccessToken(corpid, corpsecretMap[agentid])

	fmt.Println(partyIDs)
	// 序列化成json之后，\n会被转义，也就是变成了\\n，使用str替换，替换掉转义
	msg := messagesGroupCard(partyIDs, agentid, title, "<div class=\"gray\">"+dateTime+" | "+noticeType+"</div><br>"+detail, url)
	msg = strings.Replace(msg, "\\\\", "\\", -1)

	//  fmt.Println(strings.Replace(msg,"\\\\","\\",-1))
	SendData(accessToken, msg)
}
