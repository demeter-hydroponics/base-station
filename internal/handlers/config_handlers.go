package handlers

import (
	"base-station/internal/core"
	"base-station/internal/database/farm-config"
	"base-station/internal/utils"
	pb_panel "base-station/protobuf/generated/go/panel"
	"io"
	"net/http"
	"strings"

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
    err = utils.JSONToProto(body, &newFarmConfig)
    if err != nil {
        log.Error("There was an error unmarshalling the farm config", "err", err) 
        // send a response
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return 
    }

    log.Info("config recieved, and unmarshalled properly")

    err = core.UpdateConfig(&newFarmConfig)
    log.Info("after updateConfig")
    if err != nil {
        log.Error("There was an error setting the new config", "err", err) 
        if strings.Contains(err.Error(), "recognized") {
            log.Error("Not reporting Id not recognized errors to the front end")
            return
        }
        // send a response
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return 
    }
  	w.WriteHeader(http.StatusOK)
}

func configGetHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    farm_config.ConfigMutex.Lock()
    defer farm_config.ConfigMutex.Unlock()
    
    log.Info("Getting the current Config" ) 
    json_config, err := utils.ProtoToJSON(&farm_config.Config)
    if err != nil {
        log.Error("Error getting config", "err", err)
        http.Error(w, "Error getting config", http.StatusInternalServerError)
        return
    }
  	w.WriteHeader(http.StatusOK)
    w.Write(json_config)
}


func configDefaultGetHandler(w http.ResponseWriter, r *http.Request) {
    log.Info("Getting the default Config" ) 
    w.Header().Set("Content-Type", "application/json")
    json_config, err := utils.ProtoToJSON(&farm_config.DefaultConfig)
    if err != nil {
        log.Error("Error getting config", "err", err)
        http.Error(w, "Error getting config", http.StatusInternalServerError)
        return
    }
  	w.WriteHeader(http.StatusOK)
    w.Write(json_config)
}
