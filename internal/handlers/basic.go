package handlers

import (
    "net/http"
	//"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/charmbracelet/log"
	farm_config "base-station/internal/database/farm-config"
	pb_common "base-station/protobuf/generated/go"
	pb_column "base-station/protobuf/generated/go/column"
	pb_node "base-station/protobuf/generated/go/node"
    "encoding/json"
	"github.com/golang/protobuf/proto"
    "base-station/internal/core"
)

func check_origin(r *http.Request) bool {
    return true
}

var upgrader = websocket.Upgrader{CheckOrigin: check_origin} // use default options

var metrics_pb chan proto.Message

func Run () {
    go core.ProcessMetrics(core.MetricsChannel)
    log.Info("running server")
	http.HandleFunc("/ws", controllerHandler)
	http.HandleFunc("/config", configHandler)
	http.HandleFunc("/connected", testGetConnectedHandler)
	log.Error("Error in server:", "err", http.ListenAndServe(":12345", nil))
}

func pbToChannel(msg proto.Message) pb_common.MessageChannels {
		var channel pb_common.MessageChannels
        switch msg.(type) {
        case *pb_column.SetPumpStateCommand:
			channel = pb_common.MessageChannels_SET_PUMP_STATE_COMMAND
        case *pb_column.PumpUpdateScheduleCommand:
			channel = pb_common.MessageChannels_PUMP_UPDATE_SCHEDULE_COMMAND
        case *pb_column.SetMixingStateCommand:
            channel = pb_common.MessageChannels_SET_MIXING_STATE_COMMAND
		}
    return channel
}


func channelToPb(channel pb_common.MessageChannels) proto.Message {
	switch channel {
	case pb_common.MessageChannels_MIXING_STATS:
		return &pb_column.MixingTankStats{}
	case pb_common.MessageChannels_NODE_STATS:
		return &pb_node.NodeStats{}
	case pb_common.MessageChannels_PUMP_MANAGER_INFO:
		return &pb_column.PumpManagerInfo{}
	case pb_common.MessageChannels_PUMP_STATS:
		return &pb_column.PumpTankStats{}
	}
    return nil
}


func OnboardController(id, controllerType string) (string, bool, error) {
    // TODO check whether id was provided
    // TODO check if the controllerType is even valid
    // TODO check if the exists/is real

    return id, false, nil

}

func testGetConnectedHandler(w http.ResponseWriter, r *http.Request) {
    // get a list of connected controllers
    log.Info("Getting connected controllers")
    core.SenderChannelsMutex.Lock()
    connectedControllers := make([]string, 0, len(core.SenderChannels))
    
    for k := range core.SenderChannels {
        connectedControllers = append(connectedControllers, k)
    }

    core.SenderChannelsMutex.Unlock()

   w.Header().Set("Content-Type", "application/json")
    
    // Write status code
    w.WriteHeader(http.StatusOK)
    
    // Encode and write the JSON response
    err := json.NewEncoder(w).Encode(connectedControllers)
    if err != nil {
        http.Error(w, "Error encoding response", http.StatusInternalServerError)
        return
    } 
}
