package settings

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestNewSettings(t *testing.T) {

	tests := []struct {
		name string
		want Specification
	}{
		{
			name: "defaults",
			want: Specification{
				envSpecification{
					ServerReadTimeoutSeconds:   5 * time.Second,
					ServerWriteTimeoutSeconds:  10 * time.Second,
					ServicePort:                8080,
					SpeedThresholdMilesPerHour: 500,
					GeoIPDatabaseFileLocation:  "./GeoLite2-City.mmdb",
					StatsdNamespace:            "superman-detector",
					StatsdAddress:              "localhost:8125",
					StatsdBufferPoolSize:       1000,
					SqlDialect:                 "sqlite3",
					SqlConnectionString:        "/tmp/superman.db",
					RequestBodyLimit:           "1M",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSettings(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSettings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSettingsOverride(t *testing.T) {
	os.Setenv("DETECTOR_API_STATSD_ADDRESS", "telegraf:8125")

	tests := []struct {
		name string
		want Specification
	}{
		{
			name: "defaults",
			want: Specification{
				envSpecification{
					ServerReadTimeoutSeconds:   5 * time.Second,
					ServerWriteTimeoutSeconds:  10 * time.Second,
					ServicePort:                8080,
					SpeedThresholdMilesPerHour: 500,
					GeoIPDatabaseFileLocation:  "./GeoLite2-City.mmdb",
					StatsdNamespace:            "superman-detector",
					StatsdAddress:              "telegraf:8125",
					StatsdBufferPoolSize:       1000,
					SqlDialect:                 "sqlite3",
					SqlConnectionString:        "/tmp/superman.db",
					RequestBodyLimit:           "1M",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSettings(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSettings() = %v, want %v", got, tt.want)
			}
		})
	}
}
