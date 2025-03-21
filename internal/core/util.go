package core

import (
	"sync"
    "time"
	"github.com/golang/protobuf/proto"
)


var LastSeen = map[string]time.Time {}
var LastSeenMutex sync.Mutex

var TankLevels = map[string]map[string]float32 {}
var TankLevelsMutex sync.Mutex

var SenderChannels = map[string]chan proto.Message {}

// NOTE no need for rw mutex as its multiple writers, 1 reader
var SenderChannelsMutex sync.Mutex

var MetricsChannel chan PbMetric

type PbMetric struct {
    ControllerId string 
    Timestamp time.Time
    Pb proto.Message
}

func Init() {
	MetricsChannel = make(chan PbMetric, 1000)

	TankLevels["solution"] = make(map[string]float32, 0)
	TankLevels["mixing"] = make(map[string]float32, 0)
}
