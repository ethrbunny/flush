//
//  Pubsub envelope subscriber.
//

package main

import (
	zmq "github.com/pebbe/zmq4"
	"os"
	"fmt"
	"log"
	"flush/packet"
	"encoding/json"
	"time"
)

type gaugePacket struct {
	procInfo	*packet.ProcInfo
	gaugeVals	[]int
	startTime	int64
	endTime		int64
	tags 			map[string]string
}

type countPacket struct {
	procInfo	*packet.ProcInfo
	countSum	int
	startTime	int64
	endTime		int64
	tags 			map[string]string
}

func newCounter(pp packet.Packet) *countPacket {
	inPacket := new(countPacket)
	inPacket.procInfo = &pp.PInfo
	inPacket.countSum = pp.Incr.Increment
	inPacket.tags = pp.Tags
	inPacket.startTime = pp.CreateTime

	return inPacket
}

func newGauge(pp packet.Packet) *gaugePacket {
	inPacket := new(gaugePacket)
	inPacket.procInfo = &pp.PInfo
	inPacket.tags = pp.Tags
	inPacket.startTime = pp.CreateTime
	inPacket.gaugeVals = append(inPacket.gaugeVals, pp.Amt.GVal)
//fmt.Printf("gval: %d\n", inPacket.gaugeVals)
	return inPacket
}

func main() {
	subscriber, _ := zmq.NewSocket(zmq.SUB)
	defer subscriber.Close()
	subscriber.Bind("tcp://*:5563")
	subscriber.SetSubscribe("")

	counters := make(map[string] *countPacket)
	gauges := make(map[string] *gaugePacket)

	ticker := time.NewTicker(5 * time.Second)
	reset := time.NewTicker(20 * time.Second)
//	quit := make(chan struct{})
	go func() {
	    for {
	       select {
	        case <- ticker.C:
	            fmt.Println("tick")
							for k, _ := range counters {
								fmt.Printf("%s: counter: %d\n", k, counters[k].countSum)
							}

							for k, _ := range gauges {
								fmt.Printf("%s: gauges: %v", k, gauges[k].gaugeVals)
								fmt.Printf("\n")
							}
	        case <- reset.C:
						fmt.Println("reset")
	            counters = make(map[string] *countPacket)
							gauges = make(map[string] *gaugePacket)

	        }
	    }
	 }()

	for {
		msg, _ := subscriber.Recv(0)

		go func ()  {
			bts := []byte(msg)

 			var pp packet.Packet
			err := json.Unmarshal(bts, &pp)
			if err != nil {
		    log.Printf("error decoding  response: %v", err)
/*		    if e, ok := err.(*json.SyntaxError); ok {
		        log.Printf("syntax error at byte offset %d", e.Offset)
		    } */
			}
	//		fmt.Printf("pp: %+v\n", pp)
/*
			if pp.Incr != nil {
				ii := pp.Incr
				fmt.Printf("%s %v \n", pp.Name,  ii.Increment *4)
			}
*/
		// counter or gauge?
			packetName := pp.Name
			if pp.Incr != nil {
				if counters[packetName] == nil {
					counters[packetName] = newCounter(pp)
				} else {
					inPacket := counters[packetName]
					inPacket.countSum += pp.Incr.Increment
				}
			} else if pp.Amt != nil {
				if gauges[packetName] == nil {
					gauges[packetName] = newGauge(pp)
				} else {
					inGauge := gauges[packetName]
					inGauge.gaugeVals = append(inGauge.gaugeVals, pp.Amt.GVal)

				}
			}

			logToFile(fmt.Sprintf("pp: %+v", pp))
		} ()
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
