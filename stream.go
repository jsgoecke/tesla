package tesla

import (
	"bufio"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var (
	StreamParams = "speed,odometer,soc,elevation,est_heading,est_lat,est_lng,power,shift_state,range,est_range,heading"
	StreamingURL = "https://streaming.vn.teslamotors.com"
)

type StreamEvent struct {
	Timestamp  string `json:"timestamp"`
	Speed      string `json:"speed"`
	Odometer   string `json:"odometer"`
	Soc        string `json:"soc"`
	Elevation  string `json:"elevation"`
	EstHeading string `json:"est_heading"`
	EstLat     string `json:"est_lat"`
	EstLng     string `json:"est_lng"`
	Power      string `json:"power"`
	ShiftState string `json:"shift_state"`
	Range      string `json:"range"`
	EstRange   string `json:"est_range"`
	Heading    string `json:"heading"`
}

func (v Vehicle) Stream(c *Client) chan *StreamEvent {
	url := StreamingURL + "/stream/" + strconv.Itoa(v.VehicleID) + "/?values=" + StreamParams
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(c.Auth.Email, v.Tokens[0])
	if c.Token != nil {
		req.Header.Set("Authorization", "Bearer "+c.Token.AccessToken)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTP.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	eventChan := make(chan *StreamEvent)
	go readStream(resp, eventChan)
	return eventChan
}

func readStream(resp *http.Response, eventChan chan *StreamEvent) {
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			eventChan <- &StreamEvent{}
			break
		} else {
			streamEvent := parseStreamEvent(string(line))
			eventChan <- streamEvent
		}
	}
}

func parseStreamEvent(event string) *StreamEvent {
	data := strings.Split(event, ",")
	streamEvent := &StreamEvent{}
	streamEvent.Timestamp = data[0]
	streamEvent.Speed = data[1]
	streamEvent.Soc = data[2]
	streamEvent.Elevation = data[3]
	streamEvent.EstHeading = data[4]
	streamEvent.EstLat = data[5]
	streamEvent.EstLng = data[6]
	streamEvent.Power = data[7]
	streamEvent.ShiftState = data[8]
	streamEvent.Range = data[9]
	streamEvent.EstRange = data[10]
	streamEvent.Heading = data[11]
	return streamEvent
}
