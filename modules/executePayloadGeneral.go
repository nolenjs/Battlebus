package modules

import (
	"fmt"

	icarusClient "git.ace/icarus/icarusclient/v5"
)

func ExecutePayload(ourIFFID uint, targetIFFID uint, PayloadID int) {
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
	query := GenerateAuthQuery()

	if PayloadID == 3 {

		//ThermalLance
		executeSeq := query.ExecutePayload(int(ourIFFID), icarusClient.ThermalLance, 1, icarusClient.EmptyParams(), 0)

		responseChan, _ := query.Execute() // Pew pew
		fmt.Println("Waiting for responses:")
		response := <-responseChan
		executeResponse, ok := response.Get(executeSeq)
		if !ok {
			fmt.Println("ThermalLance execute response not found")
		}
		fmt.Println(executeResponse)

	} else if PayloadID == 4 {

		//Camera
		executeSeq := query.ExecutePayload(int(ourIFFID), icarusClient.Camera, 1, icarusClient.EmptyParams(), 0)

		responseChan, _ := query.Execute() // Pew pew
		fmt.Println("Waiting for responses:")
		response := <-responseChan
		executeResponse, ok := response.Get(executeSeq)
		if !ok {
			fmt.Println("Camera execute response not found")
		}
		fmt.Println(executeResponse)

	} else if PayloadID == 10 {

		//AntiMatterMissile
		executeSeq := query.ExecutePayload(int(ourIFFID), icarusClient.AntiMatterMissile, 1, icarusClient.EmptyParams(), int(targetIFFID))

		responseChan, _ := query.Execute() // Pew pew
		fmt.Println("Waiting for responses:")
		response := <-responseChan
		executeResponse, ok := response.Get(executeSeq)
		if !ok {
			fmt.Println("AntiMatterMissile execute response not found")
		}
		fmt.Println(executeResponse)

	} else if PayloadID == 15 {

		//SeekerMissile
		executeSeq := query.ExecutePayload(int(ourIFFID), icarusClient.SeekerMissile, 1, icarusClient.EmptyParams(), int(targetIFFID))

		responseChan, _ := query.Execute() // Pew pew
		fmt.Println("Waiting for responses:")
		response := <-responseChan
		executeResponse, ok := response.Get(executeSeq)
		if !ok {
			fmt.Println("SeekerMissile execute response not found")
		}
		fmt.Println(executeResponse)

	} else if PayloadID == 16 {

		//SeekerMissile
		executeSeq := query.ExecutePayload(int(ourIFFID), icarusClient.Supplies, 50, icarusClient.EmptyParams(), int(ourIFFID))

		responseChan, _ := query.Execute() // Pew pew
		fmt.Println("Waiting for responses:")
		response := <-responseChan
		executeResponse, ok := response.Get(executeSeq)
		if !ok {
			fmt.Println("SeekerMissile execute response not found")
		}
		fmt.Println(executeResponse)
	
	} else {
		fmt.Println("Error: Invalid Payload Number")
	}
	
}
