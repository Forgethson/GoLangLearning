package server

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	// 在线用户的列表
	OnlineMap map[string]*User
	// 锁对象
	mapLock sync.RWMutex

	//消息广播的channel缓冲区
	Message chan *Message
}

const BufferSize = 4096

// NewServer 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan *Message, 100),
	}
	return server
}

// MessageListener 监听Message消息channel的goroutine
func (this *Server) MessageListener() {
	for {
		msg := <-this.Message
		switch msg.MsgType {

		case GlobalMsg: // 广播消息
			this.mapLock.Lock()
			for _, cli := range this.OnlineMap {
				cli.C <- msg
			}
			this.mapLock.Unlock()

		case PrivateMsg:
			user, ok := this.OnlineMap[msg.Receiver]
			if !ok {
				fmt.Println("用户名不存在")
				this.SendMsg(&Message{MsgType: PrivateMsg, Sender: msg.Sender, Content: "用户名不存在", Receiver: msg.Sender})
				break
			}
			user.C <- msg

		case SearchAllOnlineClientMsg: // 查询所有在线用户消息
			this.mapLock.Lock()
			user := this.OnlineMap[msg.Sender]
			for _, cli := range this.OnlineMap {
				newMsg := CopyMessage(msg)
				newMsg.Content = "用户" + cli.Name + "在线"
				user.C <- newMsg
			}
			this.mapLock.Unlock()

		case RenameMsg:
			this.mapLock.Lock()
			user := this.OnlineMap[msg.Sender]
			user.C <- msg
			this.mapLock.Unlock()
		default:
			fmt.Println("unknown msg type err")
		}
	}
}

// SendMsg 发送消息
func (this *Server) SendMsg(msg *Message) {
	this.Message <- msg
}

// Handler 客户端连接处理器
func (this *Server) Handler(conn net.Conn) {
	// ...当前链接的业务
	fmt.Println("链接建立成功")
	user := NewUser(conn, this)
	user.Online()
	// 接收缓冲区
	buf := make([]byte, BufferSize)
	for {
		n, err := conn.Read(buf)
		// 客户端Socket合法关闭
		if n == 0 {
			user.Offline()
			break
		}

		if err != nil && err != io.EOF {
			fmt.Println("Conn Read err:", err)
			break
		}
		// 得到用户消息，去除末尾的'\n'
		msg := string(buf[:n-1])
		// 处理用户消息
		user.doMeg(msg)
	}
	//当前handler阻塞
	//select {}
}

// Start 启动服务器
func (this *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	// close listen socket
	defer listener.Close()

	// 启动监听广播消息的goroutine
	go this.MessageListener()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}

		//do handler
		go this.Handler(conn)
	}
}
