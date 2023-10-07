package server

import (
	"encoding/json"
	"fmt"
	"net"
)

type User struct {
	Name   string
	Addr   string
	C      chan *Message
	conn   net.Conn
	server *Server // User关联对应的Server
}

// NewUser 创建一个用户的API（工厂方法）
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan *Message),
		conn:   conn,
		server: server,
	}
	//启动监听当前user channel消息的goroutine
	go user.MessageListener()
	return user
}

// MessageListener 监听当前 User 的 channel,一旦有消息，就直接发送给对端客户端
func (this *User) MessageListener() {
	for {
		msg := <-this.C
		//this.conn.Write([]byte(msg.Content + "\n"))
		marshal, err := json.Marshal(msg)
		marshalStr := string(marshal)
		fmt.Println(marshalStr)
		if err != nil {
			fmt.Printf("msg type = %T\n", msg)
			fmt.Printf("msg = %#v", msg)
			panic(err)
		}
		this.conn.Write([]byte(marshalStr + "\n"))
	}
}

func (this *User) Online() {
	// 用户上线，将用户加入到onlineMap中
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	// 广播当前用户上线消息
	sendMsg := &Message{MsgType: GlobalMsg, Content: this.Name + ":" + "已上线"}
	this.server.SendMsg(sendMsg)
}

func (this *User) Offline() {
	// 用户下线，将用户从onlineMap中去掉
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	// 广播当前用户上线消息
	sendMsg := &Message{MsgType: GlobalMsg, Content: this.Name + ":" + "已下线"}
	this.server.SendMsg(sendMsg)
}

func (this *User) Rename(name string) {
	delete(this.server.OnlineMap, this.Name)
	this.Name = name
	this.server.OnlineMap[this.Name] = this
}

func (this *User) doMeg(msgStr string) {
	var msg *Message
	err := json.Unmarshal([]byte(msgStr), &msg)
	if err != nil {
		fmt.Println("msgStr:", msgStr)
		fmt.Println("msg unmarshal err")
		return
	}

	switch msg.MsgType {

	case GlobalMsg: // 广播消息
		sendMsg := &Message{MsgType: GlobalMsg, Sender: this.Name, Content: this.Name + ":" + msg.Content}
		this.server.SendMsg(sendMsg)

	case SearchAllOnlineClientMsg: // 查询所有当前在线用户
		sendMsg := &Message{MsgType: SearchAllOnlineClientMsg, Sender: this.Name}
		this.server.SendMsg(sendMsg)

	case PrivateMsg: // 私聊
		sendMsg := &Message{MsgType: PrivateMsg, Sender: this.Name, Content: this.Name + ":" + msg.Content, Receiver: msg.Receiver}
		this.server.SendMsg(sendMsg)

	case RenameMsg: // 修改用户名
		oriName := this.Name
		// Todo: 校验用户名
		this.Rename(msg.Content)
		sendMsg := &Message{MsgType: RenameMsg, Sender: this.Name, Content: "用户" + oriName + "将用户名修改为：" + msg.Content}
		this.server.SendMsg(sendMsg)
	default:
		fmt.Println("unknown msg type err")
	}

}
