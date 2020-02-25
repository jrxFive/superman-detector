package settings

import (
	"reflect"
	"testing"
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
					ServicePort:                8080,
					SpeedThresholdMilesPerHour: 500,
					GeoIPDatabaseFileLocation:  "./GeoLite2-City.mmdb",
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
