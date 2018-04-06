package main

import (
//    "encoding/json"
    "time"
    "flush/packet"
    "math/rand"
)

func main() {
  rand.Seed(time.Now().UnixNano())

for i:= 0; i < 5; i++ {
  packet.Counter("some name", "key0:item0, key1:item1", rand.Intn(10))
//    time.Sleep(2 * time.Second)
  packet.Gauge("name0" , "key0:item0, key1:item1", rand.Intn(100))
  packet.Gauge("name1" , "key0:item0, key1:item1", rand.Intn(100))
  packet.Gauge("name2", "key0:item0, key1:item1", rand.Intn(100))

  }
  time.Sleep(2 * time.Second)

}
