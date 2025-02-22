package handlers

import (
    "net/http"
	"github.com/gorilla/websocket"
	"github.com/charmbracelet/log"
	pb_common "base-station/protobuf/generated/go"
	pb_column "base-station/protobuf/generated/go/column"
	pb_node "base-station/protobuf/generated/go/node"
	"github.com/golang/protobuf/proto"
)

func check_origin(r *http.Request) bool {
    return true
}

var upgrader = websocket.Upgrader{CheckOrigin: check_origin} // use default options


func Run () {
    log.Info("running server")
		http.HandleFunc("/ws", Controller_Handler)
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
