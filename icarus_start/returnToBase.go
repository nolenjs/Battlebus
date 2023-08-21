package main

import (
	"fmt"

	icarus "git.ace/icarus/icarusclient/v5"
)

type MissionTarget struct {
	VehicleID                   int
	Latitude, Longitude         float64
	Altitude, Heading, Velocity float32
	Linger                      int
}

type RadarPing struct {
	VehicleID           int
	VehicleType         int
	Latitude, Longitude float64
	Altitude, Heading   float32
}

func runRTB() {
	fmt.Println("Insert Vehicle ID Number")
	var VehicleId int
	fmt.Scanln(&VehicleId)
	fmt.Println("Insert Latitude")
	var Latitude float64
	fmt.Scanln(&Latitude)
	fmt.Println("Insert Longitude")
	var Longitude float64
	fmt.Scanln(&Longitude)
	newMissionTarget := MissionTarget{}
	newMissionTarget.VehicleID = VehicleId
	newMissionTarget.Latitude = Latitude
	newMissionTarget.Longitude = Longitude
	newMissionTarget.Altitude = 100
	newMissionTarget.Heading = 0
	newMissionTarget.Velocity = 100
	fmt.Println("hello")
	runMission(newMissionTarget.VehicleID, newMissionTarget)
}

func runMission(vehicleID int, target MissionTarget) {
	//Create a new query pointed at the IcarusServer instance
	fmt.Println("hello 2")
	query := icarus.NewQuery("10.51.153.243", "443")
	resp, ok := query.Authenticate("pat.whartenby", "J!&Geu7PD+64SZ3")
	if !ok {
		fmt.Println("Unable to authenticate to IcarusServer:", resp)
		return
	}

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

	//Navigate to home base
	cmdList := getCmdList(target)

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

func getCmdList(target MissionTarget) []icarus.Cmd {
	var cmdList []icarus.Cmd = nil
	var cmdType icarus.CmdType
	if target.Linger > 0 {
		cmdType = icarus.LOITER
	} else {
		cmdType = icarus.GOTO
	}
	cmdList = icarus.AddCmd(cmdList, cmdType, target.Latitude, target.Longitude, target.Altitude, target.Velocity, 0, uint32(target.Linger), target.Heading)
	return cmdList

}
