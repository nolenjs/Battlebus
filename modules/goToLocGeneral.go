package modules

import (
	"fmt"

	icarus "git.ace/icarus/icarusclient/v5"
)

func GoToLocGeneral(Latitude float64, Longitude float64, VehicleId int, Velocity float32, Altitude float32, query icarus.QueryPackage) {
	newMissionTarget := MissionTarget{}
	newMissionTarget.VehicleID = VehicleId
	newMissionTarget.Latitude = Latitude
	newMissionTarget.Longitude = Longitude
	newMissionTarget.Altitude = Altitude
	newMissionTarget.Heading = 0
	newMissionTarget.Velocity = Velocity
	navSeq := query.SetNavMode(VehicleId, icarus.NAVIGATION)
	responseChan, _ := query.Execute()
	response := <-responseChan
	navResponse, ok := response.Get(navSeq)
	if !ok {
		fmt.Println("NAVIGATE response not found")
		return
	}
	if navResponse.Ok {
	} else {
		fmt.Println("Error during mode change:", navResponse.Message)
	}
	//Clear the mode change query from the queue
	query.ClearQueries()
	
	cmdList := GetCmdList(newMissionTarget)
	gotoSeq := query.Goto(VehicleId, cmdList)
	responseChan, _ = query.Execute()
	response = <-responseChan
	gotoResponse, ok := response.Get(gotoSeq)
	if !ok {
		fmt.Println("Go to response not found")
		return
	}

	if !gotoResponse.Ok {
		fmt.Println("Error during navigation:", gotoResponse.Message)
	}
}

func StopDrone(VehicleID int, query icarus.QueryPackage) {
	droneInfo := GetDroneInfo(VehicleID, query)
	droneTelem := droneInfo.Vehicles[0].Telem
	GoToLocGeneral(droneTelem.Latitude, droneTelem.Longitude, VehicleID, 0, droneTelem.Altitude, query)
}
