package tesla

import (
	"encoding/json"
	"errors"
	"strconv"
)

// Response from the Tesla API after POSTing a command
type CommandResponse struct {
	Response struct {
		Reason string `json:"reason"`
		Result bool   `json:"result"`
	} `json:"response"`
}

// TBD based on Github issue #7
// Toggles defrost on and off, locations values are 'front' or 'rear'
// func (v Vehicle) Defrost(location string, state bool) error {
// 	command := location + "_defrost_"
// 	if state {
// 		command += "on"
// 	} else {
// 		command += "off"
// 	}
// 	apiUrl := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/" + command
// 	fmt.Println(apiUrl)
// 	_, err := sendCommand(apiUrl, nil)
// 	return err
// }

// Wakes up the vehicle when it is powered off
func (v Vehicle) Wakeup() (*Vehicle, error) {
	apiUrl := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/wake_up"
	body, err := sendCommand(apiUrl, nil)
	if err != nil {
		return nil, err
	}
	vehicleResponse := &VehicleResponse{}
	err = json.Unmarshal(body, vehicleResponse)
	if err != nil {
		return nil, err
	}
	return vehicleResponse.Response, nil
}

// Opens the charge port so you may insert your charging cable
func (v Vehicle) OpenChargePort() error {
	apiUrl := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/charge_port_door_open"
	_, err := sendCommand(apiUrl, nil)
	return err
}

// Resets the PIN set for valet mode, if set
func (v Vehicle) ResetValetPIN() error {
	apiUrl := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/reset_valet_pin"
	_, err := sendCommand(apiUrl, nil)
	return err
}

// Sets the charge limit to the standard setting
func (v Vehicle) SetChargeLimitStandard() error {
	apiUrl := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/charge_standard"
	_, err := sendCommand(apiUrl, nil)
	return err
}

// Sets the charge limit to the max limit
func (v Vehicle) SetChargeLimitMax() error {
	apiUrl := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/charge_max_range"
	_, err := sendCommand(apiUrl, nil)
	return err
}

// Set the charge limit to a custom percentage
func (v Vehicle) SetChargeLimit(percent int) error {
	apiUrl := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/set_charge_limit"
	theJson := `{"percent": ` + strconv.Itoa(percent) + `}`
	_, err := ActiveClient.post(apiUrl, []byte(theJson))
	return err
}

// Starts the charging of the vehicle after you have inserted the
// charging cable
func (v Vehicle) StartCharging() error {
	apiUrl := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/charge_start"
	_, err := sendCommand(apiUrl, nil)
	return err
}

// Stop the charging of the vehicle
func (v Vehicle) StopCharging() error {
	apiUrl := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/charge_stop"
	_, err := sendCommand(apiUrl, nil)
	return err
}

// Flashes the lights of the vehicle
func (v Vehicle) FlashLights() error {
	apiUrl := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/flash_lights"
	_, err := sendCommand(apiUrl, nil)
	return err
}

// Honks the horn of the vehicle
func (v *Vehicle) HonkHorn() error {
	apiUrl := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/honk_horn"
	_, err := sendCommand(apiUrl, nil)
	return err
}

// Unlock the car's doors
func (v Vehicle) UnlockDoors() error {
	apiUrl := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/door_unlock"
	_, err := sendCommand(apiUrl, nil)
	return err
}

// Locks the doors of the vehicle
func (v Vehicle) LockDoors() error {
	apiUrl := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/door_lock"
	_, err := sendCommand(apiUrl, nil)
	return err
}

// Sets the temprature of the vehicle, where you may set the driver
// zone and the passenger zone to seperate temperatures
func (v Vehicle) SetTemperature(driver float64, passenger float64) error {
	apiUrl := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/set_temps"
	body, err := json.Marshal(map[string]interface{}{
		"driver_temp":    driver,
		"passenger_temp": passenger,
	})
	if err != nil {
		return err
	}
	_, err = ActiveClient.post(apiUrl, body)
	return err
}

// Starts the air conditioning in the car
func (v Vehicle) StartAirConditioning() error {
	url := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/auto_conditioning_start"
	_, err := sendCommand(url, nil)
	return err
}

// Stops the air conditioning in the car
func (v Vehicle) StopAirConditioning() error {
	apiUrl := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/auto_conditioning_stop"
	_, err := sendCommand(apiUrl, nil)
	return err
}

// The desired state of the panoramic roof. The approximate percent open
// values for each state are open = 100%, close = 0%, comfort = 80%, vent = %15, move = set %
func (v Vehicle) MovePanoRoof(state string, percent int) error {
	apiUrl := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/sun_roof_control"
	theJson := `{"state": "` + state + `", "percent":` + strconv.Itoa(percent) + `}`
	_, err := ActiveClient.post(apiUrl, []byte(theJson))
	return err
}

// Starts the car by turning it on, requires the password to be sent
// again
func (v Vehicle) Start(password string) error {
	apiUrl := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/remote_start_drive?password=" + password
	_, err := sendCommand(apiUrl, nil)
	return err
}

// Opens the trunk, where values may be 'front' or 'rear'
func (v Vehicle) OpenTrunk(trunk string) error {
	apiUrl := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/trunk_open" // ?which_trunk=" + trunk
	theJson := `{"which_trunk": "` + trunk + `"}`
	_, err := ActiveClient.post(apiUrl, []byte(theJson))
	return err
}

// Sends a command to the vehicle
func sendCommand(url string, reqBody []byte) ([]byte, error) {
	body, err := ActiveClient.post(url, reqBody)
	if err != nil {
		return nil, err
	}
	if len(body) > 0 {
		response := &CommandResponse{}
		err = json.Unmarshal(body, response)
		if err != nil {
			return nil, err
		}
		if response.Response.Result != true && response.Response.Reason != "" {
			return nil, errors.New(response.Response.Reason)
		}
	}
	return body, nil
}
