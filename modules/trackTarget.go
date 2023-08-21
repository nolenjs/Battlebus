// package missions
package modules

import (
	"fmt"
	//"liam-tool/modules"
	"liam-tool/utils"
)

const STOP_DISTANCE = 50
const PAYLOAD_RANGE = 300
const SLOWDOWN_DISTANCE = 1000

/*func InputTrackTarget() {
	var trackerID, targetID int
	var velocity, altitude float32
	fmt.Println("Enter the Tracker ID:")
	fmt.Scanln(&trackerID)
	fmt.Println("Enter the Target ID:")
	fmt.Scanln(&targetID)
	fmt.Println("Enter the Velocity:")
	fmt.Scanln(&velocity)
	fmt.Println("Enter the Altitude:")
	fmt.Scanln(&altitude)

	TrackTarget(rackerID, targetID, velocity, altitude)
}*/

func TrackTarget(TrackerIDMap map[int]int, TargetID int, MaxVel float32, Altitude float32) {
	query := GenerateAuthQuery()
	var temp []int32
	for iff,_ := range TrackerIDMap {
		temp = append(temp, int32(iff))
	}
	var oldTargetLat, oldTargetLong float64

	currVel := MaxVel
	// currVel := float32(160)
	for {
		for drone, _ := range TrackerIDMap {
			// get the radar pings and location of target drone
			pings, _, _, _ := RetrieveSomeRADAR(temp, query)

			// if there are zero pings do not try to query the radar
			if len(pings) <= 0 {
				fmt.Println("No pings from radar")
				continue
			}

			targetPing := pings[int32(TargetID)]
			targetLat, targetLong := targetPing.Latitude, targetPing.Longitude

			if targetLat != oldTargetLat || targetLong != oldTargetLong {
				fmt.Printf("Going to %f, %f\n", targetLat, targetLong)
			}
			oldTargetLat, oldTargetLong = targetLat, targetLong

			// send the bomber and ISR to the target location
			GoToLocGeneral(targetLat, targetLong, drone, currVel, Altitude, query)

			// get the ISR telem and check if it can take an image or if it should slow down
			trackerInfo := GetDroneInfo(drone, query)
			trackerTelem := trackerInfo.Vehicles[0].Telem
			trackerDist := utils.Haversine(targetLat, targetLong, trackerTelem.Latitude, trackerTelem.Longitude)
			if trackerDist <= STOP_DISTANCE {
				// fmt.Println("Tracker stopped")
				currVel = 0
			} else if trackerDist <= PAYLOAD_RANGE {
				fmt.Println("Tracker in camera range")
				currVel = 30
			} else if trackerDist <= SLOWDOWN_DISTANCE {
				// fmt.Println("Tracker slowing down")
				// fmt.Printf("Going to targetLat: %f, targetLong: %f\n", targetLat, targetLong)
				currVel = 70
			} else {
				currVel = MaxVel
			}
		}
	}
}
