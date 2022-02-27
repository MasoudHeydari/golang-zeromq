package main

func main() {
	StartDestinationServer()
	destinationSocket := NewDestinationSocket()
	destinationSocket.ReceiveMessages()
}
