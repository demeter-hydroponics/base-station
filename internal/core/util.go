package core

import (
	"sync"
    "time"
	"github.com/golang/protobuf/proto"
)

var SenderChannels map[string]chan proto.Message

// NOTE no need for rw mutex as its multiple writers, 1 reader
var SenderChannelsMutex sync.Mutex

var MetricsChannel chan PbMetric

type PbMetric struct {
    ControllerId string 
    Timestamp time.Time
    Pb proto.Message
}

