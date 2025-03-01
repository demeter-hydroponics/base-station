package core

import (
	//	"github.com/golang/protobuf/proto"
	"errors"

	"github.com/charmbracelet/log"

	"base-station/internal/database/farm-config"
	pb_panel "base-station/protobuf/generated/go/panel"
)

// TODO metrics collection function
func ProcessMetrics(metrics_pb <- chan PbMetric) {
    log.Info("Starting Metrics Processing Goroutine")
    for msg := range metrics_pb {
        // TODO do a type check here and handle it
        log.Info("Message recieved!", "time", msg.Timestamp.String())
    }
}


// NOTE not sure if this is needed
func ValidateConfig() error {
    // TODO validate the configuration makes sense
    return nil
}

func UpdateConfig(new_config *pb_panel.FarmConfig) error {
    if farm_config.ConfigMutex.TryLock() {
        return errors.New("cannot process config, another config is being processed right now")
    }
    defer farm_config.ConfigMutex.Unlock() 

    // TODO apply the changes in the config
    for _, column := range new_config.Columns {
        _ = column
        id := *column.Id
        log.Infof("setting controller %s", id) 

        // apply the pump configs
        UpdatePumpState(id, *column.PrimaryPumpState, *column.SecondaryPumpState)
        // apply the mixing configs
        UpdateMixingState(id, *column.MixingState)

        // TODO implement node configs
        for _, node := range column.Nodes {
            _ = node
            // TODO i think for now we arent actually configuring nodes 
        }
    }

    for _, node := range new_config.UnsetNodes {
        _ = node
        // TODO configure nodes
    }
    return nil
}
