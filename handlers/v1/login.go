package v1

import (
	"net/http"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	"github.com/jrxfive/superman-detector/internal/pkg/settings"
	"github.com/jrxfive/superman-detector/models"
	"github.com/jrxfive/superman-detector/pkg/calculations"
	"github.com/jrxfive/superman-detector/pkg/geoip"
	"github.com/jrxfive/superman-detector/schemas"
	"github.com/labstack/echo/v4"
)

var (
	loginPrecedingQuery      = "username = ? and timestamp < ?"
	loginPrecedingQueryOrder = "timestamp desc"

	loginSubsequentQuery      = "username = ? and timestamp > ?"
	loginSubsequentQueryOrder = "timestamp"
)

type Login struct {
	db           *gorm.DB
	geoDB        geoip.Locator
	settings     settings.Specification
	statsdClient statsd.ClientInterface
}

func NewLogin(db *gorm.DB, geoDB geoip.Locator, statsdClient statsd.ClientInterface, settings settings.Specification) Login {
	if !db.HasTable(&models.LoginEvent{}) {
		db.CreateTable(&models.LoginEvent{})
	}

	return Login{
		db:           db,
		geoDB:        geoDB,
		settings:     settings,
		statsdClient: statsdClient,
	}
}

func (l *Login) suspiciousTravel(speed int) bool {
	return speed > l.settings.SpeedThresholdMilesPerHour
}

// Given the current POST'd event and schema.IPAccess assign values if a record is found from the query results. This
// probably could have used reflection but would have certainly increased the complexity of this function.
func (l *Login) setIPAccessRequestRecord(currentEvent models.LoginEvent, whereQueryClause, orderQueryClause string, responseIPAccess *schemas.IPAccess) bool {
	queriedAccessEvent := &models.LoginEvent{}

	ipAccessRequest := l.db.
		Where(whereQueryClause, currentEvent.Username, currentEvent.Timestamp).
		Order(orderQueryClause).
		Limit(1).
		Find(queriedAccessEvent)

	if !ipAccessRequest.RecordNotFound() {
		distanceMiles, err := calculations.
			CoordinatesDistance(currentEvent.Latitude, currentEvent.Longitude, queriedAccessEvent.Latitude, queriedAccessEvent.Longitude)

		if err != nil {
			return false
		}

		speed := calculations.SpeedMilesPerHour(currentEvent.Timestamp, queriedAccessEvent.Timestamp, distanceMiles)

		responseIPAccess.Ip = queriedAccessEvent.IPAddress
		responseIPAccess.Speed = speed
		responseIPAccess.Geo = schemas.Geo{
			Lat:    queriedAccessEvent.Latitude,
			Lon:    queriedAccessEvent.Longitude,
			Radius: queriedAccessEvent.Radius,
		}
		responseIPAccess.Timestamp = queriedAccessEvent.Timestamp
		return true
	}

	return false
}

func (l *Login) PostLogin(c echo.Context) error {
	loginEvent := &schemas.LoginEvent{}

	if err := c.Bind(loginEvent); err != nil {
		c.Logger().Error("failed to bind request body to struct")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	_, err := govalidator.ValidateStruct(loginEvent)
	if err != nil {
		c.Logger().Error("error validating struct")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	requestGeoDetails, err := l.geoDB.Locate(loginEvent.IPAddress.IP())
	if err != nil {
		c.Logger().Errorf("geo ip:%s lookup failed", loginEvent.IPAddress.IP().String())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	currentIPAccessEvent := &models.LoginEvent{
		Username:  loginEvent.Username,
		Timestamp: loginEvent.TimeStamp,
		EventUuid: loginEvent.EventUuid,
		IPAddress: loginEvent.IPAddress.IP().String(),
		Latitude:  requestGeoDetails.Latitude,
		Longitude: requestGeoDetails.Longitude,
		Radius:    requestGeoDetails.Radius,
	}

	result := l.db.Create(currentIPAccessEvent)
	if result.Error != nil {
		c.Logger().Error("failed to insert model struct")
		return echo.NewHTTPError(http.StatusServiceUnavailable)
	}

	response := schemas.ProcessedLoginEvent{
		CurrentGeo: schemas.Geo{
			Lat:    requestGeoDetails.Latitude,
			Lon:    requestGeoDetails.Longitude,
			Radius: requestGeoDetails.Radius,
		},
		TravelToCurrentGeoSuspicious:   false,
		TravelFromCurrentGeoSuspicious: false,
		PrecedingIpAccess:              &schemas.IPAccess{},
		SubsequentIpAccess:             &schemas.IPAccess{},
	}

	//If no record was found set struct value to nil so it does not appear in response
	if updated := l.setIPAccessRequestRecord(*currentIPAccessEvent, loginPrecedingQuery, loginPrecedingQueryOrder, response.PrecedingIpAccess); !updated {
		response.PrecedingIpAccess = nil
	} else {
		response.TravelToCurrentGeoSuspicious = l.suspiciousTravel(response.PrecedingIpAccess.Speed)
	}

	if updated := l.setIPAccessRequestRecord(*currentIPAccessEvent, loginSubsequentQuery, loginSubsequentQueryOrder, response.SubsequentIpAccess); !updated {
		response.SubsequentIpAccess = nil
	} else {
		response.TravelFromCurrentGeoSuspicious = l.suspiciousTravel(response.SubsequentIpAccess.Speed)
	}

	return c.JSON(http.StatusCreated, response)
}
