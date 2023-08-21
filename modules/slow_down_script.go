package modules

import (
	"fmt"
	utils "liam-tool/utils"

	icarusClient "git.ace/icarus/icarusclient/v5"
)

func GoToLocSlow(telem icarusClient.Telemetry, Latitude, Longitude float64, VehicleId int, maxVel, Altitude float32, query icarusClient.QueryPackage) {
	curLat := telem.Latitude
	curLon := telem.Longitude
	var Velocity float32

	const distanceThreshold = 5000
	distance := utils.Haversine(curLat, curLon, Latitude, Longitude)
	if distance < distanceThreshold {
		temp := float32(distance)
		Velocity = maxVel * (temp / distanceThreshold)
	}
	if Velocity < 15 {
		Velocity = 15
	}

	newMissionTarget := MissionTarget{}
	newMissionTarget.VehicleID = VehicleId
	newMissionTarget.Latitude = Latitude
	newMissionTarget.Longitude = Longitude
	newMissionTarget.Altitude = Altitude
	newMissionTarget.Heading = 0
	newMissionTarget.Velocity = Velocity
	navSeq := query.SetNavMode(VehicleId, icarusClient.NAVIGATION)
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
	fmt.Println(response)
	if !ok {
		fmt.Println("Go to response not found")
		return
	}

	if gotoResponse.Ok {
	} else {
		fmt.Println("Error during navigation:", gotoResponse.Message)
	}
}
