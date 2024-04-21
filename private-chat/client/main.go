package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"go-websocket/private-chat/model"

	"github.com/gorilla/websocket"
	"github.com/spf13/cast"
)

const (
	PrivateChatUrl = "ws://127.0.0.1:9900/privateChat"
)

var (
	user model.User
	name string
	port int
)

func init() {
	flag.StringVar(&name, "name", "", "user name")
	flag.IntVar(&port, "port", 8801, "server port")
}

func main() {
	flag.Parse()

	c, _, err := websocket.DefaultDialer.Dial(PrivateChatUrl, nil)
	if err != nil {
		log.Fatal("dial:", err)
		return
	}

	//注册
	u := model.Chat{Event: model.Register, User: model.User{UserName: name}}
	bytes, _ := json.Marshal(u)
	_ = c.WriteMessage(websocket.TextMessage, bytes)

	//读取消息监听
	go func() {
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Printf("read message err:%s\n", err)
				continue
			}
			var res model.Chat
			if err := json.Unmarshal(msg, &res); err != nil {
				log.Printf("unmarshal err:%s\n", err)
				continue
			}
			switch res.Event {
			case model.Register:
				user = res.User
				fmt.Printf("register success , now user %v \n", user)
			case model.SendMsg:
				fmt.Printf("用户 %s 在 %s 给你发送了一条消息:%s \n",
					res.Message.SendUser.UserName, res.Message.CreateTime, res.Message.Content)
			}
		}
	}()

	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		uid := r.URL.Query().Get("uid")
		content := r.URL.Query().Get("content")
		msg, _ := json.Marshal(model.Chat{
			Event: model.SendMsg,
			Message: model.Message{
				SendUser: &user,
				Receiver: &model.User{UserId: cast.ToInt64(uid)},
				Content:  content,
			},
		})
		_ = c.WriteMessage(websocket.TextMessage, msg)
		_, _ = w.Write([]byte("ok"))
	})
	_ = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
