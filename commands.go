package tesla

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// CommandResponse is the response from the Tesla API after POSTing a command.
type CommandResponse struct {
	Response struct {
		Reason string `json:"reason"`
		Result bool   `json:"result"`
	} `json:"response"`
}

// AutoParkRequest are the required elements to POST an Autopark/Summon request for the vehicle.
type AutoParkRequest struct {
	VehicleID uint64  `json:"vehicle_id,omitempty"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	Action    string  `json:"action,omitempty"`
}

// SentryData shows whether Sentry is on.
type SentryData struct {
	Mode string `json:"on"`
}

// AutoparkAbort causes the vehicle to abort the Autopark request.
func (v *Vehicle) AutoparkAbort() error {
	return v.autoPark("abort")
}

// AutoparkForward causes the vehicle to pull forward.
func (v *Vehicle) AutoparkForward() error {
	return v.autoPark("start_forward")
}

// AutoparkReverse causes the vehicle to go in reverse.
func (v *Vehicle) AutoparkReverse() error {
	return v.autoPark("start_reverse")
}

// Performs the actual auto park/summon request for the vehicle
func (v *Vehicle) autoPark(action string) error {
	apiURL := v.commandPath("autopark_request")
	driveState, _ := v.DriveState()
	autoParkRequest := &AutoParkRequest{
		VehicleID: v.VehicleID,
		Lat:       driveState.Latitude,
		Lon:       driveState.Longitude,
		Action:    action,
	}
	body, _ := json.Marshal(autoParkRequest)

	_, err := v.sendCommand(apiURL, body)
	return err
}

// EnableSentry enables Sentry Mode
func (v *Vehicle) EnableSentry() error {
	apiURL := v.commandPath("set_sentry_mode")
	sentryRequest := &SentryData{
		Mode: "true",
	}

	body, _ := json.Marshal(sentryRequest)
	_, err := v.sendCommand(apiURL, body)
	return err
}

// TBD based on Github issue #7
// Toggles defrost on and off, locations values are 'front' or 'rear'
// func (v *Vehicle) Defrost(location string, state bool) error {
// 	command := location + "_defrost_"
// 	if state {
// 		command += "on"
// 	} else {
// 		command += "off"
// 	}
// 	apiURL := v.c.URL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/" + command
// 	fmt.Println(apiURL)
// 	_, err := v.sendCommand(apiURL, nil)
// 	return err
// }

// TriggerHomelink opens and closes the configured Homelink garage door of the vehicle
// keep in mind this is a toggle and the garage door state is unknown
// a major limitation of Homelink.
func (v *Vehicle) TriggerHomelink() error {
	apiURL := v.commandPath("trigger_homelink")
	driveState, _ := v.DriveState()
	autoParkRequest := &AutoParkRequest{
		Lat: driveState.Latitude,
		Lon: driveState.Longitude,
	}
	body, _ := json.Marshal(autoParkRequest)

	_, err := v.sendCommand(apiURL, body)
	return err
}

// Wakeup wakes up the vehicle when it is powered off.
func (v *Vehicle) Wakeup() (*Vehicle, error) {
	apiURL := v.wakePath()
	body, err := v.sendCommand(apiURL, nil)
	if err != nil {
		return nil, err
	}
	vehicleResponse := &VehicleResponse{}
	if err := json.Unmarshal(body, vehicleResponse); err != nil {
		return nil, err
	}
	vehicleResponse.Response.c = v.c
	return vehicleResponse.Response, nil
}

// OpenChargePort opens the charge port so you may insert your charging cable.
func (v *Vehicle) OpenChargePort() error {
	apiURL := v.commandPath("charge_port_door_open")
	_, err := v.sendCommand(apiURL, nil)
	return err
}

// ResetValetPIN resets the PIN set for valet mode, if set.
func (v *Vehicle) ResetValetPIN() error {
	apiURL := v.commandPath("reset_valet_pin")
	_, err := v.sendCommand(apiURL, nil)
	return err
}

// SetChargeLimitStandard sets the charge limit to the standard setting.
func (v *Vehicle) SetChargeLimitStandard() error {
	apiURL := v.commandPath("charge_standard")
	_, err := v.sendCommand(apiURL, nil)
	return err
}

// SetChargeLimitMax sets the charge limit to the max limit.
func (v *Vehicle) SetChargeLimitMax() error {
	apiURL := v.commandPath("charge_max_range")
	_, err := v.sendCommand(apiURL, nil)
	return err
}

// SetChargeLimit set the charge limit to a custom percentage.
func (v *Vehicle) SetChargeLimit(percent int) error {
	apiURL := v.commandPath("set_charge_limit")
	payload := `{"percent": ` + strconv.Itoa(percent) + `}`
	_, err := v.c.post(apiURL, []byte(payload))
	return err
}

// SetChargingAmps set the charging amps to a specific value.
func (v *Vehicle) SetChargingAmps(amps int) error {
	apiURL := v.commandPath("set_charging_amps")
	payload := `{"charging_amps": ` + strconv.Itoa(amps) + `}`
	_, err := v.c.post(apiURL, []byte(payload))
	return err
}

// StartCharging starts the charging of the vehicle after you have inserted the charging cable.
func (v *Vehicle) StartCharging() error {
	apiURL := v.commandPath("charge_start")
	_, err := v.sendCommand(apiURL, nil)
	return err
}

// StopCharging stops the charging of the vehicle.
func (v *Vehicle) StopCharging() error {
	apiURL := v.commandPath("charge_stop")
	_, err := v.sendCommand(apiURL, nil)
	return err
}

// FlashLights flashes the lights of the vehicle.
func (v *Vehicle) FlashLights() error {
	apiURL := v.commandPath("flash_lights")
	_, err := v.sendCommand(apiURL, nil)
	return err
}

// HonkHorn honks the horn of the vehicle.
func (v *Vehicle) HonkHorn() error {
	apiURL := v.commandPath("honk_horn")
	_, err := v.sendCommand(apiURL, nil)
	return err
}

// UnlockDoors unlock the vehicle's doors.
func (v *Vehicle) UnlockDoors() error {
	apiURL := v.commandPath("door_unlock")
	_, err := v.sendCommand(apiURL, nil)
	return err
}

// LockDoors locks the doors of the vehicle.
func (v *Vehicle) LockDoors() error {
	apiURL := v.commandPath("door_lock")
	_, err := v.sendCommand(apiURL, nil)
	return err
}

type tempRequest struct {
	DriverTemp    string `json:"driver_temp"`
	PassengerTemp string `json:"passenger_temp"`
}

// SetTemperature sets the temperature of the vehicle, where you may set the driver
// zone and the passenger zone to separate temperatures.
func (v *Vehicle) SetTemperature(driver float64, passenger float64) error {
	driveTemp := strconv.FormatFloat(driver, 'f', -1, 32)
	passengerTemp := strconv.FormatFloat(passenger, 'f', -1, 32)
	apiURL := v.commandPath("set_temps")
	b, err := json.Marshal(&tempRequest{driveTemp, passengerTemp})
	if err != nil {
		return err
	}
	_, err = v.c.post(apiURL, b)
	return err
}

// StartAirConditioning starts the air conditioning in the vehicle.
func (v *Vehicle) StartAirConditioning() error {
	url := v.commandPath("auto_conditioning_start")
	_, err := v.sendCommand(url, nil)
	return err
}

// StopAirConditioning stops the air conditioning in the vehicle.
func (v *Vehicle) StopAirConditioning() error {
	apiURL := v.commandPath("auto_conditioning_stop")
	_, err := v.sendCommand(apiURL, nil)
	return err
}

// SetSeatHeater sets the specified seat's heater level.
func (v *Vehicle) SetSeatHeater(heater int, level int) error {
	url := v.commandPath("remote_seat_heater_request")
	payload := fmt.Sprintf(`{"heater":%d, "level":%d}`, heater, level)
	_, err := v.c.post(url, []byte(payload))
	return err
}

// SetSteeringWheelHeater turns steering wheel heater on or off.
func (v *Vehicle) SetSteeringWheelHeater(on bool) error {
	url := v.commandPath("remote_steering_wheel_heater_request")
	payload := fmt.Sprintf(`{"on":%t}`, on)
	_, err := v.c.post(url, []byte(payload))
	return err
}

// MovePanoRoof sets the desired state of the panoramic roof. The approximate percent open
// values for each state are open = 100%, close = 0%, comfort = 80%, vent = %15, move = set %.
func (v *Vehicle) MovePanoRoof(state string, percent int) error {
	apiURL := v.commandPath("sun_roof_control")
	payload := `{"state": "` + state + `", "percent":` + strconv.Itoa(percent) + `}`
	_, err := v.c.post(apiURL, []byte(payload))
	return err
}

// Controls the windows. Will vent or close all windows simultaneously. command can be "vent" or "close".
// lat and lon values must be near the current location of the car for close operation to succeed.
// For vent, the lat and lon values are ignored, and may both be 0 (which has been observed from the app itself).
func (v *Vehicle) WindowControl(command string, lat, lon float64) error {
	apiURL := v.commandPath("window_control")
	payload := fmt.Sprintf(`{"command":"%s", "lat": %f, "lon": %f}`, command, lat, lon)
	_, err := v.c.post(apiURL, []byte(payload))
	return err
}

// Start starts the car by turning it on, requires the password to be sent again.
func (v *Vehicle) Start(password string) error {
	apiURL := v.commandPath("remote_start_drive?password=" + password)
	_, err := v.sendCommand(apiURL, nil)
	return err
}

// OpenTrunk opens the trunk, where values may be 'front' or 'rear'.
func (v *Vehicle) OpenTrunk(trunk string) error {
	apiURL := v.commandPath("actuate_trunk")
	payload := `{"which_trunk": "` + trunk + `"}`
	_, err := v.c.post(apiURL, []byte(payload))
	return err
}

// Sends a command to the vehicle
func (v *Vehicle) sendCommand(url string, reqBody []byte) ([]byte, error) {
	body, err := v.c.post(url, reqBody)
	if err != nil {
		return nil, err
	}
	if len(body) > 0 {
		response := &CommandResponse{}
		if err := json.Unmarshal(body, response); err != nil {
			return nil, err
		}
		if !response.Response.Result && response.Response.Reason != "" {
			return nil, errors.New(response.Response.Reason)
		}
	}
	return body, nil
}
