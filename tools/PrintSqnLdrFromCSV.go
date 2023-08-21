package main

import (
	"encoding/csv"
	"fmt"
	modulesImport "liam-tool/modules"
	"os"
	"strconv"
)

/*
	func main() {
		query := modulesImport.GenerateAuthQuery()
		var drones map[int32]string
		drones = make(map[int32]string)
		allDrones := modulesImport.GetAllIff(query) // allDrones is a map that contains every drone IFF as the keys

		for {
			modulesImport.SqnLdrDroneTable(drones, query)
		}
	}
*/
func main() {
	query := modulesImport.GenerateAuthQuery()

	var drones map[int32]string
	drones = make(map[int32]string)
	allDrones := modulesImport.GetAllIff(query)
	//open CSV file
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
		iff, err := strconv.Atoi(element)
		var iff32 = int32(iff)
		if err != nil {
			fmt.Println("Converting IFFs to Ints")
		}
		_, ok := allDrones[iff32]
		if !ok {
			fmt.Println("IFF not found in the allDrones map.")
			continue
		}
		if ok && err == nil {
			drones[iff32] = "empty"
		}
	}
	for {
		modulesImport.SqnLdrDroneTable(drones, query)
	}
}
