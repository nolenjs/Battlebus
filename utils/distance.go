package utils

import "math"

func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const r = 6371 // km
	sdLat := math.Sin(Radians(lat2-lat1) / 2)
	sdLon := math.Sin(Radians(lon2-lon1) / 2)
	a := sdLat*sdLat + math.Cos(Radians(lat1))*math.Cos(Radians(lat2))*sdLon*sdLon
	d := 2 * r * math.Asin(math.Sqrt(a))
	return d * 1000 //Distance in Meters
}

func Radians(d float64) float64 {
	return d * math.Pi / 180
}
