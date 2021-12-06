package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"ddl/auth"
	"ddl/query"
)

func main() {
	http.HandleFunc("/", index) // index 为向 url发送请求时，调用的函数
	http.HandleFunc("/api/login", auth.LoginHandler)
	http.HandleFunc("/api/logout", auth.LogoutHandler)
	http.HandleFunc("/api/get_list", query.GetListHandler)
	http.HandleFunc("/api/query_single", query.QuerySingleHandler)

	http.HandleFunc("/WW_verify_udfdZsIBL9yNi4SN.txt", WWVerify)

	http.HandleFunc("/api/auth_demo", auth.AuthDemoHandler)
	http.HandleFunc("/api/set_permission_demo", setPermission)
	log.Fatal(http.ListenAndServe(":8000", nil)) //防火墙并没有拦截:)

}
func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "DDL系统API")
}

//微信扫码授权验证
func WWVerify(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "udfdZsIBL9yNi4SN")
}

func setPermission(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if values.Get("studentID") == "" {
		fmt.Fprint(w, "参数错误")
		return
	}
	studentID, _ := strconv.ParseInt(values.Get("studentID"), 10, 32)
	permission, _ := strconv.ParseInt(values.Get("permission"), 10, 64)
	userInfo := auth.SetUserPermission(int32(studentID), permission)
	fmt.Fprintf(w, "%d权限修改为%d", userInfo.StudentID, userInfo.Permission)
}
