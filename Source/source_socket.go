package main

import (
	"encoding/json"
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"log"
	"time"
)

const (
	sourceChannelName = "source_channel"
	endPoint          = "tcp://*:5556"
)

type SourceSocket struct {
	soc          *zmq.Socket
	msgGenerator *MessageGenerator
	totalSentMsg uint
}

// NewSourceSocket create new SourceSocket and return the reference
func NewSourceSocket(msgGenerator *MessageGenerator) *SourceSocket {

	return &SourceSocket{
		msgGenerator: msgGenerator,
		soc:          createNewSourceS(),
		totalSentMsg: 0,
	}
}

func createNewSourceS() *zmq.Socket {
	zmqSocket, err := zmq.NewSocket(zmq.PUB)
	if err != nil {
		log.Fatalln("failed to create source socket - error: ", err)
		return nil
	}
	err = zmqSocket.Bind(endPoint)
	if err != nil {
		log.Fatalln("failed to bind to this end point source socket - error: ", err)
		return nil
	}
	// make sure subscriber connection has time to complete
	time.Sleep(time.Second)
	return zmqSocket
}

// GenerateAndSend send messages to the number of totalRequests to the Destination
func (sourceSocket *SourceSocket) GenerateAndSend(totalRequests int) {
	for i := 0; i < totalRequests; i++ {
		randomMsg := sourceSocket.msgGenerator.GetRandomMessage()
		err := sourceSocket.Send(randomMsg)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// Send send a message to Destination. used in GenerateAndSend function
func (sourceSocket *SourceSocket) Send(newMsg *Message) error {
	requestBody, err := json.Marshal(newMsg)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = sourceSocket.soc.Send(sourceChannelName, zmq.SNDMORE)
	if err != nil {
		fmt.Println("error while sending to destination channel - error: ", err)
		return err
	}
	_, err = sourceSocket.soc.Send(string(requestBody), 0)
	if err != nil {
		fmt.Println("error while sending new message to destination - error: ", err)
		return err
	}

	sourceSocket.totalSentMsg++
	fmt.Println("total sent messages: ", sourceSocket.totalSentMsg)
	return nil
}
