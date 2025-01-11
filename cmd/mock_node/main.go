package main

import (
	pb_common "base-station/protobuf/generated/go/common"
	pb_pump "base-station/protobuf/generated/go/pump"
	"github.com/charmbracelet/log"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/websocket"
	"time"
)

func main() {
	// config the address and origin of the websocket server
	origin := "http://localhost/"
	url := "ws://localhost:12345/ws"

	// connect to the server
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()
	for {

		var headerBytes []byte
		var statBytes []byte

		tds_ppm := float32(0.24)
		tds_v := float32(0.25)

		ph_mol := float32(0.42)
		ph_v := float32(0.52)

		sv := pb_common.SensorValidity_VALID

		stats := pb_pump.MixingTankStats{
			TDSSense: &pb_pump.TDSSensor{TDSSensePPM: &tds_ppm, AnalogVoltage: &tds_v, Validity: &sv},
			PHSense:  &pb_pump.PHSensor{PhSenseMolPerL: &ph_mol, AnalogVoltage: &ph_v, Validity: &sv},
		}
		if statBytes, err = proto.Marshal(&stats); err != nil {
			log.Fatal("failure in marshalling mixing stats message", "err", err)
		}

		stamp := uint64(time.Now().Unix())
		size := uint32(len(statBytes))
		channel := pb_common.MessageChannels_MIXING_STATS
		header := pb_common.MessageHeader{
			Channel:   &channel,
			Timestamp: &stamp,
			Length:    &size,
		}

		if headerBytes, err = proto.Marshal(&header); err != nil {
			log.Fatal("failure in marshalling header packet", "err", err)
		}
		log.Info("length of a header packet: ", "length", len(headerBytes))

		if _, err := ws.Write(append(headerBytes, statBytes...)); err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second)
	}
}
