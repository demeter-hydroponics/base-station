package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"

	//	pb_common "base-station/protobuf/generated/go"
	pb_column "base-station/protobuf/generated/go/column"
	//	pb_node "base-station/protobuf/generated/go/node"
	"bytes"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"io"
	"net/http"
)

type TXMessage struct {
	proto_type   string
	json_content string
}

var (
	TxMessages = make(chan *TXMessage, 100)
	model      = NewModel()
	prog       = tea.NewProgram(model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	Commands = map[string]proto.Message{
		"SetPumpStateCommand": &pb_column.SetPumpStateCommand{
			SelectedPump: pb_column.PumpType_BOTH.Enum(),
			State:        pb_column.PumpState_PUMP_OFF.Enum(),
		},
		"PumpUpdateScheduleCommand": &pb_column.PumpUpdateScheduleCommand{
			OnTimeSeconds: UInt32(0),
			PeriodSeconds: UInt32(0),
		},
	}
)

func main() {
	log.SetOutput(io.Discard)
	// start up the server
	go func() {
		http.HandleFunc("/ws", MockServer)
		log.Error("Error in server:", "err", http.ListenAndServe(":12345", nil))
	}()

	if _, err := prog.Run(); err != nil {
		log.Error("There was an error with TUI", "err", err)
	}
}

// JSONToProto converts JSON bytes to a proto2 message
func JSONToProto(data []byte, pb proto.Message) error {
	unmarshaler := &jsonpb.Unmarshaler{
		AllowUnknownFields: false, // More forgiving JSON parsing
	}
	return unmarshaler.Unmarshal(bytes.NewReader(data), pb)
}

// ProtoToJSON converts a proto2 message to JSON bytes
func ProtoToJSON(pb proto.Message) ([]byte, error) {
	marshaler := &jsonpb.Marshaler{
		EmitDefaults: true, // Include fields even if they're at default values
		OrigName:     true, // Use camelCase naming in JSON
		Indent:       "  ",
	}
	json, err := marshaler.MarshalToString(pb)
	return []byte(json), err
}

func UInt32(x uint32) *uint32 {
	return &x
}
