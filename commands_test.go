package tesla

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	CommandResponseJSON  = `{"response":{"reason":"","result":true}}`
	WakeupResponseJSON   = `{"response":{"color":null,"display_name":"Macak","id":900,"option_codes":"MS04,RENA,AU01,BC0R,BP01,BR01,BS00,CDM0,CH00,PBSB,CW02,DA02,DCF0,DRLH,DSH7,DV4W,FG02,HP00,IDPB,IX01,LP01,ME02,MI00,PA00,PF01,PI01,PK00,PS01,PX00,PX4D,QNEB,RFP2,SC01,SP00,SR01,SU01,TM00,TP03,TR01,UTAB,WTSG,WTX0,X001,X003,X007,X011,X013,X019,X024,X027,X028,X031,X037,X040,YF01,COUS","user_id":789,"vehicle_id":456,"vin":"abc123","tokens":["1","2"],"state":"online","id_s":"123","remote_start_enabled":true,"calendar_enabled":true,"notifications_enabled":true,"backseat_token":null,"backseat_token_updated_at":null}}`
	ChargeAlreadySetJSON = `{"response":{"reason":"already_standard","result":false}}`
	ChargedJSON          = `{"response":{"reason":"complete","result":false}}`
)

func TestCommandsSpec(t *testing.T) {
	ts := serveHTTP(t)
	defer ts.Close()

	client := NewTestClient(ts)
	vehicles, err := client.Vehicles()
	if err != nil {
		t.Fatal(err)
	}
	vehicle := vehicles[0]

	Convey("Should auto park abort Autopark", t, func() {
		err := vehicle.AutoparkAbort()
		So(err, ShouldBeNil)
	})

	Convey("Should auto park car forward", t, func() {
		err := vehicle.AutoparkForward()
		So(err, ShouldBeNil)
	})

	Convey("Should auto park car reverse", t, func() {
		err := vehicle.AutoparkReverse()
		So(err, ShouldBeNil)
	})

	Convey("Should turn on sentry mode", t, func() {
		err := vehicle.EnableSentry()
		So(err, ShouldBeNil)
	})

	Convey("Should toggle the garage door based on Homelink", t, func() {
		err := vehicle.TriggerHomelink()
		So(err, ShouldBeNil)
	})

	Convey("Should wakeup the car", t, func() {
		_, err := vehicle.Wakeup()
		So(err, ShouldBeNil)
	})

	Convey("Should flash lights", t, func() {
		err := vehicle.FlashLights()
		So(err, ShouldBeNil)
	})

	Convey("Should honk the horn", t, func() {
		err := vehicle.HonkHorn()
		So(err, ShouldBeNil)
	})

	Convey("Should open the charge port", t, func() {
		err := vehicle.OpenChargePort()
		So(err, ShouldBeNil)
	})

	Convey("Should reset the valet pin", t, func() {
		err := vehicle.ResetValetPIN()
		So(err, ShouldBeNil)
	})

	Convey("Should set the car charge limit", t, func() {
		err := vehicle.SetChargeLimit(50)
		So(err, ShouldBeNil)
	})

	Convey("Should set the car to standard charge level", t, func() {
		err := vehicle.SetChargeLimitStandard()
		So(err.Error(), ShouldEqual, "already_standard")
	})

	Convey("Should attempt to charge the car", t, func() {
		err := vehicle.StartCharging()
		So(err.Error(), ShouldEqual, "complete")
	})

	Convey("Should attempt to stop charging the car", t, func() {
		err := vehicle.StopCharging()
		So(err, ShouldBeNil)
	})

	Convey("Should set charge limit maximum", t, func() {
		err := vehicle.SetChargeLimitMax()
		So(err, ShouldBeNil)
	})

	Convey("Should set air conditioning on", t, func() {
		err := vehicle.StartAirConditioning()
		So(err, ShouldBeNil)
	})

	Convey("Should set air conditioning off", t, func() {
		err := vehicle.StopAirConditioning()
		So(err, ShouldBeNil)
	})

	Convey("Should unlock the doors", t, func() {
		err := vehicle.UnlockDoors()
		So(err, ShouldBeNil)
	})

	Convey("Should lock the doors", t, func() {
		err := vehicle.LockDoors()
		So(err, ShouldBeNil)
	})

	Convey("Should set the temperature", t, func() {
		err := vehicle.SetTemperature(72.0, 72.0)
		So(err, ShouldBeNil)
	})

	Convey("Should start the car", t, func() {
		err := vehicle.Start("foo")
		So(err, ShouldBeNil)
	})

	Convey("Should move the Pano Roof around", t, func() {
		Convey("Should vent the pano roof", func() {
			err := vehicle.MovePanoRoof("vent", 0)
			So(err, ShouldBeNil)
		})
		Convey("Should open the pano roof", func() {
			err := vehicle.MovePanoRoof("open", 0)
			So(err, ShouldBeNil)
		})
		Convey("Should move the pano roof to 50", func() {
			err := vehicle.MovePanoRoof("move", 50)
			So(err, ShouldBeNil)
		})
		Convey("Should close the pano roof", func() {
			err := vehicle.MovePanoRoof("close", 0)
			So(err, ShouldBeNil)
		})
	})
}
