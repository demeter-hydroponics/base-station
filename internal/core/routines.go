package core

import (
	//	"github.com/golang/protobuf/proto"
	"errors"

	metrics_db "base-station/internal/database/metrics-db"
	"base-station/internal/utils"

	"github.com/charmbracelet/log"

	farm_config "base-station/internal/database/farm-config"
	pb_node "base-station/protobuf/generated/go/node"
	pb_panel "base-station/protobuf/generated/go/panel"
)

func ProcessMetrics(metrics_pb <-chan PbMetric) {
	log.Info("Starting Metrics Processing Goroutine")

	for msg := range metrics_pb {
		payload, err := utils.ProtoToJSON(msg.Pb)
		if err != nil {
			log.Error("there was an error in marshalling the protobuff", "err", err)
			continue
		}
		log.Info("incoming data", "msg", msg.Pb.String())
		go metrics_db.SendMetric(msg.ControllerId, string(payload))
		if pb, ok := msg.Pb.(*pb_node.GrowLightSectionStats); ok {
			log.Info("recieved grow section stats", "msg", pb.String())
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
		if UpdatePumpState(id, *column.PrimaryPumpState, *column.SecondaryPumpState) != nil {
			log.Info("Error setting Column Pump State", "id", id)
		}
		// apply the mixing configs
		if UpdateMixingState(id, *column.MixingState) != nil {
			log.Info("Error setting Column Mixing State", "id", id)
		}
		// apply the water level configs
		if UpdateWaterlevelState(id, *column.WaterLevelState) != nil {
			log.Info("Error setting Column water level State", "id", id)
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
