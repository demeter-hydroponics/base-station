package core

import (
	//	"github.com/golang/protobuf/proto"
	"errors"
	"time"

	metrics_db "base-station/internal/database/metrics-db"
	"base-station/internal/utils"

	"github.com/charmbracelet/log"

	farm_config "base-station/internal/database/farm-config"
//	pb_node "base-station/protobuf/generated/go/node"
	pb_column "base-station/protobuf/generated/go/column"
	pb_panel "base-station/protobuf/generated/go/panel"
)

func ProcessMetrics(metrics_pb <-chan PbMetric) {
	log.Info("Starting Metrics Processing Goroutine")

	for msg := range metrics_pb {
		// register the lastSeen 
		LastSeenMutex.Lock()
		LastSeen[msg.ControllerId] = time.Now()
		LastSeenMutex.Unlock()

		payload, err := utils.ProtoToJSON(msg.Pb)
		if err != nil {
			log.Error("there was an error in marshalling the protobuff", "err", err)
			continue
		}
		log.Info("incoming data", "msg", msg.Pb.String())
		go metrics_db.SendMetric(msg.ControllerId, string(payload))

//		if pb, ok := msg.Pb.(*pb_node.GrowLightSectionStats); ok {
//			log.Info("recieved grow section stats", "msg", pb.String())
//		}
		if pb, ok := msg.Pb.(*pb_column.PumpTankStats); ok {
			TankLevelsMutex.Lock()
			// extract the relevant info
			// grab the solution reservoir levels
			TankLevels["solution"][msg.ControllerId] = 0.18 - *pb.SolutionReservoirLevel.TankFluidVolume_L
			// grab the mixing reservoir levels
			TankLevels["mixing"][msg.ControllerId] = 0.12 - *pb.MixingReservoirLevel.TankFluidVolume_L
			TankLevelsMutex.Unlock()
		}

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
        if err := UpdatePumpState(id, *column.PrimaryPumpState, *column.SecondaryPumpState); err != nil {
			log.Info("Error setting Column Pump State", "id", id, "err", err)
		}
		// apply the mixing configs
        if err := UpdateMixingState(id, *column.MixingState); err != nil {
			log.Info("Error setting Column Mixing State", "id", id, "err", err)
		}
		// apply the water level configs
        if err := UpdateWaterlevelState(id, *column.WaterLevelState); err != nil {
			log.Info("Error setting Column water level State", "id", id, "err", err)
		}

		for _, node := range column.Nodes {
            if err := UpdatePPFD(node.GetId(), node.GetPPFD()); err != nil {
				log.Info("Error setting Node PPFD", "id", node.GetId(), "err", err)
			}
		}
	}

	for _, node := range new_config.UnsetNodes {
        if err := UpdatePPFD(node.GetId(), node.GetPPFD()); err != nil {
			log.Info("Error setting Node PPFD", "id", node.GetId(), "err", err)
		}
	}

	farm_config.Config = *new_config
	return nil
}
