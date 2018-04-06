package main

import (
	zmq "go-zmq"

	"fmt"
)

func main() {
	// Socket to talk to clients
	context, _ := zmq.NewContext()
	defer context.Close()

	responder, _ := zmq.NewSocket(zmq.Sub)
	responder.Subscribe([]byte(""))
	defer responder.Close()

	responder.Bind("tcp://*:5555")

	for {
		fmt.Println("start")
		// Wait for next request from client
		request, _, _ := responder.RecvPart()
		fmt.Println("next")
		fmt.Printf("Received request: [%s]\n", string(request))

		// Do some 'work'
	//	time.Sleep(1 * time.Second)

		// Send reply back to client
	//	responder.SendPart([]byte("World"), false)
	}
}
