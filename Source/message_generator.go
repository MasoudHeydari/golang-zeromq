package main

import (
	"math/rand"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var (
	min int
	max int
)

type Message struct {
	Id      uint32 `json:"id"`
	Content string `json:"content"`
}

type MessageGenerator struct {
	totalMsgNo uint32
}

// NewMessageGenerator instantiates a new message generator
func NewMessageGenerator(lowerBound, upperBound int) *MessageGenerator {
	min = lowerBound
	max = upperBound
	rand.Seed(time.Now().UnixNano())
	return &MessageGenerator{
		totalMsgNo: 0,
	}
}

// GetRandomMessage generate a random message
func (msgGenerator *MessageGenerator) GetRandomMessage() *Message {
	randomMsg := &Message{}
	randomMsg.Id = msgGenerator.totalMsgNo
	randomMsg.Content = *getRandomString()
	msgGenerator.totalMsgNo++ // increase total number of generated message
	return randomMsg
}

// getRandomString generate a random string based on provided lower and upper bound
func getRandomString() *string {
	randInt := getRandInt()
	return randomStringBytes(randInt)
}

// randomStringBytes generate a random string to the length of provided strLength
func randomStringBytes(strLength int) *string {
	strLengthBytes := make([]byte, strLength)
	lettersLen := len(letters)

	for i := range strLengthBytes {
		strLengthBytes[i] = letters[rand.Int63()%int64(lettersLen)]
	}

	resultStr := string(strLengthBytes)
	return &resultStr
}

// getRandInt generate a random int between min and max
func getRandInt() int {
	return min + rand.Intn(max-min)
}
