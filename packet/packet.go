package packet

import (
  "time"
"strings"
  zmq "github.com/pebbe/zmq4"
  "fmt"
  "encoding/json"
  "os"
)

type ProcInfo struct {
  Pid int
  AppName string
  HostName  string
}

type CountVal struct {
  Increment int
}

type GaugeVal struct {
  GVal int
}

// TODO getters for all of these - make them private
type Packet struct {
  Name string
  Tags map[string] string
  CreateTime int64
  PInfo ProcInfo
  Incr *CountVal
  Amt *GaugeVal
}

func newPacket(name string, tags string) *Packet {
  p := new(Packet)
  p.CreateTime = time.Now().Unix()
  p.Tags = make(map[string]string)
  p.Name = name

  // process info
  p.PInfo.Pid = os.Getpid()
  p.PInfo.HostName, _ = os.Hostname()
  p.PInfo.AppName = os.Args[0]

  tagSplit := strings.Split(tags, ",")
  for _, tagItem := range tagSplit {
    tagPair := strings.Split(strings.TrimSpace(tagItem), ":")
    p.Tags[tagPair[0]] = tagPair[1]
  }
  return p
}

func sendPacket(pkt *Packet) {
  go func() {
    publisher, _ := zmq.NewSocket(zmq.PUB)
    defer publisher.Close()
    publisher.Connect("tcp://localhost:5563")

    time.Sleep(time.Second)

    b, err := json.Marshal(pkt)
  	if err != nil {
  	   fmt.Println("error:", err)
       return
    }
  //  os.Stdout.Write(b)
    publisher.Send(string(b[:]), 0)
  }()
}

// TODO error checking
func Gauge(name string, tags string, gval int) {
  go func() {
    p := newPacket(name, tags)
    p.Amt = new(GaugeVal)
    p.Amt.GVal = gval

    sendPacket(p)
  } ()
}

func Counter(name string, tags string, iment int) {
  go func() {
    p := newPacket(name, tags)
    p.Incr = new(CountVal)
    p.Incr.Increment = iment

    sendPacket(p)
  } ()
}
