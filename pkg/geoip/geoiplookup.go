package geoip

import (
	"net"

	"github.com/jrxfive/superman-detector/internal/pkg/settings"
	"github.com/oschwald/geoip2-golang"
)

//Locate based on IP Address
type Locator interface {
	Locate(ip net.IP) (*Location, error)
	Close() error
}

type Location struct {
	Latitude  float64 `maxminddb:"latitude" json:"lat"`
	Longitude float64 `maxminddb:"longitude" json:"lon"`
	Radius    uint16  `maxminddb:"accuracy_radius" json:"radius"`
}

type DefaultLocator struct {
	db *geoip2.Reader
}

// New Default locator based on github.com/oschwald/geoip2-golang requires
// local copy of GeoLite2-City.mmdb
func NewDefaultLocator(s settings.Specification) (Locator, error) {
	geoDB, err := geoip2.Open(s.GeoIPDatabaseFileLocation)
	if err != nil {
		return nil, err
	}

	return &DefaultLocator{db: geoDB}, nil
}

// Locate latitude, longitude
func (d *DefaultLocator) Locate(ip net.IP) (*Location, error) {
	requestGeoDetails, err := d.db.City(ip)
	if err != nil {
		return nil, err
	}

	return &Location{
		Radius:    requestGeoDetails.Location.AccuracyRadius,
		Latitude:  requestGeoDetails.Location.Latitude,
		Longitude: requestGeoDetails.Location.Longitude,
	}, nil

}

// Close DB
func (d *DefaultLocator) Close() error {
	return d.db.Close()
}
