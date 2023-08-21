package main

import (
	"fmt"
	"liam-tool/modules"
	"os"
	"strconv"
	"strings"
)

func main2() {
	query := modules.GenerateAuthQuery()
	if os.Args[1] == "-h" || os.Args[2] == "--help" || len(os.Args) < 2 {
		fmt.Println("Usage: ./DCCLI <cmd> <opt1> <opt2> <opt3> <...>")
		fmt.Println("=========================================\n")
		fmt.Println("-m, --move\t<IFF1,IFF2> <lat> <lon> <alt> <vel>")
		fmt.Println("---------------------------------------------------")
		fmt.Println("Example: -m 2103 45.261 -62.514 100 160")
		fmt.Println("Example: -m 2103,2014 45.261 -62.514 100 160")
		fmt.Println("If you are controlling multiple drones, enter all the IFFs split by a comma (no spaces). " +
			"Battlebus will handle deconfliction (different altitudes)\n")
		fmt.Println("-t, --takeoff\t<IFF1,IFF2>")
		fmt.Println("---------------------------")
		fmt.Println("Example -t 2103")
		fmt.Println("Example -t 2103,2104\n")
		fmt.Println("-l, --land\t<IFF1,IFF2>")
		fmt.Println("---------------------------")
	}

	// Process the command-line arguments starting from index 1
	// If you want to convert arguments to integers, you can use strconv.Atoi
	// Example:
	// num, err := strconv.Atoi(arg)
	// if err != nil {
	//     fmt.Printf("Argument %d is not a valid integer: %s\n", i+1, err)
	//     continue
	// }
	// fmt.Printf("Argument %d (as integer): %d\n", i+1, num)
	if os.Args[1] == "-m" || os.Args[1] == "--move" {
		drones := strings.Split(os.Args[2], ",")
		temp, _ := strconv.ParseFloat(os.Args[3], 64)
		lat := float64(temp)
		temp, _ = strconv.ParseFloat(os.Args[4], 64)
		lon := float64(temp)
		temp, _ = strconv.ParseFloat(os.Args[5], 32)
		alt := float32(temp)
		temp, _ = strconv.ParseFloat(os.Args[6], 32)
		vel := float32(temp)
		fmt.Println(drones)
		fmt.Printf("%f,%f,%f,%f\n", lat, lon, alt, vel)
		var count float32 = 0
		for _, id := range drones {
			vid, _ := strconv.Atoi(id)
			modules.GoToLocGeneral(lat, lon, vid, vel, alt+(count*10), query)
			count++
		}
	} else if os.Args[1] == "-t" || os.Args[1] == "--takeoff" {
		drones := strings.Split(os.Args[2], ",")
		for i, val := range drones {
			id, _ := strconv.Atoi(val)
			modules.RunTakeoffGen(id, float32(10+i*10))
		}
	} else if os.Args[1] == "-l" || os.Args[1] == "--land" {
		drones := strings.Split(os.Args[2], ",")
		for _, val := range drones {
			id, _ := strconv.Atoi(val)
			modules.RunLandGen(id)
		}
	}
}
