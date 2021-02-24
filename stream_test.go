package tesla

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	StreamEventString    = `1460905367,65,9550.3,88,10,76,30.493001,-100.457018,,,227,184,75`
	BadStreamEventString = `1460905367    9550.3,88    76,30.493001,-100.457018,,,227,184,75`
)

func TestStreamSpec(t *testing.T) {
	ts := serveHTTP(t)
	defer ts.Close()

	client := NewTestClient(ts)

	vehicle := &Vehicle{
		c:         client,
		VehicleID: 123,
		Tokens:    []string{"456", "789"},
	}

	Convey("Should get stream events", t, func() {
		eventChan, errChan, err := vehicle.Stream("elon@tesla.com")
		So(err, ShouldBeNil)

		Convey("2 good, 1 bad", func() {
			select {
			case event := <-eventChan:
				So(event.Speed, ShouldEqual, 65)
			case err = <-errChan:
				So(err, ShouldBeNil)
			}
			select {
			case event := <-eventChan:
				So(event.Speed, ShouldEqual, 65)
			case err = <-errChan:
				So(err, ShouldBeNil)
			}
			select {
			case event := <-eventChan:
				So(event, ShouldBeNil)
			case err = <-errChan:
				So(err.Error(), ShouldEqual, "bad message from Tesla API stream")
			}
			select {
			case event := <-eventChan:
				So(event, ShouldBeNil)
			case err = <-errChan:
				So(err.Error(), ShouldEqual, "HTTP stream closed")
			}
		})
	})
}
