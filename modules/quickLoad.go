package modules

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	icarus "git.ace/icarus/icarusclient/v5"
)

func RunLoadGen(VehicleId int) {
	newMissionTarget := MissionTarget{}
	newMissionTarget.VehicleID = VehicleId
	RunLoadMission(newMissionTarget.VehicleID, newMissionTarget)
}

func RunLoadMission(vehicleID int, target MissionTarget) {
logFile, err := os.OpenFile(strconv.Itoa(vehicleID)+"_log.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer logFile.Close()
	printWithTimestamp(logFile, strconv.Itoa(vehicleID), "Running: quickLoad.go")

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
	if (GetDroneInfo(vehicleID, query).Vehicles[0].VConfig.Role == icarus.Fighter || GetDroneInfo(vehicleID, query).Vehicles[0].VConfig.Role == icarus.Multi) {
		configs = icarus.AddPayloadConfig(nil, "AntiMatterMissile", icarus.AntiMatterMissile, 4, true)
		//Assumes Anti-Matter missile payloads
		configSeq = query.ConfigurePayloads(vehicleID, configs)
		//Create the chanels to get execute responses

		responseChan, _ := query.Execute()
		fmt.Println("Loading missiles...")
		response := <-responseChan

		configResponse, ok := response.Get(configSeq)
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
	if (GetDroneInfo(vehicleID, query).Vehicles[0].VConfig.Role == icarus.Bomber || GetDroneInfo(vehicleID, query).Vehicles[0].VConfig.Role == icarus.Multi) {
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
}