package main

import (
	"fmt"
	icarus "git.ace/icarus/icarusclient/v5"
	modulesImport "liam-tool/modules"
)

func UnloadCargo(vehicleID int) {
	/*
		const (
			InvalidPayload      PayloadType = 0
			AllPayloads                     = 0
			ThermalLance                    = 3
			Camera                          = 4
			Fuel                            = 5
			Phosphex                        = 7
			PhosphexRemediation             = 8
			AirRadar                        = 9
			AntiMatterMissile               = 10
			AllRadar                        = 11
			GroundRadar                     = 12
			SAM                             = 13
			Cargo                           = 14
			SeekerMissile                   = 15
			Supplies                        = 16
			G2ARadar                        = 18
		)
	*/

	query := modulesImport.GenerateAuthQuery()
	configs := icarus.AddPayloadConfig(nil, "Fuel", icarus.Fuel, 100, true)
	configSeq := query.ConfigurePayloads(vehicleID, configs)
	//query.ShowQuery()

	responseChan, _ := query.Execute()
	fmt.Println("Loading fuel...")
	response := <-responseChan
	configResponse, ok := response.Get(configSeq)
	if !ok {
		fmt.Println("Refuel response not found")
		return
	}
	fmt.Println(configResponse)
	if configResponse.Ok {
		fmt.Println("Refueling complete")
	} else {
		fmt.Println("Error during refueling:", configResponse.Message)
	}
	query.ClearQueries()
	fuck, _ := icarus.UnloadCargo(icarus.Fuel, 1)
	fmt.Println(fuck)
	/*
		configs2 := icarus.AddPayloadConfig(nil, "payload1", 16, 1, true)
		configSeq = query.ConfigurePayloads(vehicleID, configs2)
		query.ShowQuery()
	*/

	executeSeq := query.ExecutePayload(vehicleID, icarus.Cargo, 1, fuck, 0)
	fmt.Println(executeSeq)
	respChan, _ := query.Execute()
	fmt.Println("Unloading Supplies...")
	resp := <-respChan
	fmt.Println(resp)
	executeResponse, ok := resp.Get(configSeq)
	if !ok {
		fmt.Println("Supply response not found")
	}
	if executeResponse.Ok {
		fmt.Println("Loading complete")
	} else {
		fmt.Println("Error during loading:", executeResponse.Message)
	}
	query.ClearQueries()
	fmt.Println(modulesImport.GetDroneInfo(vehicleID, query).Vehicles[0].PayStatus[14])
}

func main() {
	cargo := []int{2540, 2537}
	for _, iff := range cargo {
		UnloadCargo(iff)
	}
}
