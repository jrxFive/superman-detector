package v1

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/jrxfive/superman-detector/internal/pkg/settings"
	"github.com/jrxfive/superman-detector/pkg/geoip"
	"github.com/jrxfive/superman-detector/schemas"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// If these tests look familiar it is because they are. Due to amount of time I had and did see previous implementations on github.
// I replicated one particular test suite and figured this would be a legitimately good test to verify I had the same expected functionality.
// I hope this does not come off as dishonest or against the spirit of this exercise.
// https://github.com/NEPDAVE/supermanDetector/blob/master/http_test.go
// Unlike the above link TravelFromCurrentGeoSuspicious, TravelFromCurrentGeoSuspicious, PrecedingIpAccess, SubsequentIpAccess
// will only appear in the response if it exists. Along with speed is in miles per hour.
func TestLogin_PostLogin(t *testing.T) {
	e := echo.New()
	defer e.Close()

	s := settings.NewSettings()
	s.GeoIPDatabaseFileLocation = "../../GeoLite2-City.mmdb"

	statsdClient := &statsd.NoOpClient{}

	db, err := gorm.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	defer db.Close()

	geoDB, err := geoip.NewDefaultLocator(s)
	assert.NoError(t, err)
	defer geoDB.Close()

	l := NewLogin(db, geoDB, statsdClient, s)

	tests := []struct {
		name             string
		request          schemas.LoginEvent
		expectedResponse schemas.ProcessedLoginEvent
		statusCode       int
	}{
		{
			name: "'Bob' New York IP Access at 1514764001",
			request: schemas.LoginEvent{
				Username:  "bob",
				TimeStamp: 1514764001,
				EventUuid: "2850064a-5787-11ea-a0e4-c4b301c8961b",
				IPAddress: schemas.IPAddress(net.ParseIP("206.81.252.6")),
			},
			expectedResponse: schemas.ProcessedLoginEvent{
				CurrentGeo: schemas.Geo{
					Lat:    39.211,
					Lon:    -76.8362,
					Radius: 5,
				},
			},
		},
		{
			name: "'Bob' Hong Kong IP Access 1514764000",
			request: schemas.LoginEvent{
				Username:  "bob",
				TimeStamp: 1514764000,
				EventUuid: "85ad929a-db03-4bf4-9541-8f728fa12e40",
				IPAddress: schemas.IPAddress(net.ParseIP("119.28.48.231")),
			},
			expectedResponse: schemas.ProcessedLoginEvent{
				CurrentGeo: schemas.Geo{
					Lat:    22.25,
					Lon:    114.1667,
					Radius: 50,
				},
				TravelFromCurrentGeoSuspicious: true,
				SubsequentIpAccess: &schemas.IPAccess{
					Ip:        "206.81.252.6",
					Speed:     29266323,
					Timestamp: 1514764001,
					Geo: schemas.Geo{
						Lat:    39.211,
						Lon:    -76.8362,
						Radius: 5,
					},
				},
			},
		},
		{
			name: "'Bob' Moscow IP Access at 1514764002",
			request: schemas.LoginEvent{
				Username:  "bob",
				TimeStamp: 1514764002,
				EventUuid: "85ad929a-db03-4bf4-9541-8f728fa12e42",
				IPAddress: schemas.IPAddress(net.ParseIP("31.173.221.5")),
			},
			expectedResponse: schemas.ProcessedLoginEvent{
				CurrentGeo: schemas.Geo{
					Lat:    42.9753,
					Lon:    47.5022,
					Radius: 1000,
				},
				TravelToCurrentGeoSuspicious: true,
				PrecedingIpAccess: &schemas.IPAccess{
					Ip:        "206.81.252.6",
					Speed:     20794646,
					Timestamp: 1514764001,
					Geo: schemas.Geo{
						Lat:    39.211,
						Lon:    -76.8362,
						Radius: 5,
					},
				},
			},
		},
		{
			name: "'Bob' Sydney IP Access at 1514764006",
			request: schemas.LoginEvent{
				Username:  "bob",
				TimeStamp: 1514764006,
				EventUuid: "85ad929a-db03-4bf4-9541-8f728fa12e46",
				IPAddress: schemas.IPAddress(net.ParseIP("203.2.218.214")),
			},
			expectedResponse: schemas.ProcessedLoginEvent{
				CurrentGeo: schemas.Geo{
					Lat:    -33.8919,
					Lon:    151.1554,
					Radius: 500,
				},
				TravelToCurrentGeoSuspicious: true,
				PrecedingIpAccess: &schemas.IPAccess{
					Ip:        "31.173.221.5",
					Speed:     7558027,
					Timestamp: 1514764002,
					Geo: schemas.Geo{
						Lat:    42.9753,
						Lon:    47.5022,
						Radius: 1000,
					},
				},
			},
		},
		{
			name: "'Bob' New York IP Access at 1514764005",
			request: schemas.LoginEvent{
				Username:  "bob",
				TimeStamp: 1514764005,
				EventUuid: "85ad929a-db03-4bf4-9541-8f728fa12e45",
				IPAddress: schemas.IPAddress(net.ParseIP("206.81.252.7")),
			},
			expectedResponse: schemas.ProcessedLoginEvent{
				CurrentGeo: schemas.Geo{
					Lat:    39.211,
					Lon:    -76.8362,
					Radius: 5,
				},
				TravelToCurrentGeoSuspicious:   true,
				TravelFromCurrentGeoSuspicious: true,
				PrecedingIpAccess: &schemas.IPAccess{
					Ip:        "31.173.221.5",
					Speed:     6931548,
					Timestamp: 1514764002,
					Geo: schemas.Geo{
						Lat:    42.9753,
						Lon:    47.5022,
						Radius: 1000,
					},
				},
				SubsequentIpAccess: &schemas.IPAccess{
					Ip:        "203.2.218.214",
					Speed:     35197415,
					Timestamp: 1514764006,
					Geo: schemas.Geo{
						Lat:    -33.8919,
						Lon:    151.1554,
						Radius: 500,
					},
				},
			},
		},
		{
			name: "'Bob' Hong Kong IP Access at 1514764003",
			request: schemas.LoginEvent{
				Username:  "bob",
				TimeStamp: 1514764003,
				EventUuid: "85ad929a-db03-4bf4-9541-8f728fa12e43",
				IPAddress: schemas.IPAddress(net.ParseIP("119.28.48.231")),
			},
			expectedResponse: schemas.ProcessedLoginEvent{
				CurrentGeo: schemas.Geo{
					Lat:    22.25,
					Lon:    114.1667,
					Radius: 50,
				},
				TravelToCurrentGeoSuspicious:   true,
				TravelFromCurrentGeoSuspicious: true,
				PrecedingIpAccess: &schemas.IPAccess{
					Ip:        "31.173.221.5",
					Speed:     14483739,
					Timestamp: 1514764002,
					Geo: schemas.Geo{
						Lat:    42.9753,
						Lon:    47.5022,
						Radius: 1000,
					},
				},
				SubsequentIpAccess: &schemas.IPAccess{
					Ip:        "206.81.252.7",
					Speed:     14633161,
					Timestamp: 1514764005,
					Geo: schemas.Geo{
						Lat:    39.211,
						Lon:    -76.8362,
						Radius: 5,
					},
				},
			},
		},
		{
			name: "'Bob' Sydney IP Access at 1514764007",
			request: schemas.LoginEvent{
				Username:  "bob",
				TimeStamp: 1514764007,
				EventUuid: "85ad929a-db03-4bf4-9541-8f728fa12e47",
				IPAddress: schemas.IPAddress(net.ParseIP("203.2.218.214")),
			},
			expectedResponse: schemas.ProcessedLoginEvent{
				CurrentGeo: schemas.Geo{
					Lat:    -33.8919,
					Lon:    151.1554,
					Radius: 500,
				},
				TravelToCurrentGeoSuspicious:   false,
				TravelFromCurrentGeoSuspicious: false,
				PrecedingIpAccess: &schemas.IPAccess{
					Ip:        "203.2.218.214",
					Speed:     0,
					Timestamp: 1514764006,
					Geo: schemas.Geo{
						Lat:    -33.8919,
						Lon:    151.1554,
						Radius: 500,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			requestBody, err := json.Marshal(tt.request)
			assert.NoError(t, err)

			expectedResponse, err := json.Marshal(tt.expectedResponse)
			assert.NoError(t, err)

			expectedResponseString := string(expectedResponse) + "\n"

			req := httptest.NewRequest(http.MethodGet, "/", bytes.NewBuffer(requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if assert.NoError(t, l.PostLogin(c)) {
				assert.Equal(t, http.StatusCreated, rec.Code)
				assert.Equal(t, expectedResponseString, rec.Body.String())
			}
		})
	}
}
