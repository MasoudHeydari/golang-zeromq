package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

type FileStorage struct {
	fileName                       string
	destinationAvailabilityChannel chan bool
}

// NewFileStorage create a FileStorage
func NewFileStorage(storageName string) *FileStorage {
	return &FileStorage{
		fileName:                       storageName,
		destinationAvailabilityChannel: make(chan bool),
	}
}

// StoreMessages listens to fileStoreChannel and if receive a new message, store it in file storage
func (db *FileStorage) StoreMessages(fileStoreChannel <-chan *Message, fileRetrieveChannel chan<- *Message) {
	go func() {
		for {
			select {
			case newMsgToStore := <-fileStoreChannel:
				_ = db.StoreMessageToFile(newMsgToStore)
			case isDestAvailable := <-db.destinationAvailabilityChannel:
				db.sendStoredMessagesToDestination(fileRetrieveChannel, isDestAvailable)
			}
		}
	}()
}

// sendStoredMessagesToDestination fetch stored messages from file storage and send it to fileRetrieveChannel to be sent to Destination
func (db *FileStorage) sendStoredMessagesToDestination(fileRetrieveChannel chan<- *Message, isDestAvailable bool) {
	go func() {
		if !db.IsFileStorageEmpty() && isDestAvailable {
			topMsg, err := db.ReadMsgFromFile()
			if err != nil {
				log.Println(err)
				return
			}
			fileRetrieveChannel <- topMsg
		}
	}()
}

// StoreMessageToFile receive a Message pointer and store it in file storage.
// used Mutex for locking other goroutines to work with file storage at a same time
func (db *FileStorage) StoreMessageToFile(msg *Message) error {
	// I think there is no pre-condition to be concerned about, store the message to end of file
	var mu sync.Mutex
	mu.Lock()

	file, err := os.OpenFile(db.fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}

	newStrReader := strings.NewReader(msg.ToString())
	_, err = io.Copy(file, newStrReader)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = file.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}

	mu.Unlock()
	return nil
}

// ReadMsgFromFile read a message from top of file storage and decode it to Message type and return the reference
func (db *FileStorage) ReadMsgFromFile() (*Message, error) {
	var mu sync.Mutex
	mu.Lock()
	fmt.Println("read from file")
	topMsgBytes, err := db.getTopMessageBytes()
	if err != nil {

		log.Println("baaaad", err)
		return nil, err
	}
	topMsg := new(Message)
	//err = json.Unmarshal(topMsgBytes, topMsg)
	//if err != nil {
	//	log.Println("wwwww", err)
	//	return nil, err
	//}

	fmt.Println(string(topMsgBytes))

	mu.Unlock()
	return topMsg, nil
}

// getTopMessageBytes read a message from top of file storage and delete it from file storage
func (db *FileStorage) getTopMessageBytes() ([]byte, error) {
	file, err := os.OpenFile(db.fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	topMsg, err := popTopMessage(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
	//fmt.Println("pop:", string(topMsg))

	return topMsg, nil
}

func popTopMessage(f *os.File) ([]byte, error) {
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if fi.Size() == 0 {
		return nil, errors.New("file is empty")
	}

	buf := bytes.NewBuffer(make([]byte, 0, fi.Size()))

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		fmt.Println("error 1")
		return nil, err
	}
	_, err = io.Copy(buf, f)
	if err != nil {
		fmt.Println("error 2")
		return nil, err
	}

	topMsgBytes, err := buf.ReadBytes('$') // delim: '\n'
	if err != nil && err != io.EOF {
		fmt.Println("error 3")
		return nil, err
	}

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		fmt.Println("error 4")
		return nil, err
	}
	nw, err := io.Copy(f, buf)
	if err != nil {
		fmt.Println("error 1")
		return nil, err
	}
	err = f.Truncate(nw)
	if err != nil {
		return nil, err
	}
	err = f.Sync()
	if err != nil {
		return nil, err
	}

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	return topMsgBytes[0 : len(topMsgBytes)-1], nil
}

// IsFileStorageEmpty check the status of file storage. if file is empty, return true otherwise return false
func (db *FileStorage) IsFileStorageEmpty() bool {
	fileStatus, _ := os.Stat(db.fileName)
	fileSize := fileStatus.Size()
	//fmt.Println("file size is: ", fileSize)
	return fileSize == 0
}
