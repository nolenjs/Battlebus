package main

import (
	"os"
	"strconv"
	"strings"
	"fmt"
	"sync"
	"liam-tool/modules"
)

const ERROR_GOTO = "usage: goto <float lat> <float long> <float alt> <float velocity>"

func gotogroup() {
	var offsetMap = make(map[int][]float64)
	offsetMap[0] = []float64{0,0}
	offsetMap[1] = []float64{-0.0100, 0.0100}
	offsetMap[2] = []float64{-0.0100, -0.0100}
	offsetMap[3] = []float64{0.0100, 0.0100}
	offsetMap[4] = []float64{0.0100, -0.0100}
	offsetMap[5] = []float64{-0.0200, 0}
	offsetMap[6] = []float64{0.0200, 0}
	offsetMap[7] = []float64{0, 0.0200}
	offsetMap[8] = []float64{0, -0.0200}

	args := os.Args[1:]
	lat, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		fmt.Println(ERROR_GOTO)
		return
	}

	long, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		fmt.Println(ERROR_GOTO)
		return
	}

	alt, err := strconv.ParseFloat(args[2], 32)
	if err != nil {
		fmt.Println(ERROR_GOTO)
		return
	}

	velo, err := strconv.ParseFloat(args[3], 32)
	if err != nil {
		fmt.Println(ERROR_GOTO)
		return
	}

	var dronesUnparsed = args[4]
	splitByComma := strings.Split(dronesUnparsed, ",")
	parsedDrones := make([]int, 0, len(splitByComma))
	fmt.Println(splitByComma)

	for i, s := range splitByComma {
		if i == len(splitByComma)-1{
			break
		}
		trimmedStr := strings.Trim(s, ",")
		fmt.Println(len(trimmedStr))
		droneid, err := strconv.Atoi(trimmedStr)
		if err != nil {
			fmt.Println("Error converting strings to ints")
			return
		}
		parsedDrones = append(parsedDrones, droneid)
	}
	fmt.Println(parsedDrones)

	var wg sync.WaitGroup
	var counter = 0
	for _, DroneID := range parsedDrones {
		query := modules.GenerateAuthQuery()
		if modules.GetDroneInfo(DroneID, query).Vehicles[0].Telem.Altitude == 0 {
			fmt.Println("Takeoff the damn drone first.")
			break
		} else {
			index := counter % 9
			offsets := offsetMap[index]
			/*_, ok := droneList[DroneID]
			if !ok {
				fmt.Printf("[!] IFFID %d not in used drones...", DroneID)
				return
			}*/
			wg.Add(1)
			go func(wg *sync.WaitGroup, DroneID int, alt float64) {
				defer wg.Done()
				modules.SlowApproachGeneral(DroneID, lat + offsets[0], long + offsets[1], float32(alt), float32(velo))
			}(&wg, DroneID, alt)
			alt += 5
			counter +=1
		}
		
	}
	wg.Wait()
}