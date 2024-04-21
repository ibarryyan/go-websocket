package model

type Event int

const (
	Register Event = 1 //注册事件
	SendMsg        = 2 //发送消息事件
)

type Chat struct {
	Event   Event
	Message Message
	User    User
}

type Message struct {
	SendUser   *User
	Receiver   *User
	Content    string
	CreateTime string
}

type User struct {
	UserId   int64
	UserName string
	Address  string
}
