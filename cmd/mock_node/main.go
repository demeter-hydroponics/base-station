package main

import (
	pb_common "base-station/protobuf/generated/go"
	pb_node "base-station/protobuf/generated/go/node"
	"github.com/charmbracelet/log"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
    "base-station/internal/handlers"
	"time"
)



func ConvertMessageToFrame(msg proto.Message) ([]byte, error) {
		var headerBytes []byte
		var msgBytes []byte
        var err error
		if msgBytes, err = proto.Marshal(msg); err != nil {
			log.Error("failure in marshalling message", "err", err)
            return []byte{}, err
		}

		stamp := uint64(time.Now().UnixMicro())
		size := uint32(len(msgBytes))
		channel := handlers.PbToChannel(msg) 
		header := pb_common.MessageHeader{
			Channel:   &channel,
			Timestamp: &stamp,
			Length:    &size,
		}

		if headerBytes, err = proto.Marshal(&header); err != nil {
			log.Error("failure in marshalling header packet", "err", err)
            return []byte{}, err
		}
		log.Info("length of a header packet: ", "length", len(headerBytes))
        log.Info("sent data", "header", headerBytes, "msg", msgBytes)

    return append(headerBytes, msgBytes...), nil

}

func main() {
	// config the address and origin of the websocket server
	url := "ws://localhost:12345/ws?id=0a7dbbd3-d42b-4593-9135-8509c2ed520e&type=0"

	// connect to the server
	ws,_, err := websocket.DefaultDialer.Dial(url, nil )
	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()
	for {


		sv := pb_common.SensorValidity_VALID
        index0 := uint32(0)
        index1 := uint32(1)
        ppfd := float32(1.0)
        current := float32(0.5)
        lightSense := pb_node.LightSensorStats{
            SensedPPFD: &ppfd, 
            Validity: &sv,
        }
        growLightMetrics := pb_node.GrowLightStats{
            SetPPFD: &ppfd,
            Current: &current,
            CurrentValid: &sv,
        }
        stats := pb_node.GrowLightSectionStats{
            GrowLightIndex: &index0,
            LightSense: &lightSense,
            GrowLightMetrics: &growLightMetrics,
        }
        stats2 := pb_node.GrowLightSectionStats{
            GrowLightIndex: &index1,
            LightSense: &lightSense,
            GrowLightMetrics: &growLightMetrics,
        }

        statBytes, err := ConvertMessageToFrame(&stats)
        if err != nil {
            log.Error("there was an error with converting the message frame", "err", err)
            return 
        }

        log.Info(statBytes)

		// TODO make this use gorrilla and send as a binary message
		if  err := ws.WriteMessage(websocket.BinaryMessage,statBytes); err != nil {
			log.Fatal(err)
            return
		}
        stat2Bytes, err := ConvertMessageToFrame(&stats2)
        if err != nil {
            log.Error("there was an error with converting the message frame", "err", err)
            return 
        }

		// TODO make this use gorrilla and send as a binary message
		if  err := ws.WriteMessage(websocket.BinaryMessage, stat2Bytes); err != nil {
			log.Fatal(err)
            return 
		}
		time.Sleep(time.Second)
	}
}
