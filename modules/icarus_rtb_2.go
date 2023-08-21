package modules

import (
	"fmt"
	"liam-tool/utils"

	icarus "git.ace/icarus/icarusclient/v5"
)

type RadarPing struct {
	VehicleID           int
	VehicleType         int
	Latitude, Longitude float64
	Altitude, Heading   float32
}

type MissionTarget struct {
	VehicleID                   int
	Latitude, Longitude         float64
	Altitude, Heading, Velocity float32
	Linger                      int
}

func RTB() {
	var VehicleId int
	newMissionTarget := MissionTarget{}
	fmt.Println("Insert Vehicle ID Number")
	fmt.Scanln(&VehicleId)
	newMissionTarget.VehicleID = VehicleId
	newMissionTarget.Latitude = 47.6275
	newMissionTarget.Longitude = -61.4160
	newMissionTarget.Altitude = 25.0
	newMissionTarget.Heading = 0
	fmt.Println("Insert Velocity")
	var Velocity float32
	fmt.Scanln(&Velocity)
	newMissionTarget.Velocity = Velocity
	fmt.Println("hello")
	runMission(newMissionTarget.VehicleID, newMissionTarget)
}

func RunRTB(VehicleId int, Altitude float32, Velocity float32, Latitude float64, Longitude float64) {
	newMissionTarget := MissionTarget{}
	newMissionTarget.VehicleID = VehicleId
	newMissionTarget.Latitude = Latitude
	newMissionTarget.Longitude = Longitude
	newMissionTarget.Altitude = Altitude
	newMissionTarget.Heading = 0
	newMissionTarget.Velocity = Velocity
	runMission(newMissionTarget.VehicleID, newMissionTarget)
}

func RunRTBWithSlowdown(droneTelem icarus.Telemetry, VehicleId int, Altitude float32, Velocity float32, Latitude float64, Longitude float64) {
	//Altitude = 50
	distance := utils.Haversine(droneTelem.Latitude, droneTelem.Longitude, Latitude, Longitude)
	if distance <= 200 && droneTelem.Velocity <= 15 {
		Altitude = 0
	}
	if distance < 3000 {
		Velocity = Velocity * (float32(distance) / 3000)
		if Velocity < 15 {
			Velocity = 15
		}
	}
	fmt.Printf("%d: %fm/s\t\t%f remaining\n", VehicleId, Velocity, distance)
	newMissionTarget := MissionTarget{}
	newMissionTarget.VehicleID = VehicleId
	newMissionTarget.Latitude = Latitude
	newMissionTarget.Longitude = Longitude
	newMissionTarget.Altitude = Altitude
	newMissionTarget.Heading = 0
	newMissionTarget.Velocity = Velocity
	runMission(newMissionTarget.VehicleID, newMissionTarget)
}

func runMission(vehicleID int, target MissionTarget) {
	//Create a new query pointed at the IcarusServer instance
	query := GenerateAuthQuery()

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
	cmdList := GetCmdList(target)

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

func GetCmdList(target MissionTarget) []icarus.Cmd {
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
