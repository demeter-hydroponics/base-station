package core

import (
	farm_config "base-station/internal/database/farm-config"
	pb_column "base-station/protobuf/generated/go/column"
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


// Column Update function
// TODO func for updating a specific Controller with a specific param
// NOTE These functions should construct a proto message and send it to the appropriate channel
// NOTE farm config msg -> pb message
func UpdatePumpState(id string, primary, secondary farm_config.PumpState) error {

    // convert the pumpstate to a pb command
    pumpType := pb_column.PumpType_PRIMARY
    if primary == secondary {
        pumpType = pb_column.PumpType_BOTH
    }
    pumpState := pb_column.PumpState(conversionMap[string(primary)])  

    SenderChannelsMutex.Lock()
    // send the command
    SenderChannels[id] <- &pb_column.SetPumpStateCommand{SelectedPump: &pumpType, State: &pumpState} 
    SenderChannelsMutex.Unlock()

    if pumpType == pb_column.PumpType_BOTH {
	    return nil
    }

    // convert the pumpstate to a pb command
    pumpTypeSec := pb_column.PumpType_SECONDARY
    pumpStateSec := pb_column.PumpState(conversionMap[string(secondary)])  
    SenderChannelsMutex.Lock()
    // send the command
    SenderChannels[id] <- &pb_column.SetPumpStateCommand{SelectedPump: &pumpTypeSec, State: &pumpStateSec} 
    SenderChannelsMutex.Unlock()

	return nil
}

func UpdateMixingState(id string, state farm_config.MixingState) error {
    mixingState := pb_column.MixingOverrideState(conversionMap[string(state)])  
    SenderChannelsMutex.Lock()
    // send the command
    SenderChannels[id] <- &pb_column.SetMixingStateCommand{State: &mixingState} 
    SenderChannelsMutex.Unlock()

	return nil
}
