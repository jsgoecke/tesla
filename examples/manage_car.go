package main

import (
	"fmt"
	"os"

	"github.com/jsgoecke/tesla"
)

func main() {
	client, err := tesla.NewClient(
		&tesla.Auth{
			ClientID:     os.Getenv("TESLA_CLIENT_ID"),
			ClientSecret: os.Getenv("TESLA_CLIENT_SECRET"),
			Email:        os.Getenv("TESLA_USERNAME"),
			Password:     os.Getenv("TESLA_PASSWORD"),
		})
	if err != nil {
		panic(err)
	}

	vehicles, err := client.Vehicles()
	if err != nil {
		panic(err)
	}
	vehicle := vehicles[0]
	status, err := vehicle.MobileEnabled()

	if err != nil {
		panic(err)
	}

	fmt.Println(status)
	fmt.Println(vehicle.ChargeState())
	fmt.Println(vehicle.ClimateState())
	fmt.Println(vehicle.DriveState())
	fmt.Println(vehicle.GuiSettings())
	fmt.Println(vehicle.VehicleState())
	fmt.Println(vehicle.HonkHorn())
	fmt.Println(vehicle.FlashLights())
	fmt.Println(vehicle.Wakeup())
	fmt.Println(vehicle.OpenChargePort())
	fmt.Println(vehicle.SetChargeLimitStandard())
	fmt.Println(vehicle.StartCharging())
	fmt.Println(vehicle.StopCharging())
	fmt.Println(vehicle.SetChargeLimitMax())
	fmt.Println(vehicle.StartAirConditioning())
	fmt.Println(vehicle.StopAirConditioning())
	fmt.Println(vehicle.LockDoors())
	fmt.Println(vehicle.SetTemprature(72.0, 72.0))
	fmt.Println(vehicle.Start(os.Getenv("TESLA_PASSWORD")))
	fmt.Println(vehicle.OpenTrunk("rear"))
	fmt.Println(vehicle.OpenTrunk("front"))
}
