package tesla

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	StreamEventString = `1460905367,65,9550.3,88,10,76,30.493001,-100.457018,,,227,184,75`
)

func TestStreamSpec(t *testing.T) {
	ts := serveHTTP(t)
	defer ts.Close()
	previousAuthURL := AuthURL
	previousURL := BaseURL
	AuthURL = ts.URL + "/oauth/token"
	BaseURL = ts.URL + "/api/1"
	vehicle := &Vehicle{}
	vehicle.VehicleID = 123
	vehicle.Tokens = []string{"456", "789"}

	previousStreamingURL := StreamingURL
	StreamingURL = ts.URL

	Convey("Should get stream events", t, func() {
		eventChan := vehicle.Stream()
		event := <-eventChan
		So(event.Speed, ShouldEqual, 65)
		event = <-eventChan
		So(event.Speed, ShouldEqual, 65)
		event = <-eventChan
		So(event, ShouldBeNil)
	})

	AuthURL = previousAuthURL
	BaseURL = previousURL
	StreamingURL = previousStreamingURL
}
