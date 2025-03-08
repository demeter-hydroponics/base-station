package core

import (
	//	"github.com/golang/protobuf/proto"
	"errors"

	"base-station/internal/database/metrics-db"
	"base-station/internal/utils"

	"github.com/charmbracelet/log"

	"base-station/internal/database/farm-config"
	pb_column "base-station/protobuf/generated/go/column"
	pb_node "base-station/protobuf/generated/go/node"
	pb_panel "base-station/protobuf/generated/go/panel"
)

// TODO metrics collection function
func ProcessMetrics(metrics_pb <- chan PbMetric) {
    log.Info("Starting Metrics Processing Goroutine")

    primaryPumpState := pb_column.BinaryLoadState_DISABLED
    mixingState := pb_column.BinaryLoadState_ENABLED
    var phSensorVal float32 = 0.0
    var tdsSensePPM float32 = 0.0
    var solutionResVolume float32 = 0.0
    var waterFeedResVolume float32 = 0.0
    var mixingResLevelVolume float32 = 0.0

    for msg := range metrics_pb {

        payload, err := utils.ProtoToJSON(msg.Pb)
        if err != nil {
            log.Error("there was an error in marshalling the protobuff", "err", err)
            continue
        }
//        log.Infof("sending metrics")
        log.Info("incoming data", "msg", msg.Pb.String())
        go metrics_db.SendMetric(msg.ControllerId, string(payload))
        if pb,ok := msg.Pb.(*pb_node.GrowLightSectionStats); ok {
            log.Info("recieved grow section stats", "msg",pb.String())
        }

        // TODO do a type check here and handle it
        if pb,ok := msg.Pb.(*pb_column.PumpTankStats); ok {
            //log.Info("Message recieved!", "time", msg.Timestamp.String())
            //log.Info("Pump Status", "status", pb.PrimaryPump.State)
            primaryPumpState = *pb.PrimaryPump.State
            solutionResVolume = *pb.SolutionReservoirLevel.TankFluidVolume_L
            waterFeedResVolume = *pb.FeedReservoirLevel.TankFluidVolume_L
            mixingResLevelVolume = *pb.MixingReservoirLevel.TankFluidVolume_L
        } 
        if pb, ok := msg.Pb.(*pb_column.MixingTankStats); ok {
            phSensorVal = *pb.PHSense.PhSenseMolPerL
            tdsSensePPM = *pb.TDSSense.TDSSensePPM
            mixingState = *pb.MixingValveStats.State
        }

        log.Infof("ph: %.2f, tds: %.2f, soln: %.2f, water feed: %.2f, mixing res: %.2f, PumpState: %s, Mixing State: %s", 
             phSensorVal, tdsSensePPM, 
             solutionResVolume, waterFeedResVolume, mixingResLevelVolume,primaryPumpState, mixingState)
    }
}


// NOTE not sure if this is needed
func ValidateConfig() error {
    // TODO validate the configuration makes sense
    return nil
}

func UpdateConfig(new_config *pb_panel.FarmConfig) error {
    if !farm_config.ConfigMutex.TryLock() {
        return errors.New("cannot process config, another config is being processed right now")
    }
    defer farm_config.ConfigMutex.Unlock() 

    // TODO apply the changes in the config
    for _, column := range new_config.Columns {
        _ = column
        id := *column.Id
        log.Infof("setting controller %s", id) 

        // apply the pump configs
        if UpdatePumpState(id, *column.PrimaryPumpState, *column.SecondaryPumpState) != nil {
                log.Info("Error setting Column Pump State", "id", id) 
        }
        // apply the mixing configs
        if UpdateMixingState(id, *column.MixingState) != nil {
                log.Info("Error setting Column Mixing State", "id", id) 
        }

        for _, node := range column.Nodes {
            if UpdatePPFD(node.GetId(), node.GetPPFD()) != nil {
                log.Info("Error setting Node PPFD", "id", node.GetId()) 
            }
        }
    }

    for _, node := range new_config.UnsetNodes {
        if UpdatePPFD(node.GetId(), node.GetPPFD()) != nil {
            log.Info("Error setting Node PPFD", "id", node.GetId()) 
        }
    }

    farm_config.Config = *new_config
    return nil
}
