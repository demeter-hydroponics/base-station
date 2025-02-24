package farm_config

import (
	"github.com/google/uuid"
	//	pb_column "base-station/protobuf/generated/go/column"
	"sync"
)

/*
This Might be temporary, might be a bad idea.
I want a single function/handler to be able to update the config at once, others will be rejected temporarily until the current one is processed.
*/

var ConfigMutex sync.Mutex

type ControllerType string

const (
	NodeController   ControllerType = "node"
	ColumnController ControllerType = "column"
)

type ControllerConfig struct {
	Type ControllerType
	Id   uuid.UUID
	Name string
}

type NodeConfig struct {
	CtrlCfg ControllerConfig
	// TODO something about DLI or schedules
}

type PumpState string

const (
	PumpOff      PumpState = "off"
	PumpOn       PumpState = "on"
	PumpSchedule PumpState = "schedule"
)

type MixingState string

const (
	NoValveOverride  MixingState = "no_override"
	OverrideValveOn  MixingState = "override_valve_on"
	OverrideValveOff MixingState = "override_valve_off"
)

type ColumnConfig struct {
	Name    string
	CtrlCfg ControllerConfig
	Nodes   []NodeConfig

	PrimaryPumpState   PumpState
	SecondaryPumpState PumpState
	MixingState        MixingState
}

type FarmConfig struct {
	Columns    []ColumnConfig
	UnsetNodes []NodeConfig
}
