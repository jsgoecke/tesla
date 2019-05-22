package tesla

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	TrueJSON         = `{"response":true}`
	ChargeStateJSON  = `{"response":{"charging_state":"Complete","charge_limit_soc":90,"charge_limit_soc_std":90,"charge_limit_soc_min":50,"charge_limit_soc_max":100,"charge_to_max_range":false,"battery_heater_on":null,"not_enough_power_to_heat":null,"max_range_charge_counter":0,"fast_charger_present":null,"fast_charger_type":"\u003Cinvalid\u003E","battery_range":235.92,"est_battery_range":200.46,"ideal_battery_range":304.73,"battery_level":90,"usable_battery_level":90,"battery_current":null,"charge_energy_added":19.94,"charge_miles_added_rated":64.5,"charge_miles_added_ideal":83.0,"charger_voltage":null,"charger_pilot_current":null,"charger_actual_current":null,"charger_power":null,"time_to_full_charge":0.0,"trip_charging":null,"charge_rate":0.0,"charge_port_door_open":null,"motorized_charge_port":true,"scheduled_charging_start_time":null,"scheduled_charging_pending":false,"user_charge_enable_request":null,"charge_enable_request":true,"eu_vehicle":false,"charger_phases":null,"charge_port_latch":"\u003Cinvalid\u003E","charge_current_request":40,"charge_current_request_max":40,"managed_charging_active":false,"managed_charging_user_canceled":false,"managed_charging_start_time":null}}`
	ClimateStateJSON = `{"response":{"inside_temp":null,"outside_temp":null,"driver_temp_setting":22.0,"passenger_temp_setting":22.0,"left_temp_direction":17,"right_temp_direction":17,"is_auto_conditioning_on":null,"is_front_defroster_on":false,"is_rear_defroster_on":false,"fan_status":null,"is_climate_on":false,"min_avail_temp":15,"max_avail_temp":28,"seat_heater_left":0,"seat_heater_right":0,"seat_heater_rear_left":0,"seat_heater_rear_right":0,"seat_heater_rear_center":0,"seat_heater_rear_right_back":0,"seat_heater_rear_left_back":0,"smart_preconditioning":false}}`
	DriveStateJSON   = `{"response":{"shift_state":null,"speed":null,"latitude":35.1,"longitude":20.2,"heading":57,"gps_as_of":1452491619}}`
	GuiSettingsJSON  = `{"response":{"gui_distance_units":"mi/hr","gui_temperature_units":"F","gui_charge_rate_units":"mi/hr","gui_24_hour_time":true,"gui_range_display":"Rated"}}`
	VehicleStateJSON = `{"response":{"api_version":3,"calendar_supported":true,"car_type":"s","car_version":"2.9.12","center_display_state":0,"dark_rims":false,"df":0,"dr":0,"exterior_color":"Black","ft":0,"has_spoiler":true,"locked":true,"notifications_supported":true,"odometer":3738.84633,"parsed_calendar_supported":true,"perf_config":"P2","pf":0,"pr":0,"rear_seat_heaters":1,"remote_start":false,"remote_start_supported":true,"rhd":false,"roof_color":"None","rt":0,"seat_type":1,"sun_roof_installed":2,"sun_roof_percent_open":0,"sun_roof_state":"unknown","third_row_seats":"None","valet_mode":false,"vehicle_name":"Macak","wheel_type":"Super21Gray"}}`
)

func TestStatesSpec(t *testing.T) {
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

	Convey("Should get mobile enabled status", t, func() {
		vehicles, err := client.Vehicles()
		vehicle := vehicles[0]
		status, err := vehicle.MobileEnabled()
		So(err, ShouldBeNil)
		So(status, ShouldBeTrue)
	})

	Convey("Should get charge state", t, func() {
		vehicles, err := client.Vehicles()
		vehicle := vehicles[0]
		status, err := vehicle.ChargeState()
		So(err, ShouldBeNil)
		So(status.BatteryLevel, ShouldEqual, 90)
		So(status.ChargeRate, ShouldEqual, 0)
		So(status.ChargingState, ShouldEqual, "Complete")
	})

	Convey("Should get climate state", t, func() {
		vehicles, err := client.Vehicles()
		vehicle := vehicles[0]
		status, err := vehicle.ClimateState()
		So(err, ShouldBeNil)
		So(status.DriverTempSetting, ShouldEqual, 22.0)
		So(status.PassengerTempSetting, ShouldEqual, 22.0)
		So(status.IsRearDefrosterOn, ShouldBeFalse)
	})

	Convey("Should get drive state", t, func() {
		vehicles, err := client.Vehicles()
		vehicle := vehicles[0]
		status, err := vehicle.DriveState()
		So(err, ShouldBeNil)
		So(status.Latitude, ShouldEqual, 35.1)
		So(status.Longitude, ShouldEqual, 20.2)
	})

	Convey("Should get GUI settings", t, func() {
		vehicles, err := client.Vehicles()
		vehicle := vehicles[0]
		status, err := vehicle.GuiSettings()
		So(err, ShouldBeNil)
		So(status.GuiDistanceUnits, ShouldEqual, "mi/hr")
		So(status.GuiTemperatureUnits, ShouldEqual, "F")
	})

	Convey("Should get Vehicle state", t, func() {
		vehicles, err := client.Vehicles()
		vehicle := vehicles[0]
		status, err := vehicle.VehicleState()
		So(err, ShouldBeNil)
		So(status.APIVersion, ShouldEqual, 3)
		So(status.CalendarSupported, ShouldBeTrue)
		So(status.Rt, ShouldEqual, 0)
	})

	AuthURL = previousAuthURL
	BaseURL = previousURL
}
