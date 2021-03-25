package tesla

import (
	"context"
	"encoding/json"
	"fmt"
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
		baseURL:      ts.URL + "/api/1",
		streamingURL: ts.URL,
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

var testMux = &http.ServeMux{}

func serveHTTP(t *testing.T) *httptest.Server {
	return httptest.NewServer(testMux)
}

func serveJSON(j string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Header.Get("Accept") != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
		if req.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(j))
	}
}

func serveCheck(c func(req *http.Request, body []byte) error) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Header.Get("Accept") != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
		if req.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := c(req, body); err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func init() {
	testMux.HandleFunc("/oauth/token", serveJSON("{\"access_token\": \"ghi789\"}"))
	testMux.HandleFunc("/api/1/vehicles", serveJSON(VehiclesJSON))
	testMux.HandleFunc("/api/1/vehicles/1234", serveJSON(VehicleJSON))
	testMux.HandleFunc("/api/1/vehicles/1234/command/auto_conditioning_start", serveJSON(CommandResponseJSON))
	testMux.HandleFunc("/api/1/vehicles/1234/command/auto_conditioning_stop", serveJSON(CommandResponseJSON))
	testMux.HandleFunc("/api/1/vehicles/1234/command/charge_max_range", serveJSON(CommandResponseJSON))
	testMux.HandleFunc("/api/1/vehicles/1234/command/charge_port_door_open", serveJSON(CommandResponseJSON))
	testMux.HandleFunc("/api/1/vehicles/1234/command/charge_standard", serveJSON(ChargeAlreadySetJSON))
	testMux.HandleFunc("/api/1/vehicles/1234/command/charge_start", serveJSON(ChargedJSON))
	testMux.HandleFunc("/api/1/vehicles/1234/command/charge_stop", serveJSON(CommandResponseJSON))
	testMux.HandleFunc("/api/1/vehicles/1234/command/door_lock", serveJSON(CommandResponseJSON))
	testMux.HandleFunc("/api/1/vehicles/1234/command/door_unlock", serveJSON(CommandResponseJSON))
	testMux.HandleFunc("/api/1/vehicles/1234/command/flash_lights", serveJSON(CommandResponseJSON))
	testMux.HandleFunc("/api/1/vehicles/1234/command/honk_horn", serveJSON(CommandResponseJSON))
	testMux.HandleFunc("/api/1/vehicles/1234/command/reset_valet_pin", serveJSON(CommandResponseJSON))
	testMux.HandleFunc("/api/1/vehicles/1234/data_request/charge_state", serveJSON(ChargeStateJSON))
	testMux.HandleFunc("/api/1/vehicles/1234/data_request/climate_state", serveJSON(ClimateStateJSON))
	testMux.HandleFunc("/api/1/vehicles/1234/data_request/drive_state", serveJSON(DriveStateJSON))
	testMux.HandleFunc("/api/1/vehicles/1234/data_request/gui_settings", serveJSON(GuiSettingsJSON))
	testMux.HandleFunc("/api/1/vehicles/1234/data_request/service_data", serveJSON(ServiceDataJSON))
	testMux.HandleFunc("/api/1/vehicles/1234/data_request/vehicle_state", serveJSON(VehicleStateJSON))
	testMux.HandleFunc("/api/1/vehicles/1234/mobile_enabled", serveJSON(TrueJSON))
	testMux.HandleFunc("/api/1/vehicles/1234/wake_up", serveJSON(WakeupResponseJSON))

	testMux.HandleFunc("/api/1/vehicles/1234/command/remote_start_drive", func(w http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if got, want := req.FormValue("password"), "foo"; got != want {
			http.Error(w, "password expected to be foo", http.StatusPreconditionFailed)
			return
		}
		serveJSON(CommandResponseJSON)
	})

	testMux.HandleFunc("/api/1/vehicles/1234/command/set_temps", serveCheck(func(req *http.Request, body []byte) error {
		if string(body) != `{"driver_temp":"72","passenger_temp":"72"}` {
			return fmt.Errorf("unexpected body %s", body)
		}
		return nil
	}))

	testMux.HandleFunc("/api/1/vehicles/1234/command/set_charge_limit", serveCheck(func(req *http.Request, body []byte) error {
		if string(body) != `{"percent": 50}` {
			return fmt.Errorf("unexpected body %s", body)
		}
		return nil
	}))

	testMux.HandleFunc("/api/1/vehicles/1234/command/autopark_request", serveCheck(func(req *http.Request, body []byte) error {
		apr := &AutoParkRequest{}
		if err := json.Unmarshal(body, apr); err != nil {
			return err
		}
		switch apr.Action {
		case "start_forward", "start_reverse", "abort":
		default:
			return fmt.Errorf("The Autopark command should pass start_forward, start_reverse or abort")
		}
		if g, w := apr.VehicleID, uint64(456); g != w {
			return fmt.Errorf("unexpected vehicle id: got %d want %d", g, w)
		}
		if g, w := apr.Lat, 35.1; g != w {
			return fmt.Errorf("unexpected lat: got %f want %f", g, w)
		}
		if g, w := apr.Lon, 20.2; g != w {
			return fmt.Errorf("unexpected lon: got %f want %f", g, w)
		}
		return nil
	}))

	testMux.HandleFunc("/api/1/vehicles/1234/command/trigger_homelink", serveCheck(func(req *http.Request, body []byte) error {
		apr := &AutoParkRequest{}
		if err := json.Unmarshal(body, apr); err != nil {
			return err
		}
		if g, w := apr.Lat, 35.1; g != w {
			return fmt.Errorf("unexpected lat: got %f want %f", g, w)
		}
		if g, w := apr.Lon, 20.2; g != w {
			return fmt.Errorf("unexpected lon: got %f want %f", g, w)
		}
		return nil
	}))

	testMux.HandleFunc("/api/1/vehicles/1234/command/sun_roof_control", serveCheck(func(req *http.Request, body []byte) error {
		switch string(body) {
		case `{"state": "vent", "percent":0}`:
		case `{"state": "open", "percent":0}`:
		case `{"state": "move", "percent":50}`:
		case `{"state": "close", "percent":0}`:
		default:
			return fmt.Errorf("unknown request %s", body)
		}
		return nil
	}))

	testMux.HandleFunc("/api/1/vehicles/1234/command/set_sentry_mode", serveCheck(func(req *http.Request, body []byte) error {
		switch string(body) {
		case `{"on":"true"}`:
		default:
			return fmt.Errorf("unknown request %s", body)
		}
		return nil
	}))
}
