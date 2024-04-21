package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go-websocket/private-chat/model"

	"github.com/gorilla/websocket"
)

var (
	ws      = websocket.Upgrader{}
	userMap map[int64]*websocket.Conn
)

func init() {
	userMap = make(map[int64]*websocket.Conn)
}

func main() {
	http.HandleFunc("/privateChat", privateChat)
	_ = http.ListenAndServe(":9900", nil)
}

func privateChat(w http.ResponseWriter, r *http.Request) {
	c, err := ws.Upgrade(w, r, nil) //升级将 HTTP 服务器连接升级到 WebSocket 协议
	if err != nil {
		log.Printf("upgrade err:%s\n", err)
		return
	}

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("read message err:%s\n", err)
			continue
		}

		var chat model.Chat
		if err := json.Unmarshal(message, &chat); err != nil {
			log.Printf("unmarshal err:%s\n", err)
			continue
		}

		switch chat.Event {
		case model.Register:
			//设置userID
			chat.User.UserId = time.Now().Unix()
			//保存用户连接
			userMap[chat.User.UserId] = c
			break
		case model.SendMsg:
			chat.Message.CreateTime = time.Now().Format("2006-01-02 15:04:05")
			//拿到用户连接
			ok := false
			c, ok = userMap[chat.Message.Receiver.UserId]
			if !ok {
				//如果没有，拿到发送方用户的连接，告诉他不行
				c, _ = userMap[chat.Message.SendUser.UserId]
				chat.Message.Receiver = chat.Message.SendUser
				chat.Message.Content = "发送失败"
			}
			break
		default:
			c, _ = userMap[chat.Message.SendUser.UserId]
			chat.Message.Receiver = chat.Message.SendUser
			chat.Message.Content = "消息类型不对"
		}

		log.Printf("now chat : %+v \n", chat)

		//响应数据
		bytes, err := json.Marshal(chat)
		if err != nil {
			log.Printf("marshal err:%s\n", err)
			continue
		}
		if err := c.WriteMessage(mt, bytes); err != nil {
			log.Printf("write message err:%s\n", err)
			continue
		}
	}
}