package tesla

import (
	"encoding/json"
	"strconv"
)

// Contains the current charge states that exist within the vehicle
type ChargeState struct {
	ChargingState               string      `json:"charging_state"`
	ChargeLimitSoc              int         `json:"charge_limit_soc"`
	ChargeLimitSocStd           int         `json:"charge_limit_soc_std"`
	ChargeLimitSocMin           int         `json:"charge_limit_soc_min"`
	ChargeLimitSocMax           int         `json:"charge_limit_soc_max"`
	ChargeToMaxRange            bool        `json:"charge_to_max_range"`
	BatteryHeaterOn             bool        `json:"battery_heater_on"`
	NotEnoughPowerToHeat        bool        `json:"not_enough_power_to_heat"`
	MaxRangeChargeCounter       int         `json:"max_range_charge_counter"`
	FastChargerPresent          bool        `json:"fast_charger_present"`
	FastChargerType             string      `json:"fast_charger_type"`
	BatteryRange                float64     `json:"battery_range"`
	EstBatteryRange             float64     `json:"est_battery_range"`
	IdealBatteryRange           float64     `json:"ideal_battery_range"`
	BatteryLevel                int         `json:"battery_level"`
	UsableBatteryLevel          int         `json:"usable_battery_level"`
	BatteryCurrent              interface{} `json:"battery_current"`
	ChargeEnergyAdded           float64     `json:"charge_energy_added"`
	ChargeMilesAddedRated       float64     `json:"charge_miles_added_rated"`
	ChargeMilesAddedIdeal       float64     `json:"charge_miles_added_ideal"`
	ChargerVoltage              interface{} `json:"charger_voltage"`
	ChargerPilotCurrent         interface{} `json:"charger_pilot_current"`
	ChargerActualCurrent        interface{} `json:"charger_actual_current"`
	ChargerPower                interface{} `json:"charger_power"`
	TimeToFullCharge            float64     `json:"time_to_full_charge"`
	TripCharging                interface{} `json:"trip_charging"`
	ChargeRate                  float64     `json:"charge_rate"`
	ChargePortDoorOpen          bool        `json:"charge_port_door_open"`
	MotorizedChargePort         bool        `json:"motorized_charge_port"`
	ScheduledChargingStartTime  interface{} `json:"scheduled_charging_start_time"`
	ScheduledChargingPending    bool        `json:"scheduled_charging_pending"`
	UserChargeEnableRequest     interface{} `json:"user_charge_enable_request"`
	ChargeEnableRequest         bool        `json:"charge_enable_request"`
	EuVehicle                   bool        `json:"eu_vehicle"`
	ChargerPhases               interface{} `json:"charger_phases"`
	ChargePortLatch             string      `json:"charge_port_latch"`
	ChargeCurrentRequest        int         `json:"charge_current_request"`
	ChargeCurrentRequestMax     int         `json:"charge_current_request_max"`
	ManagedChargingActive       bool        `json:"managed_charging_active"`
	ManagedChargingUserCanceled bool        `json:"managed_charging_user_canceled"`
	ManagedChargingStartTime    interface{} `json:"managed_charging_start_time"`
}

// Contains the current climate states availale from the vehicle
type ClimateState struct {
	InsideTemp              float64     `json:"inside_temp"`
	OutsideTemp             float64     `json:"outside_temp"`
	DriverTempSetting       float64     `json:"driver_temp_setting"`
	PassengerTempSetting    float64     `json:"passenger_temp_setting"`
	LeftTempDirection       float64     `json:"left_temp_direction"`
	RightTempDirection      float64     `json:"right_temp_direction"`
	IsAutoConditioningOn    bool        `json:"is_auto_conditioning_on"`
	IsFrontDefrosterOn      bool        `json:"is_front_defroster_on"`
	IsRearDefrosterOn       bool        `json:"is_rear_defroster_on"`
	FanStatus               interface{} `json:"fan_status"`
	IsClimateOn             bool        `json:"is_climate_on"`
	MinAvailTemp            float64     `json:"min_avail_temp"`
	MaxAvailTemp            float64     `json:"max_avail_temp"`
	SeatHeaterLeft          int         `json:"seat_heater_left"`
	SeatHeaterRight         int         `json:"seat_heater_right"`
	SeatHeaterRearLeft      int         `json:"seat_heater_rear_left"`
	SeatHeaterRearRight     int         `json:"seat_heater_rear_right"`
	SeatHeaterRearCenter    int         `json:"seat_heater_rear_center"`
	SeatHeaterRearRightBack int         `json:"seat_heater_rear_right_back"`
	SeatHeaterRearLeftBack  int         `json:"seat_heater_rear_left_back"`
	SmartPreconditioning    bool        `json:"smart_preconditioning"`
}

// Contains the current drive state of the vehicle
type DriveState struct {
	ShiftState interface{} `json:"shift_state"`
	Speed      float64     `json:"speed"`
	Latitude   float64     `json:"latitude"`
	Longitude  float64     `json:"longitude"`
	Heading    int         `json:"heading"`
	GpsAsOf    int64       `json:"gps_as_of"`
}

// Contains the current GUI settings of the vehicle
type GuiSettings struct {
	GuiDistanceUnits    string `json:"gui_distance_units"`
	GuiTemperatureUnits string `json:"gui_temperature_units"`
	GuiChargeRateUnits  string `json:"gui_charge_rate_units"`
	Gui24HourTime       bool   `json:"gui_24_hour_time"`
	GuiRangeDisplay     string `json:"gui_range_display"`
}

// Contains the current state of the vehicle
type VehicleState struct {
	APIVersion              int     `json:"api_version"`
	AutoParkState           string  `json:"autopark_state"`
	AutoParkStateV2         string  `json:"autopark_state_v2"`
	CalendarSupported       bool    `json:"calendar_supported"`
	CarType                 string  `json:"car_type"`
	CarVersion              string  `json:"car_version"`
	CenterDisplayState      int     `json:"center_display_state"`
	DarkRims                bool    `json:"dark_rims"`
	Df                      int     `json:"df"`
	Dr                      int     `json:"dr"`
	ExteriorColor           string  `json:"exterior_color"`
	Ft                      int     `json:"ft"`
	HasSpoiler              bool    `json:"has_spoiler"`
	Locked                  bool    `json:"locked"`
	NotificationsSupported  bool    `json:"notifications_supported"`
	Odometer                float64 `json:"odometer"`
	ParsedCalendarSupported bool    `json:"parsed_calendar_supported"`
	PerfConfig              string  `json:"perf_config"`
	Pf                      int     `json:"pf"`
	Pr                      int     `json:"pr"`
	RearSeatHeaters         int     `json:"rear_seat_heaters"`
	RemoteStart             bool    `json:"remote_start"`
	RemoteStartSupported    bool    `json:"remote_start_supported"`
	Rhd                     bool    `json:"rhd"`
	RoofColor               string  `json:"roof_color"`
	Rt                      int     `json:"rt"`
	SeatType                int     `json:"seat_type"`
	SpoilerType             string  `json:"spoiler_type"`
	SunRoofInstalled        int     `json:"sun_roof_installed"`
	SunRoofPercentOpen      int     `json:"sun_roof_percent_open"`
	SunRoofState            string  `json:"sun_roof_state"`
	ThirdRowSeats           string  `json:"third_row_seats"`
	ValetMode               bool    `json:"valet_mode"`
	VehicleName             string  `json:"vehicle_name"`
	WheelType               string  `json:"wheel_type"`
}

// Represents the request to get the states of the vehicle
type StateRequest struct {
	Response struct {
		*ChargeState
		*ClimateState
		*DriveState
		*GuiSettings
		*VehicleState
	} `json:"response"`
}

// The response when a state is requested
type Response struct {
	Bool bool `json:"response"`
}

// Returns if the vehicle is mobile enabled for Tesla API control
func (v *Vehicle) MobileEnabled() (bool, error) {
	body, err := ActiveClient.get(BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/mobile_enabled")
	if err != nil {
		return false, err
	}
	response := &Response{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return false, err
	}
	return response.Bool, nil
}

// Returns the charge state of the vehicle
func (v *Vehicle) ChargeState() (*ChargeState, error) {
	stateRequest, err := fetchState("/charge_state", v.ID)
	if err != nil {
		return nil, err
	}
	return stateRequest.Response.ChargeState, nil
}

// Returns the climate state of the vehicle
func (v Vehicle) ClimateState() (*ClimateState, error) {
	stateRequest, err := fetchState("/climate_state", v.ID)
	if err != nil {
		return nil, err
	}
	return stateRequest.Response.ClimateState, nil
}

func (v Vehicle) DriveState() (*DriveState, error) {
	stateRequest, err := fetchState("/drive_state", v.ID)
	if err != nil {
		return nil, err
	}
	return stateRequest.Response.DriveState, nil
}

// Returns the GUI settings of the vehicle
func (v Vehicle) GuiSettings() (*GuiSettings, error) {
	stateRequest, err := fetchState("/gui_settings", v.ID)
	if err != nil {
		return nil, err
	}
	return stateRequest.Response.GuiSettings, nil
}

func (v Vehicle) VehicleState() (*VehicleState, error) {
	stateRequest, err := fetchState("/vehicle_state", v.ID)
	if err != nil {
		return nil, err
	}
	return stateRequest.Response.VehicleState, nil
}

// A utility function to fetch the appropriate state of the vehicle
func fetchState(resource string, id int64) (*StateRequest, error) {
	stateRequest := &StateRequest{}
	body, err := ActiveClient.get(BaseURL + "/vehicles/" + strconv.FormatInt(id, 10) + "/data_request" + resource)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, stateRequest)
	if err != nil {
		return nil, err
	}
	return stateRequest, nil
}
