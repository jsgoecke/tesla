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
	Timestamp  string  `json:"timestamp"`
	Speed      int     `json:"speed"`
	Odometer   float64 `json:"odometer"`
	Soc        int     `json:"soc"`
	Elevation  int     `json:"elevation"`
	EstHeading int     `json:"est_heading"`
	EstLat     float64 `json:"est_lat"`
	EstLng     float64 `json:"est_lng"`
	Power      int     `json:"power"`
	ShiftState string  `json:"shift_state"`
	Range      int     `json:"range"`
	EstRange   int     `json:"est_range"`
	Heading    int     `json:"heading"`
}

func (v Vehicle) Stream(c *Client) chan *StreamEvent {
	url := StreamingURL + "/stream/" + strconv.Itoa(v.VehicleID) + "/?values=" + StreamParams
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(c.Auth.Email, v.Tokens[0])
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
			eventChan <- nil // &StreamEvent{}
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
	streamEvent.Speed, _ = strconv.Atoi(data[1])
	streamEvent.Odometer, _ = strconv.ParseFloat(data[2], 64)
	streamEvent.Soc, _ = strconv.Atoi(data[3])
	streamEvent.Elevation, _ = strconv.Atoi(data[4])
	streamEvent.EstHeading, _ = strconv.Atoi(data[5])
	streamEvent.EstLat, _ = strconv.ParseFloat(data[6], 64)
	streamEvent.EstLng, _ = strconv.ParseFloat(data[7], 64)
	streamEvent.Power, _ = strconv.Atoi(data[8])
	streamEvent.ShiftState = data[9]
	streamEvent.Range, _ = strconv.Atoi(data[10])
	streamEvent.EstRange, _ = strconv.Atoi(data[11])
	streamEvent.Heading, _ = strconv.Atoi(data[12])
	return streamEvent
}
