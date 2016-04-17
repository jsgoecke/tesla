package tesla

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	StreamEventString = `1460828294924,65,9550.3,88,10,76,30.493001,-100.457018,,,227,184,75`
)

func TestStreamSpec(t *testing.T) {
	ts := serveHTTP(t)
	defer ts.Close()
	previousAuthURL := AuthURL
	previousURL := BaseURL
	AuthURL = ts.URL + "/oauth/token"
	BaseURL = ts.URL + "/api/1"

	auth := &Auth{
		GrantType:    "password",
		ClientID:     "abc123",
		ClientSecret: "def456",
		Email:        "elon@tesla.com",
		Password:     "go",
	}
	client, _ := NewClient(auth)
	vehicle := &Vehicle{}
	vehicle.VehicleID = 123
	vehicle.Tokens = []string{"456", "789"}

	previousStreamingURL := StreamingURL
	StreamingURL = ts.URL

	Convey("Should get stream events", t, func() {
		eventChan := vehicle.Stream(client)
		event := <-eventChan
		So(event.Timestamp, ShouldEqual, "1460828294924")
		event = <-eventChan
		So(event.Timestamp, ShouldEqual, "1460828294924")
		event = <-eventChan
		So(event.Timestamp, ShouldEqual, "")
	})

	AuthURL = previousAuthURL
	BaseURL = previousURL
	StreamingURL = previousStreamingURL
}
