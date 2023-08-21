package main

import (
	"fmt"

	modules "liam-tool/modules"

	icarus "git.ace/icarus/icarusclient/v5"
)

func runLand() {
	newMissionTarget := MissionTarget{}
	newMissionTarget.VehicleID = 2120
	newMissionTarget.Latitude = 46.1160
	newMissionTarget.Longitude = -64.6753
	newMissionTarget.Altitude = 25.0
	newMissionTarget.Heading = 0
	newMissionTarget.Velocity = 100
	fmt.Println("hello")
	runLandMission(newMissionTarget.VehicleID)
}

func runLandMission(vehicleID int) {

	fmt.Println("hello 2")
	query := modules.GenerateAuthQuery()
	resp, ok := query.Authenticate("pat.whartenby", "J!&Geu7PD+64SZ3")
	if !ok {
		fmt.Println("Unable to authenticate to IcarusServer:", resp)
		return
	}

	fmt.Printf("Vehicle %d at home base\n", vehicleID)

	//Clear the navigation query from the queue
	query.ClearQueries()

	landSeq := query.SetNavMode(vehicleID, icarus.LAND_NOW)
	responseChan, _ := query.Execute()

	fmt.Println("Landing...")
	response := <-responseChan
	landResponse, ok := response.Get(landSeq)
	if !ok {
		fmt.Println("Land response not found")
		return
	}
	if landResponse.Ok {
		fmt.Println("UAV landing initialized")
	} else {
		fmt.Println("Unable to land")
	}
}
