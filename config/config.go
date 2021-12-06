package config

//企业微信
const CORPID string = "ww8a5308483ff283cc"

//企业微信应用
const CORPSECRET string = "cahx1b3DGRqCMGumztUHgCpj9B8slCHGSHs_wZeknLg"

//数据库用户名密码
const DB_USER_PW string = "squidward_top:789456"

//数据库主机
const DB_HOST string = "lllccc.top:20570"

// const DB_HOST string = "localhost"

//测试密钥，实际环境请修改
const ENCRYPT_KEY string = "1234567890abcdef"

//密码登录
var AdminKey = map[string]int64{
	"admin":    267387135,
	"12345678": 0,
}
