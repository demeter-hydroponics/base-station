package farm_config

import (
    "sync"
	pb_panel "base-station/protobuf/generated/go/panel"
)

/*
This Might be temporary, might be a bad idea.
I want a single function/handler to be able to update the config at once, others will be rejected temporarily until the current one is processed.
*/
var ConfigMutex sync.Mutex
var Config pb_panel.FarmConfig 
