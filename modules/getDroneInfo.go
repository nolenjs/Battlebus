package modules

import (
	"fmt"
	"strings"

	icarus "git.ace/icarus/icarusclient/v5"
)

func GetDroneInfo(VehicleID int, query icarus.QueryPackage) icarus.IcarusResponse {
	query.ClearQueries()
	qID := query.GetVehicleStatus(VehicleID)
	respChannel, _ := query.Execute()
	responses := <-respChannel

	resp, ok := responses.Get(qID)

	if !ok {
		fmt.Println("[!] Invalid response")
		return icarus.IcarusResponse{}
	}

	return resp
}

func GetAllDroneInfo(query icarus.QueryPackage) icarus.IcarusResponse {
	query.ClearQueries()
	qID := query.GetAllVehicleStatus()
	respChannel, _ := query.Execute()
	responses := <-respChannel

	resp, ok := responses.Get(qID)

	if !ok {
		fmt.Println("[!] Invalid response")
		return icarus.IcarusResponse{}
	}
	return resp
}

func GetAllIff(query icarus.QueryPackage) map[int32]string {
	allDrones := GetAllDroneInfo(query).Vehicles
	var IFFs map[int32]string
	IFFs = make(map[int32]string)
	for _, drone := range allDrones {
		droneType := strings.Split(drone.VehicleCallsign, "-")[0]
		IFFs[int32(drone.VehicleId)] = droneType
	}
	return IFFs
}
