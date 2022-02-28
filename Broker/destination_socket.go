package main

import (
	"encoding/json"
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"log"
	"net/http"
	"time"
)

const (
	destinationUrl         = "http://localhost:8080"
	destinationChannelName = "destination_channel"
)

type DestinationSocket struct {
	totalSentMessages uint
	endPoint          string
	soc               *zmq.Socket
}

// NewDestinationSocket create new DestinationSocket
func NewDestinationSocket(destinationEndPoint string) *DestinationSocket {
	publisher := createNewDestinationSocket(destinationEndPoint)

	return &DestinationSocket{
		totalSentMessages: 0,
		endPoint:          destinationEndPoint,
		soc:               publisher,
	}
}

func createNewDestinationSocket(destinationEndPoint string) *zmq.Socket {
	publisher, err := zmq.NewSocket(zmq.PUB)
	if err != nil {
		log.Fatalln("failed to create destination socket - error: ", err)
		return nil
	}

	err = publisher.Bind(destinationEndPoint)
	if err != nil {
		log.Fatalln("err while creating zmq.socket: ", err)
		return nil
	}
	// make sure subscriber connection has time to complete
	time.Sleep(time.Second)
	return publisher
}

// SendMessages listens to destinationChannel and everytime that receive a new Message,
// redirect it to Destination_Service via sendMessage function
func (destinationSocket *DestinationSocket) SendMessages(destinationChannel <-chan *Message) {
	go func() {
		for {
			select {
			case newMsgToSend := <-destinationChannel:
				err := destinationSocket.sendMessage(newMsgToSend)
				if err != nil {
					log.Println("error occurred while sending message to destination_service - error: ", err)
				}
			}
		}
	}()
}

// sendMessage receives a Message pointer and send it to Destination_Service
func (destinationSocket *DestinationSocket) sendMessage(newMsg *Message) error {
	requestBody, err := json.Marshal(newMsg)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = destinationSocket.soc.Send(destinationChannelName, zmq.SNDMORE)
	if err != nil {
		fmt.Println("error while sending to destination channel - error: ", err)
		return err
	}

	_, err = destinationSocket.soc.Send(string(requestBody), 0)
	if err != nil {
		fmt.Println("error while sending new message to destination - error: ", err)
		return err
	}

	return nil
}

// IsDestinationAvailable checks the Destination_Service is running or not
func (destinationSocket *DestinationSocket) IsDestinationAvailable() bool {
	resp, err := http.Get(destinationUrl)
	if err != nil {
		//log.Println("http GET failed - err: ", err)
		return false
	}

	return resp.StatusCode == http.StatusOK
}
