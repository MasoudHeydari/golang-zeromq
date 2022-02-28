package main

const (
	destinationEndPoint = "tcp://*:5557"
	sourceEndPoint      = "tcp://localhost:5556"
	storageFileName     = "stored_messages.txt"
)

func main() {
	source := NewSourceSocket(sourceEndPoint)
	destination := NewDestinationSocket(destinationEndPoint)
	fileStorage := NewFileStorage(storageFileName)

	broker := NewBroker(fileStorage, source, destination)
	broker.Run()
}
