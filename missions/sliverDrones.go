package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	icarusClient "git.ace/icarus/icarusclient/v5"
	"liam-tool/modules"
	"liam-tool/utils"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

const ERROR_USE = "usage: use <int Drone-ID 1> <int Drone-ID 2> ...>"
const ERROR_REMOVE = "usage: rm <int Drone-ID 1> <int Drone-ID 2> ...>"
const ERROR_TAKEOFF = "usage: takeoff <float alt>"

// const ERROR_GOTO = "usage: goto <float lat> <float long> <float alt> <float velocity>"
const ERROR_RTB = "choose base usage: rtb <float velocity> <string baseName>" +
	"go to closest usage: rtb <float velocity> <string m or h for Malazan or Halcyon>"
const ERROR_LAND = "usage: land"
const ERROR_RADAR_SWEEP = "usage: radarsweep"
const ERROR_TRACK = "usage: goto <int targetID> <float alt> <float velocity>"
const ERROR_LOAD = "usage: load"
const ERROR_PICTURE = "usage: picture"
const ERROR_BOMB = "usage: bomb <targetIFF>"
const ERROR_MISSILE = "usage: shoot <targetIFF>"

func main() {
	sorties, dronesUsed, availableDrones := AutomaticFunctions()
	for {
		input := GetNewCommand(dronesUsed)

		command := strings.ToLower(input[0])
		args := input[1:]

		if command == "exit" {
			return
		} else if command == "use" {
			DoUse(dronesUsed, args, availableDrones, sorties)
		} else if command == "rm" {
			DoRemove(dronesUsed, args)
		} else if command == "takeoff" {
			DoTakeoff(dronesUsed, args)
		} else if command == "goto" {
			DoGoto(dronesUsed, args)
		} else if command == "rtb" {
			DoRTB(dronesUsed, args)
		} else if command == "land" {
			DoLand(dronesUsed, args)
		} else if command == "track" {
			DoTrack(dronesUsed, args)
		} else if command == "drones" {
			PrintInitialDrones(availableDrones, args)
		} else if command == "clear" {
			DoClear(dronesUsed, args)
		} else if command == "sortie" || command == "sorties" {
			DoChangeSortie(sorties, dronesUsed, availableDrones, args)
		} else if command == "help" {
			DoShowHelpMenu()
		} else if command == "bomb" {
			DoBomb(dronesUsed, args)
		} else if command == "picture" {
			DoPicture(dronesUsed, args)
		} else if command == "shoot" {
			DoMissile(dronesUsed, args)
		} else if command == "load" {
			DoLoad(dronesUsed, args)
		} else if command == "radar" {
			DoPrintRadar(dronesUsed, args)
		} else if command == "boof" {
			Boof()
		} else if command == "loadcargo" {
			//loadCargo(dronesUsed, args)
		} else if command == "table" {
			reopenTable(args)
		}
	}
}

func reopenTable(args []string) {
	if len(args) != 0 && (args[0] == "help" || args[0] == "h" || args[0] == "-h" || args[0] == "-help") {
		fmt.Println("This command reopens the Drone Table after it inevitably crashes due to the shit show that is this place's internet")
		return
	}
	if len(args) > 0{
		fmt.Println("no args for this command\n")
		return
	}
	var cmd *exec.Cmd
	cmd = exec.Command("x-terminal-emulator", "-e", "./../binaries/MyDroneInfo")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Start()
	if err != nil {
		fmt.Println("Error starting new terminal:", err)
	}
}

func AutomaticFunctions() (map[string][]int, map[int]int, map[int]int) {
	//Print the sick banner
	RandomizeBanner()
	//fmt.Println(BANNER)
	//Prepare available drones for pilot by reading in their IFFs from a CSV file
	availableDrones := GetAvailableDrones()

	//automatically open the drone info script that Alan made
	var cmd *exec.Cmd
	cmd = exec.Command("x-terminal-emulator", "-e", "./../binaries/MyDroneInfo")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Start()
	if err != nil {
		fmt.Println("Error starting new terminal:", err)
	}

	//create the empty dronesUsed array (doesnt really need to be here, main was just cluttered)
	dronesUsed := make(map[int]int)
	sorties := make(map[string][]int)
	//These two commands prints the drones available prior to the pilot starting anything
	var empty []string //since theres a possible help argument, this was the best way to handle this.
	PrintInitialDrones(availableDrones, empty)
	return sorties, dronesUsed, availableDrones
}

func GetAvailableDrones() map[int]int {
	availableDrones := make(map[int]int)
	filePath := "../resources/iffs.csv"
	iffCSV, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Failure opening CSV file")
	}
	defer iffCSV.Close()
	reader := csv.NewReader(iffCSV)
	record, err := reader.Read()
	if err != nil {
		fmt.Println("Failure reading in available drones")
	}
	for _, element := range record {
		// Convert the element to an integer
		value, err := strconv.Atoi(element)
		availableDrones[value] = value
		if err != nil {
			fmt.Println("WTF")
		}
	}
	return availableDrones

}

func GetNewCommand(droneList map[int]int) []string {
	// prompt for a new command with a list of drones included
	fmt.Printf("\n")
	for _, droneID := range droneList {
		fmt.Printf("(%d) ", droneID)
	}
	fmt.Printf("> ")

	// get the command after the prompt
	var command string
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		command = scanner.Text()
	}

	return strings.Fields(command)
}

func DoUse(droneList map[int]int, args []string, availableDrones map[int]int, sorties map[string][]int) {
	if len(args) != 0 && (args[0] == "help" || args[0] == "h" || args[0] == "-h" || args[0] == "-help") {
		fmt.Println("\nThis command adds the specified IFFs to your selection. You can add all available IFFs to your selection by using \"use *\".\n" +
			"You can input a sortie instead of a list of IFFs to rapidly switch between groups you specified before. Think of a control group from an RTS like Starcraft/\n" +
			"Examples: \"use 2103 2104\"    OR      \"use fighters\"")
		return
	}
	if len(args) == 0 {
		return
	}
	//check if user is asking for a Sortie
	iffs, sortieExists := sorties[args[0]]
	if len(args) == 1 && sortieExists == true {
		//remove all drones from dronesUsed map
		for key := range droneList {
			delete(droneList, key)
		}
		//select all IFFs in Sortie
		for _, id := range iffs {
			droneList[id] = id
		}
		return
	}

	if len(args) != 0 && args[0] == "*" {
		for key, value := range availableDrones {
			droneList[key] = value
		}
		return
	}
	for _, arg := range args {
		id, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Println(ERROR_USE)
			return
		}
		_, ok := availableDrones[id]
		if ok {
			droneList[id] = id
			fmt.Printf("[+] Now using IFFID %d...\n", id)
		} else {
			fmt.Println("Fuck off, ", id, " is not your drone")
		}
	}
}

/*func loadCargo(droneList map[int]int,  args []string) {
	//query := modules.GenerateAuthQuery
	var wg sync.WaitGroup
	for _, DroneID := range droneList {
		_, ok := droneList[DroneID]
		if !ok {
			fmt.Printf("[!] IFFID %d not in used drones...", DroneID)
			return
		}

		wg.Add(1)
		go func(wg *sync.WaitGroup, DroneID int) {
			defer wg.Done()
			LoadCargo(DroneID)

		}(&wg, DroneID)
	}
	wg.Wait()
}*/

func DoRemove(droneList map[int]int, args []string) {
	if len(args) != 0 && (args[0] == "help" || args[0] == "h" || args[0] == "-h" || args[0] == "-help") {
		fmt.Println("\nThis command the specified IFFs from your selection. You can remove all IFFs from your selection by using \"rm *\". It acts the exact same as \"clear\"\n" +
			"Example: rm 2103 2104")
		return
	}
	if len(args) != 0 && args[0] == "*" {
		for key := range droneList {
			delete(droneList, key)
		}
		return
	}
	for _, arg := range args {
		id, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Println(ERROR_REMOVE)
			return
		}
		_, ok := droneList[id]
		if !ok {
			fmt.Printf("[!] IFFID %d not in used drones...", id)
			return
		}

		fmt.Printf("[*] Removing IFFID %d from drone list...\n", id)
		delete(droneList, id)
	}
}

func DoTakeoff(droneList map[int]int, args []string) {
	if len(args) == 1 && (args[0] == "help" || args[0] == "h" || args[0] == "-h" || args[0] == "-help") {
		fmt.Println("This command will tell all selected drones to takeoff and linger starting at the altitude specified, with increments of 5 meters between each drone upwards.\nThis command only needs altitude as a argument.\n\nExample use: \"takeoff 100\"")
		return
	}
	if len(args) != 1 {
		fmt.Println(ERROR_TAKEOFF)
		return
	}

	alt, err := strconv.ParseFloat(args[0], 32)
	if err != nil {
		fmt.Println(ERROR_TAKEOFF)
		return
	}

	var wg sync.WaitGroup
	for _, DroneID := range droneList {
		_, ok := droneList[DroneID]
		if !ok {
			fmt.Printf("[!] IFFID %d not in used drones...", DroneID)
			return
		}

		wg.Add(1)
		go func(wg *sync.WaitGroup, DroneID int, alt float64) {
			defer wg.Done()
			modules.RunTakeoffGen(DroneID, float32(alt))
		}(&wg, DroneID, alt)
		alt += 5
	}
	wg.Wait()
}

// used for passing a shit load of drones as one argument to new terminals
func droneListToString(droneList map[int]int) string {
	droneString := ""
	for key, _ := range droneList {
		droneString = fmt.Sprintf("%d,%s", key, droneString)
	}
	return droneString
}
func checkForStringElement(str [8]string, element string) (bool) {
	for _, v := range str {
		if v == element {
			return true
		}
	}
	return false
}

func DoGoto(droneList map[int]int, args []string) {
	if len(args) == 1 && (args[0] == "help" || args[0] == "h" || args[0] == "-h" || args[0] == "-help") {
		fmt.Println("\nThis command will tell all selected drones to travel to the location specified, with increments of 5 meters between each drone upwards.\nThe drones will form a grid pattern when they arrive rather than stack in one column.\nIt is highly recommended not to include more than 9 drones for a Goto command.\n\ngoto <float lat> <float long> <float alt> <float velocity>\n\nExample use: \"goto 46.0000 -60.0000 100 160\"")
		return
	}
	if len(args) != 4 && len(args) != 2 {
		//fmt.Println(ERROR_GOTO)
		return
	}


	var stringArray [8]string
	stringArray[0] = "n"
	stringArray[1] = "ne"
	stringArray[2] = "e"
	stringArray[3] = "se"
	stringArray[4] = "s"
	stringArray[5] = "sw"
	stringArray[6] = "w"
	stringArray[7] = "nw"
	if checkForStringElement(stringArray, strings.ToLower(args[0])) {
		var offset [2]float64
		switch(strings.ToLower(args[0])) {
		case "n":
			offset[0] = 0.5
			offset[1] = 0.0
		case "ne": 
			offset[0] = 0.5
			offset[1] = 0.5
		case "e":
			offset[0] = 0.0
			offset[1] = 0.5
		case "se":
			offset[0] = -0.5
			offset[1] = -0.5
		case "s":
			offset[0] = -0.5
			offset[1] = 0.0
		case "sw":
			offset[0] = -0.5
			offset[1] = 0.5
		case "w":
			offset[0] = 0.0
			offset[1] = -0.5
		case "nw":
			offset[0] = 0.5
			offset[1] = -0.5
		default:
		}
		query := modules.GenerateAuthQuery()
		var droneTelem icarusClient.Telemetry
		for _, DroneID := range droneList {
			droneTelem = modules.GetDroneInfo(DroneID, query).Vehicles[0].Telem
			break
		}
		offsetLat := strconv.FormatFloat(droneTelem.Latitude+offset[0], 'f', 2, 64)
		offsetLon := strconv.FormatFloat(droneTelem.Longitude+offset[1], 'f', 2, 64)
		altitude := strconv.FormatFloat(float64(droneTelem.Altitude), 'f', 2, 64)
		var cmd *exec.Cmd
		cmd = exec.Command("x-terminal-emulator", "-e", "./GoToSliverBus", offsetLat, offsetLon, altitude , args[1], droneListToString(droneList))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err := cmd.Start()
		if err != nil {
			fmt.Println("Error starting new terminal:", err)
		}
		return
	}
	



	//var speed = 0
	/*query := modules.GenerateAuthQuery()
	for vehicleID, _ := droneList {
		if GetDroneInfo(vehicleID, query).Vehicles[0].VConfig.Role == icarus.Bomber {
			if int(args[3]) > 105
			speed = 105
		}
		if GetDroneInfo(vehicleID, query).Vehicles[0].VConfig.Role == icarus.Fighter {
			if int(args[3]) > 190
			speed = 190
		}

	}*/
	if len(args) == 4{
		var cmd *exec.Cmd
		cmd = exec.Command("x-terminal-emulator", "-e", "./GoToSliverBus", args[0], args[1], args[2], args[3], droneListToString(droneList))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		err := cmd.Start()
		if err != nil {
			fmt.Println("Error starting new terminal:", err)
		}	
		return
	}

	fmt.Println("usage: goto <float lat> <float long> <float alt> <float velocity>")
	fmt.Println("usage: goto <n, s, se, etc> <float velocity> ")
	return
}

func DoRTB(droneList map[int]int, args []string) {
	//malazan bases
	var mBaseCoords = make(map[string][]float64)
	mBaseCoords["galladi"] = []float64{46.116, -64.6753}
	mBaseCoords["gast"] = []float64{47.0074, -65.4493}
	mBaseCoords["geleen"] = []float64{45.8693, -66.5328}
	mBaseCoords["bell"] = []float64{46.16, -60.046}
	mBaseCoords["borgin"] = []float64{46.165, -60.7836}
	mBaseCoords["broadmoor"] = []float64{45.6114, -62.6224}
	mBaseCoords["galain"] = []float64{46.14, -65.9042}
	mBaseCoords["burbage"] = []float64{45.6515, -61.3699}
	mBaseCoords["gruntle"] = []float64{47.6275, -61.416}
	mBaseCoords["serc"] = []float64{49.2305, -62.34}
	//halcyon bases
	var hBaseCoords = make(map[string][]float64)
	hBaseCoords["hope"] = []float64{49.1355, -61.7975}
	hBaseCoords["groundbreaker"] = []float64{49.8368, -64.2868}
	hBaseCoords["unreliable"] = []float64{47.4232, -61.7744}

	var offsetMap = make(map[int][]float64)
	offsetMap[0] = []float64{0, 0}
	offsetMap[1] = []float64{-0.0010, 0.0010}
	offsetMap[2] = []float64{-0.0010, -0.0010}
	offsetMap[3] = []float64{0.0010, 0.0010}
	offsetMap[4] = []float64{0.0010, -0.0010}
	offsetMap[5] = []float64{-0.0010, 0}
	offsetMap[6] = []float64{0.0020, 0}
	offsetMap[7] = []float64{0, 0.0020}
	offsetMap[8] = []float64{0, -0.0020}

	if len(args) == 1 && (args[0] == "help" || args[0] == "h" || args[0] == "-h" || args[0] == "-help") {
		fmt.Println("This command has 3 variations. ")
		fmt.Println("Type m to RTB to nearest Malazan base")
		fmt.Println("Type h to RTB to nearest Halcyon base\n")
		fmt.Println("usage:  rtb <float velocity> <name of base>\n")
		fmt.Println("usage:  rtb <float velocity> <m or h flag>")
		return
	}
	query := modules.GenerateAuthQuery()
	var userInput = false
	var halcyonFlag = false
	var malazanFlag = false

	if len(args) != 2 {
		fmt.Println(ERROR_RTB)
	}

	velo, err := strconv.ParseFloat(args[0], 32)
	if err != nil {
		fmt.Println(ERROR_RTB)
		return
	}
	//User input airbase
	var base = ""
	if args[1] != "h" && args[1] != "m" {
		base = strings.ToLower(args[1])
		userInput = true
	}

	if args[1] == "h" {
		halcyonFlag = true
	} else if args[1] == "m" {
		malazanFlag = true
	}
	var counter = 0
	var wg sync.WaitGroup
	for _, DroneID := range droneList {
		_, ok := droneList[DroneID]
		if !ok {
			fmt.Printf("[!] IFFID %d not in used drones...", DroneID)
			return
		}
		var targetLat = 0.0
		var targetLon = 0.0
		droneTelem := modules.GetDroneInfo(DroneID, query).Vehicles[0].Telem

		if userInput == true {

			baseLoc, exists := mBaseCoords[base]
			if exists != true {
				baseLoc, exists = hBaseCoords[base]
				if exists != true {
					fmt.Println("This base does not exist, dumbass")
					return
				}
			}
			targetLat = baseLoc[0]
			targetLon = baseLoc[1]
			targetAlt := droneTelem.Altitude
			fmt.Println(baseLoc)
			fmt.Println(targetLat)
			fmt.Println(targetLon)
			fmt.Println(targetAlt)
		}

		lat1 := droneTelem.Latitude
		lon1 := droneTelem.Longitude
		var distance = 1000000.0
		if halcyonFlag {
			var closestBase = ""
			for key, value := range hBaseCoords {
				iDistance := utils.Haversine(lat1, lon1, value[0], value[1])
				if iDistance < distance {
					closestBase = key
					distance = iDistance
				}
			}
			fmt.Println("RTBing to " + closestBase + " AB")
			targetLat = hBaseCoords[closestBase][0]
			targetLon = hBaseCoords[closestBase][1]
		}

		if malazanFlag {
			var closestBase = ""
			for key, value := range mBaseCoords {
				iDistance := utils.Haversine(lat1, lon1, value[0], value[1])
				if iDistance < distance {
					closestBase = key
					distance = iDistance
				}
			}
			fmt.Println("RTBing to " + closestBase)
			targetLat = mBaseCoords[closestBase][0]
			targetLon = mBaseCoords[closestBase][1]
		}

		targetAlt := float32(droneTelem.Altitude)

		index := counter % 9
		offsets := offsetMap[index]
		wg.Add(1)
		//count := 0
		go func(wg *sync.WaitGroup, DroneID int, alt float32) {
			defer wg.Done()

			modules.RunRTBWithSlowdown(droneTelem, DroneID, float32(targetAlt), float32(velo), float64(targetLat)+ offsets[0], float64(targetLon)+ offsets[1])
			//count++
		}(&wg, DroneID, targetAlt)
		//alt += 5
		counter +=1
	}
	wg.Wait()
}

func DoLand(droneList map[int]int, args []string) {
	var wg sync.WaitGroup
	for _, DroneID := range droneList {
		_, ok := droneList[DroneID]
		if !ok {
			fmt.Printf("[!] IFFID %d not in used drones...", DroneID)
			return
		}

		wg.Add(1)
		go func(wg *sync.WaitGroup, DroneID int) {
			defer wg.Done()
			modules.RunLandGen(DroneID)
		}(&wg, DroneID)
	}
	wg.Wait()
}

/*
	func DoRadarSweep(droneList map[int]int, args []string) {
		var wg sync.WaitGroup
		for _, DroneID := range droneList {
			_, ok := droneList[DroneID]
			if !ok {
				fmt.Printf("[!] IFFID %d not in used drones...", DroneID)
				return
			}*/
func DoGetRadar(droneList map[int]int) (map[int32]icarusClient.RadarPing){
	query := modules.GenerateAuthQuery()
	droneArr := make([]int32, 0, len(droneList)) // Preallocate the slice with capacity
	for droneID := range droneList {
		droneArr = append(droneArr, int32(droneID))
	}
	_, radar, _, _ := modules.RetrieveSomeRADAR(droneArr, query)
	return radar
}

func DoPrintRadar(droneList map[int]int, args []string) {
	query := modules.GenerateAuthQuery()
	droneArr := make([]int32, 0, len(droneList)) // Preallocate the slice with capacity
	for droneID := range droneList {
		droneArr = append(droneArr, int32(droneID))
	}
	_, radar, _, _ := modules.RetrieveSomeRADAR(droneArr, query)
	if len(radar) == 0 {
		fmt.Println("No radar pings")
		return
	}
	for key, value := range radar {
		var typeStr string
		if value.Type == 1 {
			typeStr = "FIGHTER"
		} else if value.Type == 2{
			typeStr = "BOMBER"
		} else if value.Type == 4 {
			typeStr = "ISR"
		} else if value.Type == 6 {
			typeStr = "MULTI"
		} else if value.Type == 12 {
			typeStr = "AIRLINER (DONT SHOOT)"
		} else {
			typeStr = "Something else"
		}
		var distanceClosest = 1000000.0
		for _, friendly := range droneArr {
			telem := modules.GetDroneInfo(int(friendly), query).Vehicles[0].Telem
			distance := utils.Haversine(value.Latitude, value.Longitude, telem.Latitude, telem.Longitude)
			if distance < distanceClosest {
				distanceClosest = distance
			}
		}
		
		fmt.Printf("IFF: %d  || ROLE: %s\t || Distance from Closest Drone: %.1f\n", key, typeStr, distanceClosest)
	}
	return
}



func DoTrack(droneList map[int]int, args []string) {
	if len(args) == 1 && (args[0] == "help" || args[0] == "h" || args[0] == "-h"){
		fmt.Println("This command tracks a specified target, continually telling the selected drones to go to the lat and lon of the target. \n Example use: \"track <target iffid> <float alt> <float velocity>\"")
		return
	}
	target, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println(ERROR_TRACK)
		return
	}

	alt, err := strconv.ParseFloat(args[1], 32)
	if err != nil {
		fmt.Println(ERROR_TRACK)
		return
	}

	velo, err := strconv.ParseFloat(args[2], 32)
	if err != nil {
		fmt.Println(ERROR_TRACK)
		return
	}

	var wg sync.WaitGroup
	for _, DroneID := range droneList {
		_, ok := droneList[DroneID]
		if !ok {
			fmt.Printf("[!] IFFID %d not in used drones...", DroneID)
			return
		}

		wg.Add(1)
		go func(wg *sync.WaitGroup, DroneID int) {
			defer wg.Done()
			modules.TrackTarget(droneList, target, float32(velo), float32(alt))

		}(&wg, DroneID)
	}
	wg.Wait()
}

func DoClear(droneList map[int]int, args []string) {
	if len(args) != 0 && (args[0] == "help" || args[0] == "h" || args[0] == "-h" || args[0] == "-help") {
		fmt.Println("\nThis command removes all IFFs from your selection. It acts the exact same as \"rm *\"")
		return
	}
	for key := range droneList {
		delete(droneList, key)
	}
}

// prints a nicely formatted list of all available drones
func PrintInitialDrones(availableDrones map[int]int, args []string) {
	//help section
	if len(args) != 0 && (args[0] == "help" || args[0] == "h" || args[0] == "-h") {
		fmt.Println("\nThis command prints out all the drone IFFs included in your iffs.csv file. It's a quick reminder of what you have!")
		return
	}
	//actual code
	var str = ""
	keys := make([]int, 0, len(availableDrones))
	for key := range availableDrones {
		keys = append(keys, key)
	}
	if len(keys) == 1 {
		fmt.Println("\n\n ||Your drone is " + strconv.Itoa(keys[0]) + "||")
	}
	for i, key := range keys {
		if i == len(keys)-1 {
			stringifiedID := strconv.Itoa(key)
			str = str + "and " + stringifiedID
			break
		}
		stringifiedID := strconv.Itoa(key)
		str = str + stringifiedID + ", "

	}
	if len(keys) > 1 {
		fmt.Println("\n ||Your drones are " + str + "||")
	}
	return
}

func DoChangeSortie(sorties map[string][]int, droneList map[int]int, availableDrones map[int]int, args []string) map[string][]int {
	if len(args) == 1 && (args[0] == "help" || args[0] == "h" || args[0] == "-h" || args[0] == "-help") {
		fmt.Println("\nThis command, used with the argument add or rm, will allow the creation and deletion of control groups within the set of drones you control.")
		fmt.Println("Use case: you have 2 fighters and 2 bombers. Use the command \"sortie add bombers 2217 2218\" to create a sortie named \"bombers\".")
		fmt.Println("You could then input \"use bombers\", which will select 2217 and 2218.\n\n")
		fmt.Println("sortie <add or rm> <name> <list of IFFIDs separated by spaces>\n")
	}
	if (args[0] != "ls") && (args[0] != "rm") && (args[0] != "add") {
		fmt.Println(args[0])
		fmt.Println("exiting")
		return sorties
	}
	if args[0] == "ls" && len(args) == 1 {
		for key, value := range sorties {
			fmt.Printf("Sortie: %s, IFFs: %d\n", key, value)
		}
	}
	if args[0] == "rm" && len(args) >= 2 {
		var selectedSortie = args[1]
		_, exists := sorties[selectedSortie]
		if exists == true {
			delete(sorties, selectedSortie)
			fmt.Println("Sortie " + selectedSortie + " deleted")
		} else {
			fmt.Println("That is not a valid sortie")
		}
	}

	if args[0] == "add" && len(args) >= 2 {
		//get the user's sortie name from args
		var sortieName = args[1]
		//make sure the pilot isnt breaking shit with their name
		if sortieName == "*" {
			fmt.Println("Dont name it that, dumbass\n")
			return sorties
		}
		//make sure the sortie doesnt already exist by that name
		_, exists := sorties[sortieName]
		if exists == true {
			fmt.Println("Cannot duplicate sorties. Try a different name.\n")
		}
		//get the iffs from args at strings
		iffsSTRs := args[2:]
		//convert the iffs to an array of integers to add into the sortie map
		iffs := make([]int, len(iffsSTRs))
		for i, s := range iffsSTRs {
			if i == len(iffsSTRs) {
				break
			}
			num, err := strconv.Atoi(s)
			if err != nil {
				fmt.Println("Failure to convert string to Integer")
			}
			_, exists := availableDrones[num]
			if exists == true {
				//adds the new IFF integer into the array
				iffs[i] = num
			} else {
				fmt.Println(s + " is not a valid IFF.\n")
			}
		}
		iffsCleaned := removeValueFromArray(iffs, 0)
		sorties[sortieName] = iffsCleaned
		fmt.Println("Sortie " + sortieName + " created!")
	}
	return sorties
}

// To remove 0 from any sorties when that occurs
func removeValueFromArray(arr []int, target int) []int {
	var result []int
	for _, value := range arr {
		if value != target {
			result = append(result, value)
		}
	}
	return result
}

func DoBomb(dronesUsed map[int]int, args []string) {
	if len(args) == 1 && (args[0] == "help" || args[0] == "h" || args[0] == "-h" || args[0] == "-help") {
		fmt.Println("This command will drop one bomb (thermal lance) each time its used. You must specify the IFFID you want to bomb. \n Example use: \"bomb 3333\"")
		return
	} else if len(args) > 1 || len(args) ==0 {
		fmt.Println(ERROR_BOMB)
		return
	}
	if len(dronesUsed) != 1 {
		fmt.Println("Please only use one drone to bomb.")
		return
	}
	enemyIFF, err := strconv.ParseFloat(args[0], 32)
	if err != nil {
		fmt.Println(err)
	}
	for friendly, _ := range dronesUsed {
		modules.ExecutePayload(uint(friendly), uint(enemyIFF), 3)
		break // Exit the loop after the first iteration
	}
	return
}

func DoMissile(dronesUsed map[int]int, args []string) {
	query := modules.GenerateAuthQuery()
	if len(args) == 1 && (args[0] == "help" || args[0] == "h" || args[0] == "-h" || args[0] == "-help") {
		fmt.Println("This command shoot one A2A missile each time its used. You must specify the IFFIDs you want to shoot at. \n Example use: \"shoot 3333 3332\"")
		return
	} else if len(args) > 0 {
		fmt.Println(ERROR_MISSILE)
	}


	var empty []string
	DoPrintRadar(dronesUsed, empty)
	var targetInput string
	fmt.Println("Select your target(s):")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	targetInput = scanner.Text()
	enemyIFFStrings := strings.Split(targetInput, " ")
	
	enemyIFFs := make([]int, len(enemyIFFStrings))
	for i, s := range enemyIFFStrings {
		if i == len(enemyIFFStrings) {
			break
		}
		num, err := strconv.Atoi(string(s))
		if err != nil {
			fmt.Println("Conversion error")
		}
		enemyIFFs[i] = num
	}	
	
	radarPings := DoGetRadar(dronesUsed)
	for _, enemyIFF := range enemyIFFs {
		var shotAt = false
		for friendly, _ := range dronesUsed {
			fmt.Println(friendly)
			friendlyInfo := modules.GetDroneInfo(friendly, query).Vehicles[0].Telem
			distanceToTarget := utils.Haversine(radarPings[int32(enemyIFF)].Latitude, radarPings[int32(enemyIFF)].Longitude, friendlyInfo.Latitude, friendlyInfo.Longitude)
			if distanceToTarget < 5000 /*&& hasMissile == true*/{
				executeSeq := query.ExecutePayload(int(friendly), icarusClient.AntiMatterMissile, 1, icarusClient.EmptyParams(), int(enemyIFF))
			
				responseChan, _ := query.Execute() // Pew pew
				fmt.Println("Waiting for responses:")
				response := <-responseChan
				_, ok := response.Get(executeSeq)
				if !ok {
					fmt.Println("AntiMatterMissile execute response not found")
				}
				shotAt = true
			} else {
				if distanceToTarget > 5000 {
					fmt.Println("Out of range: "+ strconv.FormatFloat(distanceToTarget, 'f', -1, 64))
				} 
			}
			if shotAt == true {
				break
			}
		}
	}
	
}

func DoPicture(dronesUsed map[int]int, args []string) {
	if len(args) == 1 && (args[0] == "help" || args[0] == "h" || args[0] == "-h" || args[0] == "-help") {
		fmt.Println("This command will take a picture immediately upon use. Only one ISR may use this at a time. \n Example use: \"picture\"")
		return
	} else if len(args) > 0 {
		fmt.Println(ERROR_PICTURE)
		return
	}
	if len(dronesUsed) != 1 {
		fmt.Println("Please only use one drone to bomb.")
		return
	}
	for friendly, _ := range dronesUsed {
		query := modules.GenerateAuthQuery()
		executeSeq := query.ExecutePayload(int(friendly), icarusClient.Camera, 1, icarusClient.EmptyParams(), 0)
		responseChan, _ := query.Execute() // Pew pew
		response := <-responseChan
		_, ok := response.Get(executeSeq)
		if !ok {
			fmt.Println("Camera execute response not found")
		}
		if ok {
			fmt.Println("Picture Taken!")
		}
		return
	}
}

func DoLoad(droneList map[int]int, args []string) {
	if len(args) == 1 && (args[0] == "help" || args[0] == "h" || args[0] == "-h" || args[0] == "-help") {
		fmt.Println("This command will tell all selected drones to load fuel and payloads. \nThis command takes no arguments.\n\nExample use: \"load\"")
		return
	} else if len(args) > 0 {
		fmt.Println(ERROR_LOAD)
		return
	}

	var wg sync.WaitGroup
	for _, DroneID := range droneList {
		_, ok := droneList[DroneID]
		if !ok {
			fmt.Printf("[!] IFFID %d not in used drones...", DroneID)
			return
		}

		wg.Add(1)
		go func(wg *sync.WaitGroup, DroneID int) {
			defer wg.Done()
			modules.RunLoadGen(DroneID)
		}(&wg, DroneID)
	}
	wg.Wait()
}

func DoShowHelpMenu() {
	fmt.Println("All commands:\n===========")
	fmt.Println("use \nrm \ntakeoff \ngoto \nrtb \nland \ntrack \ndrones \nclear \nsortie \nhelp \nbomb \nshoot \npicture \nload \table \boof")
}


/*func getFastestGroupSpeed(droneList map[int]int) (int) {

}*/

func Boof() {
	//A nice easter egg function
	var boof1 = "" +
		"((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((" +
		"((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((" +
		"((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((((" +
		"((((((((((((((((((((((((((((#&&&&&&&&&&&&&&&%###((((((((((((((((((((((((((((((((" +
		"((((((((((((((((((((((((#&&&&&&%%%%&&&&&&&&&@@@@&&%##(((((((((((((((((((((((((((" +
		"((((((((((((((((((((((&&&@&@&&&&&&&&@@@&&&&&@@@@@@@@@%##((((((((((((((((((((((((" +
		"(((((((((((((((((((%%&&&&&&%%&&&%%%&&&&&&@@&&&@@&&@@@@&%##((((((((((((((((((((((" +
		"((((((((((((((((((%%%%%%%#((/(((##(((((((#%%%%&&@&@@@@@@@##(((((((((((((((((((((" +
		"(((((((((((((((((%%&%#%%#******////((((((((((((#&&&@@@@@@&##((((((((((((((((((((" +
		"(((((((((((((((((&&%%%%(********//////((((((((((#%##%%&@@@&##(((((((((((((((((((" +
		"(((((((((((((((((%&%%%#/*******///////((((((((((((#####%&@@###((((((((((((((((((" +
		"((((((((((((((((#&&&%%(/******////////((((((((((######%%&&&%##((((((((((((((((((" +
		"(((((((((((((((((&&&%#(/****/((((#####(((((((((#######%%&&&%##((((((((((((((((((" +
		"(((((((((((((((((%&&&%(/***////((((((((((((##%%%####%%%%&&@###((((((((((((((((((" +
		"/((((((((/((((((/##%&%(****/(((((%%(##((//#%%#(&%%%%#%%%&@&##(((((((((((((((((((" +
		"/((/((((((((((((////#((*****////((((((//*((%%#######%#%%&@##((((((((((((((((((((" +
		"/////((/(((/(((/////(/********///////////(#%%#########%&@##(((((((((((((((((((((" +
		"//(//((((((//(/(/(//*///******//////(****((##%########%#(#((((((((((((((((((((((" +
		"/////(///(((((((/((/((((*****/////(//(#((##%###((####%&(((((((((((((((((((((((((" +
		"////////((((((((((//((((*/////**/((/////(((####((####%#(((((((((((((((((((((((((" +
		"///////////(/(////////((*///***//(#%%%%%%%%%########%%((((((((((((((((((((((((((" +
		"/////////(///(//////////**///**///////((###########%%(((((((((((((((((((((((((((" +
		"////////////////////////**////////////((((########%#((((((((((((((((((((((((((((" +
		"///////////////////////,**//*//////////((((#####%#((((((((((((((((((((((((((((((" +
		"//////////////////////...///////(((//((((((##%%%%(((((((((((((((((((((((((((((((" +
		"///////////////////&@......,(////((((#%%%%%%%%%%%(((((((((((((((((((((((((((((((" +
		"////////////////(&&&&*....,,,,,/((((((((###%##%#/(###%((((((((((/(((((((((((((((" +
		"/////////////#&&&&&&&&*.,,,,,,,,,,,((((((#####(/(&&@@@@(((((((((((((((((((((((((" +
		"//////(&&&&&&&&&&&&&&&&/,,,,,,,,,/(/,*/#####//&@@&@@@@@@@@@#((((((((((((((((((((" +
		"#&&&&&&&&&&&&&&&&&&&&&&&(*,,,,.@@&@@@@@@@@@@@@@&@&&&&&&&@@@@@@@@&#((((((((/(((((" +
		"&&&&&&&&&&&&&&&&&&&&&&&&&%*,,,,,@@&@@&@@@@&@@&@@&&&&&&&&@@@@@@@@@@@@@@@@&((/////" +
		"&&&&&&&&&&&&&&&&&&&&&&&&&&%*,,,,,,@@&&&&&@&#////***@@@@@@@@@@@@@@@@@@@@@@@@&((//" +
		"                                                                                " +
		"                         -- A L A N    T H E    S O N --                        "

	var boof2 = "" +
		"@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@&@&@@&&&&&&@&&&@@@@@@@@@@&&@&@@@@@" +
		"@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@" +
		"@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@%*.(%/(.@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@" +
		"@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@#&&##(((##%,,#/*%%&@@@@@@@@@@@@@@@@@@@@@@@@@@@@@" +
		"@@@@@@@@@@@@@@@@@@@@@@@@@@&@@/%%&/*,,*//,#(##*../**,**%@@@@@@@@@@@@@@@@@@@@@@@@@" +
		"@@@@@@@@@@@@@@@@@@@@@@@@@@**.,,//**//********/#,/(#,(##*(&@@@@@@@@@@@@@@@@@@@@@@" +
		"@@@@@@@@@@@@@@@@@@@@@@@%,/%/,/,/(((/***/***//,****(*((%*/*/%&@@@@@@@@@@@@@@@@@@@" +
		"@@@@@@@@@@@@@@@@@@@@@@(/#/***/(/////**************#/((#(#*/%(#&@@@@@@@@@@@@@@@@@" +
		"@@@@@@@@@@@@@@@@@@@@@%(#***//*/*******,***,,,,,,***///(((/&%#%%@@@@@@@@@@@@@@@@@" +
		"@@@@@@@@@@@@@@@@@@@@@#&((///****,,,.,,,,,,,,,,,,*****///(##%%&@@@@@@@@@@@@@@@@@@" +
		"@@@@@@@@@@@@@@@@@@@@&%%#///***********,*,*,*,,,,,,,**///(((%##&@@@@@@@@@@@@@@@@@" +
		"@@@@@@@@@@@@@@@@@@@@&%%(//*(/%%%%%##((////(#%#%%%##/*/((((#&&%&@@@@@@@@@@@@@@@@@" +
		"@@@@@@@@@@@@@@@@@@@@@%%(((((##/*###%#(*,*(##%###(/%#(((((#%&&@%@@@@@@@@@@@@@@@@@" +
		"@@@@@@@@@@@@@@@@@@@@@@/@%(((%%(%%(##((/%#///#(%#%/%((/(%/#&@&&@@@@@@@@@@@@@@@@@@" +
		"@@@@@@@@@@@@@@@@@@@@@%&&&(/(((##(/*@#*,,*%&*///##(////(,##&&%##@@@@@@@@@@@@@@@@@" +
		"@@@@@@@@@@@@@@@@@@@@@%#,&//*******/*******&//*,*,****/(((#&&/(%@@@@@@@@@@@@@@@@@" +
		"@@@@@@@@@@@@@@@@@@@@@/(%#///////(((%%(((%&&((/(/****//((##&#((#@@@@@@@@@@@@@@@@@" +
		"@@@@@@@@@@@@@@@@@@@@@@/(#(////((/*(####((#((/**/*(//((((#%(((#@@@@@@@@@@@@@@@@@@" +
		"@@@@@@@@@@@@@@@@@@@@@@@(#///(((/(##%######%#%##(((////((#(/(&@@@@@@@@@@@@@@@@@@@" +
		"@@@@@@@@@@@@@@@@@@@@@@@@@(/*/(#(/(((((#((((##/((##((//(##@@@@@@@@@@@@@@@@@@@@@@@" +
		"@@@@@@@@@@@@@@@@@@@@@@@@@#(//////////(((////***/((/(/((#@@@@@@@@@@@@@@@@@@@@@@@@" +
		"@@@@@@@@@@@@@@@@@@@@@@@@@@#(((//********,*****/(((#((#%@@@@@@@@@@@@@@@@@@@@@@@@@" +
		"@@@@@@@@@@@@@@@@@@@@@@@@@@(%#((#///////***///(##(##%#%#@@@@@@@@@@@@@@@@@@@@@@@@@" +
		"@@@@@@@@@@@@@@@@@@@@@@@@@@(####/#%&##%#%%%%&%%%%#%%%#/#@@@@@@@@@@@@@@@@@@@@@@@@@" +
		"@@@@@@@@@@@@@@@@@@@@@@@@@@(//###(%##%#%%#%#%##%##%((((@@@@@@@@@@@@@@@@@@@@@@@@@@" +
		"@@@@@@@@@@@@@@@@@@@@@@@@@@(((/(##(###%#%&%###(##((//((@@@@@@@@@@@@@@@@@@@@@@@@@@" +
		"@@@@@@@@@@@@@@@@@@@@@@@@&@/(/(((((((##(##((#(((((((/(#@@@@@@@@@@@@@@@@@@@@@@@@@@" +
		"@@@@@&&&@@@@@@@@@@@@@&&&@(/////(((/((///((((/(/((((((#@@@@@@@@@@@@@@@@@@@@@@@@@@" +
		"&&&&&@@&&&&@@@@@@&&&&&&@@(///////((/((/((((///(////((#@@@@@@@@@@@@@@@@&&&@@@@@@@" +
		"&&&&&&&&&&&&@@@&&@@&&&&@#(//////////(/////(///////((((#@@&@@@@@@@&@@@&&@&&&@&&&&" +
		"&&&&&&&%&&&&@@@&&@@&&&&&#((/(///////////////////((((((#@&&&@@@@@@@&@@@@&&%%&&&&&" +
		"&&&&&&&&&&&@@@&&&&&&&@&@%#((((//////(//////////(/((((&&&&&@@@&@@&@@@&&&&&&&%&&&&" +
		"                                                                                " +
		"                     --  C O L E   T H E   B A P T I S T  --                    "

	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 2
	var num = rand.Intn(max-min+1) + min
	if num == 1 {
		fmt.Println(boof1)
	} else if num == 2 {
		fmt.Println(boof2)
	} /*else if num == 3 {
		fmt.Println(randomtext3)
	} else if num == 4 {
		fmt.Println(randomtext4)*/

	return
}

func RandomizeBanner() {
	var BANNER = "" +
		" ____        _   _   _      _                \n" +
		"| __ )  __ _| |_| |_| | ___| |__  _   _ ___  \n" +
		"|  _ \\ / _` | __| __| |/ _ \\ '_ \\| | | / __| \n" +
		"| |_) | (_| | |_| |_| |  __/ |_) | |_| \\__ \\ \n" +
		"|____/ \\__,_|\\__|\\__|_|\\___|_.__/ \\__,_|___/\n" +
		"===========================================  FULL RELEASE v1.1"
	var AUTHORS = "------------------------------------------\n" +
		"Authors: The Holy Trinity and Cole the Baptist\n"

	var randomtext1 = "    F  U  C  K     V  A  L  I  N  O  R"
	var randomtext2 = "  V A L I N O R   D E L E N D A  E S T"
	var randomtext3 = "         B O O F    F O R E V E R"
	var randomtext4 = "   I F  I T  F L I E S , I T  D I E S   "

	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 4
	var num = rand.Intn(max-min+1) + min
	fmt.Println(BANNER)
	if num == 1 {
		fmt.Println(randomtext1)
	} else if num == 2 {
		fmt.Println(randomtext2)
	} else if num == 3 {
		fmt.Println(randomtext3)
	} else if num == 4 {
		fmt.Println(randomtext4)
	}
	fmt.Println(AUTHORS)
	return

}