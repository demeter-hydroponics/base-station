package handlers


import (
	"time"

	"github.com/charmbracelet/log"

	"errors"

	"github.com/gorilla/websocket"

	pb_common "base-station/protobuf/generated/go"
	"io"

	"net/http"

	"github.com/golang/protobuf/proto"

)

func ReadMessage(reader io.Reader, buf [1024]byte) error {
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

	var msg proto.Message = channelToPb(*header.Channel)
    if msg == nil {
        return errors.New("header is invalid, or unsupported")
    }

//    // ADD NEW TYPES HERE
//	switch *header.Channel {
//	case pb_common.MessageChannels_MIXING_STATS:
//		msg = &pb_column.MixingTankStats{}
//	case pb_common.MessageChannels_NODE_STATS:
//		msg = &pb_node.NodeStats{}
//	case pb_common.MessageChannels_PUMP_MANAGER_INFO:
//		msg = &pb_column.PumpManagerInfo{}
//	case pb_common.MessageChannels_PUMP_STATS:
//		msg = &pb_column.PumpTankStats{}

	err = proto.Unmarshal(buf[0:msgSize], msg)
	if err != nil {
		log.Error("error in unmarshalling message", "err", err)
		return err
	}
	log.Info("msg recieved!", "msg", msg.String())
    // TODO send this struct through to some goroutine for processing stats if its a stat

	return nil
}

func SenderRoutine(c *websocket.Conn, toSend <-chan proto.Message, quit <-chan bool) {
	for msg := range toSend {
		var pb proto.Message
		var channel = pbToChannel(msg)
        var err error

		// convert protobuf to bytes
		var pb_bytes []byte
		if pb_bytes, err = proto.Marshal(pb); err != nil {
			log.Error("There was an error marshalling the protobuf","err", err.Error())
			continue
		}
		// make the header
		stamp := uint64(time.Now().UnixMicro())
		size := uint32(len(pb_bytes))
		header := pb_common.MessageHeader{
			Channel:   &channel,
			Timestamp: &stamp,
			Length:    &size,
		}

		var header_bytes []byte
		if header_bytes, err = proto.Marshal(&header); err != nil {
			log.Error("There was an error marshalling the header","err", err.Error())
			continue
		}
		if err := c.WriteMessage(websocket.BinaryMessage, append(header_bytes, pb_bytes...)); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
                log.Error("ERR: connection closed by client", "err", err.Error())
			} else {
                log.Error("ERR: could not write outgoing msg", "err", err.Error())
			}
			continue
		}
	}
}

// this will get called on a request from a controller
func Controller_Handler(w http.ResponseWriter, r *http.Request) {
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
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
				log.Info("socket has been closed by client")
				return
			}
			log.Error("There was an error with loading the next message:", "err", err)
			continue
		}
		log.Info("Got a new frame")

		// NOTE do a check on message type, this should be a Text or skip it
		if messageType != websocket.TextMessage {
			// this could be a ping/pong, a close or a text message. ignore all for now
			_, err := io.ReadAll(reader)
			if err != nil {
				log.Error("there was an error clearing the reader", "err", err)
			}
			log.Info("Got a non binary message", "msg type", messageType)
			continue
		}

        // TODO check for connection closed somewhere so you can trigger the send routine to die

		for {
			log.Info("Processing a Message")
			err = ReadMessage(reader, buf)
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
