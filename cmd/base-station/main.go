package main

import (
    "base-station/internal/handlers"
	"github.com/charmbracelet/log"

//    "base-station/internal/database/metrics-db"
//	pb_node "base-station/protobuf/generated/go/node"
//    "base-station/internal/utils"

)

func main() {
	log.Info("Starting up Base Station")


//    ppfd := float32(0.5)
//    ref := pb_node.SetPPFDReferenceCommand{PPFD: &ppfd}
//
//    pb_json,err := utils.ProtoToJSON(&ref)
//    if err != nil {
//        log.Error("there was an err", "err", err)
//    }
//    log.Infof("json form: %s", string(pb_json))
//
//    metrics_db.SendMetric("5a0c4ea4-b1c5-4c97-ae43-7e01627dc688&", string(pb_json))

    handlers.Run()
}
