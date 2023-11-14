package tesla

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	VehiclesJSON = `{"response":[{"color":null,"display_name":"Macak","id":1234,"option_codes":"MS04,RENA,AU01,BC0R,BP01,BR01,BS00,CDM0,CH00,PBSB,CW02,DA02,DCF0,DRLH,DSH7,DV4W,FG02,HP00,IDPB,IX01,LP01,ME02,MI00,PA00,PF01,PI01,PK00,PS01,PX00,PX4D,QNEB,RFP2,SC01,SP00,SR01,SU01,TM00,TP03,TR01,UTAB,WTSG,WTX0,X001,X003,X007,X011,X013,X019,X024,X027,X028,X031,X037,X040,YF01,COUS","vehicle_id":456,"vin":"abc123","tokens":["1","2"],"state":"online","id_s":"789","remote_start_enabled":true,"calendar_enabled":true,"notifications_enabled":true,"backseat_token":null,"backseat_token_updated_at":null,"vehicle_config":{"timestamp":1614069716042}}],"count":1}`
	VehicleJSON  = `{"response":{"color":null,"display_name":"Macak","id":1234,"option_codes":"MS04,RENA,AU01,BC0R,BP01,BR01,BS00,CDM0,CH00,PBSB,CW02,DA02,DCF0,DRLH,DSH7,DV4W,FG02,HP00,IDPB,IX01,LP01,ME02,MI00,PA00,PF01,PI01,PK00,PS01,PX00,PX4D,QNEB,RFP2,SC01,SP00,SR01,SU01,TM00,TP03,TR01,UTAB,WTSG,WTX0,X001,X003,X007,X011,X013,X019,X024,X027,X028,X031,X037,X040,YF01,COUS","vehicle_id":456,"vin":"abc123","tokens":["1","2"],"state":"online","id_s":"789","remote_start_enabled":true,"calendar_enabled":true,"notifications_enabled":true,"backseat_token":null,"backseat_token_updated_at":null},"count":1}`
)

func TestVehiclesSpec(t *testing.T) {
	ts := serveHTTP(t)
	defer ts.Close()

	client := NewTestClient(ts)

	Convey("Should get vehicles", t, func() {
		vehicles, err := client.Vehicles()
		So(err, ShouldBeNil)
		So(vehicles[0].DisplayName, ShouldEqual, "Macak")
		So(vehicles[0].CalendarEnabled, ShouldBeTrue)
	})
}

func TestVehicle(t *testing.T) {
	ts := serveHTTP(t)
	defer ts.Close()

	client := NewTestClient(ts)

	Convey("Should get vehicle", t, func() {
		vehicle, err := client.Vehicle(1234)
		So(err, ShouldBeNil)
		So(vehicle.DisplayName, ShouldEqual, "Macak")
		So(vehicle.CalendarEnabled, ShouldBeTrue)
	})
}

func TestVehicle_CommandPath(t *testing.T) {
	ts := serveHTTP(t)
	defer ts.Close()

	client := NewTestClient(ts)
	v := &Vehicle{ID: 1, c: client}

	Convey("Should have a URL with /command/", t, func() {
		So(v.commandPath("honk_horn"), ShouldEndWith, "/api/1/vehicles/1/command/honk_horn")
	})
}
