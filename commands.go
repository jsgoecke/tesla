package tesla

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type CommandResponse struct {
	Response struct {
		Reason string `json:"reason"`
		Result bool   `json:"result"`
	} `json:"response"`
}

func (v Vehicle) Wakeup() (*Vehicle, error) {
	url := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/wake_up"
	body, err := sendCommand(url)
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

func (v Vehicle) OpenChargePort() error {
	url := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/charge_port_door_open"
	_, err := sendCommand(url)
	return err
}

func (v Vehicle) SetChargeLimitStandard() error {
	url := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/charge_standard"
	_, err := sendCommand(url)
	return err
}

func (v Vehicle) SetChargeLimitMax() error {
	url := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/charge_max_range"
	_, err := sendCommand(url)
	return err
}

// func (v Vehicle) SetChargeLimit(limit int) error {
// 	url := BaseURL + "/vehicles/" + strconv.Itoa(v.VehicleID) + "/command/set_charge_limit?=" + strconv.Itoa(limit)
// 	_, err := v.Client.postURLEncoded(url, nil)
// 	return err
// }

func (v Vehicle) StartCharging() error {
	url := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/charge_start"
	_, err := sendCommand(url)
	return err
}

func (v Vehicle) StopCharging() error {
	url := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/charge_stop"
	_, err := sendCommand(url)
	return err
}

func (v Vehicle) FlashLights() error {
	url := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/flash_lights"
	_, err := sendCommand(url)
	return err
}

func (v *Vehicle) HonkHorn() error {
	url := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/honk_horn"
	_, err := sendCommand(url)
	return err
}

// func (v Vehicle) UnlockDoors() error {
// 	url := BaseURL + "/vehicles/" + strconv.Itoa(v.VehicleID) + "/command/unlock_doors"
// 	_, err := v.Client.postURLEncoded(url, nil)
// 	return err
// }

func (v Vehicle) LockDoors() error {
	url := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/door_lock"
	_, err := sendCommand(url)
	return err
}

func (v Vehicle) SetTemprature(driver float64, passenger float64) error {
	driveTemp := strconv.FormatFloat(driver, 'f', -1, 32)
	passengerTemp := strconv.FormatFloat(passenger, 'f', -1, 32)
	url := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/set_temps?driver_temp=" + driveTemp + "&passenger_temp=" + passengerTemp
	fmt.Println(url)
	body, err := ActiveClient.post(url, nil)
	fmt.Println(string(body))
	fmt.Println(err)
	return err
}

func (v Vehicle) StartAirConditioning() error {
	url := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/auto_conditioning_start"
	_, err := sendCommand(url)
	return err
}

func (v Vehicle) StopAirConditioning() error {
	url := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/auto_conditioning_stop"
	_, err := sendCommand(url)
	return err
}

// func (v Vehicle) MovePanoRoof(state string, percent int) error {
// 	url := BaseURL + "/vehicles/" + strconv.Itoa(v.VehicleID) + "/command/sun_roof_control?"
// 	_, err := v.Client.postURLEncoded(url, nil)
// 	return err
// }

func (v Vehicle) Start(password string) error {
	url := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/remote_start_drive?password=" + password
	_, err := sendCommand(url)
	return err
}

func (v Vehicle) OpenTrunk(trunk string) error {
	url := BaseURL + "/vehicles/" + strconv.FormatInt(v.ID, 10) + "/command/trunk_open" // ?which_trunk=" + trunk
	fmt.Println(url)
	theJson := `{"which_trunk": "` + trunk + `"}`
	fmt.Println(theJson)
	body, err := ActiveClient.post(url, []byte(theJson))
	fmt.Println(body)
	return err
}

func sendCommand(url string) ([]byte, error) {
	fmt.Println(url)
	body, err := ActiveClient.post(url, nil)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(body))
	response := &CommandResponse{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}
	if response.Response.Result != true && response.Response.Reason != "" {
		return nil, errors.New(response.Response.Reason)
	}
	return body, nil
}
