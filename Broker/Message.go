package main

import (
	"fmt"
)

type Message struct {
	Id      uint32 `json:"id"`
	Content string `json:"content"`
}

func (msg *Message) ToString() string {
	return fmt.Sprintf("{id:%v,content:\"%s\"}$\n", msg.Id, msg.Content)
}
