package main

import (
	"encoding/json"

	"github.com/charmbracelet/log"

	"errors"

	"github.com/gorilla/websocket"

	pb_common "base-station/protobuf/generated/go/common"
	pb_node "base-station/protobuf/generated/go/metrics"
	pb_pump "base-station/protobuf/generated/go/pump"
	"io"

	"net/http"

	"github.com/golang/protobuf/proto"
)

func check_origin(r *http.Request) bool {
	return true
}
var upgrader = websocket.Upgrader{CheckOrigin:  check_origin} // use default options

func ProcessMessage(reader io.Reader, buf [1024]byte) error {
	n, err := io.ReadFull(reader, buf[0:16])
	if err == io.EOF {
		return err
	}
	if err != nil {
		log.Error("there was an error reading a header in","err", err)
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
		log.Error("error in unmarshalling header","err", err)
		return err
	}
	// determine the size of the message and read in the message
	msgSize := *header.Length
	log.Info("header recieved!", "header", header.String())

	// read in the message
	n, err = io.ReadFull(reader, buf[0:msgSize])
	if err != nil {
		log.Error("there was an error reading a header in","err", err)
		return err
	}
	if uint32(n) != msgSize {
		log.Error("incorrect number of bytes were read in", "bytes read", n, "bytes expected", msgSize)
		return errors.New("incorrect number of bytes were read in for message")
	}

	// TODO but in a switch statement on msg type
	switch *header.Channel {
	case pb_common.MessageChannels_MIXING_STATS:
		msg := pb_pump.MixingTankStats{}
		err = proto.Unmarshal(buf[0:msgSize], &msg)
		if err != nil {
			log.Error("error in unmarshalling message", "err",err)
			return err
		}
		log.Info("msg recieved!", "msg", msg.String())
		marshalled, err := json.MarshalIndent(msg, "", "  ")
		if err != nil {
			log.Error("could not pretty print")
			return nil
		}
		log.Info("pretty print", "msg", string(marshalled))
	case pb_common.MessageChannels_NODE_STATS:
		msg := pb_node.NodeStats{}
		err = proto.Unmarshal(buf[0:msgSize], &msg)
		if err != nil {
			log.Error("error in unmarshalling message","err", err)
			return err
		}
		log.Info("msg recieved!", "msg", msg.String())
		marshalled, err := json.MarshalIndent(msg, "", "  ")
		if err != nil {
			log.Error("could not pretty print")
			return nil
		}
		log.Info("pretty print", "msg", string(marshalled))

	}

	return nil

}

// TODO gorilla websocket supports 1 reader and 1 writer concurrently so implement this as a goroutine and make another one for sending messages
// Echo the data received on the WebSocket.
func EchoServer(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("upgrade :","err", err)
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
			// this may not be recoverable
			continue
		}
		log.Info("Got a new frame")

		if messageType != websocket.TextMessage {
			// this could be a ping/pong, a close or a text message. ignore all for now
			_, err := io.ReadAll(reader)
			if err != nil {
				log.Error("there was an error clearing the reader","err", err)
			}
			log.Info("Got a non binary message", "msg type", messageType)
			continue
		}

		for {
			// TODO check for EOF error and break out of this loop
			log.Info("Processing a Message")
			err = ProcessMessage(reader, buf)
			if err != nil {
				if err == io.EOF {
					log.Info("EOF found, Message over")
					break
				}
				log.Error( "encountered an error in reading message", "err",err)

			}
		}

	}
}

func main() {
	log.Info("starting server on :12345")
	log.Info("Connect on ws://<ip>:12345/ws")
	http.HandleFunc("/ws", EchoServer)
	log.Error("Error in server:","err", http.ListenAndServe(":12345", nil))
}
