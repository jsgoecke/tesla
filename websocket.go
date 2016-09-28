package tesla

import (
	"encoding/base64"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	// WebSocketServer is the host to connect to for the Tesla websocket stream.
	WebSocketServer = "streaming.vn.teslamotors.com"
	// WebSocketResource is the HTTP resource to connect to.
	WebSocketResource = "/connect/"
)

type homelinkCommand struct {
	MsgType   string  `json:"msg_type"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// WebSocket encapsulates a controlling websocket to a vehicle.
type WebSocket struct {
	Vehicle *Vehicle
	Output  <-chan map[string]interface{}

	conn *websocket.Conn

	stateWB        sync.WaitGroup
	autoparkState  string
	homelinkNearby bool
}

// Close closes the underlying connection.
func (s *WebSocket) Close() error {
	return s.conn.Close()
}

// ActivateHomelink triggers homelink via this connection.
func (s *WebSocket) ActivateHomelink() error {
	driveState, err := s.Vehicle.DriveState()
	if err != nil {
		return err
	}

	cmd := homelinkCommand{
		MsgType:   "homelink:cmd_trigger",
		Latitude:  driveState.Latitude,
		Longitude: driveState.Longitude,
	}

	log.Printf("Tesla: WriteJSON(%+v)", cmd)
	return s.conn.WriteJSON(cmd)
}

// AutoparkState waits for the state to be loaded over the socket, then returns it.
func (s *WebSocket) AutoparkState() string {
	s.stateWB.Wait()
	return s.autoparkState
}

// HomelinkNearby wait for the state to be loaded over the socket, then returns it.
func (s *WebSocket) HomelinkNearby() bool {
	s.stateWB.Wait()
	return s.homelinkNearby
}

// Returns a WebSocket connected to the vehicle.
func (v *Vehicle) WebSocket() (*WebSocket, error) {
	sockURL := url.URL{Scheme: "wss", Host: WebSocketServer, Path: WebSocketResource + strconv.Itoa(v.VehicleID)}

	data := []byte(v.client.Auth.Email + ":" + v.Tokens[0])
	encodedToken := base64.StdEncoding.EncodeToString(data)
	headers := http.Header{}
	headers.Add("Authorization", "Basic "+encodedToken)

	pipe := make(chan map[string]interface{})
	sock := &WebSocket{
		Vehicle: v,
		Output:  (<-chan map[string]interface{})(pipe),
	}
	// autopark state and homelink nearby
	sock.stateWB.Add(2)

	var err error
	sock.conn, _, err = websocket.DefaultDialer.Dial(sockURL.String(), headers)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			msg := map[string]interface{}{}
			err := sock.conn.ReadJSON(&msg)
			if err != nil {
				close(pipe)
				return
			}
			switch msg["msg_type"] {
			case "autopark:status":
				sock.autoparkState = msg["autopark_state"].(string)
				sock.stateWB.Done()
			case "homelink:status":
				sock.homelinkNearby = msg["homelink_nearby"].(bool)
				sock.stateWB.Done()
			}
			pipe <- msg
		}
	}()

	return sock, nil
}
