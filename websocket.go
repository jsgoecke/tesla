package tesla

import (
	"encoding/base64"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	// WebSocketServer is the host to connect to for the Tesla websocket stream.
	WebSocketServer = "streaming.vn.teslamotors.com"
	// WebSocketResource is the HTTP resource to connect to.
	WebSocketResource = "/connect/"
)

type autoparkCommand struct {
	MsgType   string  `json:"msg_type"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type heartbeatCommand struct {
	MsgType   string `json:"msg_type"`
	Timestamp int64  `json:"timestamp"`
}

// WebSocket encapsulates a controlling websocket to a vehicle.
type WebSocket struct {
	Vehicle *Vehicle
	Output  <-chan map[string]interface{}

	conn *websocket.Conn

	autoparkStateWG    sync.WaitGroup
	autoparkState      string
	autoparkStateInit  bool
	homelinkNearbyWG   sync.WaitGroup
	homelinkNearby     bool
	homelinkNearbyInit bool
}

// Close closes the underlying connection.
func (s *WebSocket) Close() error {
	return s.conn.Close()
}

func (s *WebSocket) Write(i interface{}) error {
	log.Printf("Tesla: WriteJSON(%+v)", i)
	return s.conn.WriteJSON(i)
}

// AutoparkReverse triggers autopark reverse via this connection.
func (s *WebSocket) AutoparkReverse() error {
	driveState, err := s.Vehicle.DriveState()
	if err != nil {
		return err
	}

	cmd := autoparkCommand{
		MsgType:   "autopark:cmd_reverse",
		Latitude:  driveState.Latitude,
		Longitude: driveState.Longitude,
	}

	return s.Write(cmd)
}

// AutoparkAbort aborts autopark via this connection.
func (s *WebSocket) AutoparkAbort() {
	s.Write(map[string]interface{}{
		"msg_type": "autopark:cmd_abort",
	})
}

// AutoparkForward triggers autopark forward via this connection.
func (s *WebSocket) AutoparkForward() error {
	driveState, err := s.Vehicle.DriveState()
	if err != nil {
		return err
	}

	cmd := autoparkCommand{
		MsgType:   "autopark:cmd_forward",
		Latitude:  driveState.Latitude,
		Longitude: driveState.Longitude,
	}

	return s.Write(cmd)
}

// ActivateHomelink triggers homelink via this connection.
func (s *WebSocket) ActivateHomelink() error {
	driveState, err := s.Vehicle.DriveState()
	if err != nil {
		return err
	}

	cmd := autoparkCommand{
		MsgType:   "homelink:cmd_trigger",
		Latitude:  driveState.Latitude,
		Longitude: driveState.Longitude,
	}

	return s.Write(cmd)
}

// AutoparkState waits for the state to be loaded over the socket, then returns it.
func (s *WebSocket) AutoparkState() string {
	s.autoparkStateWG.Wait()
	return s.autoparkState
}

// HomelinkNearby wait for the state to be loaded over the socket, then returns it.
func (s *WebSocket) HomelinkNearby() bool {
	s.homelinkNearbyWG.Wait()
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
	sock.homelinkNearbyWG.Add(1)
	sock.autoparkStateWG.Add(1)

	var err error
	sock.conn, _, err = websocket.DefaultDialer.Dial(sockURL.String(), headers)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			msg := map[string]interface{}{}
			err := sock.conn.ReadJSON(&msg)
			log.Printf("Tesla: ReadJSON: %+v, %v", msg, err)
			if err != nil {
				close(pipe)
				return
			}
			switch msg["msg_type"] {
			case "control:hello":
				freq := msg["autopark"].(map[string]interface{})["heartbeat_frequency"].(float64)
				go func() {
					for _ = range time.Tick(time.Millisecond * time.Duration(freq)) {
						if err = sock.Write(heartbeatCommand{
							MsgType:   "autopark:heartbeat_app",
							Timestamp: time.Now().UnixNano() / int64(time.Second),
						}); err != nil {
							select {
							case _, ok := <-pipe:
								if ok {
									close(pipe)
								}
							default:
								close(pipe)
							}
							return
						}
					}
				}()
			case "autopark:status":
				sock.autoparkState = msg["autopark_state"].(string)
				if !sock.autoparkStateInit {
					sock.autoparkStateWG.Done()
					sock.autoparkStateInit = true
				}

			case "homelink:status":
				sock.homelinkNearby = msg["homelink_nearby"].(bool)
				if !sock.homelinkNearbyInit {
					sock.homelinkNearbyWG.Done()
					sock.homelinkNearbyInit = true
				}
			}
			pipe <- msg
		}
	}()

	return sock, nil
}
