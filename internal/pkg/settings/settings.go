package settings

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Configuration settings from environment variables
type envSpecification struct {
	ServerReadTimeoutSeconds   time.Duration `default:"5s" split_words:"true"`
	ServerWriteTimeoutSeconds  time.Duration `default:"10s" split_words:"true"`
	ServicePort                int           `default:"8080" split_words:"true"`
	SpeedThresholdMilesPerHour int           `default:"500" split_words:"true"`
	GeoIPDatabaseFileLocation  string        `default:"./GeoLite2-City.mmdb" split_words:"true"`
}

type Specification struct {
	envSpecification
}

func NewSettings() Specification {
	e := envSpecification{}
	envconfig.MustProcess("detector_api", &e)
	return Specification{e}
}