package util

import "math"

// Cr. https://www.geodatasource.com/developers/go
// Example:
//fmt.Printf("%f Miles\n", distance(32.9697, -96.80322, 29.46786, -98.53506, "M"))
//fmt.Printf("%f Kilometers\n", distance(32.9697, -96.80322, 29.46786, -98.53506, "K"))
//fmt.Printf("%f Nautical Miles\n", distance(32.9697, -96.80322, 29.46786, -98.53506, "N"))
//:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
//:::                                                                         :::
//:::  This routine calculates the distance between two points (given the     :::
//:::  latitude/longitude of those points). It is being used to calculate     :::
//:::  the distance between two locations using GeoDataSource (TM) prodducts  :::
//:::                                                                         :::
//:::  Definitions:                                                           :::
//:::    South latitudes are negative, east longitudes are positive           :::
//:::                                                                         :::
//:::  Passed to function:                                                    :::
//:::    lat1, lon1 = Latitude and Longitude of point 1 (in decimal degrees)  :::
//:::    lat2, lon2 = Latitude and Longitude of point 2 (in decimal degrees)  :::
//:::    unit = the unit you desire for results                               :::
//:::           where: 'M' is statute miles (default)                         :::
//:::                  'K' is kilometers                                      :::
//:::                  'N' is nautical miles                                  :::
//:::                                                                         :::
//:::  Worldwide cities and other features databases with latitude longitude  :::
//:::  are available at https://www.geodatasource.com                         :::
//:::                                                                         :::
//:::  For enquiries, please contact sales@geodatasource.com                  :::
//:::                                                                         :::
//:::  Official Web site: https://www.geodatasource.com                       :::
//:::                                                                         :::
//:::               GeoDataSource.com (C) All Rights Reserved 2018            :::
//:::                                                                         :::
//:::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
func CoordinatorDistance(lat1 float64, lng1 float64, lat2 float64, lng2 float64, unit ...string) float64 {
	const PI float64 = 3.141592653589793
	radlat1 := PI * lat1 / 180
	radlat2 := PI * lat2 / 180
	theta := lng1 - lng2
	radtheta := PI * theta / 180
	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)
	if dist > 1 {
		dist = 1
	}
	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515
	if len(unit) > 0 {
		if unit[0] == "K" {
			dist = dist * 1.609344
		} else if unit[0] == "N" {
			dist = dist * 0.8684
		}
	}
	return dist
}

// https://gist.github.com/cdipaolo/d3f8db3848278b49db68
//// haversin(θ) function
//func hsin(theta float64) float64 {
//	return math.Pow(math.Sin(theta/2), 2)
//}
//
//// Distance function returns the distance (in meters) between two points of
////     a given longitude and latitude relatively accurately (using a spherical
////     approximation of the Earth) through the Haversin Distance Formula for
////     great arc distance on a sphere with accuracy for small distances
////
//// point coordinates are supplied in degrees and converted into rad. in the func
////
//// distance returned is METERS!!!!!!
//// http://en.wikipedia.org/wiki/Haversine_formula
//func Distance(lat1, lon1, lat2, lon2 float64) float64 {
//	// convert to radians
//	// must cast radius as float to multiply later
//	var la1, lo1, la2, lo2, r float64
//	la1 = lat1 * math.Pi / 180
//	lo1 = lon1 * math.Pi / 180
//	la2 = lat2 * math.Pi / 180
//	lo2 = lon2 * math.Pi / 180
//
//	r = 6378100 // Earth radius in METERS
//
//	// calculate
//	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)
//
//	return 2 * r * math.Asin(math.Sqrt(h))
//}
