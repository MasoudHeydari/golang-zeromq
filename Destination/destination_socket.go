package main

import (
	"encoding/json"
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"log"
)

const (
	destinationChannelName    = "destination_channel"
	destinationSocketEndPoint = "tcp://localhost:5557"
)

type DestinationSocket struct {
	totalReceivedMsg uint
	soc              *zmq.Socket
}

func NewDestinationSocket() *DestinationSocket {
	return &DestinationSocket{
		totalReceivedMsg: 0,
		soc:              createNewDestinationSocket(),
	}
}

func createNewDestinationSocket() *zmq.Socket {
	subscriber, _ := zmq.NewSocket(zmq.SUB)
	err := subscriber.Connect(destinationSocketEndPoint)
	if err != nil {
		log.Fatalln("failed to connect to specified end point - err: ", err)
	}
	subscriber.SetSubscribe(destinationChannelName)
	return subscriber
}

func (destinationSocket *DestinationSocket) ReceiveMessages() {
	for {
		_, err := destinationSocket.soc.Recv(0)
		if err != nil {
			panic(err)
		}
		//log.Print("received a message from address: ", address)

		if receivedMsg, err := destinationSocket.soc.Recv(0); err != nil {
			log.Println(err)
		} else {
			message := new(Message)
			if err := json.Unmarshal([]byte(receivedMsg), message); err != nil {
				panic(err)
			}

			destinationSocket.totalReceivedMsg++
			fmt.Println("received new msg - total messages: ", destinationSocket.totalReceivedMsg)

		}
	}
}
