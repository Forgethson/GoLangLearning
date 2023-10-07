package main

import (
	"IM-System/server"
	"encoding/json"
	"fmt"
)

func main() {
	msg := &server.Message{MsgType: server.GlobalMsg, Content: "hello world!", Sender: "Wang", Receiver: []string{"Lee", "xiaoming"}}
	jsonBytes, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("to json err")
		panic(err)
	}
	jsonStr := string(jsonBytes)
	fmt.Println(jsonStr)

	jsonStr2 := "{\"msgType\":0,\"content\":\"hello world!\",\"sender\":null,\"receivers\":null}"
	fmt.Println(jsonStr2)
	var msg2 *server.Message
	err = json.Unmarshal([]byte(jsonStr), &msg2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", msg2)
}
