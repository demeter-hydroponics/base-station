package metrics_db

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	"time"
    "net/http"
    "bytes"
    "io"
)

var LOKI_ENDPOINT string = "http://100.123.35.94:3100/loki/api/v1/push"

type Stream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

type Metric struct {
	Streams []Stream `json:"streams"`
}

func SendMetric(id string, payload string) error {
	epoch := fmt.Sprintf("%d", time.Now().UnixNano())

    stream := Stream{
        Stream: map[string]string{"app": "demeter", "controller": id},
        Values: [][]string{ []string{epoch, payload} },
    }
    metric := Metric{Streams: []Stream{stream}}

	jsonData, err := json.Marshal(metric)
	if err != nil {
		log.Errorf("Error marshaling to JSON: %v", err)
		return err
	}

	// Convert byte slice to string and print
	//log.Info(string(jsonData))


    resp, err := http.Post(
		LOKI_ENDPOINT,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		//log.Errorf("Error making request: %v", err)
        return err
	}
	defer resp.Body.Close()
	
	// Read response
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		//log.Errorf("Error reading response: %v", err)
        return err
	}

    //log.Info("body response", "body", string(body))

	return nil
}
