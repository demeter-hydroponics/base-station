package core

import (
	farm_config "base-station/internal/database/farm-config"
	pb_column "base-station/protobuf/generated/go/column"
)

var conversionMap = map[string]int32{
    string(farm_config.PumpOff) : int32(pb_column.PumpState_PUMP_OFF),
    string(farm_config.PumpOn) : int32(pb_column.PumpState_PUMP_ON),
    string(farm_config.PumpSchedule) : int32(pb_column.PumpState_PUMP_SCHEDULE),
    string(farm_config.NoValveOverride) : int32(pb_column.MixingOverrideState_NO_OVERRIDE),
    string(farm_config.OverrideValveOn) : int32(pb_column.MixingOverrideState_OVERRIDE_VALVE_ON),
    string(farm_config.OverrideValveOff) : int32(pb_column.MixingOverrideState_OVERRIDE_VALVE_OFF),
}
