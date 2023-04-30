package tesla

import (
	"strconv"
	"strings"
)

// VehiclePartialResponse represents the vehicle response root data as returned from the Tesla API.
type VehiclePartialResponse struct {
	Color                  interface{} `json:"color"`
	DisplayName            string      `json:"display_name"`
	ID                     int64       `json:"id"`
	OptionCodes            string      `json:"option_codes"`
	VehicleID              uint64      `json:"vehicle_id"`
	Vin                    string      `json:"vin"`
	Tokens                 []string    `json:"tokens"`
	State                  string      `json:"state"`
	IDS                    string      `json:"id_s"`
	RemoteStartEnabled     bool        `json:"remote_start_enabled"`
	CalendarEnabled        bool        `json:"calendar_enabled"`
	NotificationsEnabled   bool        `json:"notifications_enabled"`
	BackseatToken          interface{} `json:"backseat_token"`
	BackseatTokenUpdatedAt interface{} `json:"backseat_token_updated_at"`
	AccessType             string      `json:"access_type"`
	InService              bool        `json:"in_service"`
	APIVersion             int         `json:"api_version"`
	CommandSigning         string      `json:"command_signing"`
}

// Vehicle represents the vehicle as returned from the Tesla API.
type Vehicle struct {
	Color                  interface{} `json:"color"`
	DisplayName            string      `json:"display_name"`
	ID                     int64       `json:"id"`
	OptionCodes            string      `json:"option_codes"`
	VehicleID              uint64      `json:"vehicle_id"`
	Vin                    string      `json:"vin"`
	Tokens                 []string    `json:"tokens"`
	State                  string      `json:"state"`
	IDS                    string      `json:"id_s"`
	RemoteStartEnabled     bool        `json:"remote_start_enabled"`
	CalendarEnabled        bool        `json:"calendar_enabled"`
	NotificationsEnabled   bool        `json:"notifications_enabled"`
	BackseatToken          interface{} `json:"backseat_token"`
	BackseatTokenUpdatedAt interface{} `json:"backseat_token_updated_at"`
	AccessType             string      `json:"access_type"`
	InService              bool        `json:"in_service"`
	APIVersion             int         `json:"api_version"`
	CommandSigning         string      `json:"command_signing"`

	c *Client
}

// VehicleConfig represents the configuration of a vehicle.
type VehicleConfig struct {
	CanAcceptNavigationRequests bool     `json:"can_accept_navigation_requests"`
	CanActuateTrunks            bool     `json:"can_actuate_trunks"`
	CarSpecialType              string   `json:"car_special_type"`
	CarType                     string   `json:"car_type"`
	ChargePortType              string   `json:"charge_port_type"`
	DefaultChargeToMax          bool     `json:"default_charge_to_max"`
	DriverAssist                string   `json:"driver_assist"`
	EceRestrictions             bool     `json:"ece_restrictions"`
	EfficiencyPackage           string   `json:"efficiency_package"`
	EUVehicle                   bool     `json:"eu_vehicle"`
	ExteriorColor               string   `json:"exterior_color"`
	ExteriorTrim                string   `json:"exterior_trim"`
	HasAirSuspension            bool     `json:"has_air_suspension"`
	HasLudicrousMode            bool     `json:"has_ludicrous_mode"`
	InteriorTrimType            string   `json:"interior_trim_type"`
	KeyVersion                  int      `json:"key_version"`
	MotorizedChargePort         bool     `json:"motorized_charge_port"`
	PerformancePackage          string   `json:"performance_package"`
	Plg                         bool     `json:"plg"`
	RearDriveUnit               string   `json:"rear_drive_unit"`
	RearSeatHeaters             int      `json:"rear_seat_heaters"`
	RearSeatType                int      `json:"rear_seat_type"`
	Rhd                         bool     `json:"rhd"`
	RoofColor                   string   `json:"roof_color"`
	SeatType                    int      `json:"seat_type"`
	SpoilerType                 string   `json:"spoiler_type"`
	SunRoofInstalled            int      `json:"sun_roof_installed"`
	ThirdRowSeats               string   `json:"third_row_seats"`
	Timestamp                   timeMsec `json:"timestamp"`
	TrimBadging                 string   `json:"trim_badging"`
	UseRangeBadging             bool     `json:"use_range_badging"`
	WheelType                   string   `json:"wheel_type"`
}

// VehicleResponse contains the vehicle details from the Tesla API.
type VehicleResponse struct {
	Response *Vehicle `json:"response"`
	Count    int      `json:"count"`
}

// VehiclesResponse contains a slice of Vehicles from the Tesla API.
type VehiclesResponse struct {
	Response []*Vehicle `json:"response"`
	Count    int        `json:"count"`
}

// Vehicles fetches the vehicles associated to a Tesla account via the API.
func (c *Client) Vehicles() ([]*Vehicle, error) {
	vehiclesResponse := &VehiclesResponse{}
	if err := c.getJSON(c.baseURL+"/vehicles", vehiclesResponse); err != nil {
		return nil, err
	}
	for _, v := range vehiclesResponse.Response {
		v.c = c
	}
	return vehiclesResponse.Response, nil
}

// Vehicle fetches the vehicle by ID associated to a Tesla account via the API.
func (c *Client) Vehicle(vehicleID int64) (*Vehicle, error) {
	resp := &VehicleResponse{}
	if err := c.getJSON(c.baseURL+"/vehicles/"+strconv.FormatInt(vehicleID, 10), resp); err != nil {
		return nil, err
	}
	resp.Response.c = c
	return resp.Response, nil
}

func (v *Vehicle) basePath() string {
	return strings.Join([]string{v.c.baseURL, "vehicles", strconv.FormatInt(v.ID, 10)}, "/")
}

func (v *Vehicle) commandPath(command string) string {
	return strings.Join([]string{v.basePath(), "command", command}, "/")
}

func (v *Vehicle) wakePath() string {
	return strings.Join([]string{v.basePath(), "wake_up"}, "/")
}
