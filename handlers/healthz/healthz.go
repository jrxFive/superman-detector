package healthz

import (
	"net/http"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/labstack/echo/v4"
)

// Use to verify dependencies are responding to simple ACKs
type Pinger interface {
	Ping() error //Returns error if Ping fails
}

type Healthz struct {
	p            Pinger
	statsdClient statsd.ClientInterface
}

func NewHealthz(p Pinger, statsdClient statsd.ClientInterface) Healthz {
	return Healthz{
		p:            p,
		statsdClient: statsdClient,
	}
}

// Verifies API can communicate with given DB. Success/Passing returns 204,
// Failing returns 503.
func (h *Healthz) HeadHealthz(c echo.Context) error {
	if err := h.p.Ping(); err != nil {
		c.Logger().Error(err.Error())
		return echo.NewHTTPError(http.StatusServiceUnavailable)
	}

	return c.NoContent(http.StatusNoContent)
}

// Verifies API can communicate with given DB. Success/Passing returns 200 with a response body,
// Failing returns 503.
func (h *Healthz) GetHealthz(c echo.Context) error {
	if err := h.p.Ping(); err != nil {
		c.Logger().Error(err.Error())
		return echo.NewHTTPError(http.StatusServiceUnavailable)
	}

	return c.String(http.StatusOK, `{"status": "ok"}`)
}
