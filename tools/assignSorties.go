package main

import (
	"encoding/csv"
	"fmt"
	modulesImport "liam-tool/modules"
	"os"
	"strconv"
)

func readCSVs() ([]string, int) {
	sortieExists := false

	sortieCount := 1
	iffs, err := readData("sortie" + strconv.Itoa(sortieCount) + ".csv")
	droneIFFs := iffs
	for {
		if err != nil {
			if !sortieExists { // If the first sortie doesn't even exist
				fmt.Println("No sorties have been made yet.")
			}
			break
		}
		droneIFFs = append(droneIFFs, iffs...)
		sortieExists = true
		sortieCount += 1
		iffs, err = readData("sortie" + strconv.Itoa(sortieCount) + ".csv")
	}
	return droneIFFs, sortieCount
}

func readData(fileName string) ([]string, error) {

	f, err := os.Open(fileName)

	if err != nil {
		return []string{}, err
	}

	defer f.Close()

	r := csv.NewReader(f)

	// // skip first line
	// if _, err := r.Read(); err != nil {
	// 	return [][]string{}, err
	// }

	IFFs, err := r.Read() //r.ReadAll() would return array of string arrays

	if err != nil {
		return []string{}, err
	}

	return IFFs, nil
}

func isAssigned(iff string, drones []string) bool {
	for _, d := range drones {
		if d == iff {
			return true
		}
	}
	return false
}

func shouldBreak(e error) bool {
	if e != nil {
		fmt.Println("Error reading integer:", e)
		return true
	}
	return false
}

func main() {
	var ISR int
	var FTR int
	var BMR int
	var MLT int
	var CRG int
	// var WMD int
	fmt.Print("Enter how many ISR needed in this sortie: ")
	_, err := fmt.Scan(&ISR)
	if shouldBreak(err) {
		return
	}
	fmt.Print("Enter how many Fighters needed in this sortie: ")
	_, err = fmt.Scan(&FTR)
	if shouldBreak(err) {
		return
	}
	fmt.Print("Enter how many Bombers needed in this sortie: ")
	_, err = fmt.Scan(&BMR)
	if shouldBreak(err) {
		return
	}
	fmt.Print("Enter how many Multi-roles needed in this sortie: ")
	_, err = fmt.Scan(&MLT)
	if shouldBreak(err) {
		return
	}
	fmt.Print("Enter how many Cargos needed in this sortie: ")
	_, err = fmt.Scan(&CRG)
	if shouldBreak(err) {
		return
	}
	// fmt.Print("Enter how many WMDeez Nutz needed in this sortie: ")
	// _, err = fmt.Scan(&WMD)
	// if err != nil {
	// 	fmt.Println("Error reading integer:", err)
	// 	return
	// }

	ogISR, ogFTR, ogBMR, ogMLT, ogCRG := ISR, FTR, BMR, MLT, CRG
	//ogWMD := WMD
	query := modulesImport.GenerateAuthQuery()
	drones := modulesImport.GetAllIff(query)

	dronesAssigned, sortieInt := readCSVs()
	var sortie []string

	for i, droneType := range drones {
		iff := strconv.Itoa(int(i))
		if droneType == "ENVY" && ISR > 0 { //If ISR isn't full
			if !isAssigned(iff, dronesAssigned) {
				sortie = append(sortie, iff)
				ISR--
			}
		}
		if droneType == "OLAR" && FTR > 0 { //
			if !isAssigned(iff, dronesAssigned) {
				sortie = append(sortie, iff)
				FTR--
			}
		}
		if droneType == "SPITE" && BMR > 0 {
			if !isAssigned(iff, dronesAssigned) {
				sortie = append(sortie, iff)
				BMR--
			}
		}
		if droneType == "MALICE" && MLT > 0 {
			if !isAssigned(iff, dronesAssigned) {
				sortie = append(sortie, iff)
				MLT--
			}
		}
		if droneType == "BROOD" && CRG > 0 {
			if !isAssigned(iff, dronesAssigned) {
				sortie = append(sortie, iff)
				CRG--
			}
		}
	}
	//If vehicles still remain to be added to sortie but not enough of a type?
	if (ISR + FTR + BMR + MLT + CRG) > 0 {
		fmt.Print("WARNING! Not all aircraft available for sorty!\n\t")
		fmt.Print(strconv.Itoa(ogISR-ISR) + " ISR\n\t")
		fmt.Print(strconv.Itoa(ogFTR-FTR) + " Fighters\n\t")
		fmt.Print(strconv.Itoa(ogBMR-BMR) + " Bombers\n\t")
		fmt.Print(strconv.Itoa(ogMLT-MLT) + " Multi-roles\n\t")
		fmt.Print(strconv.Itoa(ogCRG-CRG) + " Cargos\n")
		fmt.Print("are included in your new sortie\n")
	}

	// Create new sortie csv file
	f, err := os.Create("sortie" + strconv.Itoa(sortieInt) + ".csv")
	defer f.Close() // defer until the end of the program

	if err != nil {
		fmt.Println("Error creating new file")
		return
	}
	w := csv.NewWriter(f)
	err = w.Write(sortie) // writes sorties to buffer
	w.Flush()             // Writes buffer to file

	if err != nil {
		fmt.Println("Error writing to new file")
	}
}
