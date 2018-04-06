//
//  Pubsub envelope subscriber.
//

package main

import (
	zmq "github.com/pebbe/zmq4"
	"os"
	"fmt"
	"log"
	"myzmq/packet"
	"encoding/json"
)


func main() {
//s	foo()
//	ch := make(chan<- string)
	//  Prepare our subscriber
	subscriber, _ := zmq.NewSocket(zmq.SUB)
	defer subscriber.Close()
	subscriber.Bind("tcp://*:5563")
	subscriber.SetSubscribe("")

	for {
		msg, _ := subscriber.Recv(0)

		go func ()  {
			bts := []byte(msg)

 			var pp packet.Packet
			err := json.Unmarshal(bts, &pp)
			if err != nil {
		    log.Printf("error decoding  response: %v", err)
		    if e, ok := err.(*json.SyntaxError); ok {
		        log.Printf("syntax error at byte offset %d", e.Offset)
		    }
			}
			fmt.Printf("pp: %+v", pp)

			logToFile(fmt.Sprintf("pp: %+v", pp))
		} ()

		fmt.Println("?")
	}
}


func logToFile(msg string) {
	f, err := os.OpenFile("test.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
	    log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println(msg)
}
