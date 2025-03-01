package core

import (
	pb_column "base-station/protobuf/generated/go/column"

	//"github.com/golang/protobuf/proto"
    "errors"
)

// Column Update function
// NOTE These functions should construct a proto message and send it to the appropriate channel
// NOTE farm config msg -> pb message
func UpdatePumpState(id string, primary, secondary pb_column.PumpState) error {

    // convert the pumpstate to a pb command
    pumpType := pb_column.PumpType_PRIMARY
    if primary == secondary {
        pumpType = pb_column.PumpType_BOTH
    }

    SenderChannelsMutex.Lock()
    // send the command
    if channel, exists := SenderChannels[id]; exists {
        channel <- &pb_column.SetPumpStateCommand{SelectedPump: &pumpType, State: &primary} 
    } else {
        SenderChannelsMutex.Unlock()
        return errors.New("Id not recognized")
    }
    SenderChannelsMutex.Unlock()

    if pumpType == pb_column.PumpType_BOTH {
	    return nil
    }

    // convert the pumpstate to a pb command
    pumpTypeSec := pb_column.PumpType_SECONDARY
    SenderChannelsMutex.Lock()
    // send the command
    if channel, exists := SenderChannels[id]; exists {
        channel <- &pb_column.SetPumpStateCommand{SelectedPump: &pumpTypeSec, State: &secondary} 
    } else {
        SenderChannelsMutex.Unlock()
        return errors.New("Id not recognized")
    }
    SenderChannelsMutex.Unlock()

	return nil
}

func UpdateMixingState(id string, state pb_column.MixingOverrideState) error {
    SenderChannelsMutex.Lock()
    // send the command
    if channel, exists := SenderChannels[id]; exists {
        channel <- &pb_column.SetMixingStateCommand{State: &state} 
    } else {
        SenderChannelsMutex.Unlock()
        return errors.New("Id not recognized")
    }
    SenderChannelsMutex.Unlock()

	return nil
}
