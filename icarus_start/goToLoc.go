package main

import (
	"fmt"
	"liam-tool/modules"

	icarus "git.ace/icarus/icarusclient/v5"
)

func runGoToLoc() {
	fmt.Println("Insert Vehicle ID Number")
	var VehicleId int
	fmt.Scanln(&VehicleId)
	fmt.Println("Insert Latitude")
	var Latitude float64
	fmt.Scanln(&Latitude)
	fmt.Println("Insert Longitude")
	var Longitude float64
	fmt.Scanln(&Longitude)
	fmt.Println("Insert Velocity")
	var Velocity float32
	fmt.Scanln(&Velocity)
	newMissionTarget := MissionTarget{}
	newMissionTarget.VehicleID = VehicleId
	newMissionTarget.Latitude = Latitude
	newMissionTarget.Longitude = Longitude
	newMissionTarget.Altitude = 100
	newMissionTarget.Heading = 0
	newMissionTarget.Velocity = Velocity
	runGoTo(newMissionTarget.VehicleID, newMissionTarget)
}

func runGoTo(vehicleID int, target MissionTarget) {
	//Create a new query pointed at the IcarusServer instance
	query := modules.GenerateAuthQuery()

	//Change mode to Navigate
	navSeq := query.SetNavMode(vehicleID, icarus.NAVIGATION)
	responseChan, _ := query.Execute()
	fmt.Println("Entering Navigate mode...")
	response := <-responseChan
	navResponse, ok := response.Get(navSeq)
	if !ok {
		fmt.Println("NAVIGATE response not found")
		return
	}

	if navResponse.Ok {
		fmt.Println("UAV ready to navigate")
	} else {
		fmt.Println("Error during mode change:", navResponse.Message)
	}
	//Clear the mode change query from the queue
	query.ClearQueries()

	cmdList := getCmdList(target)

	for {
		gotoSeq := query.Goto(vehicleID, cmdList)
		responseChan, _ = query.Execute()
		fmt.Println("Navigating to waypoint...")
		response = <-responseChan
		gotoResponse, ok := response.Get(gotoSeq)
		if !ok {
			fmt.Println("Go to response not found")
			return
		}

		if gotoResponse.Ok {
			fmt.Println("Navigating to waypoint")
		} else {
			fmt.Println("Error during navigation:", gotoResponse.Message)
		}
	}
}
