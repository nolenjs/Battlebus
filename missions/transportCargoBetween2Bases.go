// package missions
package main

import (
	"fmt"
	icarus "git.ace/icarus/icarusclient/v5"
	modulesImport "liam-tool"
	"liam-tool/modules"
	"liam-tool/utils"
	"os"
	"os/exec"
)

func handleError() {
	if err := recover(); err != nil {
		fmt.Println("Error occurred:", err)
	}
}

// func transportCargoBetween2Bases() {
func notmain() {
	defer handleError()
	instructionGiven := false
	query := modulesImport.GenerateAuthQuery()
	var IFF int
	fmt.Print("Enter an integer for IFFID: ")
	_, err := fmt.Scan(&IFF)
	if err != nil {
		fmt.Println("Error reading IFFID:", err)
		return
	}
	cargo1 := IFF
	//cargo1Telem := modules.GetDroneInfo(cargo1, query).Vehicles[0].Telem
	//cargo1 := 2535
	//cargo2 := 2536
	//cargo_array := [...]int{2212, 2213}

	var gruntle_lat float64
	fmt.Print("Enter destination Latitude: ")
	_, err1 := fmt.Scan(&gruntle_lat)
	if err1 != nil {
		fmt.Println("Error reading Latitude:", err1)
		return
	}
	var gruntle_lon float64
	fmt.Print("Enter destination Latitude: ")
	_, err2 := fmt.Scan(&gruntle_lon)
	if err2 != nil {
		fmt.Println("Error reading Longitude:", err2)
		return
	}


	burbage_lat := 45.6515
	burbage_lon := -61.3699

	var maxVel float32 = 70
	var Velocity float32 = maxVel
	//cargo_remaining := 200

	//cargo1_alive := true
	//cargo2_alive := true

	cargo1State := 0
	//cargo2State := 0

	fuck := 0

	configs := icarus.AddPayloadConfig(nil, "Fuel", icarus.Fuel, 50, true)
	configSeq := query.ConfigurePayloads(cargo1, configs)
	query.ShowQuery()

	responseChan, _ := query.Execute()
	fmt.Println("Loading fuel...")
	response := <-responseChan
	configResponse, ok := response.Get(configSeq)
	if !ok {
		fmt.Println("Refuel response not found")
		return
	}
	if configResponse.Ok {
		fmt.Println("Refueling complete")
	} else {
		fmt.Println("Error during refueling:", configResponse.Message)
	}
	query.ClearQueries()

	for fuck < 1 {
		defer handleError()
		cargo1Telem := modules.GetDroneInfo(cargo1, query).Vehicles[0].Telem
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
		fmt.Printf("Segment %d\n", cargo1State)
		if cargo1State == 0 {
			// Take off
			if cargo1Telem.Altitude >= 100 {
				cargo1State = 1
				instructionGiven = false
			} else if !instructionGiven {
				distance := utils.Haversine(cargo1Telem.Latitude, cargo1Telem.Longitude, gruntle_lat, gruntle_lon)
				const distanceThreshold = 2000
				if distance < distanceThreshold {
					Velocity = 15
				}
				modulesImport.GoToLocGeneral(gruntle_lat, gruntle_lon, cargo1, Velocity, 100, query)
				instructionGiven = true
			}
		} else if cargo1State == 1 {
			// Fly to Gruntle
			distance := utils.Haversine(cargo1Telem.Latitude, cargo1Telem.Longitude, gruntle_lat, gruntle_lon)
			fmt.Printf("Distance Left: %f", distance)
			if distance < 50 {
				cargo1State = 2
				instructionGiven = false
			} else {

				const distanceThreshold = 7500
				if distance < distanceThreshold {
					var temp float32 = float32(distance)
					Velocity = maxVel * (temp / distanceThreshold)
				}
				if Velocity < 15 {
					Velocity = 15
				}
				modulesImport.GoToLocGeneral(gruntle_lat, gruntle_lon, cargo1, Velocity, 110, query)
			}
		} else if cargo1State == 2 {
			// Land at Gruntle
			if cargo1Telem.Altitude == 0 {
				// CHANGE ME
				cargo1State = 4
				// CHANGE ME
				instructionGiven = false
				modulesImport.GoToLocGeneral(gruntle_lat, gruntle_lon, cargo1, 0, 0, query)
				fuck++
			} else if !instructionGiven {
				fmt.Printf("Starting Segment %d\n", cargo1State)
				modulesImport.GoToLocGeneral(gruntle_lat, gruntle_lon, cargo1, 10, 0, query)
				instructionGiven = true
			}
		} else if cargo1State == 3 {
			// Unload
			/*
				if unloaded {
					fuck++
				} else {
					unload
				}
			*/
		} else if cargo1State == 4 {
			// Take off
			if cargo1Telem.Altitude == 100 {
				cargo1State = 5
				instructionGiven = false
			} else if !instructionGiven {
				distance := utils.Haversine(cargo1Telem.Latitude, cargo1Telem.Longitude, burbage_lat, burbage_lon)
				const distanceThreshold = 2000
				if distance < distanceThreshold {
					Velocity = 15
				}
				fmt.Printf("Starting Segment %d\n", cargo1State)
				modulesImport.GoToLocGeneral(burbage_lat, burbage_lon, cargo1, Velocity, 100, query)
				instructionGiven = true
			}
		} else if cargo1State == 5 {
			// Fly to Burbage
			distance := utils.Haversine(cargo1Telem.Latitude, cargo1Telem.Longitude, burbage_lat, burbage_lon)
			fmt.Printf("%fm from target\n", distance)
			if distance < 50 {
				cargo1State = 6
				instructionGiven = false
			} else if !instructionGiven {
				const distanceThreshold = 5000
				if distance < distanceThreshold {
					var temp float32 = float32(distance)
					Velocity = maxVel * (temp / distanceThreshold)
				}
				fmt.Printf("Starting Segment %d\n", cargo1State)
				modulesImport.GoToLocGeneral(burbage_lat, burbage_lon, cargo1, Velocity, 100, query)
				instructionGiven = true
			}
		} else if cargo1State == 6 {
			// Land at Burbage
			if cargo1Telem.Altitude == 0 {
				// CHANGE ME
				cargo1State = 0
				fuck++
				// CHANGE ME
				instructionGiven = false
			} else if !instructionGiven {
				fmt.Printf("Starting Segment %d\n", cargo1State)
				modulesImport.GoToLocGeneral(burbage_lat, burbage_lon, cargo1, 10, 0, query)
				instructionGiven = true
			}
		} else if cargo1State == 7 {
			// Restock
			/*
				if stocked {
					cargo1State = 0
				} else {
					stock
				}
			*/
		}
	}
}
