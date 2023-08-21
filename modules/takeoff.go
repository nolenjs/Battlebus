package modules

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	icarus "git.ace/icarus/icarusclient/v5"
)

func RunTakeoff() {
	var VehicleId int
	var Altitude float32

	newMissionTarget := MissionTarget{}
	fmt.Println("Insert Vehicle ID Number")
	fmt.Scanln(&VehicleId)
	newMissionTarget.VehicleID = VehicleId
	fmt.Println("Insert Takeoff Altitude")
	fmt.Scanln(&Altitude)
	newMissionTarget.Altitude = Altitude

	RunTakeoffMission(newMissionTarget.VehicleID, newMissionTarget)
}

func RunTakeoffGen(VehicleId int, Altitude float32) {
	newMissionTarget := MissionTarget{}
	newMissionTarget.VehicleID = VehicleId
	newMissionTarget.Altitude = Altitude
	RunTakeoffMission(newMissionTarget.VehicleID, newMissionTarget)
}

func RunTakeoffMission(vehicleID int, target MissionTarget) {
	logFile, err := os.OpenFile(strconv.Itoa(vehicleID)+"_log.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer logFile.Close()

	printWithTimestamp(logFile, strconv.Itoa(vehicleID), "Running: takeoff.go")

	//Create a new query pointed at the IcarusServer instance
	query := GenerateAuthQuery()

	respl := GetDroneInfo(vehicleID, query)
	nameString := respl.Vehicles[0].VehicleCallsign

	//Load 50 units of Fuel
	configs := icarus.AddPayloadConfig(nil, "Fuel", icarus.Fuel, 100, true)

	//Assumes ISR vehicle, may add other payloads here using above syntax if munitions are needed
	configSeq := query.ConfigurePayloads(vehicleID, configs)
	//Create the chanels to get execute responses

	responseChan, _ := query.Execute()
	fmt.Println("Loading fuel...")
	response := <-responseChan
	// vResp := query.GetVehicleStatus(vehicleID)
	// response.ShowResponse()

	configResponse, ok := response.Get(configSeq)
	if !ok {
		fmt.Println("Refuel response not found")
		printWithTimestamp(logFile, strconv.Itoa(vehicleID), "Refuel response not found")
		return
	}

	if configResponse.Ok {
		fmt.Println("Refueling complete")
		printWithTimestamp(logFile, strconv.Itoa(vehicleID), "Refueling complete")

	} else {
		fmt.Println("Error during refueling:", configResponse.Message)
		printWithTimestamp(logFile, strconv.Itoa(vehicleID), "Error during refueling:")
	}

	// load other payloads
	// fighter = 1
	// bomber = 2
	// isr = 4
	// multi = 6
	// wmd = 7

	// air-to-air
	if GetDroneInfo(vehicleID, query).Vehicles[0].VConfig.Role == icarus.Fighter || GetDroneInfo(vehicleID, query).Vehicles[0].VConfig.Role == icarus.Multi {
		fmt.Println("Loading missiles")
		configs = icarus.AddPayloadConfig(nil, "AntiMatterMissile", icarus.AntiMatterMissile, 4, true)
		//Assumes Anti-Matter missile payloads
		configSeq = query.ConfigurePayloads(vehicleID, configs)
		//Create the chanels to get execute responses

		responseChan, _ := query.Execute()
		fmt.Println("Loading missiles...")
		response := <-responseChan

		configResponse, ok := response.Get(configSeq)
		fmt.Println(configResponse.PayloadResponse.Type)
		if !ok {
			fmt.Println("Missile response not found")
			printWithTimestamp(logFile, strconv.Itoa(vehicleID), "Missile response not found")

			return
		}

		if configResponse.Ok {
			fmt.Println("Loading complete")
			printWithTimestamp(logFile, strconv.Itoa(vehicleID), "loading complete: missiles")
		} else {
			fmt.Println("Error during loading:", configResponse.Message)
			printWithTimestamp(logFile, strconv.Itoa(vehicleID), "Error during loading: missiles")
		}
	}

	// air-to-ground
	if GetDroneInfo(vehicleID, query).Vehicles[0].VConfig.Role == icarus.Bomber || GetDroneInfo(vehicleID, query).Vehicles[0].VConfig.Role == icarus.Multi {
		configs = icarus.AddPayloadConfig(nil, "ThermalLance", icarus.ThermalLance, 4, true)
		//Assumes Thermal Lance payloads
		configSeq = query.ConfigurePayloads(vehicleID, configs)
		//Create the chanels to get execute responses

		responseChan, _ := query.Execute()
		fmt.Println("Loading bombs...")
		response := <-responseChan

		configResponse, ok := response.Get(configSeq)
		if !ok {
			fmt.Println("Bomb response not found")
			printWithTimestamp(logFile, strconv.Itoa(vehicleID), "Bomb response not found")

			return
		}

		if configResponse.Ok {
			fmt.Println("Loading complete")
			printWithTimestamp(logFile, strconv.Itoa(vehicleID), "loading complete: bombs")

		} else {
			fmt.Println("Error during loading:", configResponse.Message)
			printWithTimestamp(logFile, strconv.Itoa(vehicleID), "Error during loading: bombs")

		}
	}

	// super air-to-ground
	if strings.Contains(nameString, "DRACONUS") {
		configs = icarus.AddPayloadConfig(nil, "Phosphex", icarus.Phosphex, 2, true)
		//Assumes Phosphex payloads
		configSeq = query.ConfigurePayloads(vehicleID, configs)
		//Create the chanels to get execute responses

		responseChan, _ := query.Execute()
		fmt.Println("Loading phosphex...")
		response := <-responseChan

		configResponse, ok := response.Get(configSeq)
		if !ok {
			fmt.Println("Phosphex response not found")
			printWithTimestamp(logFile, strconv.Itoa(vehicleID), "Phosphex response not found")

			return
		}

		if configResponse.Ok {
			fmt.Println("Loading complete")
			printWithTimestamp(logFile, strconv.Itoa(vehicleID), "Loading complete: phosphex (god help us all)")

		} else {
			fmt.Println("Error during loading:", configResponse.Message)
			printWithTimestamp(logFile, strconv.Itoa(vehicleID), "Phosphex response not found"+configResponse.Message)

		}
	}

	//Clear the refuel query from the queue
	query.ClearQueries()
	time.Sleep(5 * time.Second)
	//Change mode to Take off
	takeOffSeq := query.SetNavMode(vehicleID, icarus.TAKE_OFF)
	responseChan, _ = query.Execute()
	fmt.Println("Taking off...")
	printWithTimestamp(logFile, strconv.Itoa(vehicleID), "Taking off ...")
	response = <-responseChan
	takeOffResponse, ok := response.Get(takeOffSeq)
	if !ok {
		fmt.Println("TAKE_OFF response not found")
		printWithTimestamp(logFile, strconv.Itoa(vehicleID), "TAKE_OFF response not found")
		return
	}

	if takeOffResponse.Ok {
		fmt.Println("UAV launch complete")
	} else {
		fmt.Println("Error during take off:", takeOffResponse.Message)
		printWithTimestamp(logFile, strconv.Itoa(vehicleID), "Error during takeoff")
	}
	//Clear the take off query from the queue
	query.ClearQueries()

	//Sleep to ensure UAV is in air before navigating
	time.Sleep(1 * time.Second)

	//Change mode to Navigate
	navSeq := query.SetNavMode(vehicleID, icarus.NAVIGATION)

	// get the current telem of the drone to take off in the right spot
	droneTelem := respl.Vehicles[0].Telem
	target.Latitude = droneTelem.Latitude
	target.Longitude = droneTelem.Longitude

	responseChan, _ = query.Execute()
	fmt.Println("Entering Navigate mode...")
	response = <-responseChan
	navResponse, ok := response.Get(navSeq)
	if !ok {
		fmt.Println("NAVIGATE response not found")
		return
	}

	if navResponse.Ok {
		fmt.Println("UAV ready to navigate")
	} else {
		fmt.Println("Error during mode change:", navResponse.Message)
		printWithTimestamp(logFile, strconv.Itoa(vehicleID), "Error during mode change")

	}
	//Clear the mode change query from the queue
	query.ClearQueries()

	//Navigate to point of interest and linger for 60 seconds then navigate to secondary location
	//The vehicle will loop through the given targets until it is told to do something else
	cmdList := GetCmdList(target)

	gotoSeq := query.Goto(vehicleID, cmdList)
	responseChan, _ = query.Execute()
	fmt.Println("Navigating to waypoint...")

	response = <-responseChan
	gotoResponse, ok := response.Get(gotoSeq)
	if !ok {
		fmt.Println("Go to response not found")
		printWithTimestamp(logFile, strconv.Itoa(vehicleID), "Go to response not found")
		return
	}

	if gotoResponse.Ok {
		fmt.Println("Navigating to waypoint")

	} else {
		fmt.Println("Error during navigation:", gotoResponse.Message)
		printWithTimestamp(logFile, strconv.Itoa(vehicleID), "Error during navigation:"+gotoResponse.Message)

	}
	//Clear the navigation query from the queue
	query.ClearQueries()

	// //Sleep to ensure UAV is stopped before landing
	time.Sleep(1 * time.Second)

	// //Change mode to Navigate
	// navSeq = query.SetNavMode(vehicleID, icarus.LAND_WAYPOINT)

	// responseChan, _ = query.Execute()
	// fmt.Println("Entering Land_Waypoint mode...")
	// response = <-responseChan
	// navResponse, ok = response.Get(navSeq)
	// if !ok {
	// 	fmt.Println("LAND_WAYPOINT response not found")
	// 	return
	// }

	// if navResponse.Ok {
	// 	fmt.Println("UAV ready to land")
	// } else {
	// 	fmt.Println("Error during mode change:", navResponse.Message)
	// }
	// //Clear the mode change query from the queue
	// query.ClearQueries()
}
