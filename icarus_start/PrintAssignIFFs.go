package main

import (
	"fmt"
	icarusClient "git.ace/icarus/icarusclient/v5"
	modules "liam-tool/modules"
)

func main() {
	query := modules.GenerateAuthQuery()
	var sorties map[string][]int
	sorties = make(map[string][]int)

	var fighters, bombers, isr, multi, cargo, strategic []int
	var remF, remB, remI, remM, remC, remS int

	for _, drone := range modules.GetAllDroneInfo(query).Vehicles {
		if drone.VConfig.Role == icarusClient.Fighter {
			fighters = append(fighters, int(drone.VehicleId))
		} else if drone.VConfig.Role == icarusClient.Bomber {
			bombers = append(bombers, int(drone.VehicleId))
		} else if drone.VConfig.Role == icarusClient.ISR {
			isr = append(isr, int(drone.VehicleId))
		} else if drone.VConfig.Role == icarusClient.Multi {
			multi = append(multi, int(drone.VehicleId))
		} else if drone.VConfig.Role == icarusClient.Cargo {
			cargo = append(cargo, int(drone.VehicleId))
		} else if drone.VConfig.Role == icarusClient.WMD {
			strategic = append(strategic, int(drone.VehicleId))
		} else {
			fmt.Printf("Unknown type %d for IFF %d\n", drone.VConfig.Role, drone.VehicleId)
		}
	}
	remF = len(fighters)
	remB = len(bombers)
	remI = len(isr)
	remC = len(cargo)
	remM = len(multi)
	remS = len(strategic)
	fmt.Println(fighters)
	fmt.Println(bombers)
	fmt.Println(isr)
	fmt.Println(multi)
	fmt.Println(cargo)
	fmt.Println(strategic)

	var numSorties int
	fmt.Print("Enter the number of sorties > ")
	fmt.Scan(&numSorties)

	// Print the table header
	// Print the table rows for each sortie

	/*


		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
		fmt.Println("Sortie Name\tFighters\tISR\tBombers\tMulti\tCargo\tStrat")
		// Print the table rows for each sortie
		for key, item := range sorties {
			fmt.Printf("%-12s\t%-8s\t%-4s\t%-7s\t%-5s\t%-5s\t%-4s\n", key, item[0], item[1], item[2], item[3], item[4], item[5])
		}

		for key, item := range sorties {
			if remF == 0 {
				break
			}
			fmt.Printf("How many fighters in %s's sortie (there are %d remaining) > ", key, remF)
			var temp int
			fmt.Scan(&temp)
			for i := 0; i < temp; i++ {
				item = append(item, fighters[i])
				remF--
			}
			fighters[temp] = fighters[len(fighters)-1]
			fighters = fighters[:len(fighters)-1]
		}
		for key, item := range sorties {
			if remB == 0 {
				break
			}
			fmt.Printf("How many bombers in %s's sortie (there are %d remaining) > ", key, remB)
			var temp int
			fmt.Scan(&temp)
			for i := 0; i < temp; i++ {
				item = append(item, bombers[i])
				remB--
			}
			fighters[temp] = fighters[len(fighters)-1]
			fighters = fighters[:len(fighters)-1]
		}
		for key, item := range sorties {
			if remI == 0 {
				break
			}
			fmt.Printf("How many ISR in %s's sortie (there are %d remaining) > ", key, remI)
			var temp int
			fmt.Scan(&temp)
			for i := 0; i < temp; i++ {
				item = append(item, isr[i])
				remI--
			}
			fighters[temp] = fighters[len(fighters)-1]
			fighters = fighters[:len(fighters)-1]
		}
		for key, item := range sorties {
			if remM == 0 {
				break
			}
			fmt.Printf("How many multi in %s's sortie (there are %d remaining) > ", key, remM)
			var temp int
			fmt.Scan(&temp)
			for i := 0; i < temp; i++ {
				item = append(item, multi[i])
				remM--
			}
			fighters[temp] = fighters[len(fighters)-1]
			fighters = fighters[:len(fighters)-1]
		}
		for key, item := range sorties {
			if remC == 0 {
				break
			}
			fmt.Printf("How many cargo in %s's sortie (there are %d remaining) > ", key, remC)
			var temp int
			fmt.Scan(&temp)
			for i := 0; i < temp; i++ {
				item = append(item, cargo[i])
				remC--
			}
			fighters[temp] = fighters[len(fighters)-1]
			fighters = fighters[:len(fighters)-1]
		}
		for key, item := range sorties {
			if remF == 0 {
				break
			}
			fmt.Printf("How many strategic bombers in %s's sortie (there are %d remaining) > ", key, remS)
			var temp int
			fmt.Scan(&temp)
			for i := 0; i < temp; i++ {
				item = append(item, strategic[i])
				remS--
			}
			fighters[temp] = fighters[len(fighters)-1]
			fighters = fighters[:len(fighters)-1]
		}
	*/
	fmt.Printf("\nCurrent Fleet:\nF-%d\nB-%d\nI-%d\nM-%d\nC-%d\nS-%d\n\n", remF, remB, remI, remM, remC, remS)
	for i := 0; i < numSorties; i++ {
		fmt.Printf("Enter the name of the pilot for sortie %d > ", i+1)
		var temp string
		fmt.Scan(&temp)

		var tempInt int
		var request []int
		fmt.Println("Number of fighters > ")
		fmt.Scan(&tempInt)
		request = append(request, tempInt)
		fmt.Println("Number of bombers > ")
		fmt.Scan(&tempInt)
		request = append(request, tempInt)
		fmt.Println("Number of ISR > ")
		fmt.Scan(&tempInt)
		request = append(request, tempInt)
		fmt.Println("Number of multi > ")
		fmt.Scan(&tempInt)
		request = append(request, tempInt)
		fmt.Println("Number of cargo > ")
		fmt.Scan(&tempInt)
		request = append(request, tempInt)
		fmt.Println("Number of strat bombers > ")
		fmt.Scan(&tempInt)
		request = append(request, tempInt)

		sorties[temp] = tempArray
	}

	fmt.Println("Sortie Name\tFighters\tISR\tBombers\tMulti\tCargo\tStrat")
	for i := 1; i <= numSorties; i++ {
		sortieName := fmt.Sprintf("Sortie %d", i)
		fmt.Printf("%-12s\t%-8s\t%-4s\t%-7s\t%-5s\t%-5s\t%-4s\n", sortieName, "", "", "", "", "", "")
	}
}
