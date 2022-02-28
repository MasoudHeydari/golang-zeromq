package main

func main() {
	// these 3 parameters are adjustable
	totalRequests := 3000     // ~10_000 sourceSocket
	msgLengthLowerBound := 15 // 50	    	->   	// 50 Byte
	msgLengthUpperBound := 30 // 8 * 1024   ->	    // 8 K Byte

	msgGenerator := NewMessageGenerator(msgLengthLowerBound, msgLengthUpperBound)
	sourceSocket := NewSourceSocket(msgGenerator)
	sourceSocket.GenerateAndSend(totalRequests)

}
