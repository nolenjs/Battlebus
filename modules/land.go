package modules

import (
	"fmt"

	icarus "git.ace/icarus/icarusclient/v5"
)

func RunLand() {
	newMissionTarget := MissionTarget{}
	fmt.Println("Insert Vehicle ID Number")
	var VehicleId int
	fmt.Scanln(&VehicleId)
	newMissionTarget.VehicleID = VehicleId
	newMissionTarget.Latitude = 47.6275
	newMissionTarget.Longitude = -61.4160
	newMissionTarget.Altitude = 25.0
	newMissionTarget.Heading = 0
	newMissionTarget.Velocity = 10
	fmt.Println("hello")
	runLandMission(newMissionTarget.VehicleID)
}

func RunLandGen(VehicleId int) {
	query := GenerateAuthQuery()
	droneDetails := GetDroneInfo(VehicleId, query)

	newMissionTarget := MissionTarget{}
	newMissionTarget.VehicleID = VehicleId
	newMissionTarget.Latitude = droneDetails.Vehicles[0].Telem.Latitude
	newMissionTarget.Longitude = droneDetails.Vehicles[0].Telem.Longitude
	newMissionTarget.Altitude = 25
	newMissionTarget.Heading = 0
	newMissionTarget.Velocity = 10
	runLandMission(newMissionTarget.VehicleID)
}

func runLandMission(vehicleID int) {

	fmt.Println("hello 2")
	query := GenerateAuthQuery()

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
