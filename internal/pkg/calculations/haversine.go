package calculations

import (
	"fmt"
	"math"
	"time"

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

// Speed in miles per hours based on timestamp deltas and given distance in miles.
func SpeedMilesPerHour(timestamp1, timestamp2 int64, distanceMiles float64) int {
	var d time.Duration
	t1 := time.Unix(timestamp1, 0)
	t2 := time.Unix(timestamp2, 0)

	if t1.Before(t2) {
		d = t2.Sub(t1)
	} else {
		d = t1.Sub(t2)
	}

	if d.Hours() > 0 {
		speed := distanceMiles / d.Hours()
		return int(speed)
	}

	return 0
}
