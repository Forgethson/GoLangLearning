package main

import (
	server "IM-System/server"
)

func main() {
	myServer := server.NewServer("127.0.0.1", 8888)
	myServer.Start()
}

/*
sender在服务器自动添加，无需客户端添加

广播：
{"msgType":0,"content":"hello!","sender":null,"receiver":null}

私聊
{"msgType":1,"content":"hello!","sender":null,"receiver":"wjd"}

查询所有当前在线用户
{"msgType":2,"content":"null","sender":null,"receiver":null}

重命名
{"msgType":3,"content":"xiaoming","sender":null,"receiver":null}
{"msgType":3,"content":"xiaowang","sender":null,"receiver":null}
*/
