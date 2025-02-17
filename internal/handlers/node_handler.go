package handlers


import (
	"encoding/json"
	"time"

	"github.com/charmbracelet/log"

	"errors"
	"strings"

	"github.com/gorilla/websocket"

	pb_common "base-station/protobuf/generated/go"
	pb_column "base-station/protobuf/generated/go/column"
	pb_node "base-station/protobuf/generated/go/node"
	"io"

	"net/http"

	"github.com/golang/protobuf/proto"
)

func ProcessMessage(reader io.Reader, buf [1024]byte) error {
	n, err := io.ReadFull(reader, buf[0:16])
	if err == io.EOF {
		return err
	}
	if err != nil {
		log.Error("there was an error reading a header in", "err", err)
		return err
	}
	if n != 16 {
		log.Error("16 bytes were not read in", "bytes read", n)
		return errors.New("header was not all read in")
	}

	// unmarshal the header
	header := pb_common.MessageHeader{}
	err = proto.Unmarshal(buf[0:16], &header)
	if err != nil {
		log.Error("error in unmarshalling header", "err", err)
		return err
	}
	// determine the size of the message and read in the message
	msgSize := *header.Length
	//log.Info("header recieved!", "header", header.String())

	// read in the message
	n, err = io.ReadFull(reader, buf[0:msgSize])
	if err != nil {
		log.Error("there was an error reading a header in", "err", err)
		return err
	}
	if uint32(n) != msgSize {
		log.Error("incorrect number of bytes were read in", "bytes read", n, "bytes expected", msgSize)
		return errors.New("incorrect number of bytes were read in for message")
	}

	var msg proto.Message
	var msgType string

    // ADD NEW TYPES HERE
	switch *header.Channel {
	case pb_common.MessageChannels_MIXING_STATS:
		msgType = "Mixing Stats"
		msg = &pb_column.MixingTankStats{}
	case pb_common.MessageChannels_NODE_STATS:
		msgType = "Node Stats"
		msg = &pb_node.NodeStats{}
	case pb_common.MessageChannels_PUMP_MANAGER_INFO:
		msgType = "Pump Manager Info"
		msg = &pb_column.PumpManagerInfo{}
	case pb_common.MessageChannels_PUMP_STATS:
		msgType = "Pump Tank Stats"
		msg = &pb_column.PumpTankStats{}
	}

	err = proto.Unmarshal(buf[0:msgSize], msg)
	if err != nil {
		log.Error("error in unmarshalling message", "err", err)
		return err
	}
	log.Info("msg recieved!", "msg", msg.String())

	marshalled, err := json.MarshalIndent(msg, "", "  ")
	if err != nil {
		log.Error("could not pretty print")
		return nil
	}

	msg_time := time.UnixMicro(int64(header.GetTimestamp()))
	string_msg := strings.Replace(string(marshalled), "\"", "", -1)
	string_msg = msg_time.Format(time.RFC3339Nano) + " - " +
		"IN >>>>>>>>>>>>\n" +
		msgType + ": " +
		string_msg
	log.Info(string_msg)

	return nil
    // do something with the message
}

func Node_Handler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("upgrade :", "err", err)
		return
	}
	defer c.Close()

	log.Info("recieved websocket connection")
	// create a buffer for the header message
	var buf [1024]byte
	for {
		messageType, reader, err := c.NextReader()
		// TODO do a check on message type, this should be a binary or skip it
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
				log.Info("socket has been closed by client")
				return
			}
			log.Error("There was an error with loading the next message:", "err", err)
			continue
		}
		log.Info("Got a new frame")

		if messageType != websocket.TextMessage {
			// this could be a ping/pong, a close or a text message. ignore all for now
			_, err := io.ReadAll(reader)
			if err != nil {
				log.Error("there was an error clearing the reader", "err", err)
			}
			log.Info("Got a non binary message", "msg type", messageType)
			continue
		}

		for {
			log.Info("Processing a Message")
			err = ProcessMessage(reader, buf)
			if err != nil {
				if err == io.EOF {
					log.Info("EOF found, Message over")
					break
				}
				log.Error("encountered an error in reading message", "err", err)

			}
		}

	}
}
