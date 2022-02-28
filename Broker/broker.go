package main

import (
	"log"
	"sync"
	"time"
)

type Broker struct {
	fileStorage       *FileStorage
	sourceSocket      *SourceSocket
	destinationSocket *DestinationSocket

	sourceChannel       chan *Message
	destinationChannel  chan *Message
	fileStoreChannel    chan *Message
	fileRetrieveChannel chan *Message
}

func NewBroker(fileStorage *FileStorage, sourceSoc *SourceSocket, destinationSoc *DestinationSocket) *Broker {
	return &Broker{
		fileStorage:       fileStorage,
		destinationSocket: destinationSoc,
		sourceSocket:      sourceSoc,

		sourceChannel:       make(chan *Message),
		destinationChannel:  make(chan *Message),
		fileStoreChannel:    make(chan *Message),
		fileRetrieveChannel: make(chan *Message),
	}
}

func (broker *Broker) Run() {
	var wg sync.WaitGroup
	wg.Add(1)
	broker.sourceSocket.ReceiveMessages(broker.sourceChannel)
	broker.destinationSocket.SendMessages(broker.destinationChannel)
	broker.fileStorage.StoreMessages(broker.fileStoreChannel, broker.fileRetrieveChannel)
	broker.start()
	broker.checkDestinationAvailability()

	wg.Wait() // block main goroutine
}

// start make the broker to listen to sourceChannel and fileRetrieveChannel and everytime receive any
// signal from these two channels, check the Destination_Service status, if it is running, redirect current
// Message to Destination_Service without saving it in file storage, else send the Message to FileStorage to
// be stored in file storage
func (broker *Broker) start() {
	go func() {
		for {
			select {
			case newMsg := <-broker.sourceChannel:
				broker.handleSendingMsgToDestination(newMsg)
			case topStoredMsg := <-broker.fileRetrieveChannel:
				broker.destinationChannel <- topStoredMsg
			}
		}
	}()
}

func (broker *Broker) handleSendingMsgToDestination(msg *Message) {
	isDestinationAvailable := broker.destinationSocket.IsDestinationAvailable()
	if isDestinationAvailable {
		// redirect incoming message to the destination
		broker.redirectMsgToDestination(msg)
	} else {
		// destination is unreachable, store message to file
		broker.fileStoreChannel <- msg
	}
}

func (broker *Broker) redirectMsgToDestination(msg *Message) {
	// before redirection, we must check there are stored messages or not
	isFileStorageEmpty := broker.fileStorage.IsFileStorageEmpty()
	if isFileStorageEmpty {
		// redirect new message to destination
		broker.destinationChannel <- msg
	} else {
		// must store current message to end of file storage and pick first stored message and send it first
		broker.fileStoreChannel <- msg
		topStoredMsg, err := broker.fileStorage.ReadMsgFromFile()
		if err != nil {
			log.Println("error while getting a msg from top of file storage", err)
		}
		broker.destinationChannel <- topStoredMsg
	}
}

// checkDestinationAvailability checks if Destination_Service is available or not in every Millisecond
func (broker *Broker) checkDestinationAvailability() {
	go func() {
		for {
			isDestinationAvailable := broker.destinationSocket.IsDestinationAvailable()
			broker.fileStorage.destinationAvailabilityChannel <- isDestinationAvailable
			time.Sleep(time.Millisecond)
		}
	}()
}
