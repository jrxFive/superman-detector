package calculations

import (
	"fmt"
	"math"

	"github.com/umahmood/haversine"
)

const (
	maximumLatitudeDegree  = 90
	maximumLongitudeDegree = 180
)

var (
	ErrInvalidLatitude  = fmt.Errorf("latitude  must be specified as degrees between -%d and %d", maximumLatitudeDegree, maximumLatitudeDegree)
	ErrInvalidLongitude = fmt.Errorf("longitude must be specified as degrees between -%d and %d", maximumLongitudeDegree, maximumLongitudeDegree)
)

// Given a pair of latitude and longitude coordinate return distance of them based on a sphere
func CoordinatesDistance(lat1, lon1, lat2, lon2 float64) (float64, error) {

	if math.Abs(lat1) > maximumLatitudeDegree || math.Abs(lat2) > maximumLatitudeDegree {
		return 0, ErrInvalidLatitude
	}

	if math.Abs(lat1) > maximumLongitudeDegree || math.Abs(lat2) > maximumLongitudeDegree {
		return 0, ErrInvalidLongitude
	}

	c1 := haversine.Coord{Lat: lat1, Lon: lon1}
	c2 := haversine.Coord{Lat: lat2, Lon: lon2}

	miles, _ := haversine.Distance(c1, c2)
	return miles, nil
}
