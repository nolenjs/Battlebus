package modules

import (
	"fmt"
	icarusClient "git.ace/icarus/icarusclient/v5"
	"os"
	"os/exec"
	"sort"
)

func SqnLdrDroneTable(drones map[int32]string, query icarusClient.QueryPackage) {
	var allTelem map[int]icarusClient.Telemetry
	allTelem = make(map[int]icarusClient.Telemetry)
	allDrones := GetAllDroneInfo(query)

	var droneDetails map[int]map[int]icarusClient.PayloadStatus
	droneDetails = make(map[int]map[int]icarusClient.PayloadStatus)
	var droneNames map[int]string
	droneNames = make(map[int]string)
	var aliveDrones map[int]bool
	aliveDrones = make(map[int]bool)
	for _, drone := range allDrones.Vehicles {
		_, ok := drones[int32(drone.VehicleId)]
		if ok {
			allTelem[int(drone.VehicleId)] = drone.Telem
			droneDetails[int(drone.VehicleId)] = drone.PayStatus
			droneNames[int(drone.VehicleId)] = drone.VehicleCallsign
			aliveDrones[int(drone.VehicleId)] = drone.Available
		}
	}
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	// Sort the keys of the drones map
	var sortedKeys []int
	for id := range drones {
		sortedKeys = append(sortedKeys, int(id))
	}
	sort.Ints(sortedKeys)

	// Print the sorted drones
	fmt.Printf("IFF\tName\t\tLat\t\tLon\t\tAlt\t\tSpeed\t\tFuel\tMssl\tBombs\n")
	for _, id := range sortedKeys {
		curTelem := allTelem[id]
		curName := droneNames[id]
		if curTelem.Latitude > 0 && curTelem.Longitude < 0 && aliveDrones[id] {
			fmt.Printf("%d\t%s\t%f\t%f\t%f\t%f\t%d", id, curName, curTelem.Latitude, curTelem.Longitude, curTelem.Altitude, curTelem.Velocity, droneDetails[id][icarusClient.Fuel].Resources)
			if droneDetails[id][icarusClient.AntiMatterMissile].Id == icarusClient.InvalidPayload {
				fmt.Printf("\t-")
			} else {
				fmt.Printf("\t%d", droneDetails[id][icarusClient.AntiMatterMissile].Resources)
			}
			if droneDetails[id][icarusClient.ThermalLance].Id == icarusClient.InvalidPayload {
				fmt.Printf("\t-")
			} else {
				fmt.Printf("\t%d", droneDetails[id][icarusClient.ThermalLance].Resources)
			}
			fmt.Printf("\n")
		} else {
			fmt.Printf("%d\t%s\tdead\t\tdead\t\tdead\t\tdead\t\tdead\tdead\tdead\n", id, curName)
		}
	}
}
