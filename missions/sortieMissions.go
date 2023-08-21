package main

func sortie() {
	//	args := os.Args
	//	controlNumber := args[1]
	//	fmt.Println(controlNumber)
	//	numDrones, _ := strconv.Atoi(args[2])
	//	for i := 0; i < numDrones; i++ {
	//		VehicleId, _ := strconv.Atoi(args[3+i])
	//		if controlNumber == "1" {
	//			modulesImport.RunTakeoffGeneral(VehicleId, 80)
	//		} else if controlNumber == "2" {
	//			Lat, _ := strconv.ParseFloat(args[3+numDrones], 64)
	//			Long, _ := strconv.ParseFloat(args[4+numDrones], 64)
	//			Velo, _ := strconv.ParseFloat(args[5+numDrones], 32)
	//			Alt, _ := strconv.ParseFloat(args[6+numDrones], 32)
	//			newMissionTarget := modulesImport.MissionTarget{}
	//			newMissionTarget.VehicleID = VehicleId
	//			newMissionTarget.Latitude = Lat + float64(i)*.0001
	//			newMissionTarget.Longitude = Long
	//			newMissionTarget.Altitude = float32(Alt) + float32(10*i)
	//			newMissionTarget.Heading = 0
	//			newMissionTarget.Velocity = float32(Velo)
	//			modulesImport.RunGoTo(newMissionTarget.VehicleID, newMissionTarget)
	//		} else if controlNumber == "3" {
	//			Velo, _ := strconv.Atoi(args[3+numDrones])
	//			Alt, _ := strconv.Atoi(args[4+numDrones])
	//			modulesImport.RTBGeneral(VehicleId, float32(Velo), float32(Alt)+float32(10*i))
	//		} else if controlNumber == "4" {
	//			modulesImport.RunLandGeneral(VehicleId)
	//		} else {
	//			fmt.Println("fuck off")
	//		}
	//	}
}
