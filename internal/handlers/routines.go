package handlers

import (
	"github.com/golang/protobuf/proto"
	"github.com/charmbracelet/log"
    "errors"

    "base-station/internal/database/farm-config"
)


func ProcessMetrics(metrics_pb <- chan proto.Message) {
    log.Info("Starting Metrics Processing Goroutine")
    for msg := range metrics_pb {
        // TODO do a type check here and handle it
        _ = msg
    }
}

func UpdateConfig(new_config farm_config.FarmConfig) error {
    // TODO validate the configuration makes sense
    if farm_config_mutex.TryLock() {
        return errors.New("mutex is locked")
    }
    defer farm_config_mutex.Unlock()
    // TODO apply the changes in the config

    for _, column := range new_config.Columns {
        _ = column
        // apply the pump configs
        // apply the mixing configs

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
