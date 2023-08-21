package main

import (
	"fmt"
	icarus "git.ace/icarus/icarusclient/v5"
	modulesImport "liam-tool/modules"
	"os"
	"os/exec"
)

func GetDroneInfo(IFFID int, query icarus.QueryPackage) icarus.IcarusResponse {
	qID := query.GetVehicleStatus(IFFID)
	respChannel, _ := query.Execute()
	responses := <-respChannel
	resp, ok := responses.Get(qID)
	if !ok {
		fmt.Println("[!] Invalid response")
		return icarus.IcarusResponse{}
	}
	return resp
}

func PrintDroneInfo() {
	query := modulesImport.GenerateAuthQuery()

	// Ask the user for IFFID
	var IFF int
	fmt.Print("Enter an integer for IFFID: ")
	_, err := fmt.Scan(&IFF)
	if err != nil {
		fmt.Println("Error reading IFFID:", err)
		return
	}

	for {
		droneInfo := GetDroneInfo(IFF, query).Vehicles[0]
		telemInfo := droneInfo.Telem
		fuelInfo := droneInfo.PayStatus[5]
		cargoInfo := droneInfo.PayStatus[14]
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
		fmt.Printf("Lat %f\n", telemInfo.Latitude)
		fmt.Printf("Lon %f\n", telemInfo.Longitude)
		fmt.Printf("Alt %f\n", telemInfo.Altitude)
		fmt.Printf("Vel %f\n", telemInfo.Velocity)
		fmt.Printf("Fuel %d\n", fuelInfo.Resources)
		fmt.Printf("Cargo %d\n", cargoInfo.Resources)
		fmt.Println(droneInfo)
	}
}
