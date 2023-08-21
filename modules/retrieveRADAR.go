package modules

import (
	"bufio"
	icarusClient "git.ace/icarus/icarusclient/v5"
	"os"
	"strconv"
)

func GetFriendlyDrones() map[int]string {
	iffMap := make(map[int]string)
	file, err := os.Open("../resources/iffsToIgnore.txt")
	if err != nil {
		return nil
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		iff, err := strconv.Atoi(line)
		if err != nil {
			return nil
		}
		iffMap[iff] = "empty"
	}
	if err := scanner.Err(); err != nil {
		return nil
	}
	return iffMap
}

// input: array of IFF IDs
// output: map (IFF IDs, RADAR returns)
func RetrieveSomeRADAR(VehicleIDArray []int32, query icarusClient.QueryPackage) (map[int32]icarusClient.RadarPing, map[int32]icarusClient.RadarPing, map[int32]icarusClient.RadarPing, bool) {
	var payload map[int32]icarusClient.RadarPing
	payload = make(map[int32]icarusClient.RadarPing)
	var payloadAir map[int32]icarusClient.RadarPing
	payloadAir = make(map[int32]icarusClient.RadarPing)
	var payloadGnd map[int32]icarusClient.RadarPing
	payloadAir = make(map[int32]icarusClient.RadarPing)
	missile := false
	malDrones := GetFriendlyDrones()
	for _, VehicleID := range VehicleIDArray {
		drone := GetDroneInfo(int(VehicleID), query)
		if drone.Ok {
			var temp map[int32]icarusClient.RadarPing
			temp = make(map[int32]icarusClient.RadarPing)
			if drone.Vehicles[0].VConfig.Role == icarusClient.ISR {
				temp = drone.Vehicles[0].PayStatus[icarusClient.AllRadar].Radar
			} else if drone.Vehicles[0].VConfig.Role == icarusClient.Bomber || drone.Vehicles[0].VConfig.Role == icarusClient.Multi {
				temp = drone.Vehicles[0].PayStatus[icarusClient.GroundRadar].Radar
			} else if drone.Vehicles[0].VConfig.Role == icarusClient.Fighter || drone.Vehicles[0].VConfig.Role == icarusClient.Multi {
				temp = drone.Vehicles[0].PayStatus[icarusClient.AirRadar].Radar
			} else if drone.Vehicles[0].VConfig.Role == icarusClient.RADAR {
				temp = drone.Vehicles[0].PayStatus[icarusClient.G2ARadar].Radar
			}
			if len(temp) > 0 {
				for key, item := range temp {
					_, ok := malDrones[int(key)]
					if !ok {
						payload[key] = item
						if item.Altitude > 2 {
							payloadAir[key] = item
							if item.Type == icarusClient.Missile {
								missile = true
							}
						} else {
							payloadAir[key] = item
						}
					}
				}
			}
		}
	}
	return payload, payloadAir, payloadGnd, missile
}

func RetrieveAllRADAR(query icarusClient.QueryPackage) (map[int32]icarusClient.RadarPing, map[int32]icarusClient.RadarPing, map[int32]icarusClient.RadarPing, bool) {
	var payload map[int32]icarusClient.RadarPing
	payload = make(map[int32]icarusClient.RadarPing)
	var payloadAir map[int32]icarusClient.RadarPing
	payloadAir = make(map[int32]icarusClient.RadarPing)
	var payloadGnd map[int32]icarusClient.RadarPing
	payloadAir = make(map[int32]icarusClient.RadarPing)
	missile := false
	malDrones := GetFriendlyDrones()
	for VehicleID := range malDrones {
		drone := GetDroneInfo(int(VehicleID), query)
		var temp map[int32]icarusClient.RadarPing
		temp = make(map[int32]icarusClient.RadarPing)
		gndRADARCheck := drone.Vehicles[0].PayStatus[icarusClient.GroundRadar].Id
		allRADARCheck := drone.Vehicles[0].PayStatus[icarusClient.AllRadar].Id
		airRADARCheck := drone.Vehicles[0].PayStatus[icarusClient.AirRadar].Id
		g2aRADARCheck := drone.Vehicles[0].PayStatus[icarusClient.G2ARadar].Id
		if allRADARCheck == icarusClient.AllRadar {
			temp = drone.Vehicles[0].PayStatus[icarusClient.AllRadar].Radar
		} else if gndRADARCheck == icarusClient.GroundRadar {
			temp = drone.Vehicles[0].PayStatus[icarusClient.GroundRadar].Radar
		} else if airRADARCheck == icarusClient.AirRadar {
			temp = drone.Vehicles[0].PayStatus[icarusClient.AirRadar].Radar
		} else if g2aRADARCheck == icarusClient.G2ARadar {
			temp = drone.Vehicles[0].PayStatus[icarusClient.G2ARadar].Radar
		}
		if len(temp) > 0 {
			for key, item := range temp {
				_, ok := malDrones[int(key)]
				if !ok {
					payload[key] = item
					if item.Altitude > 2 {
						payloadAir[key] = item
						if item.Type == icarusClient.Missile {
							missile = true
						}
					} else {
						payloadAir[key] = item
					}
				}
			}
		}
	}
	return payload, payloadAir, payloadGnd, missile
}
