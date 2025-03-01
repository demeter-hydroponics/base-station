package handlers

import (
	"base-station/internal/core"
	pb_panel "base-station/protobuf/generated/go/panel"
	"io"
	"net/http"

	"github.com/charmbracelet/log"
	//"github.com/golang/protobuf/proto"
)


func configHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        configGetHandler(w, r)
    case http.MethodPost:
        configPostHandler(w, r)
    default:
        http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
    }
}


// gets the message body, converts protobuf struct, calls the config update function
func configPostHandler(w http.ResponseWriter, r *http.Request) {
    body, err := io.ReadAll(r.Body)
    if err != nil {
        log.Error("there was an issue with reading the request body", "err", err) 
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    var newFarmConfig pb_panel.FarmConfig
//    err = proto.Unmarshal(body, &farm_config)
    err = JSONToProto(body, &newFarmConfig)
    if err != nil {
        log.Error("There was an error unmarshalling the farm config", "err", err) 
        // send a response
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return 
    }

    err = core.UpdateConfig(&newFarmConfig)
    if err != nil {
        log.Error("There was an error setting the new config", "err", err) 
        // send a response
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return 
    }
}

func configGetHandler(w http.ResponseWriter, r *http.Request) {
}


