package main

import (
	"bufio"
	"fmt"
	icarusClient "git.ace/icarus/icarusclient/v5"
	modulesImport "liam-tool/modules"
	"liam-tool/utils"
	"os"
	"strconv"
	"time"
)

// var drones = list.New()
var drones map[int32]string
var fighters map[int32]fighter
var isrs map[int32]isr
var allPings map[int32]icarusClient.RadarPing
var airPings map[int32]icarusClient.RadarPing
var groundPings map[int32]icarusClient.RadarPing
var latA, lonA, latB, lonB float64
var weaponsFree bool

type fighter struct {
	IFF       int32
	telem     icarusClient.Telemetry
	state     droneState
	status    icarusClient.VehicleStatus
	toLat     [2]float64
	toLon     [2]float64
	altitude  float32
	maxVel    float32
	locIndex  int
	curTarget int32
	action    string
}

type isr struct {
	IFF      int32
	telem    icarusClient.Telemetry
	state    droneState
	status   icarusClient.VehicleStatus
	toLat    [2]float64
	toLon    [2]float64
	altitude float32
	maxVel   float32
	locIndex int
	closest  int32
	action   string
}

type droneState int

const (
	Init            droneState = -1
	LandedNoFuel               = 0
	LandedAndFueled            = 1
	EnRoute                    = 2
	Patrol                     = 3
	Intercept                  = 4
	FiredRetreat               = 5
	Empty                      = 6
	RTB                        = 7
	Landing                    = 8
	Done                       = 9
)

func readIFFsFromFile(filename string, query icarusClient.QueryPackage) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		iff, err := strconv.Atoi(line)
		if err != nil {
			return fmt.Errorf("invalid IFF in file: %s", err)
		}
		temp := modulesImport.GetDroneInfo(iff, query)
		if temp.Ok {
			if temp.Vehicles[0].Available {
				drones[int32(iff)] = "Waiting"
			} else {
				fmt.Printf("%d is dead\n", iff)
			}
		} else {
			fmt.Printf("%d does not exist\n", iff)
		}
	}
	return scanner.Err()
}

func readCoordinatesFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	//latA, lonA, latB, lonB float64
	var linesRead int
	for scanner.Scan() {
		line := scanner.Text()
		coord, err := strconv.ParseFloat(line, 64)
		if err != nil {
			return fmt.Errorf("invalid coordinate in file: %s", err)
		}
		switch linesRead {
		case 0:
			latA = coord
		case 1:
			lonA = coord
		case 2:
			latB = coord
		case 3:
			lonB = coord
		default:
			return fmt.Errorf("too many lines in the file")
		}

		linesRead++
	}

	if linesRead < 4 {
		return fmt.Errorf("not enough lines in the file")
	}
	return nil
}

func setup(query icarusClient.QueryPackage) {
	allPings = make(map[int32]icarusClient.RadarPing)
	airPings = make(map[int32]icarusClient.RadarPing)
	groundPings = make(map[int32]icarusClient.RadarPing)
	weaponsFree = false // Do we kill
	drones = make(map[int32]string)
	fighters = make(map[int32]fighter)
	isrs = make(map[int32]isr)
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter the filename containing IFF values: ")
	scanner.Scan()
	iffFilename := scanner.Text()

	fmt.Print("Enter the filename containing coordinate values: ")
	scanner.Scan()
	coordFilename := scanner.Text()

	// Read drone IFFs from the file
	if err := readIFFsFromFile(iffFilename, query); err != nil {
		fmt.Println("Error reading IFF values from file:", err)
		return
	}

	// Read coordinates from the file
	if err := readCoordinatesFromFile(coordFilename); err != nil {
		fmt.Println("Error reading coordinates from file:", err)
		return
	}
	ftemp, itemp := 0, 0
	for drone := range drones {
		temp := modulesImport.GetDroneInfo(int(drone), query).Vehicles[0]
		if temp.VConfig.Role == icarusClient.Fighter {
			var tempF = fighter{
				IFF:       drone,
				telem:     temp.Telem,
				state:     Init,
				status:    temp,
				toLat:     [2]float64{latA, latB},
				toLon:     [2]float64{lonA, lonB},
				altitude:  float32(100 + (10 * ftemp)),
				locIndex:  0,
				curTarget: 0,
				maxVel:    190,
				action:    "init",
			}
			fighters[drone] = tempF
			ftemp = 1 - ftemp
		} else if temp.VConfig.Role == icarusClient.ISR {
			var tempI = isr{
				IFF:      drone,
				telem:    temp.Telem,
				state:    Init,
				status:   temp,
				toLat:    [2]float64{latA, latB},
				toLon:    [2]float64{lonA, lonB},
				altitude: float32(300 + (10 * itemp)),
				maxVel:   160,
				locIndex: 0,
				closest:  0,
				action:   "init",
			}
			isrs[drone] = tempI
			itemp = 1 - itemp
		}
	}
	temp := 0
	latDif := latB - latA
	lonDif := lonB - lonA
	var count float64 = 0
	for iff, tempDrone := range fighters {
		drone := &tempDrone
		if drone.telem.Altitude == 0 {
			if drone.status.PayStatus[5].Resources < 100 {
				drone.state = LandedNoFuel
			} else {
				drone.state = LandedAndFueled
			}
		} else {
			drone.state = Patrol
		}
		if temp == 0 {
			drone.altitude = 110
			fmt.Println("Alt set to 110")
		} else {
			drone.altitude = 100
			fmt.Println("Alt set to 100")
		}
		toLatA := latA + (count * (latDif / float64(len(fighters))))
		toLatB := toLatA + latDif/float64(len(fighters))
		drone.toLat = [2]float64{toLatA, toLatB}
		toLonA := lonA + (count * (lonDif / float64(len(fighters))))
		toLonB := toLonA + lonDif/float64(len(fighters))
		drone.toLon = [2]float64{toLonA, toLonB}
		fighters[iff] = *drone
		temp = 1 - temp
		count++
	}
	fmt.Println("After fighter init")
	count = 0
	for iff, tempDrone := range isrs {
		drone := &tempDrone
		if drone.telem.Altitude == 0 {
			if drone.status.PayStatus[5].Resources < 100 {
				drone.state = LandedNoFuel
			} else {
				drone.state = LandedAndFueled
			}
		} else {
			drone.state = Patrol
		}

		if temp == 0 {
			drone.altitude = 310
			fmt.Println("Alt set to 310")
		} else {
			drone.altitude = 300
			fmt.Println("Alt set to 300")
		}
		toLatA := latA + (count * (latDif / float64(len(isrs))))
		toLatB := toLatA + latDif/float64(len(isrs))
		drone.toLat = [2]float64{toLatA, toLatB}
		toLonA := lonA + (count * (lonDif / float64(len(isrs))))
		toLonB := toLonA + lonDif/float64(len(isrs))
		drone.toLon = [2]float64{toLonA, toLonB}
		isrs[iff] = *drone
		temp = 1 - temp
		count++
	}
	fmt.Println("After ISR init")
}

// Just a function to store all means of getting knowledge for the decision
// A mixture of intel from drones, old intel, and any external files
func gatherIntel(query icarusClient.QueryPackage) {
	var temp []int
	for key := range drones {
		temp = append(temp, int(key))
	}
	var tempDroneArray []int32
	for key, _ := range drones {
		tempDroneArray = append(tempDroneArray, key)
	}
	allPings, airPings, groundPings, weaponsFree = modulesImport.RetrieveSomeRADAR(tempDroneArray, query)
}

func act(query icarusClient.QueryPackage) {
	// list potential actions here

	// takeoff
	for iff, drone := range fighters {
		tempDrone := &drone
		tempDrone.telem = modulesImport.GetDroneInfo(int(iff), query).Vehicles[0].Telem
		fighters[iff] = *tempDrone
		// patrol
		fmt.Printf("%d\n", drone.state)
		if drone.state == LandedNoFuel || drone.state == LandedAndFueled {
			modulesImport.RunTakeoffGen(int(iff), 10)
			temp := &drone
			temp.curTarget = 0
			temp.state = Patrol
			drone = *temp
		}
		if drone.state == Patrol {
			fmt.Printf("%d patrol\n", iff)
			goToLocSlow(drone.telem, drone.toLat[drone.locIndex], drone.toLon[drone.locIndex], int(iff), drone.maxVel, drone.altitude, query)
			tempIndex := checkMoveNext(drone.telem.Latitude, drone.telem.Longitude, drone.toLat[drone.locIndex], drone.toLon[drone.locIndex], drone.locIndex, len(drone.toLat))
			tempDrone.locIndex = tempIndex
			fighters[iff] = *tempDrone
		} else if drone.state == Intercept {
			tgtTemp, exists := airPings[drone.curTarget]
			if !exists {
				temp := &drone
				temp.curTarget = 0
				temp.state = Patrol
				drone = *temp
			}
			distance := utils.Haversine(drone.telem.Latitude, drone.telem.Longitude, tgtTemp.Latitude, tgtTemp.Longitude)
			if distance <= 5000 {
				if askUserToFire() {
					modulesImport.ExecutePayload(uint(drone.IFF), uint(drone.curTarget), icarusClient.AntiMatterMissile)
				}
				if modulesImport.GetDroneInfo(int(drone.IFF), query).Vehicles[0].PayStatus[icarusClient.AntiMatterMissile].Resources == 0 {
					temp := &drone
					temp.curTarget = 0
					temp.state = Empty
					drone = *temp
				}
				temp := &drone
				temp.state = FiredRetreat
				drone = *temp
			}
			tgtLat := ((distance / 5000) - 1) * (tgtTemp.Latitude - drone.telem.Latitude)
			tgtLon := ((distance / 5000) - 1) * (tgtTemp.Longitude - drone.telem.Longitude)
			goToLocSlow(drone.telem, tgtLat+drone.telem.Latitude, tgtLon+drone.telem.Longitude, int(drone.IFF), drone.maxVel, drone.altitude, query)
		} else if drone.state == FiredRetreat {
			goToLocSlow(drone.telem, drone.toLat[drone.locIndex], drone.toLon[drone.locIndex], int(iff), drone.maxVel, drone.altitude, query)
			tempIndex := checkMoveNext(drone.telem.Latitude, drone.telem.Longitude, drone.toLat[drone.locIndex], drone.toLon[drone.locIndex], drone.locIndex, len(drone.toLat))
			tempDrone.locIndex = tempIndex
			fighters[iff] = *tempDrone
			_, exists := airPings[drone.curTarget]
			if !exists {
				temp := &drone
				temp.curTarget = 0
				temp.state = Patrol
				drone = *temp
			}
		}
	}

	// shoot
	// retreat
	// rtb

	for iff, drone := range isrs {
		tempDrone := &drone
		tempDrone.telem = modulesImport.GetDroneInfo(int(iff), query).Vehicles[0].Telem
		isrs[iff] = *tempDrone
		if drone.state == LandedNoFuel || drone.state == LandedAndFueled {
			modulesImport.RunTakeoffGen(int(iff), 10)
			temp := &drone
			temp.state = Patrol
			drone = *temp
		}
		// patrol
		if drone.state == Patrol {
			fmt.Printf("%d patrol\n", iff)
			goToLocSlow(drone.telem, drone.toLat[drone.locIndex], drone.toLon[drone.locIndex], int(iff), drone.maxVel, drone.altitude, query)
			tempIndex := checkMoveNext(drone.telem.Latitude, drone.telem.Longitude, drone.toLat[drone.locIndex], drone.toLon[drone.locIndex], drone.locIndex, len(drone.toLat))
			tempDrone.locIndex = tempIndex
			isrs[iff] = *tempDrone
		}
	}
}

func goToLocSlow(telem icarusClient.Telemetry, Latitude, Longitude float64, VehicleId int, maxVel, Altitude float32, query icarusClient.QueryPackage) {
	curLat := telem.Latitude
	curLon := telem.Longitude
	var Velocity = maxVel

	const distanceThreshold = 5000
	distance := utils.Haversine(curLat, curLon, Latitude, Longitude)
	if distance < distanceThreshold {
		temp := float32(distance)
		Velocity = maxVel * (temp / distanceThreshold)
	}
	if Velocity < 15 {
		Velocity = 15
	}

	newMissionTarget := modulesImport.MissionTarget{}
	newMissionTarget.VehicleID = VehicleId
	newMissionTarget.Latitude = Latitude
	newMissionTarget.Longitude = Longitude
	newMissionTarget.Altitude = Altitude
	newMissionTarget.Heading = 0
	newMissionTarget.Velocity = Velocity
	navSeq := query.SetNavMode(VehicleId, icarusClient.NAVIGATION)
	responseChan, _ := query.Execute()
	response := <-responseChan
	navResponse, ok := response.Get(navSeq)

	if !ok {
		fmt.Println("NAVIGATE response not found")
		return
	}
	if navResponse.Ok {
	} else {
		fmt.Println("Error during mode change:", navResponse.Message)
	}
	//Clear the mode change query from the queue
	query.ClearQueries()
	cmdList := modulesImport.GetCmdList(newMissionTarget)
	gotoSeq := query.Goto(VehicleId, cmdList)
	responseChan, _ = query.Execute()
	response = <-responseChan
	gotoResponse, ok := response.Get(gotoSeq)
	if !ok {
		fmt.Println("Go to response not found")
		return
	}

	if gotoResponse.Ok {
	} else {
		fmt.Println("Error during navigation:", gotoResponse.Message)
		fmt.Printf("===%d - %f - %f===\n", VehicleId, Latitude, Longitude)
	}
}

// If we are close enough to our target, return the next index, otherwise return the original index
func checkMoveNext(curLat, curLon, Latitude, Longitude float64, index, arrayLength int) int {
	fmt.Printf("%f | %f | %f | %f\n", curLat, curLon, Latitude, Longitude)
	distance := utils.Haversine(curLat, curLon, Latitude, Longitude)
	fmt.Printf("Distance: %f | index: %d | Length: %d\n", distance, index, arrayLength)
	if distance <= 1000 {
		index++
		if index == arrayLength {
			return 0
		}
	}
	return index
}

func think() {
	var freeFighters []fighter
	var foundFightersMulti map[int32]icarusClient.RadarPing
	var foundBombers map[int32]icarusClient.RadarPing
	var foundOther map[int32]icarusClient.RadarPing
	foundFightersMulti = make(map[int32]icarusClient.RadarPing)
	foundBombers = make(map[int32]icarusClient.RadarPing)
	foundOther = make(map[int32]icarusClient.RadarPing)
	for iff, pings := range airPings {
		if pings.Type == icarusClient.Fighter || pings.Type == icarusClient.Multi {
			foundFightersMulti[iff] = pings
		}
	}
	for iff, pings := range foundFightersMulti {
		var shortest float64 = 10000
		var closest fighter
		for _, drone := range fighters {
			if drone.state == Patrol || drone.state == FiredRetreat {
				freeFighters = append(freeFighters, drone)
			}
		}
		if len(freeFighters) == 0 {
			break
		}
		for _, drone := range freeFighters {
			tempDist := utils.Haversine(pings.Latitude, pings.Longitude, drone.telem.Latitude, drone.telem.Longitude)
			if tempDist < shortest {
				fmt.Printf("\n\nSHORTES - %f\n\n", tempDist)
				closest = drone
				shortest = tempDist
			}
		}
		if closest.IFF != 0 {
			temp := &closest
			temp.curTarget = iff
			temp.state = Intercept
			closest = *temp
			fmt.Println(closest.state)
			fighters[closest.IFF] = closest
			fmt.Println(fighters[closest.IFF].state)
			fmt.Println(fighters[closest.IFF])
		}
	}
	for iff, pings := range airPings {
		if pings.Type == icarusClient.Bomber || pings.Type == icarusClient.WMD {
			foundBombers[iff] = pings
		}
	}
	for iff, pings := range foundBombers {
		var shortest float64 = 7500
		var closest fighter
		for _, drone := range fighters {
			if drone.state == Patrol || drone.state == FiredRetreat {
				freeFighters = append(freeFighters, drone)
			}
		}
		if len(freeFighters) == 0 {
			break
		}
		for _, drone := range freeFighters {
			tempDist := utils.Haversine(pings.Latitude, pings.Longitude, drone.telem.Latitude, drone.telem.Longitude)
			if tempDist < shortest {
				closest = drone
				shortest = tempDist
			}
		}
		temp := &closest
		temp.curTarget = iff
		temp.state = Intercept
		closest = *temp
	}
	for iff, pings := range airPings {
		if pings.Type == icarusClient.ISR || pings.Type == icarusClient.Cargo {
			foundOther[iff] = pings
		}
	}
	for iff, pings := range foundOther {
		var shortest float64 = 7500
		var closest fighter
		for _, drone := range fighters {
			if drone.state == Patrol || drone.state == FiredRetreat {
				freeFighters = append(freeFighters, drone)
			}
		}
		if len(freeFighters) == 0 {
			break
		}
		for _, drone := range freeFighters {
			tempDist := utils.Haversine(pings.Latitude, pings.Longitude, drone.telem.Latitude, drone.telem.Longitude)
			if tempDist < shortest {
				closest = drone
				shortest = tempDist
			}
		}
		temp := &closest
		temp.curTarget = iff
		temp.state = Intercept
		closest = *temp
	}
}

func askUserToFire() bool {
	resultChan := make(chan bool)
	// Start a goroutine to listen for user input
	go func() {
		var input string
		fmt.Print("\n\nDo you want to fire? (Press Enter to fire): ")
		fmt.Scanln(&input)
		if input == "" {
			resultChan <- true
		} else {
			resultChan <- false
		}
	}()
	select {
	case result := <-resultChan:
		return result
	case <-time.After(3 * time.Second):
		fmt.Println("Time's up! No response received.")
		return false
	}
}

func main3() {
	query := modulesImport.GenerateAuthQuery()

	setup(query)
	intelCount := 0
	gatherIntel(query)
	for {
		if intelCount >= 5 {
			gatherIntel(query)
			intelCount = -1
		}
		think()
		act(query)
		intelCount++
	}
}
