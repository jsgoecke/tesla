package tesla

import (
	"encoding/json"
	"time"
)

// Represents the vehicle as returned from the Tesla API
type Vehicle struct {
	Color                  interface{}    `json:"color"`
	DisplayName            string         `json:"display_name"`
	ID                     int64          `json:"id"`
	OptionCodes            string         `json:"option_codes"`
	VehicleID              uint64         `json:"vehicle_id"`
	Vin                    string         `json:"vin"`
	Tokens                 []string       `json:"tokens"`
	State                  string         `json:"state"`
	IDS                    string         `json:"id_s"`
	RemoteStartEnabled     bool           `json:"remote_start_enabled"`
	CalendarEnabled        bool           `json:"calendar_enabled"`
	NotificationsEnabled   bool           `json:"notifications_enabled"`
	BackseatToken          interface{}    `json:"backseat_token"`
	BackseatTokenUpdatedAt interface{}    `json:"backseat_token_updated_at"`
	AccessType             string         `json:"access_type"`
	InService              bool           `json:"in_service"`
	APIVersion             int            `json:"api_version"`
	CommandSigning         string         `json:"command_signing"`
	VehicleConfig          *VehicleConfig `json:"vehicle_config"`
}

type VehicleConfig struct {
	CanAcceptNavigationRequests bool      `json:"can_accept_navigation_requests"`
	CanActuateTrunks            bool      `json:"can_actuate_trunks"`
	CarSpecialType              string    `json:"car_special_type"`
	CarType                     string    `json:"car_type"`
	ChargePortType              string    `json:"charge_port_type"`
	DefaultChargeToMax          bool      `json:"default_charge_to_max"`
	EceRestrictions             bool      `json:"ece_restrictions"`
	EUVehicle                   bool      `json:"eu_vehicle"`
	ExteriorColor               string    `json:"exterior_color"`
	HasAirSuspension            bool      `json:"has_air_suspension"`
	HasLudicrousMode            bool      `json:"has_ludicrous_mode"`
	MotorizedChargePort         bool      `json:"motorized_charge_port"`
	Plg                         bool      `json:"plg"`
	RearSeatHeaters             int       `json:"rear_seat_heaters"`
	RearSeatType                int       `json:"rear_seat_type"`
	Rhd                         bool      `json:"rhd"`
	RoofColor                   string    `json:"roof_color"`
	SeatType                    int       `json:"seat_type"`
	SpoilerType                 string    `json:"spoiler_type"`
	SunRoofInstalled            int       `json:"sun_roof_installed"`
	ThirdRowSeats               string    `json:"third_row_seats"`
	Timestamp                   time.Time `json:"timestamp"`
	TrimBadging                 string    `json:"trim_badging"`
	UseRangeBadging             bool      `json:"use_range_badging"`
	WheelType                   string    `json:"wheel_type"`
}

// The response that contains the vehicle details from the Tesla API
type VehicleResponse struct {
	Response *Vehicle `json:"response"`
	Count    int      `json:"count"`
}

// Represents the vehicles from an account, as you could have more than
// one Tesla associated to your account
type Vehicles []struct {
	*Vehicle
}

// The response that contains the vehicles details from the Tesla API
type VehiclesResponse struct {
	Response Vehicles `json:"response"`
	Count    int      `json:"count"`
}

// Fetches the vehicles associated to a Tesla account via the API
func (c *Client) Vehicles() (Vehicles, error) {
	vehiclesResponse := &VehiclesResponse{}
	body, err := c.get(BaseURL + "/vehicles")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, vehiclesResponse)
	if err != nil {
		return nil, err
	}
	return vehiclesResponse.Response, nil
}
