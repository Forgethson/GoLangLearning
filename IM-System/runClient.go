package main

import (
	"IM-System/client"
	"flag"
	"fmt"
)

func main() {
	flag.Parse()
	newClient := client.NewClient("127.0.0.1", 8888)
	if newClient == nil {
		fmt.Println(">>>>> 链接服务器失败...")
		return
	}

	//单独开启一个goroutine去处理server的回执消息
	go newClient.DealResponse()

	fmt.Println(">>>>>链接服务器成功...")

	//启动客户端的业务
	//newClient.Run()
	newClient.Run2()
}
