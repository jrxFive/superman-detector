package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/labstack/echo/v4"
)

type Stats struct {
	statsdClient statsd.ClientInterface
}

func NewStats(statsdClient statsd.ClientInterface) *Stats {
	return &Stats{statsdClient: statsdClient}
}

//Generate statsd metrics for handler response times, counts, percentiles along with tagged status code
func (s *Stats) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		//run eventual handler
		if err := next(c); err != nil {
			c.Error(err)
		}
		elapsed := time.Since(start)

		codeTag := fmt.Sprintf("code:%d", c.Response().Status)
		handlerTag := strings.Replace(fmt.Sprintf("handler:%s.%s", strings.ToLower(c.Request().Method), c.Path()), "/", "", -1)

		_ = s.statsdClient.Incr("http.status.count", []string{codeTag, handlerTag}, 1)
		_ = s.statsdClient.Timing("http.timer", elapsed, []string{handlerTag}, 1)
		return nil
	}
}
