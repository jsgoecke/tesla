package tesla

import "encoding/json"

type Vehicle struct {
	Color                  interface{} `json:"color"`
	DisplayName            string      `json:"display_name"`
	ID                     int64       `json:"id"`
	OptionCodes            string      `json:"option_codes"`
	VehicleID              int         `json:"vehicle_id"`
	Vin                    string      `json:"vin"`
	Tokens                 []string    `json:"tokens"`
	State                  string      `json:"state"`
	IDS                    string      `json:"id_s"`
	RemoteStartEnabled     bool        `json:"remote_start_enabled"`
	CalendarEnabled        bool        `json:"calendar_enabled"`
	NotificationsEnabled   bool        `json:"notifications_enabled"`
	BackseatToken          interface{} `json:"backseat_token"`
	BackseatTokenUpdatedAt interface{} `json:"backseat_token_updated_at"`
}

type VehicleResponse struct {
	Response *Vehicle `json:"response"`
	Count    int      `json:"count"`
}

type Vehicles []struct {
	*Vehicle
}

type VehiclesResponse struct {
	Response Vehicles `json:"response"`
	Count    int      `json:"count"`
}

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
