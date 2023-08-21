package main

/*
func main() {
	query := modulesImport.GenerateAuthQuery()
	temp := modulesImport.GetAllIff(query)
	var permpings map[string]icarusClient.RadarPing
	permpings = make(map[string]icarusClient.RadarPing)
	var temp2 []int32 poop
	for drone, _ := range temp {
		temp2 = append(temp2, drone)
	}

	for {
		pings := modulesImport.RetrieveAllRADAR(query)
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
		if len(pings) == 0 {
			fmt.Println("No RADAR pings")
		} else {
			fmt.Printf("IFF\tType\tLat\t\tLon\n")
			for key, item := range pings {
				fmt.Printf("%d\t%s\t%f\t%f\n", key, icarusClient.VehicleRole.String(icarusClient.VehicleRole(item.Type)), item.Latitude, item.Longitude)
			}
		}
	}
}

func testRetriveve() {
	query := modulesImport.GenerateAuthQuery()
	temp := modulesImport.GetAllIff(query)
	var permpings map[int32]icarusClient.RadarPing
	permpings = make(map[int32]icarusClient.RadarPing)
	var temp2 []int32
	for drone, _ := range temp {
		temp2 = append(temp2, drone)
	}

	for {
		pings := modulesImport.RetrieveAllRADAR(query)
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()

		if len(pings) == 0 {
			fmt.Println("No RADAR pings")
		} else {
			fmt.Printf("IFF\tType\tLat\t\tLon\n")
			for key, item := range pings {
				fmt.Printf("%d\t%s\t%f\t%f\n", key, icarusClient.VehicleRole.String(icarusClient.VehicleRole(item.Type)), item.Latitude, item.Longitude)

				// Add new unique keys and items from pings to permpings
				if _, ok := permpings[key]; !ok {
					permpings[key] = item
				}
			}

			// Write permpings contents to a file called "radlogs.txt"
			if err := writePermpingsToFile(permpings); err != nil {
				fmt.Printf("Error writing to radlogs.txt: %s\n", err)
			}
		}

	}
}

func writePermpingsToFile(permpings map[int32]icarusClient.RadarPing) error {
	file, err := os.Create("radlogs.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the contents of permpings to the file
	for key, item := range permpings {
		_, err := fmt.Fprintf(file, "%d\t%s\t%f\t%f\n", key, icarusClient.VehicleRole.String(icarusClient.VehicleRole(item.Type)), item.Latitude, item.Longitude)
		if err != nil {
			return err
		}
	}

	return nil
}

*/
