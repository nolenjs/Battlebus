package main

import (
	"fmt"
	icarusClient "git.ace/icarus/icarusclient/v5"

	modulesImport "liam-tool/modules"
)

func main() {
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
	var IFF int
	fmt.Print("Enter an integer for IFFID: ")
	_, err := fmt.Scan(&IFF)
	if err != nil {
		fmt.Println("Error reading IFFID:", err)
		return
	}
	var enemy int
	fmt.Print("Enter an integer for enemy IFFID: ")
	_, err = fmt.Scan(&enemy)
	if err != nil {
		fmt.Println("Error reading IFFID:", err)
		return
	}
	query := modulesImport.GenerateAuthQuery()
	//executeSeq := query.ExecutePayload(IFF, icarus.Cargo, 1, oof, 0)
	executeSeq := query.ExecutePayload(IFF, icarusClient.AntiMatterMissile, 1, icarusClient.EmptyParams(), enemy)
	responseChan, _ := query.Execute() // Pew pew
	fmt.Println("Waiting for responses:")
	response := <-responseChan
	executeResponse, ok := response.Get(executeSeq)
	if !ok {
		fmt.Println("Payload execute response not found")
	}
	fmt.Println(executeResponse)
}
