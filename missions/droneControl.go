package main

import (
	"fmt"
	modulesImport "liam-tool/modules"
)

func droneControl() {
	fmt.Println("Choose flight mode: \n" +
		"1 for takeoff \n" +
		"2 for location tasking \n" +
		"3 for return to base \n" +
		"4 for landing \n" +
		"5 for tracking \n" +
		"6 for ISR Sweep")
	var controlNumber int
	fmt.Scanln(&controlNumber)
	if controlNumber == 1 {
		modulesImport.RunTakeoff()
	} else if controlNumber == 2 {
		modulesImport.SlowApproach()
	} else if controlNumber == 3 {
		modulesImport.RTB()
	} else if controlNumber == 4 {
		modulesImport.RunLand()
	//} else if controlNumber == 5 {
		//modulesImport.InputTrackTarget()
	//} else if controlNumber == 6 {
		//ISRRadarSweep()
	} else {
		fmt.Println("fuck off")
	}
}
