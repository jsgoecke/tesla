package tesla

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"golang.org/x/oauth2"
)

func NewTestClient(ts *httptest.Server) *Client {
	ctx := context.Background()
	tok := &oauth2.Token{
		AccessToken:  "refresh",
		RefreshToken: "refresh",
		Expiry:       time.Now().Add(1 * time.Hour),
	}

	config := &oauth2.Config{
		ClientID: "ownerapi",
		Endpoint: oauth2.Endpoint{
			TokenURL: ts.URL + "/oauth/token",
		},
		Scopes: []string{"openid", "email", "offline_access"},
	}

	client := &Client{
		BaseURL:      ts.URL + "/api/1",
		StreamingURL: ts.URL,
		hc:           config.Client(ctx, tok),
	}
	return client
}

func TestClientSpec(t *testing.T) {
	ts := serveHTTP(t)
	defer ts.Close()

	client := NewTestClient(ts)

	Convey("Should set the HTTP headers", t, func() {
		req, _ := http.NewRequest("GET", "http://foo.com", nil)
		client.setHeaders(req)
		So(req.Header.Get("Accept"), ShouldEqual, "application/json")
		So(req.Header.Get("Content-Type"), ShouldEqual, "application/json")
	})
}

func serveHTTP(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		body, _ := ioutil.ReadAll(req.Body)
		defer req.Body.Close()
		switch req.URL.String() {
		case "/oauth/token":
			checkHeaders(t, req)
			Convey("Request body should be set correctly", t, func() {
				/*
					auth := &Auth{}
					json.Unmarshal(body, auth)
					So(auth.ClientID, ShouldEqual, "abc123")
					So(auth.ClientSecret, ShouldEqual, "def456")
					So(auth.Email, ShouldEqual, "elon@tesla.com")
					So(auth.Password, ShouldEqual, "go")
				*/
			})
			w.WriteHeader(200)
			w.Write([]byte("{\"access_token\": \"ghi789\"}"))
		case "/api/1/vehicles":
			checkHeaders(t, req)
			w.WriteHeader(200)
			w.Write([]byte(VehiclesJSON))
		case "/api/1/vehicles/1234":
			checkHeaders(t, req)
			w.WriteHeader(200)
			w.Write([]byte(VehicleJSON))
		case "/api/1/vehicles/1234/mobile_enabled":
			checkHeaders(t, req)
			w.WriteHeader(200)
			w.Write([]byte(TrueJSON))
		case "/api/1/vehicles/1234/data_request/charge_state":
			checkHeaders(t, req)
			w.WriteHeader(200)
			w.Write([]byte(ChargeStateJSON))
		case "/api/1/vehicles/1234/data_request/climate_state":
			w.WriteHeader(200)
			w.Write([]byte(ClimateStateJSON))
		case "/api/1/vehicles/1234/data_request/drive_state":
			checkHeaders(t, req)
			w.WriteHeader(200)
			w.Write([]byte(DriveStateJSON))
		case "/api/1/vehicles/1234/data_request/gui_settings":
			checkHeaders(t, req)
			w.WriteHeader(200)
			w.Write([]byte(GuiSettingsJSON))
		case "/api/1/vehicles/1234/data_request/vehicle_state":
			checkHeaders(t, req)
			w.WriteHeader(200)
			w.Write([]byte(VehicleStateJSON))
		case "/api/1/vehicles/1234/data_request/service_data":
			checkHeaders(t, req)
			w.WriteHeader(200)
			w.Write([]byte(ServiceDataJSON))
		case "/api/1/vehicles/1234/wake_up":
			checkHeaders(t, req)
			w.WriteHeader(200)
			w.Write([]byte(WakeupResponseJSON))
		case "/api/1/vehicles/1234/command/set_charge_limit":
			w.WriteHeader(200)
			Convey("Should receive a set charge limit request", t, func() {
				So(string(body), ShouldEqual, `{"percent": 50}`)
			})
		case "/api/1/vehicles/1234/command/charge_standard":
			checkHeaders(t, req)
			w.WriteHeader(200)
			w.Write([]byte(ChargeAlreadySetJSON))
		case "/api/1/vehicles/1234/command/charge_start":
			checkHeaders(t, req)
			w.WriteHeader(200)
			w.Write([]byte(ChargedJSON))
		case "/api/1/vehicles/1234/command/charge_stop",
			"/api/1/vehicles/1234/command/charge_max_range",
			"/api/1/vehicles/1234/command/charge_port_door_open",
			"/api/1/vehicles/1234/command/flash_lights",
			"/api/1/vehicles/1234/command/honk_horn",
			"/api/1/vehicles/1234/command/auto_conditioning_start",
			"/api/1/vehicles/1234/command/auto_conditioning_stop",
			"/api/1/vehicles/1234/command/door_unlock",
			"/api/1/vehicles/1234/command/door_lock",
			"/api/1/vehicles/1234/command/reset_valet_pin",
			"/api/1/vehicles/1234/command/set_temps?driver_temp=72&passenger_temp=72",
			"/api/1/vehicles/1234/command/remote_start_drive?password=foo":
			checkHeaders(t, req)
			w.WriteHeader(200)
			w.Write([]byte(CommandResponseJSON))
		case "/stream/123/?values=speed,odometer,soc,elevation,est_heading,est_lat,est_lng,power,shift_state,range,est_range,heading":
			w.WriteHeader(200)
			events := StreamEventString + "\n" +
				StreamEventString + "\n" +
				BadStreamEventString + "\n"
			b := bytes.NewBufferString(events)
			b.WriteTo(w)
		case "/api/1/vehicles/1234/command/autopark_request":
			w.WriteHeader(200)
			Convey("Auto park request should have appropriate body", t, func() {
				autoParkRequest := &AutoParkRequest{}
				err := json.Unmarshal(body, autoParkRequest)
				So(err, ShouldBeNil)
				So(autoParkRequest.Action, shouldBeValidAutoparkCommand)
				So(autoParkRequest.VehicleID, ShouldEqual, 456)
				So(autoParkRequest.Lat, ShouldEqual, 35.1)
				So(autoParkRequest.Lon, ShouldEqual, 20.2)
			})
		case "/api/1/vehicles/1234/command/trigger_homelink":
			w.WriteHeader(200)
			Convey("Auto park request should have appropriate body", t, func() {
				autoParkRequest := &AutoParkRequest{}
				err := json.Unmarshal(body, autoParkRequest)
				So(err, ShouldBeNil)
				So(autoParkRequest.Lat, ShouldEqual, 35.1)
				So(autoParkRequest.Lon, ShouldEqual, 20.2)
			})
		case "/api/1/vehicles/1234/command/sun_roof_control":
			w.WriteHeader(200)
			Convey("Should set the Pano roof appropriately", t, func() {
				passed := false
				strBody := string(body)
				if strBody == `{"state": "vent", "percent":0}` {
					passed = true
				}
				if strBody == `{"state": "open", "percent":0}` {
					passed = true
				}
				if strBody == `{"state": "move", "percent":50}` {
					passed = true
				}
				if strBody == `{"state": "close", "percent":0}` {
					passed = true
				}
				So(passed, ShouldBeTrue)
			})
		}
	}))
}

func checkHeaders(t *testing.T, req *http.Request) {
	Convey("HTTP headers should be present", t, func() {
		So(req.Header["Accept"][0], ShouldEqual, "application/json")
		So(req.Header["Content-Type"][0], ShouldEqual, "application/json")
	})
}

func shouldBeValidAutoparkCommand(actual interface{}, expected ...interface{}) string {
	if actual == "start_forward" || actual == "start_reverse" || actual == "abort" {
		return ""
	} else {
		return "The Autopark command should pass start_forward, start_reverse or abort"
	}
}
