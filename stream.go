package tesla

import (
	"bufio"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	StreamParams = "speed,odometer,soc,elevation,est_heading,est_lat,est_lng,power,shift_state,range,est_range,heading"
	StreamingURL = "https://streaming.vn.teslamotors.com"
)

// The event returned by the vehicle by the Tesla API
type StreamEvent struct {
	Timestamp  time.Time `json:"timestamp"`
	Speed      int       `json:"speed"`
	Odometer   float64   `json:"odometer"`
	Soc        int       `json:"soc"`
	Elevation  int       `json:"elevation"`
	EstHeading int       `json:"est_heading"`
	EstLat     float64   `json:"est_lat"`
	EstLng     float64   `json:"est_lng"`
	Power      int       `json:"power"`
	ShiftState string    `json:"shift_state"`
	Range      int       `json:"range"`
	EstRange   int       `json:"est_range"`
	Heading    int       `json:"heading"`
}

// Requests a stream from the vehicle and returns a Go channel
func (v Vehicle) Stream() (chan *StreamEvent, chan error, error) {
	url := StreamingURL + "/stream/" + strconv.Itoa(v.VehicleID) + "/?values=" + StreamParams
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(ActiveClient.Auth.Email, v.Tokens[0])
	resp, err := ActiveClient.HTTP.Do(req)

	if err != nil {
		return nil, nil, err
	}

	eventChan := make(chan *StreamEvent)
	errChan := make(chan error)
	go readStream(resp, eventChan, errChan)

	return eventChan, errChan, nil
}

// Reads the stream itself from the vehicle
func readStream(resp *http.Response, eventChan chan *StreamEvent, errChan chan error) {
	reader := bufio.NewReader(resp.Body)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	defer resp.Body.Close()

	for scanner.Scan() {
		streamEvent, err := parseStreamEvent(scanner.Text())
		if err == nil {
			eventChan <- streamEvent
		} else {
			errChan <- err
		}
	}
	errChan <- errors.New("HTTP stream closed")
}

// Parses the stream event, setting all of the appropriate data types
func parseStreamEvent(event string) (*StreamEvent, error) {
	data := strings.Split(event, ",")
	if len(data) != 13 {
		return nil, errors.New("Bad message from Tesla API stream")
	}

	streamEvent := &StreamEvent{}
	timestamp, _ := strconv.ParseInt(data[0], 10, 64)
	streamEvent.Timestamp = time.Unix(0, timestamp*int64(time.Millisecond))
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
	return streamEvent, nil
}
