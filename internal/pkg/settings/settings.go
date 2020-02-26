package settings

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Configuration settings from environment variables
type envSpecification struct {
	RequestBodyLimit           string        `default:"1M"` //https://echo.labstack.com/middleware/body-limit
	ServerReadTimeoutSeconds   time.Duration `default:"5s" split_words:"true"`
	ServerWriteTimeoutSeconds  time.Duration `default:"10s" split_words:"true"`
	ServicePort                int           `default:"8080" split_words:"true"`
	SpeedThresholdMilesPerHour int           `default:"500" split_words:"true"`
	GeoIPDatabaseFileLocation  string        `default:"./GeoLite2-City.mmdb" split_words:"true"`
	StatsdNamespace            string        `default:"superman-detector" split_words:"true"`
	StatsdAddress              string        `default:"localhost:8125" split_words:"true"`
	StatsdBufferPoolSize       int           `default:"1000" split_words:"true"`
	SqlDialect                 string        `default:"sqlite3" split_words:"true"`
	SqlConnectionString        string        `default:"/tmp/superman.db" split_words:"true"`
}

type Specification struct {
	envSpecification
}

func NewSettings() Specification {
	e := envSpecification{}
	envconfig.MustProcess("detector_api", &e)
	return Specification{e}
}
