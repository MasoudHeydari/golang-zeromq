package main

import (
	"encoding/json"
	zmq "github.com/pebbe/zmq4"
	"log"
)

const sourceChannelName = "source_channel"

type SourceSocket struct {
	endPoint              string
	totalReceivedMessages uint
	soc                   *zmq.Socket
}

// NewSourceSocket create new SourceSocket
func NewSourceSocket(sourceEndPoint string) *SourceSocket {
	subscriber := createNewSourceSocket(sourceEndPoint)
	return &SourceSocket{
		totalReceivedMessages: 0,
		endPoint:              sourceEndPoint,
		soc:                   subscriber,
	}
}

func createNewSourceSocket(sourceEndPoint string) *zmq.Socket {
	subscriber, err := zmq.NewSocket(zmq.SUB)
	if err != nil {
		log.Fatalln("failed to create source socket - error: ", err)
		return nil
	}
	err = subscriber.Connect(sourceEndPoint)
	if err != nil {
		log.Fatalln("failed to connect to specified end point - err: ", err)
	}
	subscriber.SetSubscribe(sourceChannelName)
	return subscriber
}

// ReceiveMessages receives messages that Source_Service are sending to Broker_Service
// and when receives a new message, send it to sourceSocket
func (sourceSocket *SourceSocket) ReceiveMessages(sourceChannel chan<- *Message) {
	go func() {
		for {
			address, err := sourceSocket.soc.Recv(0)
			if err != nil {
				panic(err)
			}
			log.Print("received a message from address: ", address)

			if receivedMsg, err := sourceSocket.soc.Recv(0); err != nil {
				log.Println(err)
			} else {
				newReceivedMessage := new(Message)
				if err := json.Unmarshal([]byte(receivedMsg), newReceivedMessage); err != nil {
					panic(err)
				}

				sourceChannel <- newReceivedMessage
			}
		}
	}()
}
