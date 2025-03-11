package farm_config

import (
    "sync"
	pb_panel "base-station/protobuf/generated/go/panel"
	pb_column "base-station/protobuf/generated/go/column"
	pb_common "base-station/protobuf/generated/go"
)

/*
This Might be temporary, might be a bad idea.
I want a single function/handler to be able to update the config at once, others will be rejected temporarily until the current one is processed.
*/
var ConfigMutex sync.Mutex
var Config pb_panel.FarmConfig 
var DefaultConfig pb_panel.FarmConfig 

func InitFarmConfig() {
    Config = pb_panel.FarmConfig{
        Controllers: make(map[string]*pb_panel.ControllerConfig),
        Columns: make([]*pb_panel.ColumnConfig, 0),
        UnsetNodes: make([]*pb_panel.NodeConfig, 0),
    }

    column_id := "5a0c4ea4-b1c5-4c97-ae43-7e01627dc688"
    column_name := "bulbasaur"
    column_type_ctrllr := pb_common.ControllerType_COLUMN

    Config.Controllers["5a0c4ea4-b1c5-4c97-ae43-7e01627dc688"] = &pb_panel.ControllerConfig{
        Id: &column_id, 
        Type: &column_type_ctrllr,
        Name: &column_name,
    }

    node_id := "0a7dbbd3-d42b-4593-9135-8509c2ed520e"
    node_name := "lettuce"
    node_type_ctrllr := pb_common.ControllerType_NODE

    Config.Controllers["0a7dbbd3-d42b-4593-9135-8509c2ed520e"] = &pb_panel.ControllerConfig{
        Id: &node_id, 
        Type: &node_type_ctrllr,
        Name: &node_name,
    }
    
    node_2_id := "0a7dcbd3-d52b-4693-9235-8519c2ed521f"
    node_2_name := "green onions"
    node_2_type_ctrllr := pb_common.ControllerType_NODE

    Config.Controllers["0a7dcbd3-d52b-4693-9235-8519c2ed521f"] = &pb_panel.ControllerConfig{
        Id: &node_2_id, 
        Type: &node_2_type_ctrllr,
        Name: &node_2_name,
    }

    pumpState := pb_column.PumpState_PUMP_OFF
    mixingState := pb_column.MixingOverrideState_OVERRIDE_VALVE_OFF
    column := pb_panel.ColumnConfig{
        Id: &column_id,
        PrimaryPumpState: &pumpState,
        SecondaryPumpState: &pumpState,
        MixingState: &mixingState,
        WaterLevelState: &mixingState,
        Nodes: make([]*pb_panel.NodeConfig, 0),
    }

    ppfd := float32(0.0)
    node := pb_panel.NodeConfig{
        Id: &node_id,
        PPFD: &ppfd,
    }
    node_2 := pb_panel.NodeConfig{
        Id: &node_2_id,
        PPFD: &ppfd,
    }
    
    column.Nodes = append(column.Nodes, &node)
    column.Nodes = append(column.Nodes, &node_2)
    Config.Columns = append(Config.Columns,&column)
    DefaultConfig = Config
}
