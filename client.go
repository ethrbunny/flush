package main

import (
//    "encoding/json"
    "time"
    "myzmq/packet"
)

func main() {
  packet.Counter("some name", "key0:item0, key1:item1", 2)
    time.Sleep(2 * time.Second)
    packet.Gauge("some name", "key0:item0, key1:item1", 11)
    time.Sleep(2 * time.Second)
}
