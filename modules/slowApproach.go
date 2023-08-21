package modules

import (
	"fmt"
	utils "liam-tool/utils"
	"os"
	"strconv"

	icarus "git.ace/icarus/icarusclient/v5"
)

const DISTANCE_THRESHOLD = 2000
const STOPPING_THRESHOLD = 50

func SlowApproachInput() {

}

func SlowApproach() {
	query := GenerateAuthQuery()

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
	fmt.Println("Insert Altitude")
	var TargetAlt float32
	fmt.Scanln(&TargetAlt)
	newMissionTarget := MissionTarget{}
	newMissionTarget.VehicleID = VehicleId
	newMissionTarget.Latitude = Latitude
	newMissionTarget.Longitude = Longitude
	newMissionTarget.Altitude = TargetAlt
	newMissionTarget.Heading = 0
	newMissionTarget.Velocity = Velocity
	RunGoTo(newMissionTarget.VehicleID, newMissionTarget)

	logFile, err := os.OpenFile(strconv.Itoa(VehicleId)+"_log.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer logFile.Close()

	printWithTimestamp(logFile, strconv.Itoa(VehicleId), "Running: takeoff.go")

	for {
		resp2 := GetDroneInfo(newMissionTarget.VehicleID, query)
		telem := resp2.Vehicles[0].Telem
		lat := telem.Latitude
		long := telem.Longitude
		distance := utils.Haversine(lat, long, newMissionTarget.Latitude, newMissionTarget.Longitude)

		if distance < STOPPING_THRESHOLD {
			newMissionTarget.Velocity = 0
			break
		} else if distance < DISTANCE_THRESHOLD {
			newMissionTarget.Velocity = 30
			RunGoTo(newMissionTarget.VehicleID, newMissionTarget)
		} else {
			RunGoTo(newMissionTarget.VehicleID, newMissionTarget)
		}
	}
}

func SlowApproachGeneral(VehicleId int, Latitude float64, Longitude float64, TargetAlt float32, Velocity float32) {
	query := GenerateAuthQuery()

	newMissionTarget := MissionTarget{}
	newMissionTarget.VehicleID = VehicleId
	newMissionTarget.Latitude = Latitude
	newMissionTarget.Longitude = Longitude
	newMissionTarget.Altitude = TargetAlt
	newMissionTarget.Heading = 0
	newMissionTarget.Velocity = Velocity
	RunGoTo(newMissionTarget.VehicleID, newMissionTarget)

	logFile, err := os.OpenFile(strconv.Itoa(VehicleId)+"_log.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer logFile.Close()

	printWithTimestamp(logFile, strconv.Itoa(VehicleId), "Running: takeoff.go")

	for {
		resp2 := GetDroneInfo(newMissionTarget.VehicleID, query)
		telem := resp2.Vehicles[0].Telem
		lat := telem.Latitude
		long := telem.Longitude
		distance := utils.Haversine(lat, long, newMissionTarget.Latitude, newMissionTarget.Longitude)

		if distance < STOPPING_THRESHOLD {
			newMissionTarget.Velocity = 0
			break
		} else if distance < DISTANCE_THRESHOLD {
			newMissionTarget.Velocity = 30
			RunGoTo(newMissionTarget.VehicleID, newMissionTarget)
		} else {
			RunGoTo(newMissionTarget.VehicleID, newMissionTarget)
		}
	}
}

func RunGoTo(vehicleID int, target MissionTarget) {
	//Create a new query pointed at the IcarusServer instance
	query := GenerateAuthQuery()

	//Change mode to Navigate
	navSeq := query.SetNavMode(vehicleID, icarus.NAVIGATION)
	responseChan, _ := query.Execute()
	response := <-responseChan
	navResponse, ok := response.Get(navSeq)
	if !ok {
		fmt.Println("NAVIGATE response not found")
		return
	}

	if !navResponse.Ok {
		fmt.Println("Error during mode change:", navResponse.Message)
	}
	//Clear the mode change query from the queue
	query.ClearQueries()

	cmdList := GetCmdList(target)

	gotoSeq := query.Goto(vehicleID, cmdList)
	responseChan, _ = query.Execute()
	// fmt.Println("Navigating to waypoint...")
	response = <-responseChan
	gotoResponse, ok := response.Get(gotoSeq)
	if !ok {
		fmt.Println("Go to response not found")
		return
	}

	if gotoResponse.Ok {
		fmt.Printf("Drone %d going to: %f, %f\n", vehicleID, target.Latitude, target.Longitude)
	} else {
		fmt.Println("Error during navigation:", gotoResponse.Message)
	}
}
