package auth

import (
	"ddl/common"
	"ddl/config"
	"ddl/wecom"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-basic/uuid"
	"github.com/gorilla/websocket"
)

type Session struct {
	Time     time.Time
	UserInfo common.UserInfo
	Wsconn   *websocket.Conn
}

var lock sync.RWMutex
var sessions = make(map[string]Session)

var upgrader = websocket.Upgrader{
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

func QRCodeWechatLoginHandler(c *gin.Context) {
	// fmt.Println(values)
	code := c.Request.FormValue("code")
	id := c.Request.FormValue("id")
	ref := c.Request.FormValue("ref")
	if ref == "" {
		ref = "http://squidward.top/"
	}
	if code == "" {
		c.String(http.StatusBadRequest, "参数错误")
	} else {
		userID := wecom.GetWecomID(code)
		// fmt.Println(userID)
		if userID != "" {
			lock.Lock()
			session, exist := sessions[id]
			if exist {
				session.UserInfo = GetUserInfo(userID)
				sessions[id] = session
			} else {
				ref = "http://squidward.top/"
			}

			lock.Unlock()

			c.Header("Content-Type", "text/html")
			c.String(200, "<script language='javascript'>window.location.href='"+ref+"'</script>")

			// fmt.Println(userInfo)
			// fmt.Fprintf(w, "%s %d %d %d", userInfo.UserID, userInfo.StudentID, userInfo.Class, userInfo.Permission)
		}
	}
}

func QRCodeWSHandler(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("upgrade:", err)
		return
	}
	defer ws.Close()

	lock.Lock()
	id := uuid.New()
	sessions[id] = Session{Time: time.Now(), Wsconn: ws}
	lock.Unlock()

	for {
		// _, exist := sessions[id]
		// if !exist {
		// 	break
		// }
		mt, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			break
		}
		// fmt.Printf("recv: %s\n", message)
		if string(message) == "?" {
			lock.Lock()
			session := sessions[id]
			msgSend := "wait|" + id + "|" + time.Now().String()
			finished := false
			if session.UserInfo.StudentID != 0 {
				session.UserInfo.ExpTime = time.Now().AddDate(0, 0, 14)
				data, _ := json.Marshal(session.UserInfo)
				fmt.Println(string(data))
				value, _ := Encrypt(string(data), []byte(config.ENCRYPT_KEY))
				msgSend = "ok|" + value
				finished = true
			}
			sessions[id] = session

			lock.Unlock()
			err = ws.WriteMessage(mt, []byte(msgSend))
			if err != nil {
				fmt.Println("write:", err)
				break
			}

			if finished {
				break
			}

		}
	}
}

func QRCodeWSHeatBeat() {
	for {
		lock.Lock()
		now := time.Now()
		for id, s := range sessions {
			// fmt.Println(sessions)
			if now.After(s.Time.Add(time.Minute)) {
				s.Wsconn.Close()
				delete(sessions, id)
				// fmt.Println("del")
			}
		}
		lock.Unlock()
		time.Sleep(time.Second * 30)
		// fmt.Println("-")
	}

}

func QRCodeLoginHandler(c *gin.Context) {
	// fmt.Println(values)
	cookie := c.Request.FormValue("cookie")
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "UserInfo",
		Value:    cookie,
		Path:     "/",
		Domain:   "",
		MaxAge:   14 * 24 * 3600,
		Secure:   true,
		HttpOnly: false,
		SameSite: http.SameSiteNoneMode,
	})
	c.JSON(200, true)
}
